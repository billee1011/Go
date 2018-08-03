package logic

import (
	"steve/back/data"
	"steve/entity/gamelog"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
)

var redisPlayerCli *redis.Client

func init() {
	conf := mysql.Config{
		User:                 "pipi",
		Passwd:               "123456",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
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

	data.MysqlEngineGetter = func(mysqlName string) (*xorm.Engine, error) {
		return mysqlPlayerEngine, nil
	}
	data.RedisCliGetter = func(redis string, db int) (*redis.Client, error) {
		return redisPlayerCli, nil
	}
}

func TestInsertSummary(t *testing.T) {
	info := gamelog.TGameSummary{
		Sumaryid:  123,
		Deskid:    1,
		Gameid:    4,
		Levelid:   1,
		Playerids: []uint64{1, 2, 3, 4},
		Scoreinfo: []int64{90, -30, -30, -30},
		Winnerids: []uint64{1},
		Roundcurrency: []gamelog.RoundCurrency{
			gamelog.RoundCurrency{
				Settletype: 1,
				Settledetails: []gamelog.SettleDetail{
					gamelog.SettleDetail{
						Playerid:  1,
						ChangeVal: 90,
					},
					gamelog.SettleDetail{
						Playerid:  2,
						ChangeVal: -30,
					},
					gamelog.SettleDetail{
						Playerid:  3,
						ChangeVal: -30,
					},
					gamelog.SettleDetail{
						Playerid:  4,
						ChangeVal: -30,
					},
				},
			},
		},
		Createtime: time.Now(),
		Createby:   "1",
		Updatetime: time.Now(),
		Updateby:   "1",
	}
	assert.Nil(t, insertSummaryInfo(info))
}

func TestInsertDetail(t *testing.T) {
	detail := gamelog.TGameDetail{
		Sumaryid:   123,
		Playerid:   1,
		Deskid:     1,
		Gameid:     4,
		Amount:     1,
		Iswinner:   1,
		Createtime: time.Now(),
		Createby:   "1",
		Updatetime: time.Now(),
		Updateby:   "1",
	}
	assert.Nil(t, insertDetailInfo(detail))
}

func TestUpdatePlayerGame(t *testing.T) {
	detail := gamelog.TGameDetail{
		Sumaryid:   123,
		Playerid:   6,
		Deskid:     1,
		Gameid:     4,
		Amount:     -1,
		Iswinner:   1,
		Createtime: time.Now(),
		Createby:   "1",
		Updatetime: time.Now(),
		Updateby:   "1",
	}
	assert.Nil(t, updatePlayerInfo(detail))
}
