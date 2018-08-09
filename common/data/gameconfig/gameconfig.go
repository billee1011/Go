package gameconfig

import (
	"encoding/json"
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/structs"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
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

func getRedisCli(redis string, db int) (*redis.Client, error) {
	exposer := structs.GetGlobalExposer()
	redisCli, err := exposer.RedisFactory.GetRedisClient(redis, db)
	if err != nil {
		return nil, fmt.Errorf("获取 redis client 失败: %v", err)
	}
	return redisCli, nil
}

func getMysqlEngine(mysqlName string) (*xorm.Engine, error) {
	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(mysqlName)
	if err != nil {
		return nil, fmt.Errorf("获取 mysql 引擎失败：%v", err)
	}
	return engine, nil
}

// 单元测试通过这两个值修改 mysql 引擎获取和 redis cli 获取
var mysqlEngineGetter = getMysqlEngine
var redisCliGetter = getRedisCli

func getRedisField(redisName string, key string, field ...string) ([]interface{}, error) {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return nil, err
	}
	result, err := redisCli.HMGet(key, field...).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		return nil, fmt.Errorf("redis 命令执行失败: %v", err)
	}
	return result, nil
}

func setRedisFields(redisName string, key string, fields map[string]string, duration time.Duration) error {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return err
	}
	kv := make(map[string]interface{}, len(fields))
	for k, field := range fields {
		kv[k] = field
	}
	status := redisCli.HMSet(key, kv)
	if status.Err() != nil {
		return fmt.Errorf("设置失败(%v)", status.Err())
	}
	redisCli.Expire(key, duration)
	return nil
}
