package logic

import (
	"fmt"
	"steve/external/robotclient"
	"steve/server_pb/robot"
)

func startTestRobotServer() {
	pid := TestGetLRoboyPlayer()
	TestSetLRoboyPlayerState(pid)
	// TestSetLRoboyPlayerState(100002)
	// TestUpdataWinRate()
	// TestIsRobot()
}

func TestIsRobot() {
	flag, err := robotclient.IsRobotPlayer(77777)
	fmt.Println("++++++++++++++++++++++++++++")
	fmt.Printf("err(%v)\n", err)
	fmt.Printf("rsp(%v)\n", flag)
}

func TestUpdataWinRate() {
	flag, err := robotclient.UpdataRobotPlayerWinRate(100002, 3, 55, 80)
	fmt.Println("++++++++++++++++++++++++++++")
	fmt.Printf("err(%v)\n", err)
	fmt.Printf("rsp(%v)\n", flag)
}

func TestGetLRoboyPlayer() uint64 {
	req := robotclient.LeisureRobotReqInfo{
		CoinHigh:    10000,
		CoinLow:     0,
		WinRateHigh: 100,
		WinRateLow:  0,
		GameID:      2,
		LevelID:     1,
	}
	playerID, coin, winR, err := robotclient.GetLeisureRobotInfoByInfo(req)
	fmt.Println("++++++++++++++++++++++++++++")
	fmt.Printf("err(%v)\n", err)
	fmt.Printf("robotPlayerID(%v)\n", playerID)
	fmt.Printf("coin(%v)\n", coin)
	fmt.Printf("winR(%v)\n", winR)
	return playerID
}

func TestSetLRoboyPlayerState(playerID uint64) {
	fmt.Println("++++++++++++++++++++++++++++")
	NewState := uint32(robot.RobotPlayerState_RPS_MATCHING)
	OldState := uint32(robot.RobotPlayerState_RPS_IDIE)
	ServerType := uint32(robot.ServerType_ST_MATCH)
	ServerAddr := "127.0.0.1:3306"
	flag, err := robotclient.SetRobotPlayerState(uint64(playerID), OldState, NewState, ServerType, ServerAddr)
	fmt.Printf("err(%v)\n", err)
	fmt.Printf("是否更改成功(%v)\n", flag)
}
