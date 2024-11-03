package nemesis

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
	"github.com/huanghj78/jepsenFuzz/pkg/core"
	"github.com/huanghj78/jepsenFuzz/util"
	"github.com/ngaut/log"
)

type killGenerator struct {
	name string
}

func (g killGenerator) Generate(nodes []cluster.Node) []*core.NemesisOperation {
	log.Info("=========================")
	var n int
	var duration = time.Second * time.Duration(rand.Intn(120)+60)
	switch g.name {
	case "minor_kill":
		n = len(nodes)/2 - 1
	case "major_kill":
		n = len(nodes)/2 + 1
	default:
		n = 1
	}
	log.Info(g.name)
	return killNodes(nodes, n, duration)
}

func (g killGenerator) Name() string {
	return g.name
}

func killNodes(nodes []cluster.Node, n int, duration time.Duration) []*core.NemesisOperation {
	var ops []*core.NemesisOperation
	indices := shuffleIndices(len(nodes))
	if n > len(indices) {
		n = len(indices)
	}
	for i := 0; i < n; i++ {
		ops = append(ops, &core.NemesisOperation{
			Type:        core.ProcKill,
			Node:        &nodes[indices[i]],
			InvokeArgs:  nil,
			RecoverArgs: nil,
			RunTime:     duration,
		})
	}
	log.Info(ops)
	return ops
}

// NewKillGenerator creates a generator.
// Name is random_kill, minor_kill, major_kill, and all_kill.
func NewKillGenerator(name string) core.NemesisGenerator {
	return killGenerator{name: name}
}

// kill implements Nemesis
type kill struct {
}

func (k kill) Invoke(ctx context.Context, node *cluster.Node, _ ...interface{}) error {
	log.Infof("apply nemesis %s on node %s", core.ProcKill, node.IP)
	cmd := fmt.Sprintf("pkill -f gaussdb")
	output, err := util.ExecuteRemoteCommand(node.IP, "omm", "ilovedds", cmd)
	log.Info("======================================================")
	log.Info(output)
	return err
}

func (k kill) Recover(ctx context.Context, node *cluster.Node, _ ...interface{}) error {
	log.Infof("unapply nemesis %s on node %s", core.ProcKill, node.IP)
	cmd := fmt.Sprintf("cm_ctl start -n %d", node.ID)
	output, err := util.ExecuteRemoteCommand(node.IP, "omm", "ilovedds", cmd)
	log.Info(output)
	return err
}

func (kill) Name() string {
	return string(core.ProcKill)
}
