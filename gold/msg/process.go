package msg

/*
 功能：
		1. 完成从GateWay(网关）过来的所有Client的请求消息的处理。
 		2. 通过core.coreConfig配置需要处理的消息列表。
*/

/*
// ProcessMatchReq 匹配请求的处理(来自网关服)
func ProcessMatchReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.MatchReq) (ret []exchanger.ResponseMsg) {

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

	defaultMgr.addPlayer(playerID, int(req.GetGameId()))
	return
}
*/
