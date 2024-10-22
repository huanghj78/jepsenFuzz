module github.com/huanghj78/jepsenFuzz/testcase/gauss

go 1.22.8

replace github.com/huanghj78/jepsenFuzz => ../../.

replace github.com/huanghj78/jepsenFuzz/logsearch => ../../logsearch

require (
	gitee.com/opengauss/openGauss-connector-go-pq v1.0.4
	github.com/huanghj78/jepsenFuzz v0.0.0-00010101000000-000000000000
	github.com/huanghj78/jepsenFuzz/logsearch v0.0.0-00010101000000-000000000000
	github.com/juju/errors v1.0.0
	github.com/ngaut/log v0.0.0-20221012222132-f3329cba28a5
	github.com/rogpeppe/fastuuid v1.2.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/pingcap/kvproto v0.0.0-20240924080114-4a3e17f5e62d // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	k8s.io/apimachinery v0.31.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8 // indirect
)
