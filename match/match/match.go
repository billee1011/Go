package match

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
)

// 匹配请求的处理(来自网关服)
func HandleMatchReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomJoinDeskReq) (ret []exchanger.ResponseMsg) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "matchCore::handleMatch()",
	})

	response := &room.RoomJoinDeskRsp{
		ErrCode: room.RoomError_SUCCESS.Enum(),
	}
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_RSP),
		Body:  response,
	}}

	logEntry.WithField("playerID", playerID).Debugln("加入新的匹配玩家")

	defaultManager.addPlayer(playerID, int(req.GetGameId()))
	return
}
