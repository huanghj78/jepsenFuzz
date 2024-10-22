package cluster

import (
	"context"
	"log"
	"strconv"
	"strings"
)

type DockerCluster struct{}

func NewDockerClusterProvider(dbs []string) Provider {
	return &DockerClusterProvider{
		DBs: dbs,
	}
}

type DockerClusterProvider struct {
	DBs []string
}

func (d *DockerClusterProvider) SetUp(ctx context.Context, _ Specs) ([]Node, []ClientNode, error) {
	var nodes []Node
	var clientNode []ClientNode
	for _, node := range d.DBs {
		addr := strings.Split(node, ":")
		if len(addr) != 2 {
			log.Fatalf("expect format ip:port, got %s", addr)
		}
		ip := addr[0]
		port, err := strconv.Atoi(addr[1])
		if err != nil {
			log.Fatalf("illegal port %s", addr[1])
		}
		nodes = append(nodes, Node{IP: ip, Port: int32(port)})
		clientNode = append(clientNode, ClientNode{IP: ip, Port: int32(port)})
	}
	return nodes, clientNode, nil
}

func (d *DockerClusterProvider) TearDown(ctx context.Context, _ Specs) error {
	return nil
}
