package data

import (
	"encoding/json"
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/gutils"
	"steve/server_pb/user"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// idAllocObject id分配
var idAllocObject *gutils.Node

// showUID 最大展示uid
var showUID = "max_show_uid"

// redis 过期时间
var redisTimeOut = time.Hour * 24

// 玩家基本信息列表
var playerInfoList = map[int32]string{
	1: "nickname",
	2: "avatar",
	3: "gender",
	4: "name",
	5: "phone",
	6: "idCard",
}

// gameconfigList 游戏配置
var gameconfigList = map[int16]string{
	1: "gameID",
	2: "name",
	3: "type",
}

const (
	playerRedisName          = "player"
	playerMysqlName          = "player"
	playerTableName          = "t_player"
	playerCurrencyTableName  = "t_player_currency"
	playerGameTableName      = "t_player_game"
	gameconfigTableName      = "t_game_config"
	gamelevelconfigTableName = "t_game_level_config"
)

type gameConfigDetail struct {
	db.TGameConfig      `xorm:"extends"`
	db.TGameLevelConfig `xorm:"extends"`
}

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
		err = fmt.Errorf("select sql err：err=%v", err)
		return
	}
	if exist {
		playerID = dbPlayerID.ID
		if err := setRedisVal(playerRedisName, redisKey, playerID, redisTimeOut); err != nil {
			err = fmt.Errorf("save playerId into redis fail： %v", err)
		}
	}
	return
}

// GetPlayerInfo 根据玩家id获取玩家个人资料信息
func GetPlayerInfo(playerID uint64, fields ...string) (dbPlayer *db.TPlayer, err error) {
	logrus.Debugln("get player info playerId :%d, fields:%s", playerID, fields)

	dbPlayer, err = new(db.TPlayer), nil

	// 从redis获取
	dbPlayer, err = getPlayerFieldsFromRedis(playerID, fields)
	if err == nil {
		return
	}

	// 从数据库获取
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return
	}
	strCol := ""
	for _, col := range fields {
		if len(strCol) > 0 {
			strCol += ","
		}
		strCol += col
	}

	sql := fmt.Sprintf("select %s from t_player  where playerID='%d';", strCol, playerID)
	res, err := engine.QueryString(sql)
	if err != nil {
		err = fmt.Errorf("select sql err：sql=%s,err=%v", sql, err)
		return
	}
	if len(res) != 1 {
		err = fmt.Errorf("玩家存在多条信息记录： %v", err)
		return
	}
	dbPlayer, err = generateDbPlayer(playerID, res[0], fields...)
	if err != nil {
		err = fmt.Errorf("generate dbPlayer 失败(%v)", err.Error())
	}

	// 更新redis
	if err = updatePlayerFieldsToRedis(playerID, fields, dbPlayer); err != nil {
		err = fmt.Errorf("更新 redis 失败(%v)", err.Error())
	}
	return
}

// GetPlayerGameInfo 获取玩家游戏信息
func GetPlayerGameInfo(playerID uint64, gameID uint32, fields ...string) (exist bool, dbPlayerGame *db.TPlayerGame, err error) {
	logrus.Debugf("get player game info playerId :%d, gameID:%d, fields:%v", playerID, gameID, fields)

	exist, dbPlayerGame, err = true, new(db.TPlayerGame), nil

	// 从redis获取
	dbPlayerGame, err = getPlayerGameFieldsFromRedis(playerID, gameID, fields)
	if err == nil {
		return
	}

	// 从数据库获取
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return
	}
	strCol := ""
	for _, col := range fields {
		if len(strCol) > 0 {
			strCol += ","
		}
		strCol += col
	}

	sql := fmt.Sprintf("select %s from t_player_game  where playerID='%d' and gameID='%d';", strCol, playerID, gameID)
	res, err := engine.QueryString(sql)

	if err != nil {
		err = fmt.Errorf("select t_player_game sql:%s ,err：%v", sql, err)
		return
	}

	if len(res) == 0 {
		exist = false
		return
	}

	if len(res) != 1 {
		err = fmt.Errorf("玩家存在多条 gameId:%d 信息记录： %v", gameID, err)
		return
	}

	dbPlayerGame, err = generateDbPlayerGame(playerID, gameID, res[0], fields...)
	if err != nil {
		err = fmt.Errorf("generate dbPlayerGame 失败(%v)", err.Error())
	}

	// 更新redis
	if err = updatePlayerGameFieldsToRedis(playerID, gameID, fields, dbPlayerGame); err != nil {
		err = fmt.Errorf("更新 redis 失败(%v)", err.Error())
	}
	return
}

