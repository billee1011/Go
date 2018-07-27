package utils

// SendChangePlayerReq 发送换对手请求
// func SendChangePlayerReq(seat int, gameID common.GameId, deskData *DeskData) error {
// 	player := GetDeskPlayerBySeat(seat, deskData)
// 	client := player.Player.GetClient()
// 	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_CHANGE_PLAYERS_REQ), &room.RoomChangePlayersReq{
// 		GameId: &gameID,
// 	})
// 	return err
// }
