package loader

import (
	"steve/serviceloader/exchanger"
	iexchanger "steve/structs/exchanger"
	"steve/structs/sgrpc"

	"github.com/Sirupsen/logrus"
)

func createExchanger(rpcServer sgrpc.RPCServer) iexchanger.Exchanger {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "createExchanger",
	})
	e, err := exchanger.NewExchanger(rpcServer)
	if err != nil {
		logEntry.WithError(err).Panicln("创建 Exchanger 失败")
	}
	return e
}