// GetPlayerState 获取游戏状态,游戏id,ip地址
func GetPlayerState(playerID uint64, fields ...string) (pState *PlayerState, err error) {
	enrty := logrus.WithFields(logrus.Fields{
		"func_name": GetPlayerState,
		"playerID":  playerID,
	})

	pState, err = getPlayerStateFromRedis(playerID, fields...)

	if err != nil {
		enrty.WithError(err).Warningln("get player state from redis fail")
		return pState, err
	}
	return
}

// UpdatePlayerState 修改玩家游戏状态
func UpdatePlayerState(playerID uint64, oldState, newState, gameID, levelID uint32) (result bool, err error) {
	enrty := logrus.WithFields(logrus.Fields{
		"func_name": "update_player_state",
		"playerID":  playerID,
		"oldState":  oldState,
		"newState":  newState,
	})

	result, err = true, nil
	redisKey := cache.FmtPlayerIDKey(uint64(playerID))

	// 校验玩家当前状态是否一致
	val, _ := getRedisField(playerRedisName, redisKey, cache.GameState)
	if len(val) != 0 && val[0] != nil {
		state, _ := strconv.Atoi(val[0].(string))
		if oldState != uint32(state) {
			return
		}
	}

	rfields := map[string]string{
		cache.GameState: fmt.Sprintf("%d", newState),
		cache.GameID:    fmt.Sprintf("%d", gameID),
		cache.LevelID:   fmt.Sprintf("%d", levelID),
	}

	if err = setPlayerStateByWatch(playerRedisName, redisKey, oldState, rfields, redisTimeOut); err != nil {
		err = fmt.Errorf("save playerInfo  into redis fail： %v", err)
	}
	enrty.WithError(err).Warningln("update_player_state finish")
	return
}

// UpdatePlayerGateInfo 修改玩家网关服信息
func UpdatePlayerGateInfo(playerID uint64, idAddr, gateAddr string) (result bool, err error) {
	enrty := logrus.WithFields(logrus.Fields{
		"func_name": "update_player_gateInfo",
		"playerID":  playerID,
		"idAddr":    idAddr,
		"gateAddr":  gateAddr,
	})

	result, err = true, nil

	redisCli, err := redisCliGetter(playerRedisName, 0)
	if err != nil {
		return false, fmt.Errorf("获取 redis 客户端失败(%s)。", err.Error())
	}
	playerKey := cache.FmtPlayerIDKey(playerID)

	kv := map[string]interface{}{
		cache.IPAddr:   idAddr,
		cache.GateAddr: gateAddr,
	}

	status := redisCli.HMSet(playerKey, kv)
	if status.Err() != nil {
		return false, fmt.Errorf("设置失败(%v)", status.Err())
	}

	enrty.WithError(err).Warningln("update_player_gateInfo finish")
	return
}

