package demo

import (
	"context"

	"github.com/ngaut/log"

	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
	"github.com/huanghj78/jepsenFuzz/pkg/core"
)

type testDB struct{}

func (db testDB) SetUp(ctx context.Context, nodes []cluster.Node, node cluster.Node) error {
	log.Info("SetUp DB")
	// Setup DB
	return nil
}

func (db testDB) TearDown(ctx context.Context, nodes []cluster.Node, node cluster.Node) error {
	// TearDown DB
	log.Info("TearDown DB")
	return nil
}

func (db testDB) Name() string {
	return "testDB"
}

// DBCreator ...
type DBCreator struct{}

// Create ...
func (c DBCreator) Create() core.DB {
	return testDB{}
}
