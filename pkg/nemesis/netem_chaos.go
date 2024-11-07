package nemesis

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ngaut/log"

	// chaosv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/huanghj78/jepsenFuzz/util"

	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
	"github.com/huanghj78/jepsenFuzz/pkg/core"
)

type netemChaosGenerator struct {
	name string
}

// NewNetemChaos create a netem chaos.
func NewNetemChaos(name string) core.NemesisGenerator {
	return netemChaosGenerator{name: name}
}

// network loss
type loss struct {
	FaultIdMap map[string]string
}

func (l loss) apply(node *cluster.Node, args ...string) error {
	if len(args) != 2 {
		panic("args number error")
	}
	netInterface := args[0]
	lossPercent := args[1]

	cmd := fmt.Sprintf("blade create network loss --exclude-port 22 --interface %s --percent %s", netInterface, lossPercent)
	// 这里需要vim ~/.bashrc，添加export PATH="$PATH:/usr/local/bin/chaosblade-1.7.4"（对应的位置）并source ~/.bashrc
	// 请一定要排除22端口，否则sshpass难以发送指令使其复原，也会使机器难以控制
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "centos", cmd)
	if err != nil {
		log.Errorf("Error executing remote command: %v", err)
	}

	var result map[string]interface{}
	jsonOutput := strings.TrimSpace(output)
	err = json.Unmarshal([]byte(jsonOutput), &result)
	if err != nil {
		log.Error(jsonOutput)
		log.Errorf("Error unmarshalling JSON: %v", err)
	}
	l.FaultIdMap[node.IP], _ = result["result"].(string)

	log.Info("======================================================")
	log.Info(output)

	return err
}

func (l loss) defaultArgsChaosApply(node *cluster.Node) error {
	return l.apply(node, "eth0", "60")
}

func (l loss) chaosDestroy(node *cluster.Node) error {
	log.Infof("unapply nemesis %s on node %s", core.TimeChaos, node.IP)

	faultId := l.FaultIdMap[node.IP]
	delete(l.FaultIdMap, faultId)

	cmd := fmt.Sprintf("blade destroy %s", faultId)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "centos", cmd)

	log.Info("======================================================")
	log.Info(output)

	return err
}

type delay struct {
	FaultIdMap map[string]string
}

func (d delay) apply(node *cluster.Node, args ...string) error {
	if len(args) != 3 {
		panic("args number error")
	}
	netInterface := args[0]
	time := args[1]   // 延迟时间
	offset := args[2] // 延迟上下浮动的时间

	cmd := fmt.Sprintf("blade create network delay --exclude-port 22 --interface %s --time %s --offset %s", netInterface, time, offset)
	// 这里需要vim ~/.bashrc，添加export PATH="$PATH:/usr/local/bin/chaosblade-1.7.4"（对应的位置）并source ~/.bashrc
	// 请一定要排除22端口，否则sshpass难以发送指令使其复原，也会使机器难以控制
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "centos", cmd)
	if err != nil {
		log.Errorf("Error executing remote command: %v", err)
	}

	var result map[string]interface{}
	jsonOutput := strings.TrimSpace(output)
	err = json.Unmarshal([]byte(jsonOutput), &result)
	if err != nil {
		log.Error(jsonOutput)
		log.Errorf("Error unmarshalling JSON: %v", err)
	}
	d.FaultIdMap[node.IP], _ = result["result"].(string)

	log.Info("======================================================")
	log.Info(output)

	return err
}

func (d delay) defaultArgsChaosApply(node *cluster.Node) error {
	return d.apply(node, "eth0", "1000", "500")
}

func (d delay) chaosDestroy(node *cluster.Node) error {
	log.Infof("unapply nemesis %s on node %s", core.TimeChaos, node.IP)

	faultId := d.FaultIdMap[node.IP]
	delete(d.FaultIdMap, faultId)

	cmd := fmt.Sprintf("blade destroy %s", faultId)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "centos", cmd)

	log.Info("======================================================")
	log.Info(output)

	return err
}

type duplicate struct {
	FaultIdMap map[string]string
}

func (d duplicate) apply(node *cluster.Node, args ...string) error {
	if len(args) != 2 {
		panic("args number error")
	}
	netInterface := args[0]
	percent := args[1]

	cmd := fmt.Sprintf("blade create network duplicate --exclude-port 22 --interface %s --percent %s", netInterface, percent)
	// 这里需要vim ~/.bashrc，添加export PATH="$PATH:/usr/local/bin/chaosblade-1.7.4"（对应的位置）并source ~/.bashrc
	// 请一定要排除22端口，否则sshpass难以发送指令使其复原，也会使机器难以控制
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "centos", cmd)
	if err != nil {
		log.Errorf("Error executing remote command: %v", err)
	}

	var result map[string]interface{}
	jsonOutput := strings.TrimSpace(output)
	err = json.Unmarshal([]byte(jsonOutput), &result)
	if err != nil {
		log.Error(jsonOutput)
		log.Errorf("Error unmarshalling JSON: %v", err)
	}
	d.FaultIdMap[node.IP], _ = result["result"].(string)

	log.Info("======================================================")
	log.Info(output)

	return err
}

