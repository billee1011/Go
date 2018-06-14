package utils

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
)

// SendRecoverGameReq 发送恢复游戏请求
func SendRecoverGameReq(seat int, deskData *DeskData) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_RESUME_GAME_REQ), &room.RoomResumeGameReq{})
	return err
}
