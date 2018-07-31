package data

import (
	"encoding/json"
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/gutils"
	"steve/server_pb/user"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// idAllocObject id分配
var idAllocObject *gutils.Node

// redis 过期时间
var redisTimeOut = time.Hour * 24 * 30

// 需要获得的玩家就基本信息
var playerInfoList = map[int16]string{
	1: "nickname",
	2: "gender",
	3: "avatar",
	4: "name",
	5: "phone",
}

// gameconfigList 游戏配置
var gameconfigList = map[int16]string{
	1: "gameID",
	2: "name",
	3: "type",
}

// gameconfigList 类型配置
var gameLevelconfigList = map[int16]string{
	1: "levelID",
	2: "baseScores",
	3: "lowScores",
	4: "highScores",
	5: "minPeople",
	6: "maxPeople",
	7: "status",
}

const (
	playerRedisName          = "player"
	playerMysqlName          = "player"
	playerTableName          = "t_player"
	playerCurrencyTableName  = "t_player_currency"
	gameconfigTableName      = "t_game_config"
	gamelevelconfigTableName = "t_game_level_config"
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

	pinfo, err := loadPlayerInfoFromRedis(playerID, playerRedisName)
	if err == nil {
		trans2hallPlayer(&cp, pinfo)
		return
	}
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return
	}
	strCol := ""
	for _, col := range playerInfoList {
		if len(strCol) > 0 {
			strCol += ","
		}
		strCol += col
	}

	sql := fmt.Sprintf("select %s from t_player  where playerID='%d';", strCol, playerID)
	res, err := engine.QueryString(sql)
	if err != nil {
		err = fmt.Errorf("mysql 操作失败： %v", err)
		return
	}
	if len(res) != 1 {
		err = fmt.Errorf("玩家存在多条信息记录： %v", err)
		return
	}
	trans2hallPlayer(&cp, res[0])

	if err = SavePlayerInfoToRedis(playerID, res[0], playerRedisName); err != nil {
		err = fmt.Errorf("get playerInfo save redis失败： %v", err)
	}
	return
}

