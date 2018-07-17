package matchv2

import (
	"steve/client_pb/common"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/common/data/player"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// HandleMatchReq 匹配请求的处理(来自网关服)
func HandleMatchReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.MatchReq) (ret []exchanger.ResponseMsg) {

	response := &match.MatchRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_RSP),
		Body:  response,
	}}

	state := player.GetPlayerPlayState(playerID)
	if state != int(common.PlayerState_PS_IDLE) {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_MATCH_ALREADY_GAMEING))
		response.ErrDesc = proto.String("已经在游戏中了")
		return
	}

	defaultMgr.addPlayer(playerID, int(req.GetGameId()), false)
	return
}

// HandleContinueReq 处理续局请求
func HandleContinueReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.MatchDeskContinueReq) (ret []exchanger.ResponseMsg) {

	response := &match.MatchDeskContinueRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_CONTINUE_RSP),
		Body:  response,
	}}

	defaultMgr.addPlayer(playerID, int(req.GetGameId()), true)
	return
}
