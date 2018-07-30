package data

import (
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/gutils"
	"steve/server_pb/user"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper" 
)

var idAllocObject *gutils.Node

const (
	playerRedisName         = "player"
	playerMysqlName         = "player"
	playerTableName         = "t_player"
	playerCurrencyTableName = "t_player_currency"
)

// GetPlayerIDByAccountID 根据账号 ID 获取其关联的玩家 ID
func GetPlayerIDByAccountID(accountID uint64) (exist bool, playerID uint64, err error) {
	exist, playerID, err = false, 0, nil

	redisKey := cache.FmtAccountPlayerKey(accountID)
	playerID, err = getRedisUint64Val(playerRedisName, redisKey)
	if err == nil {
		return
	}
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return
	}
	where := fmt.Sprintf("accountID=%d", accountID)
	var dbPlayerID struct {
		ID uint64 `xorm:"playerID"`
	}
	exist, err = engine.Table(playerTableName).Select("playerID").Where(where).Get(&dbPlayerID)
	if err != nil {
		fmt.Println(err.Error())
		err = fmt.Errorf("mysql 操作失败： %v", err)
		return
	}
	if exist {
		playerID = dbPlayerID.ID
		if err := setRedisVal(playerRedisName, redisKey, playerID, time.Hour*24); err != nil {
			logrus.WithFields(logrus.Fields{
				"key": redisKey,
				"val": playerID,
			}).Warningln("设置 redis 失败")
		}
	}
	return
}

// GetPlayerInfoByPlayerID 根据玩家id获取玩家的基本信息
func GetPlayerInfoByPlayerID(playerID uint64) (cp cache.HallPlayer, err error) {
	cp, err = cache.HallPlayer{}, nil

	redisKey := cache.FmtPlayerKey(playerID) 
	pInfokeys := GetPlayerInfoFields()
	playerInfo := make(map[string]interface{}, 0)

	playerInfo, err = hmGetRedisFields(playerRedisName, redisKey, pInfokeys...)
	if err == nil {
		cp.NickName = playerInfo[cache.NickNameField].(string)
		cp.HeadImage = playerInfo[cache.HeadImageField].(string)
		cp.Coin = playerInfo[cache.CoinField].(uint64) 
		return
	}
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return
	}
	where := fmt.Sprintf("playerID=%d", playerID)
	var dbPlayer struct {
		NickName  string `xorm:"nickname"`
		HeadImage string `xorm:"avatar"`
	}
	err = engine.Table(playerTableName).Select("nickname,avatar").Where(where).Find(&dbPlayer)
	if err != nil {
		err = fmt.Errorf("mysql 操作失败： %v", err)
		return
	}
	cp.NickName = dbPlayer.NickName
	cp.HeadImage = dbPlayer.HeadImage

	redisFields := map[string]interface{}{
		cache.NickNameField:  cp.NickName,
		cache.HeadImageField: cp.HeadImage,
	}
	if err := hmSetRedisVal(playerRedisName, redisKey, redisFields); err != nil {
		logrus.WithFields(logrus.Fields{
			"key":    redisKey,
			"fields": redisFields,
			"val":    playerID,
		}).Warningln("设置 redis 失败")
	}
	return
}

// SetPlayerInfo 设置玩家个人信息
func SetPlayerInfo(cp cache.HallPlayer) (exist, result bool, err error) {
	exist, result, err = false, false, nil

	playerID := cp.PlayerID
	redisKey := cache.FmtPlayerKey(playerID)
	pfields := setPlayerInfoFields(cp)

	err = hmSetRedisVal(playerRedisName, redisKey, pfields)
	if err != nil {
		return
	}
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return
	}
	where := fmt.Sprintf("playerID=%d", playerID)
	var id int
	exist, err = engine.Table(playerTableName).Where(where).Get(&id)
	if err != nil {
		err = fmt.Errorf("mysql 操作失败： %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("玩家不存在： %v", playerID)
		return
	}
	dbPlayer := db.TPlayer{
		Playerid: int64(cp.PlayerID),
		Nickname: cp.NickName,
		Avatar:   cp.HeadImage,
	}
	err = UpdatePlayerData(dbPlayer)
	if err == nil {
		exist, result, err = true, true, nil
	}
	return
}

