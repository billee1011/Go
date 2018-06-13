package lb

import (
	"steve/structs/sgrpc"
	"time"

	bckd "github.com/bsm/grpclb/grpclb_backend_v1"
	"github.com/bsm/grpclb/load"
)

// RegisterLBReporter 注册负载
func RegisterLBReporter(rpcSvr sgrpc.RPCServer) error {
	return rpcSvr.RegisterService(bckd.RegisterLoadReportServer, load.NewRateReporter(5*time.Second))
}
