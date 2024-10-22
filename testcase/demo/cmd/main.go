// package
package main

import (
	"context"
	"flag"

	"github.com/huanghj78/jepsenFuzz/cmd/util"
	logs "github.com/huanghj78/jepsenFuzz/logsearch/pkg/logs"
	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
	"github.com/huanghj78/jepsenFuzz/pkg/control"
	"github.com/huanghj78/jepsenFuzz/pkg/test-infra/fixture"
	"github.com/huanghj78/jepsenFuzz/testcase/demo"
)

func main() {
	flag.Parse()
	cfg := control.Config{
		Mode:        control.ModeStandard,
		ClientCount: 1,
		RunTime:     fixture.Context.RunTime,
		RunRound:    1,
		DB:          "testDB",
	}
	suit := util.Suit{
		Config:        &cfg,
		Provider:      cluster.NewDefaultClusterProvider(),
		ClientCreator: demo.ClientCreator{},
		DBCreator:     demo.DBCreator{},
		NemesisGens:   util.ParseNemesisGenerators(fixture.Context.Nemesis),
		LogsClient:    logs.NewDiagnosticLogClient(),
	}
	suit.Run(context.Background())
}