// UpdatePlayerServerAddr 修改玩家服务端地址
func UpdatePlayerServerAddr(playerID uint64, serverType uint32, serverAddr string) (result bool, err error) {
	enrty := logrus.WithFields(logrus.Fields{
		"func_name":  "update_player_server_addr",
		"playerID":   playerID,
		"serverType": serverType,
		"serverAddr": serverAddr,
	})

	result, err = true, nil

	redisCli, err := redisCliGetter(playerRedisName, 0)
	if err != nil {
		return false, fmt.Errorf("获取 redis 客户端失败(%s)。", err.Error())
	}
	playerKey := cache.FmtPlayerIDKey(playerID)

	serverField := map[user.ServerType]string{
		user.ServerType_ST_GATE:  cache.GateAddr,
		user.ServerType_ST_MATCH: cache.MatchAddr,
		user.ServerType_ST_ROOM:  cache.RoomAddr,
	}[user.ServerType(serverType)]

	kv := map[string]interface{}{
		serverField: serverAddr,
	}
	status := redisCli.HMSet(playerKey, kv)
	if status.Err() != nil {
		return false, fmt.Errorf("设置失败(%v)", status.Err())
	}

	enrty.WithError(err).Warningln("update_player_serveraddr finish")
	return
}

// setRedisFieldByWatch 修改玩家状态（事务）
func setPlayerStateByWatch(redisName string, key string, oldState uint32, fields map[string]string, duration time.Duration) error {
	redisCli, err := redisCliGetter(redisName, 0)

	list := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		list[k] = v
	}

	err = redisCli.Watch(func(tx *redis.Tx) error {
		err := tx.HKeys(key).Err()
		if err != nil && err != redis.Nil {
			return err
		}
		stateString := tx.HGet(key, cache.GameState).Val()
		stateInt, _ := strconv.Atoi(stateString)

		if uint32(stateInt) != oldState {
			err = fmt.Errorf("修改玩家游戏状态出错，玩家当前状态不为：%d", oldState)
			return err
		}
		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.HMSet(key, list)
			return nil
		})
		return err
	}, key)
	if err == nil {
		redisCli.Expire(key, duration)
	}
	return err
}

// GetGameInfoList 获取游戏配置信息
func GetGameInfoList() (gameConfig []*db.TGameConfig, gamelevelConfig []*db.TGameLevelConfig, err error) {
	gameConfig, gamelevelConfig, err = make([]*db.TGameConfig, 0), make([]*db.TGameLevelConfig, 0), nil

	gameConfigKey := "gameconfig"
	gameLevelConfigKey := "gamelevelconfig"

	rKey := cache.FmtGameInfoConfigKey()

	// 从redis获取
	val, err := getRedisField(playerRedisName, rKey, []string{gameConfigKey, gameLevelConfigKey}...)
	if err == nil && len(val) == 2 {
		if val[0] != nil && val[0].(string) != "" {
			json.Unmarshal([]byte(val[0].(string)), gameConfig)
		}
		if val[1] != nil && val[1].(string) != "" {
			json.Unmarshal([]byte(val[1].(string)), gamelevelConfig)
		}
	}
	if len(gameConfig) != 0 && len(gamelevelConfig) != 0 {
		return
	}

	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return
	}
	strCol := "id,gameID,name,type"
	err = engine.Table(gameconfigTableName).Select(strCol).Find(&gameConfig)

	if err != nil {
		err = fmt.Errorf("select sql error： %v", err)
		return
	}

	strCol = "id,gameID,levelID,name,fee,baseScores,lowScores,highScores,minPeople,maxPeople,status,tag,remark"
	err = engine.Table(gamelevelconfigTableName).Select(strCol).Find(&gamelevelConfig)

	if err != nil {
		err = fmt.Errorf("select sql error： %v", err)
		return
	}

	// 更新redis
	gameConfigData, _ := json.Marshal(gameConfig)
	gameLevelConfigData, _ := json.Marshal(gamelevelConfig)
	rFields := map[string]string{
		gameConfigKey:      string(gameConfigData),
		gameLevelConfigKey: string(gameLevelConfigData),
	}
	if err = setRedisFields(playerRedisName, rKey, rFields, redisTimeOut); err != nil {
		err = fmt.Errorf("save game_config  into redis fail： %v", err)
	}
	return
}

