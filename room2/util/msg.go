package util

import (
	"steve/client_pb/msgid"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// SendMessageToPlayer 发送消息给玩家
func SendMessageToPlayer(playerID uint64, msgID msgid.MsgID, body proto.Message) error {
	sender := GetMessageSender()
	return sender.SendPackageByPlayerID(playerID, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgID),
	}, body)
}

// BroadCastMessageBare 向玩家广播消息
func BroadCastMessageBare(playerIDs []uint64, msgID msgid.MsgID, body []byte) error {
	sender := GetMessageSender()
	return sender.BroadcastPackageBareByPlayerID(playerIDs, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgID),
	}, body)
}