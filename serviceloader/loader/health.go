package loader

import (
	"context"
	"steve/structs/proto/health"
	"steve/structs/sgrpc"
)

type healthServerImpl struct {
}

func (hsi *healthServerImpl) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (resp *grpc_health_v1.HealthCheckResponse, err error) {
	resp = &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}
	return
}

func registerHealthServer(rpcServer sgrpc.RPCServer) error {
	return rpcServer.RegisterService(grpc_health_v1.RegisterHealthServer, &healthServerImpl{})
}