// AllocPlayerID 生成玩家 ID
func AllocPlayerID() uint64 {
	return uint64(idAllocObject.Generate().Int64())
}

// AllocShowUID 生成玩家展示id(10位数), 暂时从redis生成
func AllocShowUID() int64 {
	r, _ := redisCliGetter(playerRedisName, 0)
	if r.Exists(showUID).Val() == 0 {
		r.Set(showUID, 10000*10000*100, -1)
	}
	return r.Incr(showUID).Val()
}

// InitPlayerData 初始化玩家数据
func InitPlayerData(player db.TPlayer) error {
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return err
	}
	session := engine.Table(playerTableName)
	affected, err := session.Insert(&player)
	if err != nil || affected == 0 {
		sql, _ := session.LastSQL()
		return fmt.Errorf("insert sql error：%v， affect=%d, sql=%s", err, affected, sql)
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
		return fmt.Errorf("insert t_player_cuccency sql：%v， affect=%d", err, affected)
	}
	return nil
}

// InitPlayerState 初始化玩家状态
func InitPlayerState(playerID int64) (err error) {
	redisKey := cache.FmtPlayerIDKey(uint64(playerID))

	rfields := map[string]string{
		cache.GameState: fmt.Sprintf("%d", user.PlayerState_PS_IDIE),
		cache.IPAddr:    "",
		cache.GateAddr:  "",
	}

	if err = setRedisFields(playerRedisName, redisKey, rfields, redisTimeOut); err != nil {
		err = fmt.Errorf("save player_state into redis fail： %v", err)
	}
	return
}

// getPlayerStateFromRedis 从redis查找玩家状态信息
func getPlayerStateFromRedis(playerID uint64, fields ...string) (*PlayerState, error) {

	redisCli, err := redisCliGetter(playerRedisName, 0)
	if err != nil {
		return nil, fmt.Errorf("获取 redis 客户端失败(%s)。", err.Error())
	}
	playerKey := cache.FmtPlayerIDKey(playerID)
	cmd := redisCli.HMGet(playerKey, fields...)
	if cmd.Err() != nil {
		return nil, fmt.Errorf("执行 redis 命令失败(%s)", cmd.Err().Error())
	}
	result, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("获取 redis 结果失败(%s) fields=(%v)", err.Error(), fields)
	}
	var playerState PlayerState
	for index, field := range fields {
		v, ok := result[index].(string)
		if !ok {
			continue
		}
		if err = setPlayerStateByField(&playerState, field, v); err != nil {
			return nil, err
		}
	}
	redisCli.Expire(playerKey, redisTimeOut)
	return &playerState, nil
}

// SetPlayerFields 设置玩家指定字段值
func SetPlayerFields(playerID uint64, fields []string, dbPlayer *db.TPlayer) error {
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return fmt.Errorf("获取 mysql 引擎失败(%s)", err.Error())
	}
	_, err = engine.Table(playerTableName).Where("playerID = ?", playerID).Cols(fields...).Update(dbPlayer)
	if err != nil {
		return fmt.Errorf("更新失败 (%v)", err.Error())
	}
	if err = updatePlayerFieldsToRedis(playerID, fields, dbPlayer); err != nil {
		return fmt.Errorf("更新 redis 失败(%v)", err.Error())
	}
	return nil
}

func getPlayerFieldsFromRedis(playerID uint64, fields []string) (*db.TPlayer, error) {
	redisCli, err := redisCliGetter(playerRedisName, 0)
	if err != nil {
		return nil, fmt.Errorf("获取 redis 客户端失败(%s)。", err.Error())
	}
	playerKey := cache.FmtPlayerIDKey(playerID)
	cmd := redisCli.HMGet(playerKey, fields...)
	if cmd.Err() != nil {
		return nil, fmt.Errorf("执行 redis 命令失败(%s)", cmd.Err().Error())
	}
	result, err := cmd.Result()
	if err != nil || len(result) != len(fields) {
		return nil, fmt.Errorf("获取 redis 结果失败(%s) fields=(%v)", err.Error(), fields)
	}
	var dbPlayer db.TPlayer
	for index, field := range fields {
		v, ok := result[index].(string)
		if !ok {
			return nil, fmt.Errorf("错误的数据类型。field=%s val=%v", field, result[index])
		}
		if err = setDBPlayerByField(&dbPlayer, field, v); err != nil {
			return nil, err
		}
	}
	redisCli.Expire(playerKey, redisTimeOut)
	return &dbPlayer, nil
}

