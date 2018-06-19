package gateservice

import (
	"context"
	"steve/gateway/config"
	"steve/server_pb/gateway"

	"github.com/spf13/viper"
)

type gateService struct{}

func (gs *gateService) GetGatewayAddress(ctx context.Context, request *gateway.GetGatewayAddressRequest) (*gateway.GetGatewayAddressResponse, error) {
	response := &gateway.GetGatewayAddressResponse{
		Addr: &gateway.GatewayAddress{
			Ip:   viper.GetString(config.ListenClientAddrInquire),
			Port: int32(viper.GetInt(config.ListenClientPort)),
		},
	}
	return response, nil
}

// New 创建服务
func New() gateway.GateServiceServer {
	return &gateService{}
}
