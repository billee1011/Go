package matchv3

import (
	"steve/client_pb/common"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/external/hallclient"
	server_pb_match "steve/server_pb/match"
	"steve/server_pb/user"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HandleMatchReq 匹配请求的处理(来自网关服)
func HandleMatchReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.MatchReq) (ret []exchanger.ResponseMsg) {
	logEntry := logrus.WithFields(logrus.Fields{
		"request":  req,
		"playerID": playerID,
	})

	logEntry.Debugln("进入函数-匹配请求")

	// 默认的回复消息
	response := &match.MatchRsp{
		ErrCode: proto.Int32(int32(match.MatchError_EC_SUCCESS)),
		ErrDesc: proto.String("成功"),
		GameId:  proto.Uint32(req.GetGameId()),
		LevelId: proto.Uint32(req.GetLevelId()),
	}

	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_RSP),
		Body:  response,
	}}

	// 玩家当前状态
	rsp, err := hallclient.GetPlayerState(playerID)
	if err != nil || rsp == nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("从hall服获取玩家状态出错")

		logEntry.WithError(err).Errorln("内部错误，从hall服获取玩家状态出错")
		return
	}

	// 客户端IP地址不能为空
	/* 	rspIP := rsp.GetIpAddr()
	   	if rspIP == "" {
	   		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
	   		response.ErrDesc = proto.String("从hall服获取玩家状态时发现IP地址为空")

	   		logEntry.WithError(err).Errorln("内部错误，从hall服获取玩家状态时发现IP地址为空")
	   		return
	   	} */

	curState := rsp.GetState()
	curGameID := rsp.GetGameId()
	curLevelID := rsp.GetLevelId()
	clintIP := /* IPStringToUInt32(rspIP) */ 127010101

	// 如果处于游戏状态，返回
	if curState == user.PlayerState_PS_GAMEING {
		response.ErrCode = proto.Int32(int32(match.MatchError_EC_ALREADY_GAMEING))
		response.ErrDesc = proto.String("已经在游戏中了")

		logEntry.Warningf("匹配时发现已经在游戏状态中了，所在游戏ID:%v，所在场次ID:%v \n", curGameID, curLevelID)
		return
	}

	// 如果处于匹配状态，返回
	if curState == user.PlayerState_PS_MATCHING {
		response.ErrCode = proto.Int32(int32(match.MatchError_EC_ALREADY_MATCHING))
		response.ErrDesc = proto.String("已经在匹配中了")

		logEntry.Warningf("匹配时发现已经在匹配状态中了，正在匹配游戏ID:%v，正在匹配场次ID:%v \n", curGameID, curLevelID)
		return
	}

	// 最后：不是空闲状态说明错误
	if curState != user.PlayerState_PS_IDIE {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("不是空闲状态")

		logEntry.Errorln("内部错误，检测完不是游戏状态，不是匹配状态，最后判定仍然不是空闲状态")
		return
	}

	// 请求信息
	reqGameID := req.GetGameId()
	reqLevelID := req.GetLevelId()

	// 最终判定空闲状态

	// 分发该游戏，该场次的匹配请求通道
	errString := matchMgr.dispatchMatchReq(playerID, reqGameID, reqLevelID, uint32(clintIP))

	// 处理过程有错，回复客户端，且服务器自身报错
	if errString != "" {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = &errString

		logEntry.Errorf("处理客户端的请求匹配失败，请求匹配的游戏ID:%v，场次ID:%v \n", reqGameID, reqLevelID)
		return
	}

	// 设置状态成功
	logEntry.Debugln("离开函数-匹配请求")

	return
}

// HandleCancelMatchReq 取消匹配的处理(来自网关服)
func HandleCancelMatchReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.CancelMatchReq) (ret []exchanger.ResponseMsg) {
	logEntry := logrus.WithFields(logrus.Fields{
		"request":  req,
		"playerID": playerID,
	})

	logEntry.Debugln("进入函数-取消匹配请求")

	// 默认的回复消息
	response := &match.CancelMatchRsp{
		ErrCode: proto.Int32(int32(match.MatchError_EC_SUCCESS)),
		ErrDesc: proto.String("成功"),
	}

	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_CANCEL_MATCH_RSP),
		Body:  response,
	}}

	// 玩家当前状态
	rsp, err := hallclient.GetPlayerState(playerID)
	if err != nil || rsp == nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("从hall服获取玩家状态出错")

		logEntry.WithError(err).Errorln("内部错误，从hall服获取玩家状态出错")
		return
	}

	curState := rsp.GetState()
	curGameID := rsp.GetGameId()
	curLevelID := rsp.GetLevelId()

	// 若不是匹配状态，返回
	if curState != user.PlayerState_PS_MATCHING {
		response.ErrCode = proto.Int32(int32(match.MatchError_EC_NOT_IN_MATCHING))
		response.ErrDesc = proto.String("不在匹配中")

		logEntry.Warningf("取消匹配时发现不在匹配状态中")
		return
	}

	// 分发该游戏，该场次的匹配请求通道
	errString := matchMgr.dispatchCancelMatchReq(playerID, curGameID, curLevelID)

	// 处理过程有错，回复客户端，且服务器自身报错
	if errString != "" {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = &errString

		logEntry.Errorf("内部错误，分发客户端的取消匹配失败")
		return
	}

	// 设置状态成功
	logEntry.Debugln("离开函数-取消匹配请求")

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
	//matchMgr.addContinueApply(playerID, req.GetCancel(), int(req.GetGameId()))

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
	//matchMgr.addContinueDesk(players, int(request.GetGameId()), request.GetFixBanker(), int(request.GetBankerSeat()))

	return response
}

// ClearAllMatch 清空所有的匹配
func ClearAllMatch(req *server_pb_match.ClearAllMatchReq) *server_pb_match.ClearAllMatchRsp {

	logrus.Debugln("开始处理玩家清空所有匹配的请求")

	rsp := &server_pb_match.ClearAllMatchRsp{}

	matchMgr.ClearAllMatch()

	return rsp
}
