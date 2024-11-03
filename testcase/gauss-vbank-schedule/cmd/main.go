package main

import (
	"context"
	"flag"

	"github.com/huanghj78/jepsenFuzz/cmd/util"
	logs "github.com/huanghj78/jepsenFuzz/logsearch/pkg/logs"
	"github.com/huanghj78/jepsenFuzz/pkg/check/porcupine"
	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
	"github.com/huanghj78/jepsenFuzz/pkg/control"
	"github.com/huanghj78/jepsenFuzz/pkg/core"
	"github.com/huanghj78/jepsenFuzz/pkg/test-infra/fixture"
	"github.com/huanghj78/jepsenFuzz/pkg/verify"
	vbank "github.com/huanghj78/jepsenFuzz/testcase/gauss_vbank_schedule"
)

var (
	pkType        = flag.String("pk_type", "int", "primary key type, int,decimal,string")
	partition     = flag.Bool("partition", true, "use partitioned table")
	useRange      = flag.Bool("range", false, "use range condition")
	updateInPlace = flag.Bool("update_in_place", false, "use update in place mode")
	readCommitted = flag.Bool("read_committed", false, "use READ-COMMITTED isolation level")
	connParams    = flag.String("conn_params", "", "connection parameters")
)

// ./bin/vbank -node-addr 10.10.3.0:26000  -node-addr 10.10.4.26:26000 -node-addr 10.10.3.76:26000 -node-addr 10.10.1.9:26000 -node-addr 10.10.1.174:26000 -nemesis proc-kill

func main() {
	flag.Parse()

	cfg := control.Config{
		Mode:         control.ModeOnSchedule,
		ClientCount:  fixture.Context.ClientCount,
		RequestCount: fixture.Context.RequestCount,
		RunRound:     fixture.Context.RunRound,
		RunTime:      fixture.Context.RunTime,
		History:      fixture.Context.HistoryFile,
	}
	verifySuit := verify.Suit{
		Model:   &vbank.Model{},
		Checker: core.MultiChecker("v_bank checkers", porcupine.Checker{}),
		Parser:  &vbank.Parser{},
	}
	vbCfg := &vbank.Config{
		PKType:        *pkType,
		Partition:     *partition,
		Range:         *useRange,
		ReadCommitted: *readCommitted,
		UpdateInPlace: *updateInPlace,
		ConnParams:    *connParams,
	}
	suit := util.Suit{
		Config:           &cfg,
		ClientCreator:    vbank.NewClientCreator(vbCfg),
		ClientRequestGen: util.OnClientLoop,
		VerifySuit:       verifySuit,
		DBCreator:        vbank.DBCreator{},
		Provider:         cluster.NewDefaultClusterProvider(),
		NemesisGens:      util.ParseNemesisGenerators(fixture.Context.Nemesis),
		LogsClient:       logs.NewDiagnosticLogClient(),
	}
	suit.Run(context.Background())
}
