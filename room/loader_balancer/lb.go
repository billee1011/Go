package lb

import (
	"context"
	"steve/room/interfaces/global"
	"steve/structs/sgrpc"

	"github.com/Sirupsen/logrus"
	bckd "github.com/bsm/grpclb/grpclb_backend_v1"
)

// RegisterLBReporter 注册负载
func RegisterLBReporter(rpcSvr sgrpc.RPCServer) error {
	return rpcSvr.RegisterService(bckd.RegisterLoadReportServer, &lbReporter{})
}

type lbReporter struct{}

func (lbr *lbReporter) Load(ctx context.Context, request *bckd.LoadRequest) (response *bckd.LoadResponse, err error) {
	score := int64(global.GetDeskMgr().GetDeskCount())
	response, err = &bckd.LoadResponse{
		Score: score,
	}, nil
	logrus.WithFields(logrus.Fields{
		"func_name": "lbReporter.Load",
		"score":     score,
	}).Debugln("获取负载")
	return
}
