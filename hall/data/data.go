package data

import (
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/gutils"
	"steve/structs"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"github.com/spf13/viper"
)

var idAllocObject *gutils.Node

const (
	playerRedisName = "player"
	playerMysqlName = "player"
	playerTableName = "t_player"
)

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
	engine, err := exposer.MysqlEngineMgr.GetEngine(playerMysqlName)
	if err != nil {
		return nil, fmt.Errorf("获取 mysql 引擎失败：%v", err)
	}
	return engine, nil
}

// 单元测试通过这两个值修改 mysql 引擎获取和 redis cli 获取
var mysqlEngineGetter = getMysqlEngine
var redisCliGetter = getRedisCli

func getRedisUint64Val(redis string, key string) (uint64, error) {
	redisCli, err := redisCliGetter(redis, 0)
	if err != nil {
		return 0, err
	}
	redisCmd := redisCli.Get(key)
	if redisCmd.Err() == nil {
		playerID, err := redisCmd.Uint64()
		if err != nil {
			return 0, fmt.Errorf("获取 redis 数据失败")
		}
		return playerID, nil
	}
	return 0, fmt.Errorf("redis 命令执行失败: %v", redisCmd.Err())
}

func setRedisVal(redis string, key string, val interface{}, duration time.Duration) error {
	redisCli, err := redisCliGetter(redis, 0)
	if err != nil {
		return err
	}
	redisCmd := redisCli.Set(key, val, duration)
	if redisCmd.Err() != nil {
		return fmt.Errorf("redis 命令执行失败：%v", redisCmd.Err())
	}
	return nil
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

func init() {
	node := viper.GetInt("node")
	var err error
	idAllocObject, err = gutils.NewNode(int64(node))
	if err != nil {
		logrus.Panicf("创建 id 生成器失败: %v", err)
	}
}