func getPlayerGameFieldsFromRedis(playerID uint64, gameID uint32, fields []string) (*db.TPlayerGame, error) {
	redisCli, err := redisCliGetter(playerRedisName, 0)
	if err != nil {
		return nil, fmt.Errorf("获取 redis 客户端失败(%s)。", err.Error())
	}
	playerGameKey := cache.FmtPlayerGameInfoKey(playerID, gameID)
	cmd := redisCli.HMGet(playerGameKey, fields...)
	if cmd.Err() != nil {
		return nil, fmt.Errorf("执行 redis 命令失败(%s)", cmd.Err().Error())
	}
	result, err := cmd.Result()
	if err != nil || len(result) != len(fields) {
		return nil, fmt.Errorf("获取 redis 结果失败(%s) fields=(%v)", err.Error(), fields)
	}
	var dbPlayerGame db.TPlayerGame
	for index, field := range fields {
		v, ok := result[index].(string)
		if !ok {
			return nil, fmt.Errorf("错误的数据类型。field=%s val=%v", field, result[index])
		}
		if err = setDBPlayerGameByField(&dbPlayerGame, field, v); err != nil {
			return nil, err
		}
	}
	redisCli.Expire(playerGameKey, redisTimeOut)
	return &dbPlayerGame, nil
}

// setDBPlayerFieldByName 设置 dbPlayer 中的指定字段
func setDBPlayerByField(dbPlayer *db.TPlayer, field string, val string) error {
	switch field {
	case "id":
		dbPlayer.Id, _ = strconv.ParseInt(val, 10, 64)
	case "accountID":
		dbPlayer.Accountid, _ = strconv.ParseInt(val, 10, 64)
	case "playerID":
		dbPlayer.Playerid, _ = strconv.ParseInt(val, 10, 64)
	case "type":
		dbPlayer.Type, _ = strconv.Atoi(val)
	case "channelID":
		dbPlayer.Channelid, _ = strconv.Atoi(val)
	case "nickname":
		dbPlayer.Nickname = val
	case "gender":
		dbPlayer.Gender, _ = strconv.Atoi(val)
	case "avatar":
		dbPlayer.Avatar = val
	case "provinceID":
		dbPlayer.Provinceid, _ = strconv.Atoi(val)
	case "cityID":
		dbPlayer.Cityid, _ = strconv.Atoi(val)
	case "name":
		dbPlayer.Name = val
	case "phone":
		dbPlayer.Phone = val
	case "idCard":
		dbPlayer.Idcard = val
	case "isWhiteList":
		dbPlayer.Iswhitelist, _ = strconv.Atoi(val)
	case "zipCode":
		dbPlayer.Zipcode, _ = strconv.Atoi(val)
	case "shippingAddr":
		dbPlayer.Shippingaddr = val
	case "status":
		dbPlayer.Status, _ = strconv.Atoi(val)
	case "remark":
	case "createTime":
	case "createBy":
	case "updateTime":
	case "updateBy":
		return nil
	default:
		return fmt.Errorf("未处理的字段:%s", field)
	}
	return nil
}