func (d duplicate) defaultArgsChaosApply(node *cluster.Node) error {
	return d.apply(node, "eth0", "80")
}

func (d duplicate) chaosDestroy(node *cluster.Node) error {
	log.Infof("unapply nemesis %s on node %s", core.TimeChaos, node.IP)

	faultId := d.FaultIdMap[node.IP]
	delete(d.FaultIdMap, faultId)

	cmd := fmt.Sprintf("blade destroy %s", faultId)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "centos", cmd)

	log.Info("======================================================")
	log.Info(output)

	return err
}

type corrupt struct {
	FaultIdMap map[string]string
}

func (c corrupt) apply(node *cluster.Node, args ...string) error {
	if len(args) != 2 {
		panic("args number error")
	}
	netInterface := args[0]
	percent := args[1]

	cmd := fmt.Sprintf("blade create network duplicate --exclude-port 22 --interface %s --percent %s", netInterface, percent)
	// 这里需要vim ~/.bashrc，添加export PATH="$PATH:/usr/local/bin/chaosblade-1.7.4"（对应的位置）并source ~/.bashrc
	// 请一定要排除22端口，否则sshpass难以发送指令使其复原，也会使机器难以控制
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "centos", cmd)
	if err != nil {
		log.Errorf("Error executing remote command: %v", err)
	}

	var result map[string]interface{}
	jsonOutput := strings.TrimSpace(output)
	err = json.Unmarshal([]byte(jsonOutput), &result)
	if err != nil {
		log.Error(jsonOutput)
		log.Errorf("Error unmarshalling JSON: %v", err)
	}
	c.FaultIdMap[node.IP], _ = result["result"].(string)

	log.Info("======================================================")
	log.Info(output)

	return err
}

func (c corrupt) defaultArgsChaosApply(node *cluster.Node) error {
	return c.apply(node, "eth0", "80")
}

func (c corrupt) chaosDestroy(node *cluster.Node) error {
	log.Infof("unapply nemesis %s on node %s", core.TimeChaos, node.IP)

	faultId := c.FaultIdMap[node.IP]
	delete(c.FaultIdMap, faultId)

	cmd := fmt.Sprintf("blade destroy %s", faultId)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "centos", cmd)

	log.Info("======================================================")
	log.Info(output)

	return err
}

type netemChaos interface {
	apply(node *cluster.Node, args ...string) error
	defaultArgsChaosApply(node *cluster.Node) error
	chaosDestroy(node *cluster.Node) error
}

func selectNetem(name string) netemChaos {
	switch name {
	case "loss":
		return loss{FaultIdMap: make(map[string]string)}
	case "delay":
		return delay{FaultIdMap: make(map[string]string)}
	case "duplicate":
		return duplicate{FaultIdMap: make(map[string]string)}
	case "corrupt":
		return corrupt{FaultIdMap: make(map[string]string)}
	default:
		panic("unsupported netem action")
	}
}

// Generate will randomly generate a chaos without selecting nodes.
func (g netemChaosGenerator) Generate(nodes []cluster.Node) []*core.NemesisOperation {
	nChaos := selectNetem(g.name)
	ops := make([]*core.NemesisOperation, len(nodes))

	for idx := range nodes {
		node := nodes[idx]
		ops = append(ops, &core.NemesisOperation{
			Type:        core.NetemChaos,
			Node:        &node,
			InvokeArgs:  []interface{}{nChaos},
			RecoverArgs: []interface{}{nChaos},
			RunTime:     time.Second * time.Duration(rand.Intn(120)+60),
		})
	}

	return ops
}

func (g netemChaosGenerator) Name() string {
	return g.name
}

type netem struct {
	// k8sNemesisClient
}

func (n netem) extractChaos(node *cluster.Node, operation bool, args ...interface{}) error {
	if len(args) != 1 {
		panic("netem args number is wrong")
	}
	var nChaos netemChaos
	var ok bool

	if nChaos, ok = args[0].(netemChaos); !ok {
		panic("netem get wrong type")
	}
	if operation {
		return nChaos.defaultArgsChaosApply(node)
	} else {
		return nChaos.chaosDestroy(node)
	}
}

func (n netem) Invoke(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	log.Infof("apply netem chaos on node %s(ns:%s)", node.PodName, node.Namespace)
	return n.extractChaos(node, true, args...)
}

func (n netem) Recover(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	log.Infof("unapply netem chaos on node %s(ns:%s)", node.PodName, node.Namespace)
	return n.extractChaos(node, false, args...)
}

func (n netem) Name() string {
	return string(core.NetemChaos)
}
