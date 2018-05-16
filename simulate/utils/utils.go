package utils

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/interfaces"

	"github.com/golang/protobuf/proto"
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

// SendChupaiReq 发送出牌请求
func SendChupaiReq(deskData *DeskData, seat int, card uint32) error {
	zjPlayer := GetDeskPlayerBySeat(seat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	_, err := zjClient.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_CHUPAI_REQ), &room.RoomChupaiReq{
		Card: proto.Uint32(card),
	})
	return err
}
