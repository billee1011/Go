package user

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/common/data/player"
	"steve/hall/data"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// getPlayerState 获取玩家状态
func getPlayerState(playerID uint64) (common.PlayerState, int) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "getPlayerState",
		"player_id": playerID,
	})
	states, err := player.GetPlayerPlayStates(playerID, player.PlayStates{
		State: int(common.PlayerState_PS_IDLE),
	})
	if err != nil {
		entry.Errorln("获取玩家状态失败")
		return common.PlayerState_PS_IDLE, 0
	}
	return common.PlayerState(states.State), states.GameID
}

// HandleGetPlayerInfoReq 处理获取玩家信息请求
func HandleGetPlayerInfoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerInfoReq) (rspMsg []exchanger.ResponseMsg) {
	// 返回消息
	response := &hall.HallGetPlayerInfoRsp{}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_PLAYER_INFO_RSP),
			Body:  response,
		},
	}

	response.ErrCode = proto.Uint32(0)
	response.Coin = proto.Uint64(player.GetPlayerCoin(playerID))
	response.NickName = proto.String(player.GetPlayerNickName(playerID))
	state, gameID := getPlayerState(playerID)
	response.PlayerState = state.Enum()
	response.GameId = common.GameId(gameID).Enum()
	return
}

// HandleGetPlayerStateReq 获取玩家是否正在游戏中
func HandleGetPlayerStateReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerStateReq) (rspMsg []exchanger.ResponseMsg) {
	userData := req.GetUserData()
	response := &hall.HallGetPlayerStateRsp{
		UserData: &userData,
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP),
			Body:  response,
		},
	}
	state, gameID := getPlayerState(playerID)
	response.PlayerState = state.Enum()
	response.GameId = common.GameId(gameID).Enum()

	logrus.WithFields(logrus.Fields{
		"func_name": "HandleGetPlayerStateReq",
		"player_id": playerID,
		"response":  response,
	}).Infoln("获取玩家状态")
	return
}

// HandleGetGameInfoReq client-> 获取游戏信息列表请求
func HandleGetGameInfoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetGameListInfoReq) (rspMsg []exchanger.ResponseMsg) {
	// 返回消息
	response := &hall.HallGetGameListInfoRsp{
		ErrCode: proto.Uint32(1),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_PLAYER_INFO_RSP),
			Body:  response,
		},
	}

	// 逻辑处理
	gameInfos, gameLevelInfos, err := data.GetGameInfoList()
	if err == nil {
		response.GameConfig = ServerGameConfig2Client(gameInfos)
		response.GameLevelConfig = ServerGameLevelConfig2Client(gameLevelInfos)
	}

	logrus.WithFields(logrus.Fields{
		"func_name": "HandleGetGameInfoList",
		"player_id": playerID,
		"response":  response,
	}).Infoln("获取游戏配置")
	return
}
