package user

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/entity/cache"
	"steve/hall/data"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HandleGetPlayerInfoReq2 处理获取玩家信息请求
func HandleGetPlayerInfoReq2(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerInfoReq) (rspMsg []exchanger.ResponseMsg) {
	// 返回消息
	response := &hall.HallGetPlayerInfoRsp{}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_PLAYER_INFO_RSP),
			Body:  response,
		},
	}
	dbPlayer, err := data.GetPlayerFields(playerID, []string{"nickname", "name", "idCard"})
	if err != nil {
		response.ErrCode = proto.Uint32(uint32(common.ErrCode_EC_FAIL))
		return
	}
	if dbPlayer.Name != "" && dbPlayer.Idcard != "" {
		response.RealnameStatus = proto.Uint32(1)
	} else {
		response.RealnameStatus = proto.Uint32(0)
	}

	playerInfo, err := data.GetPlayerInfo(playerID)
	if err == nil {
		response.ErrCode = proto.Uint32(0)
		response.Coin = proto.Uint64(0)
		response.NickName = proto.String(playerInfo[cache.NickNameField])
	}

	state, gameID, err := data.GetPlayerState(playerID)
	if err == nil {
		response.PlayerState = common.PlayerState(state).Enum()
		response.GameId = common.GameId(gameID).Enum()
	}
	return
}

// HandleGetPlayerStateReq2 获取玩家是否正在游戏中
func HandleGetPlayerStateReq2(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerStateReq) (rspMsg []exchanger.ResponseMsg) {
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
	state, gameID, _ := data.GetPlayerState(playerID)
	response.PlayerState = common.PlayerState(state).Enum()
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
