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

type diskBurnGenerator struct {
	name string
}

func (g diskBurnGenerator) Generate(nodes []cluster.Node) []*core.NemesisOperation {
	ops := make([]*core.NemesisOperation, len(nodes))

	for idx := range nodes {
		node := nodes[idx]
		ops = append(ops, &core.NemesisOperation{
			Type:        core.DiskBurn,
			Node:        &node,
			InvokeArgs:  []interface{}{},
			RecoverArgs: []interface{}{},
			RunTime:     time.Second * time.Duration(rand.Intn(5)+1),
		})
	}

	return ops
}

func (g diskBurnGenerator) Name() string {
	return g.name
}

func NewDiskBurnGenerator(name string) core.NemesisGenerator {
	return diskBurnGenerator{
		name: name,
	}
}

type diskBurn struct {
	FaultIdMap map[string]string
}

func (n diskBurn) Invoke(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	log.Infof("inject disk burn on node%d", node.ID)
	var result map[string]interface{}
	cmd := fmt.Sprintf("blade create disk burn --read --write --timeout %d", Timeout)
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

func (n diskBurn) Recover(ctx context.Context, node *cluster.Node, args ...interface{}) error {
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

func (n diskBurn) Name() string {
	return string(core.DiskBurn)
}
