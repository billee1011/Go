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
		User:   "root",
		Passwd: "123456",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "player",
		Params: map[string]string{"charset": "utf8"},
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
