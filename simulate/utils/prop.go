package utils

import (
	"steve/client_pb/common"
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
)

// SendUsePropReq 发送使用道具请求
func SendUsePropReq(seat int, deskData *DeskData, toPlayerID uint64, propID common.PropType) error {
	player := GetDeskPlayerBySeat(seat, deskData)
	client := player.Player.GetClient()
	_, err := client.SendPackage(CreateMsgHead(msgid.MsgID_ROOM_USE_PROP_REQ),
		&room.RoomUsePropReq{
			PlayerId: &toPlayerID,
			PropId:   &propID,
		})
	return err
}
