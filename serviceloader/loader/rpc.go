package loader

import (
	"steve/serviceloader/structs/sgrpc"
	sgrpcinterface "steve/structs/sgrpc"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// createRPCServer 创建 RPC 服务
func createRPCServer(keyFile string, certFile string) *sgrpc.RPCServerImpl {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "createRPCServer",
		"key_file":  keyFile,
		"cert_file": certFile,
	})
	rpcOption := []grpc.ServerOption{}
	if keyFile != "" {
		cred, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			logEntry.Panicln("创建 TLS 证书失败")
		}
		rpcOption = append(rpcOption, grpc.Creds(cred))
	}
	return sgrpc.NewRPCServer(rpcOption...)
}

// createRPCClient 创建 RPC 客户端
func createRPCClient(caFile string, caServerName string) *sgrpc.ClientImpl {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":      "createRPCClient",
		"ca_file":        caFile,
		"ca_server_name": caServerName,
	})
	result := sgrpc.NewClientImpl(caFile, caServerName)
	if result == nil {
		logEntry.Panicln("创建 RPC 客户端失败")
	}
	return result
}

// runRPCServer 启动 RPC 服务
func runRPCServer(server sgrpcinterface.RPCServer, addr string, port int) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "runRPCServer",
		"addr":      addr,
		"port":      port,
	})
	if addr == "" || port == 0 {
		logEntry.Info("未配置 RPC 地址或者端口，不启动 RPC 服务")
		return
	}
	s, ok := server.(*sgrpc.RPCServerImpl)
	if !ok {
		logEntry.Panicln("转换 RPC 服务失败")
	}
	if err := s.Work(addr, port); err != nil {
		logEntry.WithError(err).Panicln("启动 RPC 服务失败")
	}
}
