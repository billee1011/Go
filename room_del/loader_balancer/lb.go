package lb

import (
	"context"
	"steve/structs/sgrpc"
	"steve/room/interfaces/global"
	"github.com/Sirupsen/logrus"
	bckd "steve/thirdpart/github.com/bsm/grpclb/grpclb_backend_v1"
)

/*
// RegisterLBReporter 注册负载
func RegisterLBReporter(rpcSvr sgrpc.RPCServer) error {
	return rpcSvr.RegisterService(bckd.RegisterLoadReportServer, &lbReporter{svr:rpcSvr})
}
*/

type lbReporter struct{
	svr sgrpc.RPCServer
}

var _ lbReporter

func (lbr *lbReporter) Load(ctx context.Context, request *bckd.LoadRequest) (response *bckd.LoadResponse, err error) {
	score := int64(global.GetDeskMgr().GetDeskCount())

	// 如果退休，将score设置成-1
	if lbr.svr.IsRetire() {
		score = -1
	}
	response, err = &bckd.LoadResponse{
		Score: score,
	}, nil
	logrus.WithFields(logrus.Fields{
		"func_name": "lbReporter.Load",
		"score":     score,
	}).Debugln("获取负载")
	return
}
