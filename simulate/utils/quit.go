package utils

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
)

// SendQuitReq 发送退出牌桌请求
func SendQuitReq(deskData *DeskData, seat int) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_DESK_QUIT_REQ), &room.RoomDeskQuitReq{})
	return err
}
