package loader
/*
功能： 服务负载报告，返回当前负载值，以及退休引起的负载值< 0
作者： SkyWang
日期： 2018-7-19
 */

import (
	"context"
	"steve/structs/sgrpc"

	bckd "github.com/bsm/grpclb/grpclb_backend_v1"
	"github.com/Sirupsen/logrus"
)

// RegisterLBReporter 注册负载
func RegisterLBReporter(rpcSvr sgrpc.RPCServer) error {
	return rpcSvr.RegisterService(bckd.RegisterLoadReportServer, &lbReporter{svr:rpcSvr})
}

type lbReporter struct{
	svr sgrpc.RPCServer		// RPCServer
	curScore int64    		// 当前负载值
}

// 获取服务当前负载值 API
func (lbr *lbReporter) Load(ctx context.Context, request *bckd.LoadRequest) (response *bckd.LoadResponse, err error) {
	// 获取当前负载值
	score := lbr.svr.GetScore()
	// 如果退休，将score设置成-1
	if lbr.svr.IsRetire() {
		score = -1
	}
	response, err = &bckd.LoadResponse{
		Score: score,
	}, nil

	if lbr.curScore != score {
		lbr.curScore = score
		logrus.WithFields(logrus.Fields{
			"func_name": "lbReporter.Load",
			"score":     score,
		}).Infoln("获取负载")
	}
	return
}
