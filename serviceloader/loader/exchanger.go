package loader

import (
	"steve/serviceloader/exchanger"
	iexchanger "steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"steve/structs/sgrpc"

	"github.com/Sirupsen/logrus"
)

func createExchanger(rpcServer sgrpc.RPCServer) iexchanger.Exchanger {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "createExchanger",
	})
	h, e := exchanger.NewMessageHandlerServer()
	if err := rpcServer.RegisterService(steve_proto_gaterpc.RegisterMessageHandlerServer, h); err != nil {
		logEntry.Panicln("注册服务失败")
	}
	return e
}
