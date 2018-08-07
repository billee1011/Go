package logic

import (
	"fmt"
	"steve/external/robotclient"
	"steve/server_pb/robot"
)

func startTestRobotServer() {
	TestGetLRoboyPlayer()
	TestSetLRoboyPlayerState()
}

func TestGetLRoboyPlayer() {
	req := robotclient.LeisureRobotReqInfo{
		CoinHigh:    8000,
		CoinLow:     5000,
		WinRateHigh: 50,
		WinRateLow:  50,
		GameID:      1,
		LevelID:     1,
	}
	playerID, coin, winR, err := robotclient.GetLeisureRobotInfoByInfo(req)
	fmt.Println("++++++++++++++++++++++++++++")
	fmt.Printf("err(%v)\n", err)
	fmt.Printf("robotPlayerID(%v)\n", playerID)
	fmt.Printf("coin(%v)\n", coin)
	fmt.Printf("coin(%v)\n", winR)
}

func TestSetLRoboyPlayerState() {
	NewState := uint32(robot.RobotPlayerState_RPS_IDIE)
	OldState := uint32(robot.RobotPlayerState_RPS_MATCHING)
	ServerType := uint32(robot.ServerType_ST_MATCH)
	ServerAddr := "127.0.0.1:3306"
	flag, err := robotclient.SetRobotPlayerState(uint64(2005), OldState, NewState, ServerType, ServerAddr)
	fmt.Printf("err(%v)\n", err)
	fmt.Printf("是否更改成功(%v)\n", flag)
	fmt.Println("++++++++++++++++++++++++++++")
}
