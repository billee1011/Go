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

// RoomCard2UInt32 Card 转 int
func RoomCard2UInt32(card *room.Card) uint32 {
	var color uint32
	if *card.Color == room.CardColor_CC_WAN {
		color = 1
	} else if *card.Color == room.CardColor_CC_TIAO {
		color = 2
	} else if *card.Color == room.CardColor_CC_TONG {
		color = 3
	} else if *card.Color == room.CardColor_CC_FENG {
		color = 4
	}
	value := color*10 + uint32(*card.Point)
	return value
}

// RoomCards2UInt32 Card 转 int
func RoomCards2UInt32(card []*room.Card) []uint32 {
	result := []uint32{}
	for _, c := range card {
		result = append(result, RoomCard2UInt32(c))
	}
	return result
}
