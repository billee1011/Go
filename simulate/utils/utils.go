package utils

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/interfaces"
)

func createMsgHead(msgID msgid.MsgID) interfaces.SendHead {
	return interfaces.SendHead{
		Head: interfaces.Head{
			MsgID: uint32(msgID),
		},
	}
}

// CreateMsgHead 创建消息头
func CreateMsgHead(msgID msgid.MsgID) interfaces.SendHead {
	return createMsgHead(msgID)
}

// MakeRoomCards 构造牌切片
func MakeRoomCards(card ...room.Card) []*room.Card {
	result := []*room.Card{}
	for i := range card {
		result = append(result, &card[i])
	}
	return result
}
