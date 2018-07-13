package utils

import (
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
)

// SendPengReq 发送碰请求
func SendPengReq(deskData *DeskData, seat int) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_XINGPAI_ACTION_REQ), &room.RoomXingpaiActionReq{
		ActionId: room.XingpaiAction_XA_PENG.Enum(),
	})
	return err
}
