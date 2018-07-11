package utils

import (
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
)

// SendChangePlayerReq 发送换对手请求
func SendChangePlayerReq(seat int, gameID room.GameId, deskData *DeskData) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_CHANGE_PLAYERS_REQ), &room.RoomChangePlayersReq{
		GameId: &gameID,
	})
	return err
}
