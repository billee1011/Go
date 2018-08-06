package data

import (
	"encoding/json"
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/gutils"
	"steve/server_pb/user"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// idAllocObject id分配
var idAllocObject *gutils.Node

// redis 过期时间
var redisTimeOut = time.Hour * 24 * 30

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

// gameconfigList 类型配置
var gameLevelconfigList = map[int16]string{
	1: "levelID",
	2: "name",
	3: "baseScores",
	4: "lowScores",
	5: "highScores",
	6: "minPeople",
	7: "maxPeople",
	8: "status",
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
		if err := setRedisVal(playerRedisName, redisKey, playerID, time.Hour*24); err != nil {
			err = fmt.Errorf("save playerId into redis fail： %v", err)
		}
	}
	return
}

// GetPlayerFields 获取玩家的指定字段值
func GetPlayerFields(playerID uint64, fields []string) (*db.TPlayer, error) {
	if dbPlayer, err := getPlayerFieldsFromRedis(playerID, fields); err == nil {
		return dbPlayer, nil
	}
	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return nil, fmt.Errorf("获取 mysql 引擎失败(%s)", err.Error())
	}
	var dbPlayer db.TPlayer
	exist, err := engine.Table(playerTableName).Where("playerID = ?", playerID).Cols(fields...).Get(&dbPlayer)
	if !exist || err != nil {
		return nil, fmt.Errorf("获取数据失败。exist=%v, err=%s", exist, err.Error())
	}
	if err = updatePlayerFieldsToRedis(playerID, fields, &dbPlayer); err != nil {
		logrus.WithFields(logrus.Fields{
			"player_id": playerID,
			"fields":    fields,
		}).WithError(err).Errorln("更新 redis 失败")
		// 因为拿到数据了，所以不返回失败
	}
	return &dbPlayer, nil
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
	redisCli.Expire(playerKey, time.Hour*24)
	return &dbPlayer, nil
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
	redisCli.Expire(playerKey, time.Hour*24)
	return nil
}

// GetPlayerInfo 根据玩家id获取玩家的基本信息
func GetPlayerInfo(playerID uint64) (info map[string]string, err error) {
	info, err = map[string]string{}, nil

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
	dbPlayer = generateDbPlayer(playerID, res[0])
	return
}

// UpdatePlayerInfo 修改玩家个人信息
func UpdatePlayerInfo(playerID uint64, nickName, avatar string, gender uint32) (exist, result bool, err error) {
	entry := logrus.WithFields(logrus.Fields{
		"opr":      "update_player_info",
		"playerID": playerID,
		"nickName": nickName,
		"avatar":   avatar,
	})
	exist, result, err = true, true, nil

	tbPlayer := db.TPlayer{
		Nickname: nickName,
		Gender:   int(gender),
		Avatar:   avatar,
	}

	engine, err := mysqlEngineGetter(playerMysqlName)
	affected, uerr := engine.Update(&tbPlayer, db.TPlayer{Playerid: int64(playerID)})

	if uerr != nil {
		entry.WithError(err).Errorln("update t_player mysql fail")
		exist, result, err = true, false, uerr
		return
	}
	if affected == 0 {
		entry.WithError(err).Errorln("update t_player playerId:%d 不存在", playerID)
		exist, result, err = false, false, nil
		return
	}

	return
}

// GetPlayerGameInfo 获取玩家游戏信息
func GetPlayerGameInfo(playerID uint64, gameID uint32) (exist bool, info *db.TPlayerGame, err error) {
	exist, info, err = false, new(db.TPlayerGame), nil

	engine, err := mysqlEngineGetter(playerMysqlName)

	where := fmt.Sprintf("playerID=%d and gameID='%d'", playerID, gameID)
	exist, err = engine.Table(playerGameTableName).Select("gameID").Where(where).Get(info)

	if err != nil {
		err = fmt.Errorf("select t_player_game sql err：%v", err)
		return
	}
	return
}

// GetPlayerState 获取游戏状态,游戏id,ip地址
func GetPlayerState(playerID uint64) (pState *PlayerState, err error) {
	enrty := logrus.WithFields(logrus.Fields{
		"func_name": GetPlayerState,
		"playerID":  playerID,
	})
	pState, err = new(PlayerState), nil

	val, err := loadFromRedis(playerID, playerRedisName)

	if err != nil {
		enrty.WithError(err).Warningln("get player state from redis fail")
		return pState, err
	}
	pState.generatePlayerState(val)
	return
}

