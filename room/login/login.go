package login

import (
	"steve/client_pb/room"
	"steve/structs/exchanger"

	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
)

// HandleLogin 处理客户端登录消息
func HandleLogin(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomLoginReq) []exchanger.ResponseMsg {
	logrus.WithField("client_id", clientID).Debugln("客户端登录")
	return nil
}
