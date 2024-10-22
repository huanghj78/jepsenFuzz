package cluster

import (
	"context"
	"fmt"
	"strings"

	"github.com/huanghj78/jepsenFuzz/pkg/test-infra/fixture"
)

// Component is the identifier of Cluster
type Component string

// Client provides useful methods about cluster
type Client struct {
	Namespace    string
	ClusterName  string
	PDMemberFunc func(ns, name string) (string, []string, error)
}

// Node is the cluster endpoint in K8s, it's maybe podIP:port or CLUSTER-IP:port
type Node struct {
	Namespace string    // Cluster k8s' namespace
	Component Component // Node component type
	PodName   string    // Pod's name
	IP        string
	Port      int32
	*Client   `json:"-"`
}

// Address returns the endpoint address of node
func (node Node) Address() string {
	return fmt.Sprintf("%s:%d", node.IP, node.Port)
}

// String ...
func (node Node) String() string {
	sb := new(strings.Builder)
	fmt.Fprintf(sb, "node[comp=%s,ip=%s:%d", node.Component, node.IP, node.Port)
	if node.Namespace != "" {
		fmt.Fprintf(sb, ",ns=%s", node.Namespace)
	}
	if node.PodName != "" {
		fmt.Fprintf(sb, ",pod=%s", node.PodName)
	}
	fmt.Fprint(sb, "]")
	return sb.String()
}

// ClientNode is TiDB's exposed endpoint, can be a nodeport, or downgrade cluster ip
type ClientNode struct {
	Namespace   string // Cluster k8s' namespace
	ClusterName string // Cluster name, use to differentiate different TiDB clusters running on same namespace
	Component   Component
	IP          string
	Port        int32
}

// Address returns the endpoint address of clientNode
func (clientNode ClientNode) Address() string {
	return fmt.Sprintf("%s:%d", clientNode.IP, clientNode.Port)
}

// String ...
func (clientNode ClientNode) String() string {
	sb := new(strings.Builder)
	fmt.Fprintf(sb, "client_node[comp=%s,ip=%s:%d", clientNode.Component, clientNode.IP, clientNode.Port)
	if clientNode.Namespace != "" {
		fmt.Fprintf(sb, ",ns=%s", clientNode.Namespace)
	}
	if clientNode.ClusterName != "" {
		fmt.Fprintf(sb, ",cluster=%s", clientNode.ClusterName)
	}
	fmt.Fprint(sb, "]")
	return sb.String()
}

// Cluster interface
type Cluster interface {
	Apply() error
	Delete() error
	GetNodes() ([]Node, error)
	GetClientNodes() ([]ClientNode, error)
}

// Specs is a cluster specification
type Specs struct {
	Cluster     Cluster
	NemesisGens []string
	Namespace   string
}

// Provider provides a collection of APIs to deploy/destroy a cluster
type Provider interface {
	// SetUp sets up cluster, returns err or all nodes info
	SetUp(ctx context.Context, spec Specs) ([]Node, []ClientNode, error)
	// TearDown tears down the cluster
	TearDown(ctx context.Context, spec Specs) error
}

// NewDefaultClusterProvider ...
func NewDefaultClusterProvider() Provider {
	return NewDockerClusterProvider(fixture.Context.NodeAddr)
}