// setDBPlayerGameFieldByName 设置 dbPlayerGame 中的指定字段
func setDBPlayerGameByField(dbPlayerGame *db.TPlayerGame, field string, val string) error {
	switch field {
	case "id":
		dbPlayerGame.Id, _ = strconv.ParseInt(val, 10, 64)
	case "userID":
		dbPlayerGame.Playerid, _ = strconv.ParseInt(val, 10, 64)
	case "gameID":
		dbPlayerGame.Gameid, _ = strconv.Atoi(val)
	case "gameName":
		dbPlayerGame.Gamename = val
	case "winningRate":
		dbPlayerGame.Winningrate, _ = strconv.Atoi(val)
	case "winningBurea":
		dbPlayerGame.Winningburea, _ = strconv.Atoi(val)
	case "totalBureau":
		dbPlayerGame.Totalbureau, _ = strconv.Atoi(val)
	case "maxWinningStream":
		dbPlayerGame.Maxwinningstream, _ = strconv.Atoi(val)
	case "maxMultiple":
		dbPlayerGame.Maxmultiple, _ = strconv.Atoi(val)
	case "createTime":
	case "createBy":
	case "updateTime":
	case "updateBy":
		return nil
	default:
		return fmt.Errorf("未处理的字段:%s", field)
	}
	return nil
}

// setPlayerStateByField 设置玩家状态
func setPlayerStateByField(playerState *PlayerState, field string, val string) error {
	switch field {
	case cache.GameState:
		playerState.State, _ = strconv.ParseUint(val, 10, 64)
	case cache.GameID:
		playerState.GameID, _ = strconv.ParseUint(val, 10, 64)
	case cache.LevelID:
		playerState.LevelID, _ = strconv.ParseUint(val, 10, 64)
	case cache.IPAddr:
		playerState.IPAddr = val
	case cache.GateAddr:
		playerState.GateAddr = val
	case cache.MatchAddr:
		playerState.MatchAddr = val
	case cache.RoomAddr:
		playerState.RoomAddr = val
	default:
		return fmt.Errorf("未处理的字段:%s", field)
	}
	return nil
}

func getDBPlayerField(field string, dbPlayer *db.TPlayer) (interface{}, error) {
	var v interface{}
	switch field {
	case "id":
		v = dbPlayer.Id
	case "accountID":
		v = dbPlayer.Accountid
	case "playerID":
		v = dbPlayer.Playerid
	case "type":
		v = dbPlayer.Type
	case "channelID":
		v = dbPlayer.Channelid
	case "nickname":
		v = dbPlayer.Nickname
	case "gender":
		v = dbPlayer.Gender
	case "avatar":
		v = dbPlayer.Avatar
	case "provinceID":
		v = dbPlayer.Provinceid
	case "cityID":
		v = dbPlayer.Cityid
	case "name":
		v = dbPlayer.Name
	case "phone":
		v = dbPlayer.Phone
	case "idCard":
		v = dbPlayer.Idcard
	case "isWhiteList":
		v = dbPlayer.Iswhitelist
	case "zipCode":
		v = dbPlayer.Zipcode
	case "shippingAddr":
		v = dbPlayer.Shippingaddr
	case "status":
		v = dbPlayer.Status
	case "remark":
		v = dbPlayer.Remark
	case "createTime":
		v = dbPlayer.Createtime
	case "createBy":
		v = dbPlayer.Createby
	case "updateTime":
		v = dbPlayer.Updatetime
	case "updateBy":
		v = dbPlayer.Updateby
	default:
		return nil, fmt.Errorf("不能识别的字段: %s", field)
	}
	return v, nil
}

func getDBPlayerGameField(field string, dbPlayerGame *db.TPlayerGame) (interface{}, error) {
	var v interface{}
	switch field {
	case "id":
		v = dbPlayerGame.Id
	case "playerID":
		v = dbPlayerGame.Playerid
	case "gameID":
		v = dbPlayerGame.Gameid
	case "gameName":
		v = dbPlayerGame.Gamename
	case "winningRate":
		v = dbPlayerGame.Winningrate
	case "winningBurea":
		v = dbPlayerGame.Winningburea
	case "totalBureau":
		v = dbPlayerGame.Totalbureau
	case "maxWinningStream":
		v = dbPlayerGame.Maxwinningstream
	case "maxMultiple":
		v = dbPlayerGame.Maxmultiple
	case "createTime":
		v = dbPlayerGame.Createtime
	case "createBy":
		v = dbPlayerGame.Createby
	case "updateTime":
		v = dbPlayerGame.Updatetime
	case "updateBy":
		v = dbPlayerGame.Updateby
	default:
		return nil, fmt.Errorf("不能识别的字段: %s", field)
	}
	return v, nil
}

