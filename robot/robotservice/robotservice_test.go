package robotservice

import (
	"context"
	"fmt"
	"log"
	"steve/server_pb/robot"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// func getinit() {
// 	conf := mysql.Config{
// 		User:                 "backuser",
// 		Passwd:               "Sdf123esdf",
// 		Net:                  "tcp",
// 		Addr:                 "192.168.7.108:3306",
// 		DBName:               "steve",
// 		Params:               map[string]string{"charset": "utf8"},
// 		AllowNativePasswords: true,
// 	}
// 	fmt.Println(conf.FormatDSN())
// 	mysqlEngine, _ := xorm.NewEngine("mysql", conf.FormatDSN())
// 	if err := mysqlEngine.Ping(); err != nil {
// 		panic(err)
// 	}
// 	data.MysqlEnginefunc = func(mysqlName string) (*xorm.Engine, error) {
// 		return mysqlEngine, nil
// 	}
// 	redisPlayerCli := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "",
// 		DB:       0,
// 	})
// 	data.RedisClifunc = func() *redis.Client {
// 		return redisPlayerCli
// 	}
// }

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
		High: 8000,
		Low:  5000,
	}
	winR := &robot.WinRateRange{
		High: 50,
		Low:  50,
	}
	game := &robot.GameConfig{
		GameId: 1,
	}
	req := &robot.GetLeisureRobotInfoReq{
		Game:         game,
		CoinsRange:   coinsR,
		WinRateRange: winR,
		NewState:     robot.RobotPlayerState_RPS_MATCHING,
	}
	rsq, err := client.GetLeisureRobotInfoByInfo(ctx, req)
	assert.Nil(t, err)
	fmt.Println(rsq)
	assert.Equal(t, rsq.GetErrCode(), int32(robot.ErrCode_EC_SUCCESS))
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
		RobotPlayerId: 2000,
		NewState:      robot.RobotPlayerState_RPS_IDIE,
		OldState:      robot.RobotPlayerState_RPS_MATCHING,
		ServerType:    robot.ServerType_ST_MATCH,
		ServerAddr:    "127.0.0.1:3306",
	}
	rsq, err := client.SetRobotPlayerState(ctx, req)
	assert.Nil(t, err)
	assert.True(t, rsq.GetResult())
	assert.Equal(t, rsq.GetErrCode(), int32(robot.ErrCode_EC_SUCCESS))
}
