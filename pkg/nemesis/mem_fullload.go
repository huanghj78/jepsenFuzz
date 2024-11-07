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

type memFullLoadGenerator struct {
	name string
}

func (g memFullLoadGenerator) Generate(nodes []cluster.Node) []*core.NemesisOperation {
	ops := make([]*core.NemesisOperation, len(nodes))

	for idx := range nodes {
		node := nodes[idx]
		ops = append(ops, &core.NemesisOperation{
			Type:        core.MemFullLoad,
			Node:        &node,
			InvokeArgs:  []interface{}{rand.Intn(10) + 80},
			RecoverArgs: []interface{}{},
			RunTime:     time.Second * time.Duration(rand.Intn(10)+20),
		})
	}

	return ops
}

func (g memFullLoadGenerator) Name() string {
	return g.name
}

func NewMemFullloadGenerator(name string) core.NemesisGenerator {
	return memFullLoadGenerator{
		name: name,
	}
}

type memFullLoad struct {
	FaultIdMap map[string]string
}

func (n memFullLoad) Invoke(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	log.Infof("inject mem full load on node%d", node.ID)
	percent, _ := args[0].(int)
	var result map[string]interface{}
	cmd := fmt.Sprintf("blade create mem load --mode ram --mem-percent %d --timeout %d", percent, Timeout)
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

func (n memFullLoad) Recover(ctx context.Context, node *cluster.Node, args ...interface{}) error {
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

func (n memFullLoad) Name() string {
	return string(core.MemFullLoad)
}
