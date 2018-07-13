package utils

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"

	"github.com/Sirupsen/logrus"
)

// ApplyJoinDesk 申请加入牌桌，从match
func ApplyJoinDesk(player interfaces.ClientPlayer, gameID room.GameId) (*room.RoomJoinDeskRsp, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ApplyJoinDesk",
		"user_id":   player.GetID(),
	})
	req := room.RoomJoinDeskReq{
		GameId: &gameID,
	}
	rsp := room.RoomJoinDeskRsp{}

	client := player.GetClient()
	err := client.Request(createMsgHead(msgid.MsgID_MATCH_REQ), &req, global.DefaultWaitMessageTime, uint32(msgid.MsgID_MATCH_RSP), &rsp)
	if err != nil {
		logEntry.WithError(err).Errorln(errRequestFailed)
		return nil, errRequestFailed
	}
	return &rsp, nil
}

// ApplyJoinDeskPlayers 多个玩家同时加入游戏
func ApplyJoinDeskPlayers(players []interfaces.ClientPlayer, gameID room.GameId) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ApplyJoinDeskPlayers",
	})
	req := room.RoomJoinDeskReq{
		GameId: &gameID,
	}
	rsp := room.RoomJoinDeskRsp{}
	for _, player := range players {
		client := player.GetClient()
		err := client.Request(createMsgHead(msgid.MsgID_MATCH_REQ), &req, global.DefaultWaitMessageTime, uint32(msgid.MsgID_MATCH_RSP), &rsp)
		if err != nil {
			logEntry.WithField("user_id", player.GetID()).WithError(err).Errorln(errRequestFailed)
			return errRequestFailed
		}
	}

	return nil
}
