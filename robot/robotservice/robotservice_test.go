package robotservice

import (
	"context"
	"fmt"
	"log"
	"steve/robot/data"
	"steve/server_pb/robot"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func getinit() {
	conf := mysql.Config{
		User:                 "root",
		Passwd:               "123456",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "steve",
		Params:               map[string]string{"charset": "utf8"},
		AllowNativePasswords: true,
	}
	fmt.Println(conf.FormatDSN())
	mysqlEngine, _ := xorm.NewEngine("mysql", conf.FormatDSN())
	if err := mysqlEngine.Ping(); err != nil {
		panic(err)
	}
	data.MysqlEnginefunc = func(mysqlName string) (*xorm.Engine, error) {
		return mysqlEngine, nil
	}
	redisPlayerCli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	data.RedisClifunc = func() *redis.Client {
		return redisPlayerCli
	}
}

func Test_GetRobotPlayerIDByInfo(t *testing.T) {
	getinit()
	coinsR := &robot.CoinsRange{
		High: 1005,
		Low:  1002,
	}
	winR := &robot.WinRateRange{
		High: 50,
		Low:  50,
	}
	req := &robot.GetRobotPlayerIDReq{
		CoinsRange:   coinsR,
		WinRateRange: winR,
	}
	playerID, err := getRobotPlayerIDByInfo(req)
	if err != nil {
		fmt.Printf("有错误 : %v", err)
	}
	fmt.Println(playerID)
}

func Test_setRobotPlayerState(t *testing.T) {
	getinit()
	req := &robot.SetRobotPlayerStateReq{
		RobotPlayerId: 13,
		State:         robot.RobotPlayerState_RPS_MATCHING,
	}
	err := setRobotPlayerState(req)
	assert.Nil(t, err)
}

func Test_grpc_clien(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:36303", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := robot.NewRobotServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coinsR := &robot.CoinsRange{
		High: 1005,
		Low:  1002,
	}
	winR := &robot.WinRateRange{
		High: 50,
		Low:  50,
	}
	req := &robot.GetRobotPlayerIDReq{
		CoinsRange:   coinsR,
		WinRateRange: winR,
	}
	rsq, err := client.GetRobotPlayerIDByInfo(ctx, req)
	assert.Nil(t, err)
	fmt.Println(rsq)
}

func Test_grpc_clien2(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:36303", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := robot.NewRobotServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &robot.SetRobotPlayerStateReq{
		RobotPlayerId: 13,
		State:         robot.RobotPlayerState_RPS_MATCHING,
	}
	rsq, err := client.SetRobotPlayerState(ctx, req)
	assert.Nil(t, err)
	fmt.Println(rsq)
}
