package user

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/external/goldclient"
	"steve/hall/data"
	"steve/server_pb/gold"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HandleGetPlayerInfoReq 处理获取玩家个人资料请求
func HandleGetPlayerInfoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerInfoReq) (rspMsg []exchanger.ResponseMsg) {
	logrus.Debugf("Handle get player info req:%v", req)

	// 默认返回消息
	response := &hall.HallGetPlayerInfoRsp{
		ErrCode: proto.Uint32(1),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_PLAYER_INFO_RSP),
			Body:  response,
		},
	}

	// 获取玩家基本个人资料
	player, err := data.GetPlayerInfo(playerID, cache.ShowUID, cache.NickName, cache.Avatar, cache.Gender, cache.Name, cache.IDCard)
	if err == nil {
		response.ErrCode = proto.Uint32(0)
		response.NickName = proto.String(player.Nickname)
		response.Avator = proto.String(player.Avatar)
		response.ShowUid = proto.Uint64(uint64(player.Showuid))
		response.Gender = proto.Uint32(uint32(player.Gender))
		if player.Name != "" && player.Idcard != "" {
			response.RealnameStatus = proto.Uint32(1)
		} else {
			response.RealnameStatus = proto.Uint32(0)
		}
	}

	// 获取玩家货币信息
	coin, err := goldclient.GetGold(playerID, int16(gold.GoldType_GOLD_COIN))
	if err == nil {
		response.Coin = proto.Uint64(uint64(coin))
	}

	// 获取玩家游戏信息
	pState, err := data.GetPlayerState(playerID, []string{cache.GameState, cache.GameID}...)
	if err == nil {
		response.PlayerState = common.PlayerState(pState.State).Enum()
		response.GameId = common.GameId(pState.GameID).Enum()
		response.ErrCode = proto.Uint32(0)
	}
	logrus.Debugf("Handle get player info rsp: %v", response)
	return
}

// HandleUpdatePlayerInoReq 处理更新玩家个人资料请求
func HandleUpdatePlayerInoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallUpdatePlayerInfoReq) (rspMsg []exchanger.ResponseMsg) {
	logrus.Debugf("Handle update player info req: %v", req)

	// 默认返回消息
	response := &hall.HallUpdatePlayerInfoRsp{
		ErrCode: proto.Uint32(1),
		// Result:  proto.Bool(false),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_UPDATE_PLAYER_INFO_RSP),
			Body:  response,
		},
	}

	// 参数校验
	dbPlayer := db.TPlayer{
		Playerid: int64(playerID),
	}
	gender := uint32(req.GetGender())
	if gender == 1 || gender == 2 {
		dbPlayer.Gender = int(req.GetGender())
	}

	if req.NickName != nil && req.GetNickName() != "" {
		dbPlayer.Nickname = req.GetNickName()
	}
	if req.Avator != nil && req.GetAvator() != "" {
		dbPlayer.Avatar = req.GetAvator()

	}

	// 逻辑处理
	fields := []string{cache.NickName, cache.Avatar, cache.Gender}
	error := data.SetPlayerFields(playerID, fields, &dbPlayer)

	if error == nil {
		response.ErrCode = proto.Uint32(0)
		response.NickName = proto.String(req.GetNickName())
		response.Gender = req.GetGender().Enum()
		response.Avator = proto.String(req.GetAvator())
	}

	logrus.Debugf("Handle update player info rsp: %v", response)
	return
}

// HandleGetPlayerStateReq 获取玩家游戏状态信息
func HandleGetPlayerStateReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerStateReq) (rspMsg []exchanger.ResponseMsg) {
	logrus.Debugf("Handle get player state req:%v", req)

	// 默认返回消息
	response := &hall.HallGetPlayerStateRsp{
		UserData: proto.Uint64(req.GetUserData()),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP),
			Body:  response,
		},
	}

	// 逻辑处理
	fields := []string{cache.GameState, cache.GameID}
	pState, err := data.GetPlayerState(playerID, fields...)

	// 返回结果
	if err == nil {
		response.PlayerState = common.PlayerState(pState.State).Enum()
		response.GameId = common.GameId(pState.GameID).Enum()
	}
	logrus.Debugf("Handle get player state rsp:%v", response)
	return
}

// HandleGetGameInfoReq client-> 获取游戏信息列表请求
func HandleGetGameInfoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetGameListInfoReq) (rspMsg []exchanger.ResponseMsg) {
	logrus.Debugf("Handle get game info req : %v", req)

	// 默认返回消息
	response := &hall.HallGetGameListInfoRsp{
		ErrCode: proto.Uint32(1),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_GAME_INFO_RSP),
			Body:  response,
		},
	}

	// 逻辑处理
	gameInfos, gameLevelInfos, err := data.GetGameInfoList()

	// 返回结果
	if err == nil {
		response.GameConfig = DBGameConfig2Client(gameInfos)
		response.GameLevelConfig = DBGamelevelConfig2Client(gameLevelInfos)
		response.ErrCode = proto.Uint32(0)

	}
	logrus.Debugf("Handle get game info rsp:%v ", response)

	return
}

// HandleGetPlayerGameInfoReq 获取玩家游戏信息
func HandleGetPlayerGameInfoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerGameInfoReq) (rspMsg []exchanger.ResponseMsg) {
	logrus.Debugf("Handle get player game info req : %v", req)

	// 传入参数
	uid := req.GetUid()
	gameID := req.GetGameId()

	// 默认返回消息
	response := &hall.HallGetPlayerGameInfoRsp{
		Uid:     proto.Uint64(uid),
		GameId:  common.GameId(gameID).Enum(),
		ErrCode: proto.Uint32(1),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_GET_PLAYER_GAME_INFO_RSP),
			Body:  response,
		},
	}

	// 逻辑处理
	fields := []string{cache.TotalBurea, cache.WinningRate, cache.MaxWinningStream, cache.MaxMultiple}
	exist, dbPlayerGame, err := data.GetPlayerGameInfo(uid, uint32(gameID), fields...)

	// 不存在直接返回
	if !exist {
		return
	}

	// 返回结果
	if err == nil {
		response.TotalBureau = proto.Uint32(uint32(dbPlayerGame.Totalbureau))
		response.WinningRate = proto.Float32(float32(dbPlayerGame.Winningrate))
		response.MaxWinningStream = proto.Uint32(uint32(dbPlayerGame.Maxwinningstream))
		response.MaxMultiple = proto.Uint32(uint32(dbPlayerGame.Maxmultiple))
	}
	logrus.Debugf("Handle get player game info rsp:%v ", response)

	return
}
