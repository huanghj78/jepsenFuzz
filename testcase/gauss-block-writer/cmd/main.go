package main

import (
	"context"
	"flag"

	"github.com/huanghj78/jepsenFuzz/cmd/util"

	_ "github.com/go-sql-driver/mysql"

	logs "github.com/huanghj78/jepsenFuzz/logsearch/pkg/logs"
	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
	"github.com/huanghj78/jepsenFuzz/pkg/control"
	"github.com/huanghj78/jepsenFuzz/pkg/test-infra/fixture"
	"github.com/huanghj78/jepsenFuzz/testcase/gauss"
)

var (
	tables      = flag.Int("tables", 10, "the number of the tables")
	concurrency = flag.Int("concurrency", 200, "concurrency of worker")
)

func main() {
	flag.Parse()
	cfg := control.Config{
		Mode:        control.ModeStandard,
		ClientCount: 1,
		RunTime:     fixture.Context.RunTime,
		RunRound:    1,
		DB:          "openGauss",
		History:     "history.log",
	}
	suit := util.Suit{
		Config:        &cfg,
		Provider:      cluster.NewDefaultClusterProvider(),
		ClientCreator: gauss.ClientCreator{TableNum: *tables, Concurrency: *concurrency},
		NemesisGens:   util.ParseNemesisGenerators(fixture.Context.Nemesis),
		DBCreator:     gauss.DBCreator{},
		LogsClient:    logs.NewDiagnosticLogClient(),
	}
	suit.Run(context.Background())
}
