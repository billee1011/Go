package utils

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	msgid "steve/client_pb/msgid"
)

// SendGetPlayerGameInfoReq 发送获取玩家信息的请求
func SendGetPlayerGameInfoReq(seat int, deskData *DeskData, toPlayerID uint64, gameID common.GameId) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_HALL_GET_PLAYER_GAME_INFO_REQ),
		&hall.HallGetPlayerGameInfoReq{
			Uid:    &toPlayerID,
			GameId: &gameID,
		})
	return err
}
