// package
package main

import (
	"context"
	"flag"
	"time"

	"github.com/huanghj78/jepsenFuzz/cmd/util"
	logs "github.com/huanghj78/jepsenFuzz/logsearch/pkg/logs"
	"github.com/huanghj78/jepsenFuzz/pkg/cluster"
	"github.com/huanghj78/jepsenFuzz/pkg/control"
	"github.com/huanghj78/jepsenFuzz/pkg/test-infra/fixture"
	"github.com/huanghj78/jepsenFuzz/testcase/gauss"
)

var (
	// case config
	retryLimit          = flag.Int("retry-limit", 2, "retry count")
	accounts            = flag.Int("accounts", 1000, "the number of accounts")
	interval            = flag.Duration("interval", 2*time.Second, "the interval")
	pessimistic         = flag.Bool("pessimistic", false, "use pessimistic transaction")
	concurrency         = flag.Int("concurrency", 200, "concurrency worker count")
	longTxn             = flag.Bool("long-txn", true, "enable long-term transactions")
	tables              = flag.Int("tables", 1, "the number of the tables")
	replicaRead         = flag.String("tidb-replica-read", "leader", "tidb_replica_read mode, support values: leader / follower / leader-and-follower, default value: leader.")
	dbname              = flag.String("dbname", "test", "name of database to test")
	tiflashDataReplicas = flag.Int("tiflash-data-replicas", 0, "the number of the tiflash data replica")
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
	gaussConfig := gauss.Config{
		EnableLongTxn:       *longTxn,
		Pessimistic:         *pessimistic,
		RetryLimit:          *retryLimit,
		Accounts:            *accounts,
		Tables:              *tables,
		Interval:            *interval,
		Concurrency:         *concurrency,
		ReplicaRead:         *replicaRead,
		DbName:              *dbname,
		TiFlashDataReplicas: *tiflashDataReplicas,
	}
	suit := util.Suit{
		Config:        &cfg,
		Provider:      cluster.NewDefaultClusterProvider(),
		ClientCreator: gauss.ClientCreator{Cfg: &gaussConfig},
		DBCreator:     gauss.DBCreator{},
		NemesisGens:   util.ParseNemesisGenerators(fixture.Context.Nemesis),
		LogsClient:    logs.NewDiagnosticLogClient(),
	}
	suit.Run(context.Background())
}
