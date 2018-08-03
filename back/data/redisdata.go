package data

import (
	"fmt"
	"steve/structs"

	"github.com/go-redis/redis"
)

const redisName = "back"

func getRedisCli(redis string, db int) (*redis.Client, error) {
	exposer := structs.GetGlobalExposer()
	redisCli, err := exposer.RedisFactory.GetRedisClient(redis, db)
	if err != nil {
		return nil, fmt.Errorf("获取 redis client 失败: %v", err)
	}
	return redisCli, nil
}

// RedisCliGetter 单元测试通过这两个值修改 mysql 引擎获取和 redis cli 获取
var RedisCliGetter = getRedisCli

// SetPlayerMaxwinningstream 储存最大连胜
func SetPlayerMaxwinningstream(key string, maxStream int) error {
	redisCli, err := RedisCliGetter(redisName, 0)
	if err != nil {
		return err
	}
	err = redisCli.Set(key, maxStream, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetPlayerMaxwinningstream 获取最大连胜
func GetPlayerMaxwinningstream(key string) (int, error) {
	redisCli, err := RedisCliGetter(redisName, 0)
	if err != nil {
		return 0, err
	}
	streamCmd := redisCli.Get(key)
	MaxStream, err := streamCmd.Int64()
	if err != nil {
		return 0, err
	}
	return int(MaxStream), nil
}

// // GetPlayerGameInfo 获取玩家游戏信息
// func GetPlayerGameInfo(playerID uint64, gameID uint32) (exist bool, info *db.TPlayerGame, err error) {
// 	exist, info, err = false, new(db.TPlayerGame), nil

// 	rKey := cache.FmtPlayerIDKey(playerID)
// 	// 从redis获取
// 	val, _ := getRedisField(playerRedisName, rKey, cache.FmtPlayerGameInfoKey(gameID))
// 	err = json.Unmarshal([]byte(val[0].(string)), info)
// 	if err == nil {
// 		exist = true
// 		return
// 	}

// 	engine, err := mysqlEngineGetter(playerMysqlName)

// 	where := fmt.Sprintf("playerID=%d and gameID='%d'", playerID, gameID)
// 	exist, err = engine.Table(playerGameTableName).Where(where).Get(info)

// 	if err != nil {
// 		err = fmt.Errorf("select t_player_game sql err：%v", err)
// 		return
// 	}
// 	// 更新redis
// 	data, _ := json.Marshal(info)
// 	rFields := map[string]string{
// 		cache.FmtPlayerGameInfoKey(gameID): string(data),
// 	}
// 	if err = setRedisFields(playerRedisName, rKey, rFields, redisTimeOut); err != nil {
// 		err = fmt.Errorf("save game_config  into redis fail： %v", err)
// 	}
// 	return
// }

// func setRedisFields(redisName string, key string, fields map[string]string, duration time.Duration) error {
// 	redisCli, err := redisCliGetter(redisName, 0)
// 	if err != nil {
// 		return err
// 	}
// 	kv := make(map[string]interface{}, len(fields))
// 	for k, field := range fields {
// 		kv[k] = field
// 	}
// 	status := redisCli.HMSet(key, kv)
// 	if status.Err() != nil {
// 		return fmt.Errorf("设置失败(%v)", status.Err())
// 	}
// 	redisCli.Expire(key, duration)
// 	return nil
// }
