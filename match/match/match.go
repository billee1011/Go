package match

import (
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HandleMatchReq 匹配请求的处理(来自网关服)
func HandleMatchReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.MatchReq) (ret []exchanger.ResponseMsg) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "matchCore::handleMatch()",
	})

	response := &match.MatchRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_RSP),
		Body:  response,
	}}

	logEntry.WithField("playerID", playerID).Debugln("加入新的匹配玩家")

	defaultManager.addPlayer(playerID, int(req.GetGameId()))
	return
}

// HandleContinueReq 处理续局请求
func HandleContinueReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.MatchDeskContinueReq) (ret []exchanger.ResponseMsg) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleContinueReq",
	})

	response := &match.MatchDeskContinueRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_RSP),
		Body:  response,
	}}

	logEntry.WithField("playerID", playerID).Debugln("续局")
	defaultManager.addPlayer(playerID, int(req.GetGameId()))
	return
}
