package loader

import (
	"context"
	"steve/structs/proto/health"
	"steve/structs/sgrpc"

	"github.com/Sirupsen/logrus"
)

type healthServerImpl struct {
}

func (hsi *healthServerImpl) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (resp *grpc_health_v1.HealthCheckResponse, err error) {
	resp = &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}
	return
}

func registerHealthServer(rpcServer sgrpc.RPCServer) {
	if err := rpcServer.RegisterService(grpc_health_v1.RegisterHealthServer, &healthServerImpl{}); err != nil {
		logrus.WithError(err).Panicln("注册健康检查服务失败")
	}
}
