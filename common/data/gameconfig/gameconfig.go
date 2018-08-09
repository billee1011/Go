package gameconfig

import (
	"encoding/json"
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
)

// redis 过期时间
var redisTimeOut = time.Hour * 24

const (
	playerRedisName          = "player"
	playerMysqlName          = "player"
	configRedisName          = "config"
	configMysqlName          = "config"
	playerTableName          = "t_player"
	playerCurrencyTableName  = "t_player_currency"
	playerGameTableName      = "t_player_game"
	gameconfigTableName      = "t_game_config"
	gamelevelconfigTableName = "t_game_level_config"
)

// GetPlayerInfo 根据玩家id获取玩家个人资料信息
func GetPlayerInfo(playerID uint64, fields ...string) (dbPlayer *db.TPlayer, err error) {
	logrus.Debugln("get player info playerId :(%d), fields:(%s)", playerID, fields)

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
		err = fmt.Errorf("select sql err：sql=(%s),err=(%v)", sql, err)
		return
	}
	if len(res) != 1 {
		err = fmt.Errorf("玩家存在多条信息记录：(%v)", err)
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
	logrus.Debugf("get player game info playerId :(%d), gameID:(%d), fields:(%v)", playerID, gameID, fields)

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
		err = fmt.Errorf("select t_player_game sql:(%s) ,err：(%v)", sql, err)
		return
	}

	if len(res) == 0 {
		exist = false
		return
	}

	if len(res) != 1 {
		err = fmt.Errorf("玩家存在多条 gameId:(%d) 信息记录： %v", gameID, err)
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

// GetGameInfoList 获取游戏配置信息
func GetGameInfoList() (gameConfig []*db.TGameConfig, gamelevelConfig []*db.TGameLevelConfig, funcErr error) {
	gameConfig, gamelevelConfig, funcErr = make([]*db.TGameConfig, 0), make([]*db.TGameLevelConfig, 0), nil

	gameConfigKey := "gameconfig"
	gameLevelConfigKey := "gamelevelconfig"

	rKey := cache.FmtGameInfoConfigKey()

	// 从redis获取
	val, rerr := getRedisField(configRedisName, rKey, []string{gameConfigKey, gameLevelConfigKey}...)
	if rerr == nil && len(val) == 2 {
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

	engine, merr := mysqlEngineGetter(configMysqlName)
	if merr != nil {
		funcErr = fmt.Errorf("get mysql enginer error：(%v)", merr.Error())
		return
	}
	strCol := "id,gameID,name,type,minPeople,maxPeople"
	funcErr = engine.Table(gameconfigTableName).Select(strCol).Find(&gameConfig)

	if funcErr != nil {
		funcErr = fmt.Errorf("select sql error：(%v)", funcErr.Error())
		return
	}

	strCol = "id,gameID,levelID,name,fee,baseScores,lowScores,highScores,realOnlinePeople,showOnlinePeople,status,tag,remark"
	funcErr = engine.Table(gamelevelconfigTableName).Select(strCol).Find(&gamelevelConfig)

	if funcErr != nil {
		funcErr = fmt.Errorf("select sql error：(%v)", funcErr.Error())
		return
	}

	// 更新redis
	gameConfigData, _ := json.Marshal(gameConfig)
	gameLevelConfigData, _ := json.Marshal(gamelevelConfig)
	rFields := map[string]string{
		gameConfigKey:      string(gameConfigData),
		gameLevelConfigKey: string(gameLevelConfigData),
	}
	if funcErr = setRedisFields(playerRedisName, rKey, rFields, redisTimeOut); funcErr != nil {
		funcErr = fmt.Errorf("save game_config  into redis fail：(%v)", funcErr.Error())
	}
	return
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
			return nil, fmt.Errorf("错误的数据类型。field=(%s) val=(%v)", field, result[index])
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
			return nil, fmt.Errorf("错误的数据类型。field=(%s) val=(%v)", field, result[index])
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
	case "showUID":
		dbPlayer.Showuid, _ = strconv.ParseInt(val, 10, 64)
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
		return fmt.Errorf("未处理的字段:(%s)", field)
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
		dbPlayerGame.Winningrate, _ = strconv.ParseFloat(val, 64)
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
		return fmt.Errorf("未处理的字段:(%s)", field)
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
		return fmt.Errorf("未处理的字段:(%s)", field)
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
	case "showUID":
		v = dbPlayer.Showuid
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
		return nil, fmt.Errorf("不能识别的字段:(%s)", field)
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
		return nil, fmt.Errorf("不能识别的字段:(%s)", field)
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