func updatePlayerFieldsToRedis(playerID uint64, fields []string, dbPlayer *db.TPlayer) error {
	redisCli, err := redisCliGetter(playerRedisName, 0)
	if err != nil {
		return fmt.Errorf("获取 redis 客户端失败(%s)。", err.Error())
	}
	playerKey := cache.FmtPlayerIDKey(playerID)
	kv := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		v, err := getDBPlayerField(field, dbPlayer)
		if err != nil {
			return err
		}
		if v == nil {
			continue
		}
		kv[field] = v
	}
	status := redisCli.HMSet(playerKey, kv)
	if status.Err() != nil {
		return fmt.Errorf("设置失败(%v)", status.Err())
	}
	redisCli.Expire(playerKey, redisTimeOut)
	return nil
}

func updatePlayerGameFieldsToRedis(playerID uint64, gameID uint32, fields []string, dbPlayerGame *db.TPlayerGame) error {
	redisCli, err := redisCliGetter(playerRedisName, 0)
	if err != nil {
		return fmt.Errorf("获取 redis 客户端失败(%s)。", err.Error())
	}
	playerGameKey := cache.FmtPlayerGameInfoKey(playerID, gameID)
	kv := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		v, err := getDBPlayerGameField(field, dbPlayerGame)
		if err != nil {
			return err
		}
		if v == nil {
			continue
		}
		kv[field] = v
	}
	status := redisCli.HMSet(playerGameKey, kv)
	if status.Err() != nil {
		return fmt.Errorf("设置失败(%v)", status.Err())
	}
	redisCli.Expire(playerGameKey, redisTimeOut)
	return nil
}

// GetPlayerTodayCharge 获取玩家今日充值数量
func GetPlayerTodayCharge(playerID uint64) (uint64, error) {
	playerKey := cache.FmtPlayerChargeKey(playerID)
	vals, err := getRedisField(playerRedisName, playerKey, cache.TodayChargeKey, cache.LastChargeTime)
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	chargeStr, _ := vals[0].(string)
	charge, _ := strconv.ParseUint(chargeStr, 10, 16)
	if charge == uint64(0) {
		return 0, nil
	}
	chargeTimeStr, _ := vals[1].(string)
	chargeTime, err := time.Parse(time.UnixDate, chargeTimeStr)
	if err == nil {
		year, month, day := chargeTime.Date()
		cyear, cmonth, cday := time.Now().Date()
		if year == cyear && month == cmonth && day == cday {
			return charge, nil
		}
	}
	return 0, nil
}

// AddPlayerTodayCharge 添加玩家今日充值数
// 非重要的数据，不保证数据一致性
func AddPlayerTodayCharge(playerID uint64, charge uint64) error {
	todayCharge, err := GetPlayerTodayCharge(playerID)
	if err != nil {
		return fmt.Errorf("获取今日充值数失败")
	}
	todayCharge += charge
	lastChargeTimeStr := time.Now().Format(time.UnixDate)
	playerKey := cache.FmtPlayerChargeKey(playerID)
	return setRedisFields(playerRedisName, playerKey, map[string]string{
		cache.TodayChargeKey: strconv.FormatUint(todayCharge, 10),
		cache.LastChargeTime: lastChargeTimeStr,
	}, time.Hour*24)
}

func init() {
	node := viper.GetInt("node")
	var err error
	idAllocObject, err = gutils.NewNode(int64(node))
	if err != nil {
		logrus.Panicf("创建 id 生成器失败: %v", err)
	}
}
