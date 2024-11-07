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

type networkPartitionGenerator struct {
	name string
}

func (g networkPartitionGenerator) Generate(nodes []cluster.Node) []*core.NemesisOperation {
	n := 1
	switch g.name {
	case "two_partition":
		if len(nodes) == 1 {
			n = 1
		} else {
			n = 2
		}
	case "multi_partition":
		if len(nodes) >= 3 {
			n = rand.Intn(len(nodes)-2) + 3 // 随机选择 3 到 len(nodes) 之间的数
		} else {
			n = len(nodes) // 如果节点数少于 3，则设为节点数
		}
	case "all_partition":
		n = len(nodes)
	default:
		n = 1
	}
	return partitionNodes(nodes, n, time.Second*time.Duration(rand.Intn(10)+10))
	// return partitionNodes(nodes, n, time.Second*time.Duration(rand.Intn(120)+60))
}

func (g networkPartitionGenerator) Name() string {
	return g.name
}

func partitionNodes(nodes []cluster.Node, n int, duration time.Duration) []*core.NemesisOperation {
	if n < 1 {
		log.Fatalf("the partition part size cannot be less than 1")
	}
	var ops []*core.NemesisOperation
	// randomly shuffle the indices and get the first n nodes to be partitioned.
	indices := shuffleIndices(len(nodes))

	// 创建分区的切片
	partitions := make([][]cluster.Node, n)

	// 将节点均匀分配到 n 个分区
	for i, index := range indices {
		partitions[i%n] = append(partitions[i%n], nodes[index])
	}

	for i, partition := range partitions {
		for _, node := range partition {
			ops = append(ops, &core.NemesisOperation{
				Type:        core.NetworkPartition,
				Node:        &node,
				InvokeArgs:  []interface{}{getOtherPartitions(partitions, i)},
				RecoverArgs: []interface{}{getOtherPartitions(partitions, i)},
				RunTime:     duration,
			})
		}
	}

	return ops
}

// getOtherPartitions 返回除了当前分区外的所有其他分区
func getOtherPartitions(partitions [][]cluster.Node, current int) []cluster.Node {
	var otherNodes []cluster.Node
	for i, partition := range partitions {
		if i != current {
			otherNodes = append(otherNodes, partition...)
		}
	}
	return otherNodes
}

// NewNetworkPartitionGenerator creates a generator.
// Name is partition-one, etc.
func NewNetworkPartitionGenerator(name string) core.NemesisGenerator {
	return networkPartitionGenerator{name: name}
}

// networkPartition implements Nemesis
type networkPartition struct {
	FaultIdMap map[string]string
}

func (n networkPartition) Invoke(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	anotherNodes := extractArgs(args...)
	var destinationIPs []string
	for _, anotherNode := range anotherNodes {
		log.Infof("inject network partition between node%d and node%d", node.ID, anotherNode.ID)
		destinationIPs = append(destinationIPs, anotherNode.IP)
	}
	destinationIPStr := strings.Join(destinationIPs, ",")
	var result map[string]interface{}
	cmd := fmt.Sprintf("blade create network loss --percent 100 --interface eth0 --timeout 300 --destination-ip %s", destinationIPStr)
	log.Debug("cmd=", cmd)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "ilovedds", cmd)
	if err != nil {
		log.Error("Execute command failed, err: %v, output: %s", err, output)
		return err
	}
	jsonOutput := strings.TrimSpace(output)
	err = json.Unmarshal([]byte(jsonOutput), &result)
	if err != nil {
		log.Errorf("Error unmarshalling JSON, err: %v, output", err, output)
		return err
	}
	n.FaultIdMap[node.IP], _ = result["result"].(string)
	return nil
}

func (n networkPartition) Recover(ctx context.Context, node *cluster.Node, args ...interface{}) error {
	// srcNode, dstNode := extractArgs(args...)
	log.Infof("recover network partition between node%d", node.ID)
	id := n.FaultIdMap[node.IP]
	// log.Debug("key=", srcNode.IP+"-"+dstNode.IP)
	log.Debug("id=", id)
	cmd := fmt.Sprintf("blade destroy %s", id)
	output, err := util.ExecuteRemoteCommand(node.IP, "root", "ilovedds", cmd)
	if err != nil {
		log.Error("Execute command failed, err: %v, output: %s", err, output)
		return err
	}
	jsonOutput := strings.TrimSpace(output)
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonOutput), &result)
	if err != nil {
		log.Errorf("Error unmarshalling JSON, err: %v, output", err, output)
	}
	log.Debug(result)
	delete(n.FaultIdMap, id)
	log.Info(output)
	return nil
}

func (n networkPartition) Name() string {
	return string(core.NetworkPartition)
}

func extractArgs(args ...interface{}) []cluster.Node {
	var anotherNodes []cluster.Node

	if len(args) != 1 {
		log.Fatalf("expect one args, got %+v", args)
	}

	anotherNodes = args[0].([]cluster.Node)
	return anotherNodes
}