// GetPlayerState 获取游戏状态
func GetPlayerState(playerID uint64) (state uint64, err error) {
	enrty := logrus.WithFields(logrus.Fields{
		"func_name": GetPlayerState,
		"playerID":  playerID,
	})
	state, err = 0, nil
	redisKey := cache.FmtPlayerKey(playerID)
	state, err = hgetRedisUint64Val(playerRedisName, redisKey, cache.PlayerStateField)
	if err == nil {
		enrty.WithError(err).Warningln("获取游戏状态 失败")
		err = fmt.Errorf("redis 操作失败： %v", err)
	}
	return
}

// SetPlayerState 设置玩家游戏状态
func SetPlayerState(playerID, oldState, newState uint64, serverType int32, serverAddr string) (result bool, err error) {
	result, err = false, nil
	redisKey := cache.FmtPlayerKey(playerID)

	currentState, _ := hgetRedisUint64Val(playerRedisName, redisKey, cache.PlayerStateField)
	if oldState != currentState {
		return
	}
	playerServerField := map[user.ServerType]string{ 
		user.ServerType_ST_GATE:  cache.GateAddrField,
		user.ServerType_ST_MATCH: cache.MatchAddrField,
		user.ServerType_ST_ROOM:  cache.RoomAddrField,
	}[user.ServerType(serverType)]

	if playerServerField == "" {
		return
	}
	redisFields := map[string]interface{}{
		cache.PlayerStateField: newState,
		playerServerField:      serverType,
	} 
	err = hmSetRedisVal(playerRedisName, redisKey, redisFields)
	if err != nil {
		err = fmt.Errorf("redis 操作失败： %v", err)
	}
	return
}

// GetGameListInfo 获取游戏列表
func GetGameListInfo() (gameID []uint64, err error) {

	return
}

// AllocPlayerID 生成玩家 ID
func AllocPlayerID() uint64 {
	return uint64(idAllocObject.Generate().Int64())
}

// InitPlayerData 初始化玩家数据
func InitPlayerData(player db.TPlayer) error {
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return err
	}
	affected, err := engine.Table(playerTableName).Insert(&player)
	if err != nil || affected == 0 {
		return fmt.Errorf("插入数据失败：%v， affect=%d", err, affected)
	}
	return nil
}

// InitPlayerCoin 初始化玩家货币信息
func InitPlayerCoin(currency db.TPlayerCurrency) error {
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return err
	}
	affected, err := engine.Table(playerCurrencyTableName).Insert(&currency)
	if err != nil || affected == 0 {
		return fmt.Errorf("插入数据失败：%v， affect=%d", err, affected)
	}
	return nil
}

// UpdatePlayerData 更新玩家数据
func UpdatePlayerData(player db.TPlayer) error {
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return err
	}
	where := fmt.Sprintf("playerID=%d", player.Playerid)
	affected, err := engine.Table(playerTableName).Where(where).Update(&player)
	if err != nil || affected == 0 {
		return fmt.Errorf("更新数据失败：%v， affect=%d", err, affected)
	}
	return nil
}

// GetPlayerInfoFields 获取palyerID属性的field值
func GetPlayerInfoFields() (playerFields []string) {
	playerFields = make([]string, 0)
	// 昵称
	playerFields = append(playerFields, cache.NickNameField)
	// 头像
	playerFields = append(playerFields, cache.HeadImageField)
	// 金币
	playerFields = append(playerFields, cache.CoinField)

	return
}

// setPlayerInfoFields 设置玩家属性值
func setPlayerInfoFields(cp cache.HallPlayer) (fields map[string]interface{}) {
	fields = map[string]interface{}{
		cache.NickNameField:  cp.NickName,
		cache.HeadImageField: cp.HeadImage,
	}
	return
}

func init() {
	node := viper.GetInt("node")
	var err error
	idAllocObject, err = gutils.NewNode(int64(node))
	if err != nil {
		logrus.Panicf("创建 id 生成器失败: %v", err)
	}
}
