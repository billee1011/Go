package utils

import (
	"steve/client_pb/common"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/simulate/interfaces"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// RequestMatch 申请匹配
func RequestMatch(player interfaces.ClientPlayer, gameID uint32, levelID uint32) (*match.MatchRsp, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "RequestMatch",
		"gameID":    gameID,
		"levelID":   levelID,
	})
	req := match.MatchReq{
		GameId:  proto.Uint32(gameID),
		LevelId: proto.Uint32(levelID),
	}
	rsp := match.MatchRsp{}

	client := player.GetClient()
	err := client.Request(createMsgHead(msgid.MsgID_MATCH_REQ), &req, global.DefaultWaitMessageTime, uint32(msgid.MsgID_MATCH_RSP), &rsp)
	if err != nil {
		logEntry.WithError(err).Errorln(errRequestFailed)
		return nil, errRequestFailed
	}
	return &rsp, nil
}

// ApplyJoinDesk 申请加入牌桌，从match
func ApplyJoinDesk(player interfaces.ClientPlayer, gameID common.GameId) (*match.MatchRsp, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ApplyJoinDesk",
		"user_id":   player.GetID(),
	})
	req := match.MatchReq{
		GameId: proto.Uint32(uint32(gameID)),
	}
	rsp := match.MatchRsp{}

	client := player.GetClient()
	err := client.Request(createMsgHead(msgid.MsgID_MATCH_REQ), &req, global.DefaultWaitMessageTime, uint32(msgid.MsgID_MATCH_RSP), &rsp)
	if err != nil {
		logEntry.WithError(err).Errorln(errRequestFailed)
		return nil, errRequestFailed
	}
	return &rsp, nil
}

// ApplyJoinDeskPlayers 多个玩家同时加入游戏
func ApplyJoinDeskPlayers(players []interfaces.ClientPlayer, gameID common.GameId) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ApplyJoinDeskPlayers",
	})
	req := match.MatchReq{
		GameId: proto.Uint32(uint32(gameID)),
	}
	rsp := match.MatchRsp{}
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
