package nemesis

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
	"github.com/huanghj78/jepsenFuzz/pkg/core"
	"github.com/huanghj78/jepsenFuzz/util"
	"github.com/ngaut/log"
)

type cpuFullLoadGenerator struct {
	name string
}

func (g cpuFullLoadGenerator) Generate(nodes []cluster.Node) []*core.NemesisOperation {
	n := 1
	switch g.name {
	case "random_cpufl":
		n = 1
	case "all_cpufl":
		n = len(nodes)
	case "major_cpufl":
		n = len(nodes)/2 + 1
	case "minor_cpufl":
		n = len(nodes) / 2
	default:
		n = 1
	}
	return cpuFullLoadNodes(nodes, n, time.Second*time.Duration(rand.Intn(10)+10), rand.Intn(10)+20)
}

func (g cpuFullLoadGenerator) Name() string {
	return g.name
}

func cpuFullLoadNodes(nodes []cluster.Node, n int, duration time.Duration, percent int) []*core.NemesisOperation {
	var ops []*core.NemesisOperation
	indices := shuffleIndices(len(nodes))
	if n > len(indices) {
		n = len(indices)
	}
	for i := 0; i < n; i++ {
		ops = append(ops, &core.NemesisOperation{
			Type:        core.CPUFullLoad,
			Node:        &nodes[indices[i]],
			InvokeArgs:  []interface{}{percent},
			RecoverArgs: []interface{}{percent},
			RunTime:     duration,
		})
	}
	return ops
}

func NewCPUFullLoadGenerator(name string) core.NemesisGenerator {
	return cpuFullLoadGenerator{name: name}
}

type cpuFullLoad struct {
	FaultIdMap map[string]string
}

func (c cpuFullLoad) Invoke(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	log.Infof("inject cpu fullload on node%d", node.ID)
	percent, _ := args[0].(int)
	cmd := fmt.Sprintf("blade create cpu load  --timeout 300 --cpu-percent %d", percent)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "ilovedds", cmd)
	if err != nil {
		log.Error("Execute command failed, err: %v, output: %s", err, output)
		return err
	}
	var result map[string]interface{}
	jsonOutput := strings.TrimSpace(output)
	err = json.Unmarshal([]byte(jsonOutput), &result)
	if err != nil {
		log.Errorf("Error unmarshalling JSON, err: %v, output: %s", err, output)
		return err
	}
	c.FaultIdMap[node.IP], _ = result["result"].(string)
	return nil
}

func (c cpuFullLoad) Recover(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	log.Infof("recover cpu fullload on node%d", node.ID)
	id := c.FaultIdMap[node.IP]
	cmd := fmt.Sprintf("blade destroy %s", id)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "ilovedds", cmd)
	if err != nil {
		log.Error("Execute command failed, err: %v, output: %s", err, output)
		return err
	}
	delete(c.FaultIdMap, id)
	log.Info(output)
	return nil
}

func (c cpuFullLoad) Name() string {
	return string(core.CPUFullLoad)
}
