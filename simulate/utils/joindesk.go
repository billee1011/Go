package utils

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/interfaces"
	"time"

	"github.com/Sirupsen/logrus"
)

// ApplyJoinDesk 申请加入牌桌
func ApplyJoinDesk(player interfaces.ClientPlayer) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ApplyJoinDesk",
		"user_id":   player.GetID(),
	})

	req := room.RoomJoinDeskReq{}
	rsp := room.RoomJoinDeskRsp{}
	client := player.GetClient()
	err := client.Request(createMsgHead(msgid.MsgID_ROOM_JOIN_DESK_REQ), &req, time.Second*5, uint32(msgid.MsgID_ROOM_JOIN_DESK_RSP), &rsp)
	if err != nil {
		logEntry.WithError(err).Errorln(errRequestFailed)
		return errRequestFailed
	}
	return nil
}