// UpdatePlayerState 修改玩家游戏状态
func UpdatePlayerState(playerID uint64, oldState, newState, reqServerType uint32, serverAddr string) (result bool, err error) {
	result, err = true, nil
	redisKey := cache.FmtPlayerIDKey(uint64(playerID))

	val, _ := getRedisField(playerRedisName, redisKey, cache.GameState)
	state, _ := strconv.Atoi(val[0].(string))

	if oldState != uint32(state) {
		return
	}

	serverType := map[user.ServerType]string{
		user.ServerType_ST_GATE:  cache.GateAddr,
		user.ServerType_ST_MATCH: cache.MatchAddr,
		user.ServerType_ST_ROOM:  cache.RoomAddr,
	}[user.ServerType(reqServerType)]

	rfields := map[string]string{
		cache.GameState: fmt.Sprintf("%d", newState),
		serverType:      serverAddr,
	}
	if err = setRedisWatch(playerRedisName, redisKey, rfields, redisTimeOut); err != nil {
		err = fmt.Errorf("save playerInfo  into redis fail： %v", err)
	}
	return
}

// GetGameInfoList 获取游戏配置信息
func GetGameInfoList() (gameInfos []*user.GameConfig, gamelevelInfos []*user.GameLevelConfig, err error) {
	gameInfos, gamelevelInfos, err = make([]*user.GameConfig, 0), make([]*user.GameLevelConfig, 0), nil

	gameConfigKey := "gameconfig"
	gameLevelConfigKey := "gamelevelconfig"

	var dbgameConfigs []db.TGameConfig
	var dbgamelevelConfigs []db.TGameLevelConfig

	gameConfigdata, err := getRedisByteVal(playerRedisName, gameConfigKey)
	if gameConfigdata != nil && len(gameConfigdata) != 0 {
		err = json.Unmarshal(gameConfigdata, &dbgameConfigs)
	}
	gameLeveldata, err := getRedisByteVal(playerRedisName, gameLevelConfigKey)
	if gameConfigdata != nil && len(gameConfigdata) != 0 {
		err = json.Unmarshal(gameLeveldata, &dbgamelevelConfigs)
	}
	if err == nil {
		dbGameConfig2serverGameConfig(dbgameConfigs)
		dbGamelevelConfig2serverGameConfig(dbgamelevelConfigs)
		return
	}

	engine, err := mysqlEngineGetter(playerMysqlName)
	if err != nil {
		return
	}
	err = engine.Table(gameconfigTableName).Find(&dbgameConfigs)

	if err != nil {
		err = fmt.Errorf("select sql error： %v", err)
		return
	}

	err = engine.Table(gamelevelconfigTableName).Find(&dbgamelevelConfigs)

	if err != nil {
		err = fmt.Errorf("select sql error： %v", err)
		return
	}
	dbGameConfig2serverGameConfig(dbgameConfigs)
	dbGamelevelConfig2serverGameConfig(dbgamelevelConfigs)
	// 写入redis
	data, _ := json.Marshal(dbgameConfigs)
	if err = setRedisVal(playerRedisName, gameConfigKey, data, redisTimeOut); err != nil {
		err = fmt.Errorf("save game_config  into redis fail： %v", err)
	}
	data, _ = json.Marshal(dbgamelevelConfigs)
	if err = setRedisVal(playerRedisName, gameLevelConfigKey, data, redisTimeOut); err != nil {
		err = fmt.Errorf("save game_level_config  into redis fail： %v", err)
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
		return fmt.Errorf("insert sql error：%v， affect=%d", err, affected)
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
		cache.IPAddr:    fmt.Sprintf("%s", "127.0.0.1"),
	}

	if err = setRedisWatch(playerRedisName, redisKey, rfields, redisTimeOut); err != nil {
		err = fmt.Errorf("save player_state into redis fail： %v", err)
	}
	return
}

// loadFromRedis 从redis查找信息
func loadFromRedis(playerID uint64, redisName string) (map[string]string, error) {

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
