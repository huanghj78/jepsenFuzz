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

type cpuFullloadGenerator struct {
	name string
}

func (g cpuFullloadGenerator) Generate(nodes []cluster.Node) []*core.NemesisOperation {
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
	return cpuFullloadNodes(nodes, n, time.Second*time.Duration(rand.Intn(10)+10), rand.Intn(10)+80)
}

func (g cpuFullloadGenerator) Name() string {
	return g.name
}

func cpuFullloadNodes(nodes []cluster.Node, n int, duration time.Duration, percent int) []*core.NemesisOperation {
	var ops []*core.NemesisOperation
	indices := shuffleIndices(len(nodes))
	if n > len(indices) {
		n = len(indices)
	}
	for i := 0; i < n; i++ {
		ops = append(ops, &core.NemesisOperation{
			Type:        core.CPUFullload,
			Node:        &nodes[indices[i]],
			InvokeArgs:  []interface{}{percent},
			RecoverArgs: []interface{}{percent},
			RunTime:     duration,
		})
	}
	return ops
}

func NewCPUFullloadGenerator(name string) core.NemesisGenerator {
	return cpuFullloadGenerator{name: name}
}

type cpuFullload struct {
	FaultIdMap map[string]string
}

func (c cpuFullload) Invoke(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	percent, _ := args[0].(int)
	cmd := fmt.Sprintf("blade create cpu load  --timeout 300 --cpu-percent %d", percent)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "ilovedds", cmd)
	if err != nil {
		log.Error(output)
		return err
	}
	var result map[string]interface{}
	jsonOutput := strings.TrimSpace(output)
	err = json.Unmarshal([]byte(jsonOutput), &result)
	if err != nil {
		log.Errorf("Error unmarshalling JSON: %v", err)
	}
	c.FaultIdMap[node.IP], _ = result["result"].(string)
	return nil
}

func (c cpuFullload) Recover(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	log.Infof("recover cpu fullload on node%d", node.ID)
	id := c.FaultIdMap[node.IP]
	cmd := fmt.Sprintf("blade destroy %s", id)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "ilovedds", cmd)
	if err != nil {
		log.Error(output)
		return err
	}
	delete(c.FaultIdMap, id)
	log.Info(output)
	return nil
}

func (c cpuFullload) Name() string {
	return string(core.CPUFullload)
}
