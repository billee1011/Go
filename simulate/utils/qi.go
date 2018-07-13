package utils

import (
	 "steve/client_pb/msgid"
	"steve/client_pb/room"
)

// SendQiReq 发送弃请求
func SendQiReq(deskData *DeskData, seat int) error {
	zjPlayer := GetDeskPlayerBySeat(seat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	_, err := zjClient.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_XINGPAI_ACTION_REQ), &room.RoomXingpaiActionReq{
		ActionId: room.XingpaiAction_XA_QI.Enum(),
	})
	return err
}
