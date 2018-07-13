package facade

import (
	"fmt"
	"steve/client_pb/msgid"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// SendMessageToPlayer 发送消息给玩家
func SendMessageToPlayer(playerID uint64, msgID msgid.MsgID, body proto.Message) error {
	sender := global.GetMessageSender()
	return sender.SendPackageByPlayerID(playerID, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgID),
	}, body)
}

// BroadCastMessageBare 向玩家广播消息
func BroadCastMessageBare(playerIDs []uint64, msgID msgid.MsgID, body []byte) error {
	sender := global.GetMessageSender()
	return sender.BroadcastPackageBareByPlayerID(playerIDs, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgID),
	}, body)
}

// BroadCastDeskMessage 广播消息给牌卓玩家
func BroadCastDeskMessage(desk interfaces.Desk, playerIDs []uint64, msgID msgid.MsgID, body proto.Message, exceptQuit bool) error {
	msgBody, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	desk.BroadcastMessage(playerIDs, msgID, msgBody, exceptQuit)
	return nil
}

func find(datas []uint64, data uint64) bool {
	for _, d := range datas {
		if d == data {
			return true
		}
	}
	return false
}

// BroadCastDeskMessageExcept 广播消息给牌桌玩家
func BroadCastDeskMessageExcept(desk interfaces.Desk, expcetPlayers []uint64, exceptQuit bool, msgID msgid.MsgID, body proto.Message) error {
	playerIDs := []uint64{}
	deskPlayers := desk.GetDeskPlayers()
	for _, deskPlayer := range deskPlayers {
		playerID := deskPlayer.GetPlayerID()
		if find(expcetPlayers, playerID) {
			continue
		}
		playerIDs = append(playerIDs, playerID)
	}
	if len(playerIDs) == 0 {
		return fmt.Errorf("没有广播玩家")
	}
	err := BroadCastDeskMessage(desk, playerIDs, msgID, body, exceptQuit)
	return err
}
