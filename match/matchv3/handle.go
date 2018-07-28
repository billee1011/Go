package matchv3

import (
	"steve/client_pb/common"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/common/data/player"
	server_pb_match "steve/server_pb/match"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HandleMatchReq 匹配请求的处理(来自网关服)
func HandleMatchReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.MatchReq) (ret []exchanger.ResponseMsg) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleMatchReq",
		"request":   req.String(),
		"playerID":  playerID,
	})

	logEntry.Debugln("进入函数")

	// 默认的回复消息
	response := &match.MatchRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}

	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_RSP),
		Body:  response,
	}}

	// 玩家当前状态
	state := player.GetPlayerPlayState(playerID)

	// 请求信息
	gameID := req.GetGameId()
	levelID := req.GetLevelId()

	// 如果处于游戏状态，返回
	if state == int(common.PlayerState_PS_GAMEING) {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_MATCH_ALREADY_GAMEING))
		response.ErrDesc = proto.String("已经在游戏中了")
		return
	}

	// 如果处于匹配状态，返回
	if state == int(common.PlayerState_PS_MATCHING) {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_MATCH_ALREADY_MATCHING))
		response.ErrDesc = proto.String("已经在匹配中了")
		return
	}

	// 最后：不是空闲状态说明错误
	if state != int(common.PlayerState_PS_IDLE) {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("不是空闲状态")

		logEntry.Errorln("内部错误，检测完不是游戏状态，不是匹配状态，最后判定仍然不是空闲状态")
		return
	}

	// 最终判定空闲状态，开始处理

	// 分发该游戏，该场次的匹配请求通道
	errString := matchMgr.dispatchMatchReq(playerID, int32(gameID), levelID)

	// 处理过程有错，回复客户端，且服务器报错
	if errString != "" {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = &errString

		logEntry.Errorln("内部错误，处理客户端的请求匹配失败，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
		return
	}

	// 设置为匹配状态，后面匹配过程中出错删除时再标记为空闲状态，匹配成功时不需处理(room服会标记为游戏状态)
	matchState := player.PlayStates{
		State:  int(common.PlayerState_PS_MATCHING), // 匹配状态
		GameID: int(gameID),                         // 游戏ID
	}

	// 如果设置状态错误，可能是客户端刚刚匹配了其他游戏，又发起了本游戏的匹配，所以这里要失败
	err := player.SetPlayerPlayStates(playerID, matchState)
	if err != nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_MATCH_ALREADY_MATCHING))
		response.ErrDesc = proto.String("刚刚匹配了其他游戏")

		logEntry.Errorln("设置匹配状态时失败，可能是客户端刚刚匹配了其他游戏，请求匹配的游戏ID:%v，场次ID:%v，玩家ID:%v", gameID, levelID, playerID)
		return
	}

	// 设置状态成功
	logEntry.Debugln("离开函数")

	return
}

// HandleContinueReq 续局请求的处理
func HandleContinueReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.MatchDeskContinueReq) (ret []exchanger.ResponseMsg) {
	logrus.WithFields(logrus.Fields{
		"func_name": "HandleContinueReq",
		"request":   req.String(),
	}).Debugln("收到玩家续局的请求")

	response := &match.MatchDeskContinueRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_CONTINUE_RSP),
		Body:  response,
	}}

	// 添加该匹配玩家（续局请求）
	matchMgr.addContinueApply(playerID, req.GetCancel(), int(req.GetGameId()))

	return
}

// AddContinueDesk 添加续局牌桌
func AddContinueDesk(request *server_pb_match.AddContinueDeskReq) *server_pb_match.AddContinueDeskRsp {
	logrus.WithFields(logrus.Fields{
		"func_name": "AddContinueDesk",
		"request":   request.String(),
	}).Debugln("收到添加续局牌桌的请求")

	response := &server_pb_match.AddContinueDeskRsp{}

	players := make([]deskPlayer, 0, len(request.GetPlayers()))

	continuePlayers := request.GetPlayers()

	// 把参数里面的player信息，转换为本地的deskPlayer
	for _, continuePlayer := range continuePlayers {
		players = append(players, deskPlayer{
			playerID: continuePlayer.GetPlayerId(),        // playerID
			robotLv:  int(continuePlayer.GetRobotLevel()), // 机器人级别
			seat:     int(continuePlayer.GetSeat()),       // 座位号
			winner:   continuePlayer.GetWin(),             // 是否是胜利者
		})
	}

	// 添加该续局牌桌
	matchMgr.addContinueDesk(players, int(request.GetGameId()), request.GetFixBanker(), int(request.GetBankerSeat()))

	return response
}
