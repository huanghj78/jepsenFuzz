// Package demo ...
package demo

import (
	"context"

	"github.com/ngaut/log"

	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
	"github.com/huanghj78/jepsenFuzz/pkg/core"
)

// TestCase ...
type testCase struct{}

// ClientCreator ...
type ClientCreator struct{}

// Create ...
func (c ClientCreator) Create(_ cluster.ClientNode) core.Client {
	return NewTestCase()
}

// NewTestCase ...
func NewTestCase() core.Client {
	return &testCase{}
}

func (t *testCase) SetUp(ctx context.Context, _ []cluster.Node, clientNodes []cluster.ClientNode, idx int) error {
	log.Info("SetUp")
	return nil
}

func (t *testCase) TearDown(ctx context.Context, nodes []cluster.ClientNode, idx int) error {
	log.Info("TearDown")
	return nil
}

func (t *testCase) Start(ctx context.Context, cfg interface{}, clientNodes []cluster.ClientNode) error {
	log.Info("Start")
	return nil
}
