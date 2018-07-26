package player

import (
	"fmt"
	"steve/entity/cache"
	"steve/hall/data"
	"steve/server_pb/hall"
	"sync"

	"github.com/Sirupsen/logrus"
)

const (
	// 玩家金币数字段名
	playerCoinField string = "coin"
	// playerGatewayAddrField 玩家网关地址字段名
	playerGatewayAddrField string = "gate_addr"
	// playerRoomAddrField 玩家所在 room 地址字段名
	playerRoomAddrField string = "room_addr"
	// playerRoomAddrField 玩家所在 match 地址字段名
	playeMatchAddrField string = "match_addr"
	// playerStateField 玩家状态字段名
	playerStateField string = "state"
	// playerGameIDField 玩家游戏 ID 字段名
	playerGameIDField string = "game_id"

	// playerNickNameField 玩家昵称字段
	playerNickNameField string = "nick_name"
	// playerHeadImageField 玩家头像字段
	playerHeadImageField string = "head_image"

	// maxPlayerIDField 玩家已分配的最大playerId
	maxPlayerIDField string = "max_player_id"
)

var mux sync.Mutex

// Login 处理玩家登录
func Login(accountID uint64) (uint64, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "Login",
		"accountID": accountID,
	})
	playerID, err := data.GetAccountPlayerID(accountID)
	if err != nil {
		entry.WithError(err).Errorln("玩家登录失败")
		return 0, err
	}
	if playerID == 0 {
		playerID = newPlayer(accountID)
	}
	return playerID, nil
}

// GetPlayerInfo 获取玩家个人信息请求
func GetPlayerInfo(playerID uint64) (cache.HallPlayer, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetPlayerInfo",
		"playerID":  playerID,
	})
	hallPlayer := new(cache.HallPlayer)
	playerInfo, err := data.GetPlayerFields(playerID, playerNickNameField, playerHeadImageField, playerStateField, playerCoinField)
	if err != nil {
		entry.WithError(err).Errorln("获取玩家信息失败")
		return *hallPlayer, err
	}
	hallPlayer.NickName = playerInfo[playerNickNameField].(string)   // 昵称
	hallPlayer.HeadImage = playerInfo[playerHeadImageField].(string) // 头像
	hallPlayer.State = playerInfo[playerStateField].(uint64)         // 玩家状态
	hallPlayer.Coin = playerInfo[playerCoinField].(uint64)           // 金币
	return *hallPlayer, nil
}

// GetPlayerState 获取玩家状态
func GetPlayerState(playerID uint64) (uint64, error) {
	return data.GetPlayerUint64Filed(playerID, playerStateField) // 玩家状态
}

// UpdatePlayerInfo 更新玩家基本信息
func UpdatePlayerInfo(player cache.HallPlayer) bool {
	playerID := player.PlayerID
	if player.NickName != "" {
		data.SetPlayerFiled(playerID, playerNickNameField, player.NickName)
	}
	if player.HeadImage != "" {
		data.SetPlayerFiled(playerID, playerHeadImageField, player.HeadImage)
	}
	return true
}

// UpdatePlayerState 更新玩家状态
func UpdatePlayerState(playerID, oldState, newState uint64, serverType int32, serverAddr string) bool {
	currentState, _ := data.GetPlayerUint64Filed(playerID, playerStateField)
	if oldState != currentState {
		return false
	}
	playerServerField := map[hall.ServerType]string{
		hall.ServerType_ST_GATE:  playerGatewayAddrField,
		hall.ServerType_ST_MATCH: playeMatchAddrField,
		hall.ServerType_ST_ROOM:  playerRoomAddrField,
	}[hall.ServerType(serverType)]

	if playerServerField == "" {
		return false
	}
	playerFields := map[string]interface{}{
		playerStateField:  newState,
		playerServerField: serverType,
	}
	err := data.SetPlayerFields(playerID, playerFields)
	return err == nil
}

// newPlayer 创建玩家
func newPlayer(accountID uint64) uint64 {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":  "newPlayer",
		"account_id": accountID,
	})
	playerID, err := data.AllocIDIncr(maxPlayerIDField)
	if err != nil {
		entry.WithError(err).Errorln("分配玩家 ID 失败")
		return 0
	}
	if err := data.NewPlayer(accountID, playerID); err != nil {
		entry.WithError(err).Errorln("创建玩家失败")
		return 0
	}
	initPlayer(playerID)
	return playerID
}

// initPlayer 初始化玩家基本信息
func initPlayer(playerID uint64) {
	playerField := map[string]interface{}{
		playerCoinField:     10000,
		playerNickNameField: fmt.Sprintf("玩家%v", playerID),
		playerStateField:    hall.PlayerState_PS_IDIE,
	}
	data.SetPlayerFields(playerID, playerField)
}
