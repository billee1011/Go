package gateservice

import (
	"context"
	"steve/client_pb/gate"
	"steve/client_pb/msgid"
	"steve/gateway/config"
	"steve/gateway/connection"
	"steve/gateway/watchdog"
	"steve/server_pb/gateway"
	"steve/structs/net"
	"steve/structs/proto/base"

	"github.com/Sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

// GateService 网关服务
// 实现 gateway.GateServiceServer
type GateService struct {
	watchDog net.WatchDog
}

var defaultObject = new(GateService)
var _ gateway.GateServiceServer = Default()

// Default 默认对象
func Default() *GateService {
	return defaultObject
}

// SetWatchDog 设置 watch dog
func (gs *GateService) SetWatchDog(dog net.WatchDog) {
	gs.watchDog = dog
}

// GetGatewayAddress 获取客户端连接地址
func (gs *GateService) GetGatewayAddress(ctx context.Context, request *gateway.GetGatewayAddressRequest) (*gateway.GetGatewayAddressResponse, error) {
	response := &gateway.GetGatewayAddressResponse{
		Addr: &gateway.GatewayAddress{
			Ip:   viper.GetString(config.ListenClientAddrInquire),
			Port: int32(viper.GetInt(config.ListenClientPort)),
		},
	}
	return response, nil
}

// AnotherLogin 顶号
func (gs *GateService) AnotherLogin(ctx context.Context, request *gateway.AnotherLoginRequest) (response *gateway.AnotherLoginResponse, err error) {
	response = &gateway.AnotherLoginResponse{}
	playerID := request.GetPlayerId()
	AnotherLogin(playerID)
	return
}

// AnotherLogin 顶号
func AnotherLogin(playerID uint64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "AnotherLogin",
		"player_id": playerID,
	})
	connMgr := connection.GetConnectionMgr()
	connection := connMgr.GetPlayerConnection(playerID)
	if connection == nil {
		entry.Infoln("玩家不在本网关")
		return
	}
	connection.DetachPlayer()
	connectionID := connection.GetClientID()

	notify := gate.GateAnotherLoginNtf{}
	data, err := proto.Marshal(&notify)
	if err != nil {
		entry.WithError(err).Errorln("消息序列化失败")
		return
	}
	dog := watchdog.Get()
	dog.SendPackage(connectionID, &base.Header{
		MsgId:   proto.Uint32(uint32(msgid.MsgID_GATE_ANOTHER_LOGIN_NTF)),
		Version: proto.String("1.0"),
	}, data)
	dog.Disconnect(connectionID)
}