// UpdatePlayerInfo 修改玩家个人信息
func UpdatePlayerInfo(playerID uint64, nickName, avatar string) (exist, result bool, err error) {
	entry := logrus.WithFields(logrus.Fields{
		"opr":      "update_pinfo",
		"playerID": playerID,
		"nickName": nickName,
		"avatar":   avatar,
	})
	exist, result, err = true, true, nil

	rfields := map[string]interface{}{
		cache.NickNameField: nickName,
		cache.AvatarField:   avatar,
	}

	if err != nil {
		return
	}
	strCol := "playerID="
	strCol += fmt.Sprintf("'%v'", playerID)
	for key, field := range rfields {
		strCol += ","
		strCol += key
		strCol += "="
		strCol += fmt.Sprintf("'%v'", field)
	}
	engine, err := mysqlEngineGetter(playerMysqlName)
	sql := fmt.Sprintf("update t_player set %s  where playerID=?;", strCol)
	res, sqlerror := engine.Exec(sql, playerID)
	if sqlerror != nil {
		entry.WithError(err).Errorln("update t_player mysql 操作失败")
		exist, result, err = true, false, sqlerror
	}
	if aff, aerr := res.RowsAffected(); aff == 0 {
		entry.WithError(err).Errorln("update t_player playerId:%d 不存在", playerID)
		exist, result, err = false, false, aerr
	}

	list := make(map[string]string, len(rfields))
	for key, field := range rfields {
		list[key] = field.(string)
	}
	if err = SavePlayerInfoToRedis(playerID, list, playerRedisName); err != nil {
		err = fmt.Errorf("get playerInfo save redis 失败： %v", err)
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
	redisKey := cache.FmtPlayerIDKey(playerID)
	state, err = hgetRedisUint64Val(playerRedisName, redisKey, cache.PlayerStateField)
	if err == nil {
		enrty.WithError(err).Warningln("get player_state 失败")
		err = fmt.Errorf("redis 操作失败： %v", err)
	}
	return
}

// UpdatePlayerState 修改玩家游戏状态
func UpdatePlayerState(playerID, oldState, newState uint64, serverType int32, serverAddr string) (result bool, err error) {
	result, err = true, nil
	redisKey := cache.FmtPlayerIDKey(playerID)

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
	list := make(map[string]string, 0)
	list[cache.PlayerStateField] = string(newState)
	list[playerServerField] = serverAddr

	if err = SavePlayerInfoToRedis(playerID, list, playerRedisName); err != nil {
		err = fmt.Errorf("get playerInfo save redis失败： %v", err)
	}
	return
}

// GetGameInfoList 获取游戏配置信息
func GetGameInfoList() (gameInfos []*user.GameInfo, err error) {
	gameInfos, err = []*user.GameInfo{}, nil

	gameDetails := make([]gameDetail, 0)

	redisKey := cache.FmtGameInfoKey()

	data, err := getRedisByteVal(playerRedisName, redisKey)
	if data != nil && len(data) != 0 {
		err = json.Unmarshal(data, &gameDetails)
		if err == nil {
			gameInfos = trans2GameInfo(gameDetails)
			return
		}
	}
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return
	}
	selectSQL := ""
	for _, gameconfig := range gameconfigList {
		selectSQL += gameconfigTableName
		selectSQL += "."
		selectSQL += gameconfig
		selectSQL += ","
	}
	selectSQL += gamelevelconfigTableName
	selectSQL += "."
	selectSQL += "gameID"
	for _, gamelevel := range gameLevelconfigList {
		selectSQL += ","
		selectSQL += gamelevelconfigTableName
		selectSQL += "."
		selectSQL += gamelevel
	}

	err = engine.Table(gameconfigTableName).Select(selectSQL).
		Join("INNER", gamelevelconfigTableName, "t_game_config.gameID = t_game_level_config.gameID ").
		Find(gameDetails)

	if err != nil {
		err = fmt.Errorf("mysql 操作失败： %v", err)
		return
	}
	if len(gameDetails) == 0 {
		err = fmt.Errorf("游戏配置数据不存在")
		return
	}
	// 写入redis
	data, err = json.Marshal(gameDetails)
	if err != nil {
		err = fmt.Errorf("写入redis 操作失败： %v", err)
	}
	if err := setRedisVal(playerRedisName, redisKey, data, redisTimeOut); err != nil {
		logrus.WithFields(logrus.Fields{
			"key": redisKey,
			"val": gameDetails,
		}).Warningln("设置 redis 失败")
	}
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
	list := make(map[string]string, 0)
	list[cache.NickNameField] = player.Nickname
	list[cache.AvatarField] = player.Avatar
	list[cache.GenderField] = string(player.Gender)
	list[cache.NameField] = player.Name

	if err = SavePlayerInfoToRedis(uint64(player.Playerid), list, playerRedisName); err != nil {
		err = fmt.Errorf("get playerInfo save redis失败： %v", err)
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

// InitPlayerState 初始化玩家状态
func InitPlayerState(playerID int64) (err error) {
	redisKey := cache.FmtPlayerIDKey(uint64(playerID))
	redisFields := map[string]interface{}{
		cache.PlayerStateField: fmt.Sprintf("%d", user.PlayerState_PS_IDIE),
	}
	if err = hmSetRedisVal(playerRedisName, redisKey, redisFields, redisTimeOut); err != nil {
		logrus.WithFields(logrus.Fields{
			"key": redisKey,
			"val": playerID,
		}).Warningln("设置 redis 失败")
	}
	return
}

// GetPlayerInfoFields 获取palyerID属性的field值
func GetPlayerInfoFields() (playerFields []string) {
	playerFields = make([]string, 0)
	// 昵称
	playerFields = append(playerFields, cache.NickNameField)
	// 头像
	playerFields = append(playerFields, cache.AvatarField)
	// 金币
	playerFields = append(playerFields, cache.CoinField)

	return
}

// loadPlayerInfoFromRedis 从redis查找用户信息
func loadPlayerInfoFromRedis(playerID uint64, redisName string) (map[string]string, error) {

	r, err := redisCliGetter(redisName, 0)
	if err != nil {
		return nil, err
	}

	redisKey := cache.FmtPlayerIDKey(playerID)

	cmd := r.HGetAll(redisKey)
	if cmd.Err() != nil {
		return nil, fmt.Errorf("get redis err:%v", cmd.Err())
	}
	m := cmd.Val()
	if len(m) == 0 {
		return nil, fmt.Errorf("redis no user: playerID=%d", playerID)
	}
	list := make(map[string]string, len(m))
	for k, v := range m {
		sp := strings.Split(k, "_")
		if len(sp) == 2 {
			k = sp[1]
		}
		list[k] = v
	}

	return list, nil
}

// SavePlayerInfoToRedis 玩家信息保存到redis
func SavePlayerInfoToRedis(playerID uint64, pinfo map[string]string, redisName string) error {
	r, err := redisCliGetter(redisName, 0)
	if err != nil {
		return err
	}

	redisKey := cache.FmtPlayerIDKey(playerID)
	list := make(map[string]interface{}, len(pinfo))
	for k, v := range pinfo {
		list[k] = v
	}
	cmd := r.HMSet(redisKey, list)
	if cmd.Err() != nil {
		return fmt.Errorf("set redis err:%v", cmd.Err())
	}
	r.Expire(redisKey, redisTimeOut)
	return nil
}

func init() {
	node := viper.GetInt("node")
	var err error
	idAllocObject, err = gutils.NewNode(int64(node))
	if err != nil {
		logrus.Panicf("创建 id 生成器失败: %v", err)
	}
}
