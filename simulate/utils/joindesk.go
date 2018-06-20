package utils

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"

	"github.com/Sirupsen/logrus"
)

// ApplyJoinDesk 申请加入牌桌
func ApplyJoinDesk(player interfaces.ClientPlayer, gameID *room.GameId) (*room.RoomJoinDeskRsp, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ApplyJoinDesk",
		"user_id":   player.GetID(),
	})
	req := room.RoomJoinDeskReq{
		GameId: gameID,
	}
	rsp := room.RoomJoinDeskRsp{}
	client := player.GetClient()
	err := client.Request(createMsgHead(msgid.MsgID_ROOM_JOIN_DESK_REQ), &req, global.DefaultWaitMessageTime, uint32(msgid.MsgID_ROOM_JOIN_DESK_RSP), &rsp)
	if err != nil {
		logEntry.WithError(err).Errorln(errRequestFailed)
		return nil, errRequestFailed
	}
	return &rsp, nil
}
