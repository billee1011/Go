package user

import (
	"bytes"
	"io/ioutil"
	"math"
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/common/data/prop"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/external/configclient"
	"steve/external/goldclient"
	"steve/hall/data"
	"steve/server_pb/gold"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// HandleGetPlayerInfoReq 处理获取玩家个人资料请求
func HandleGetPlayerInfoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerInfoReq) (rspMsg []exchanger.ResponseMsg) {
	logrus.Debugf("Handle get player info req:(%v)", req)

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
	logrus.Debugf("Handle get player info rsp: (%v)", response)
	return
}

// HandleUpdatePlayerInoReq 处理更新玩家个人资料请求
func HandleUpdatePlayerInoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallUpdatePlayerInfoReq) (rspMsg []exchanger.ResponseMsg) {
	logrus.Debugf("Handle update player info req: (%v)", req)

	// 默认返回消息
	response := &hall.HallUpdatePlayerInfoRsp{
		ErrCode: proto.Uint32(1),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_HALL_UPDATE_PLAYER_INFO_RSP),
			Body:  response,
		},
	}

	// 参数校验
	fields := []string{}
	dbPlayer := db.TPlayer{
		Playerid: int64(playerID),
	}
	gender := uint32(req.GetGender())
	if gender == 1 || gender == 2 {
		dbPlayer.Gender = int(req.GetGender())
		fields = append(fields, cache.Gender)
	}

	if req.NickName != nil && req.GetNickName() != "" {
		if !validateNickName(req.GetNickName()) {
			response.ErrCode = proto.Uint32(1)
			return
		}
		dbPlayer.Nickname = req.GetNickName()
		fields = append(fields, cache.NickName)
	}

	if req.Avator != nil && req.GetAvator() != "" {
		dbPlayer.Avatar = req.GetAvator()
		fields = append(fields, cache.Avatar)
	}

	// 逻辑处理
	error := data.SetPlayerFields(playerID, fields, &dbPlayer)

	if error == nil {
		response.ErrCode = proto.Uint32(0)
		response.NickName = proto.String(req.GetNickName())
		response.Gender = req.GetGender().Enum()
		response.Avator = proto.String(req.GetAvator())
	}

	logrus.Debugf("Handle update player info rsp: (%v)", response)
	return
}

// validateNickName 校验昵称
func validateNickName(nickName string) bool {
	if len(nickName) <= 16 {
		return true
	}
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(nickName)), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		logrus.Debugf("validateNickName utf8 transfer gbk err: (%v)", err.Error())
		return false
	}
	count := len(data)
	if count > 16 {
		return false
	}
	return true
}

// HandleGetPlayerStateReq 获取玩家游戏状态信息
func HandleGetPlayerStateReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetPlayerStateReq) (rspMsg []exchanger.ResponseMsg) {
	logrus.Debugf("Handle get player state req: (%v)", req)

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
	logrus.Debugf("Handle get player state rsp: (%v)", response)
	return
}

// HandleGetGameInfoReq client-> 获取游戏信息列表请求
func HandleGetGameInfoReq(playerID uint64, header *steve_proto_gaterpc.Header, req hall.HallGetGameListInfoReq) (rspMsg []exchanger.ResponseMsg) {
	logrus.Debugf("Handle get game info req : (%v)", req)

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

	// 游戏配置
	gameConf, err := configclient.GetGameConfigMap()
	if err != nil {
		logrus.WithError(err).Errorln("获取游戏配置失败！！")
		return
	}
	// 场次配置
	levelConf, err := configclient.GetGameLevelConfigMap()
	if err != nil {
		logrus.WithError(err).Errorln("获取游戏级别配置失败！！")
		return
	}

	// 返回结果
	if err == nil {
		response.GameConfig = DBGameConfig2Client(gameConf)
		response.GameLevelConfig = DBGamelevelConfig2Client(levelConf)
		response.ErrCode = proto.Uint32(0)
		return
	}
	logrus.Debugf("Handle get game info rsp: (%v),err :(%v) ", response, err.Error())

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
		Uid:          proto.Uint64(uid),
		GameId:       common.GameId(gameID).Enum(),
		UserProperty: make([]*common.Property, 0),
		ErrCode:      proto.Uint32(1),
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
	if !exist && playerID == uid {
		response.ErrCode = proto.Uint32(0)
		return
	}

	// 出错直接返回
	if err != nil {
		logrus.Debugf("Handle get player game info rsp:(%v),err:(%v) ", response, err.Error())
		return
	}

	if exist {
		response.TotalBureau = proto.Uint32(uint32(dbPlayerGame.Totalbureau))
		response.WinningRate = proto.Float32(float32(dbPlayerGame.Winningrate))
		response.MaxWinningStream = proto.Uint32(uint32(dbPlayerGame.Maxwinningstream))
		response.MaxMultiple = proto.Uint32(uint32(dbPlayerGame.Maxmultiple))
	}

	// 获取自己游戏信息直接返回
	if playerID == uid {
		response.ErrCode = proto.Uint32(0)
		return
	}

	// 获取玩家道具
	props, err := prop.GetPlayerAllProps(playerID)
	if err != nil {
		logrus.Debugf("Handle get player game info uid:(%d)获取玩家道具失败 err:(%v)", uid, err.Error())
		return
	}

	propIds := make([]int32, 0)
	propCount := make(map[int32]int64, len(props))
	for _, prop := range props {
		propIds = append(propIds, prop.PropID)
		propCount[prop.PropID] = prop.Count
	}

	// 获取道具属性
	propConfigs, err := prop.GetSomePropsConfig(propIds)

	logrus.Debugf("Handle get player game info propConfigs:(%v),", propConfigs)

	if err != nil {
		logrus.Debugf("Handle get player game info uid:(%d)获取玩家道具属性失败 err:(%v)", uid, err.Error())
		return
	}

	for _, propConfig := range propConfigs {
		userProperty := new(common.Property)
		userProperty.PropId = proto.Int32(propConfig.PropID)
		userProperty.PropName = proto.String(propConfig.PropName)
		userProperty.PropType = common.PropType(propConfig.Type).Enum()
		userProperty.PropCost = proto.Int64(int64(math.Abs(float64(propConfig.Value))))
		userProperty.PropCount = proto.Uint32(uint32(propCount[propConfig.PropID]))
		response.UserProperty = append(response.UserProperty, userProperty)
		response.ErrCode = proto.Uint32(0)
	}

	logrus.Debugf("Handle get player game info uid:(%d) rsp:(%v)", uid, response)
	return
}
