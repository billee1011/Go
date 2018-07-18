package gutils

import (
	"steve/client_pb/msgid"
	"steve/structs"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// SendMessage 发送消息
func SendMessage(playerID uint64, msgID msgid.MsgID, body proto.Message) {
	exchanger := structs.GetGlobalExposer().Exchanger
	exchanger.SendPackageByPlayerID(playerID, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgID),
	}, body)
}
