package facade

import (
	msgid "steve/client_pb/msgId"
	"steve/room/interfaces/global"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// SendMessageToPlayer 发送消息给玩家
func SendMessageToPlayer(playerID uint64, msgID msgid.MsgID, body proto.Message) error {
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayer(playerID)
	clientID := player.GetClientID()
	if clientID == 0 {
		return nil
	}
	sender := global.GetMessageSender()
	return sender.SendPackage(clientID, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgID),
	}, body)
}

// BroadCastMessageBare 向玩家广播消息
func BroadCastMessageBare(playerIDs []uint64, msgID msgid.MsgID, body []byte) error {
	playerMgr := global.GetPlayerMgr()

	clientIDs := []uint64{}
	for _, playerID := range playerIDs {
		player := playerMgr.GetPlayer(playerID)
		clientID := player.GetClientID()
		if clientID != 0 {
			clientIDs = append(clientIDs, clientID)
		}
	}
	if len(clientIDs) == 0 {
		return nil
	}
	sender := global.GetMessageSender()
	return sender.BroadcastPackageBare(clientIDs, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgID),
	}, body)
}
