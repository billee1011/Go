package data

import (
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
	"steve/gutils"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var redisPlayerCli *redis.Client

func init() {
	conf := mysql.Config{
		User:                 "root",
		Passwd:               "12345678",
		Net:                  "tcp",
		Addr:                 "192.168.7.108:3306",
		DBName:               "steve",
		Params:               map[string]string{"charset": "utf8"},
		AllowNativePasswords: true,
	}
	mysqlPlayerEngine, _ := xorm.NewEngine("mysql", conf.FormatDSN())

	redisPlayerCli = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	mysqlEngineGetter = func(mysqlName string) (*xorm.Engine, error) {
		return mysqlPlayerEngine, nil
	}
	redisCliGetter = func(redis string, db int) (*redis.Client, error) {
		return redisPlayerCli, nil
	}
}

func NewPlayerData(accID uint64, playerID uint64) {
	GetPlayerIDByAccountID(accID)
	InitPlayerData(db.TPlayer{
		Accountid:    int64(accID),
		Playerid:     int64(playerID),
		Type:         1,
		Channelid:    0,                                 // TODO ，渠道 ID
		Nickname:     fmt.Sprintf("player%d", playerID), // TODO,昵称
		Gender:       1,
		Avatar:       "", // TODO , 头像
		Provinceid:   0,  // TODO， 省ID
		Cityid:       0,  // TODO 市ID
		Name:         "", // TODO: 真实姓名
		Phone:        "", // TODO: 电话
		Idcard:       "", // TODO 身份证
		Iswhitelist:  0,
		Zipcode:      0,
		Shippingaddr: "",
		Status:       1,
		Remark:       "",
		Createtime:   time.Now(),
		Createby:     "",
		Updatetime:   time.Now(),
		Updateby:     "",
	})

	InitPlayerCoin(db.TPlayerCurrency{
		Playerid:       int64(playerID),
		Coins:          10000,
		Ingots:         0,
		Keycards:       0,
		Obtainingots:   0,
		Obtainkeycards: 0,
		Costingots:     0,
		Costkeycards:   0,
		Remark:         "",
		Createtime:     time.Now(),
		Createby:       "",
		Updatetime:     time.Now(),
		Updateby:       "",
	})
	InitPlayerState(int64(playerID))
	return
}

// TestInitPlayerData 初始化玩家
func TestInitPlayerData(t *testing.T) {
	viper.SetDefault("node", 200)
	playerID := AllocPlayerID()
	assert.NotZero(t, playerID)

	alloc, err := gutils.NewNode(300)
	assert.Nil(t, err)
	accID := uint64(alloc.Generate().Int64()) // 用这个

	exist, _, err := GetPlayerIDByAccountID(accID)
	assert.False(t, exist)
	assert.Nil(t, err)

	err = InitPlayerData(db.TPlayer{
		Accountid:    int64(accID),
		Playerid:     int64(playerID),
		Type:         1,
		Channelid:    0,                                 // TODO ，渠道 ID
		Nickname:     fmt.Sprintf("player%d", playerID), // TODO,昵称
		Gender:       1,
		Avatar:       "", // TODO , 头像
		Provinceid:   0,  // TODO， 省ID
		Cityid:       0,  // TODO 市ID
		Name:         "", // TODO: 真实姓名
		Phone:        "", // TODO: 电话
		Idcard:       "", // TODO 身份证
		Iswhitelist:  0,
		Zipcode:      0,
		Shippingaddr: "",
		Status:       1,
		Remark:       "",
		Createtime:   time.Now(),
		Createby:     "",
		Updatetime:   time.Now(),
		Updateby:     "",
	})
	assert.Nil(t, err)

	err = InitPlayerCoin(db.TPlayerCurrency{
		Playerid:       int64(playerID),
		Coins:          10000,
		Ingots:         0,
		Keycards:       0,
		Obtainingots:   0,
		Obtainkeycards: 0,
		Costingots:     0,
		Costkeycards:   0,
		Remark:         "",
		Createtime:     time.Now(),
		Createby:       "",
		Updatetime:     time.Now(),
		Updateby:       "",
	})
	assert.Nil(t, err)

	err = InitPlayerState(int64(playerID))
	assert.Nil(t, err)

	// 能正确拿到账号关联的玩家
	exist, _playerID, err := GetPlayerIDByAccountID(accID)
	assert.Equal(t, _playerID, playerID)
	assert.Nil(t, err)
	assert.True(t, exist)

	// redis 中有数据
	redisKey := cache.FmtAccountPlayerKey(accID)
	redisCmd := redisPlayerCli.Get(redisKey)
	assert.Nil(t, redisCmd.Err())

	_playerID, err = redisCmd.Uint64()
	assert.Nil(t, err)
	assert.Equal(t, _playerID, playerID)
}

// TestGetPlayerInfo 获取玩家信息
func TestGetPlayerInfo(t *testing.T) {
	viper.SetDefault("node", 200)
	playerID := AllocPlayerID()
	assert.NotZero(t, playerID)

	alloc, err := gutils.NewNode(300)
	assert.Nil(t, err)
	accID := uint64(alloc.Generate().Int64())

	NewPlayerData(accID, playerID)

	player, err := GetPlayerInfoByPlayerID(playerID)

	assert.Nil(t, err)
	assert.NotNil(t, player.State)

	// redis 中有数据
	redisKey := cache.FmtPlayerIDKey(playerID)
	result, err := redisPlayerCli.HMGet(redisKey, cache.PlayerStateField).Result()
	assert.Nil(t, err)
	assert.NotEmpty(t, result[0])
	fmt.Println(result)
}

// TestGetGameInfoList 获取游戏信息
func TestGetGameInfoList(t *testing.T) {
	gameInfo, err := GetGameInfoList()
	assert.Nil(t, err)
	assert.Nil(t, gameInfo)
}

// TestSetPlayerState 修改玩家状态
func TestSetPlayerState(t *testing.T) {
	viper.SetDefault("node", 200)
	playerID := AllocPlayerID()
	assert.NotZero(t, playerID)

	alloc, err := gutils.NewNode(300)
	assert.Nil(t, err)
	accID := uint64(alloc.Generate().Int64())

	NewPlayerData(accID, playerID)

	result, err := UpdatePlayerState(playerID, 0, 1, 1, "127.0.0.1")
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}

func TestUpdatePlayerInfo(t *testing.T) {
	viper.SetDefault("node", 200)
	playerID := AllocPlayerID()
	assert.NotZero(t, playerID)

	alloc, err := gutils.NewNode(300)
	assert.Nil(t, err)
	accID := uint64(alloc.Generate().Int64())

	NewPlayerData(accID, playerID)

	exists, result, err := UpdatePlayerInfo(playerID, "正是", "你好")
	assert.Nil(t, err)
	assert.Equal(t, true, result)
	assert.Equal(t, true, exists)
}
