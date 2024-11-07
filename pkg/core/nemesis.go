package core

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/atomic"

	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
)

// ChaosKind is the kind of applying chaos
type ChaosKind string

const (
	ProcKill ChaosKind = "proc-kill"
	// PodFailure Applies pod failure
	PodFailure ChaosKind = "pod-failure"
	// PodKill will random kill a pod, this will make the Node be illegal
	PodKill ChaosKind = "pod-kill"
	// ContainerKill will random kill the specified container of pod, but retain the pod
	ContainerKill ChaosKind = "container-kill"
	// NetworkPartition partitions network between nodes
	NetworkPartition ChaosKind = "network-partition"
	// NetemChaos adds corrupt or other chaos.
	NetemChaos ChaosKind = "netem-chaos"
	// TimeChaos means
	TimeChaos ChaosKind = "time-chaos"
	// PDScheduler adds scheduler
	PDScheduler ChaosKind = "pd-scheduler"
	// PDLeaderShuffler will randomly shuffle pds.
	PDLeaderShuffler ChaosKind = "pd-leader-shuffler"
	// Scaling scales cluster
	Scaling     ChaosKind = "scaling"
	CPUFullload ChaosKind = "cpu-fullload"
	DiskBurn    ChaosKind = "disk-burn"
	DiskFill    ChaosKind = "disk-fill"
)

// Nemesis injects failure and disturbs the database.
type Nemesis interface {
	// // SetUp initializes the nemesis
	// SetUp(ctx context.Context, node string) error
	// // TearDown tears down the nemesis
	// TearDown(ctx context.Context, node string) error

	// Invoke executes the nemesis
	Invoke(ctx context.Context, node *cluster.Node, args ...interface{}) error
	// Recover recovers the nemesis
	Recover(ctx context.Context, node *cluster.Node, args ...interface{}) error
	// Name returns the unique name for the nemesis
	Name() string
}

var nemesises = map[string]Nemesis{}

// RegisterNemesis registers nemesis. Not thread-safe.
func RegisterNemesis(n Nemesis) {
	name := n.Name()
	_, ok := nemesises[name]
	if ok {
		panic(fmt.Sprintf("nemesis %s is already registered", name))
	}

	nemesises[name] = n
}

// GetNemesis gets the registered nemesis.
func GetNemesis(name string) Nemesis {
	fmt.Println(nemesises)
	return nemesises[name]
}

const (
	waitForStart   = 1
	enableStart    = 2
	starting       = 3
	enableRollback = 4
)

// NemesisControl is used to operate nemesis between the control side and test client side
type NemesisControl struct {
	// 0: wait for start
	// 1: enable start
	// 2: starting
	// 3: enable rollback
	s atomic.Int32
}

// WaitForStart is used on control side to wait for enabling start nemesis
func (n *NemesisControl) WaitForStart() {
	// init state
	n.s.Store(waitForStart)
	for {
		if n.s.CAS(enableStart, starting) {
			return
		}
		time.Sleep(time.Second)
	}
}

// WaitForRollback is used on control side to wait for enabling rollback nemesis
func (n *NemesisControl) WaitForRollback(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if n.s.Load() == enableRollback {
				return
			}
			time.Sleep(time.Second)
		}
	}
}

// Start is used on client side to enable control side starting nemesis
func (n *NemesisControl) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if n.s.CAS(waitForStart, enableStart) {
				return
			}
			time.Sleep(time.Second)
		}
	}
}

// Rollback is used on client side to enable control side rollbacking nemesis
func (n *NemesisControl) Rollback(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if n.s.CAS(starting, enableRollback) {
				return
			}
			time.Sleep(time.Second)
		}
	}
}

// NemesisOperation is nemesis operation used in control.
type NemesisOperation struct {
	Type        ChaosKind     // Nemesis name
	Node        *cluster.Node // Nemesis target node, optional if it affects
	InvokeArgs  []interface{} // Nemesis invoke args
	RecoverArgs []interface{} // Nemesis recover args

	// We have two approaches to trigger recovery
	// 1. through `RunTime`
	// 2. through `NemesisControl` WaitForRollback
	RunTime        time.Duration   // Nemesis duration time
	NemesisControl *NemesisControl // Nemesis recovery signal
}

// String ...
func (n NemesisOperation) String() string {
	return fmt.Sprintf("type:%s,duration:%s,node:%s,invoke_args:%+v,recover_args:%+v", n.Type, n.RunTime, n.Node, n.InvokeArgs, n.RecoverArgs)
}

// NemesisGeneratorRecord is used to record operations generated by NemesisGenerator.Generate
type NemesisGeneratorRecord struct {
	Name string
	Ops  []*NemesisOperation
}

// NemesisGenerator is used in control, it will generate a nemesis operation
// and then the control can use it to disturb the cluster.
type NemesisGenerator interface {
	// Generate generates the nemesis operation for all nodes.
	// Every node will be assigned a nemesis operation.
	Generate(nodes []cluster.Node) []*NemesisOperation
	Name() string
}

// DelayNemesisGenerator delays nemesis generation after `Delay`
type DelayNemesisGenerator struct {
	Gen   NemesisGenerator
	Delay time.Duration
}

// Generate ...
func (d DelayNemesisGenerator) Generate(nodes []cluster.Node) []*NemesisOperation {
	time.Sleep(d.Delay)
	return d.Gen.Generate(nodes)
}

// Name ...
func (d DelayNemesisGenerator) Name() string {
	return d.Gen.Name()
}

// NemesisGenerators is a NemesisGenerator iterator
type NemesisGenerators interface {
	Next() NemesisGenerator
	HasNext() bool
	// Reset resets iterator, return false if it cannot be reset
	Reset() bool
}

// nemesisGenerators is a wrapper of []NemesisGenerator
type nemesisGenerators struct {
	idx        int
	generators []NemesisGenerator
}

func (i *nemesisGenerators) HasNext() bool {
	return i.idx < len(i.generators)
}

// Next ...
func (i *nemesisGenerators) Next() NemesisGenerator {
	gen := i.generators[i.idx]
	i.idx += 1
	return gen
}

// Reset reset if it could
func (i *nemesisGenerators) Reset() bool {
	i.idx = 0
	return true
}

// NewNemesisGenerators ...
func NewNemesisGenerators(gens []NemesisGenerator) NemesisGenerators {
	return &nemesisGenerators{
		idx:        0,
		generators: gens,
	}
}

// OneRoundNemesisGenerators is easier than nemesisGenerators, and suitable in cases that need to interact between client and control
type OneRoundNemesisGenerators struct {
	gen     NemesisGenerator
	hasNext bool
}

// HasNext ...
func (m *OneRoundNemesisGenerators) HasNext() bool {
	return m.hasNext
}

// Next ...
func (m *OneRoundNemesisGenerators) Next() NemesisGenerator {
	m.hasNext = false
	return m.gen
}

// Reset just returns false because we forbid reset for OneRoundNemesisGenerators
func (m *OneRoundNemesisGenerators) Reset() bool {
	return false
}

// NewOneRoundNemesisGenerators ...
func NewOneRoundNemesisGenerators(gen NemesisGenerator) NemesisGenerators {
	return &OneRoundNemesisGenerators{
		gen:     gen,
		hasNext: true,
	}
}
