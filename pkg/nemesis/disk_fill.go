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

type diskFillGenerator struct {
	name string
}

func (g diskFillGenerator) Generate(nodes []cluster.Node) []*core.NemesisOperation {
	ops := make([]*core.NemesisOperation, len(nodes))

	for idx := range nodes {
		node := nodes[idx]
		ops = append(ops, &core.NemesisOperation{
			Type:        core.DiskFill,
			Node:        &node,
			InvokeArgs:  []interface{}{rand.Intn(10) + 80},
			RecoverArgs: []interface{}{},
			RunTime:     time.Second * time.Duration(rand.Intn(10)+20),
		})
	}

	return ops
}

func (g diskFillGenerator) Name() string {
	return g.name
}

func NewDiskFillGenerator(name string) core.NemesisGenerator {
	return diskFillGenerator{
		name: name,
	}
}

type diskFill struct {
	FaultIdMap map[string]string
}

func (n diskFill) Invoke(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	log.Infof("inject disk fill on node%d", node.ID)
	percent, _ := args[0].(int)
	var result map[string]interface{}
	cmd := fmt.Sprintf("blade create disk fill --percent %d --timeout %d", percent, Timeout)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "ilovedds", cmd)
	if err != nil {
		log.Error("Execute command failed, err: %v, output: %s", err, output)
		return err
	}
	jsonOutput := strings.TrimSpace(output)
	err = json.Unmarshal([]byte(jsonOutput), &result)
	if err != nil {
		log.Errorf("Error unmarshalling JSON, err: %v, output: %s", err, output)
		return err
	}
	n.FaultIdMap[node.IP], _ = result["result"].(string)
	return nil
}

func (n diskFill) Recover(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	id := n.FaultIdMap[node.IP]
	cmd := fmt.Sprintf("blade destroy %s", id)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "ilovedds", cmd)
	if err != nil {
		log.Error("Execute command failed, err: %v, output: %s", err, output)
		return err
	}
	delete(n.FaultIdMap, id)
	return nil
}

func (n diskFill) Name() string {
	return string(core.DiskFill)
}
