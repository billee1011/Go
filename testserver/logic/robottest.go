package logic

import (
	"fmt"
	"steve/external/robotclient"
)

func startTestRobotServer() {
	pid := TestGetLRoboyPlayer()
	TestSetLRoboyPlayerState(pid)
	// TestUpdataWinRate()
	// TestIsRobot()
}

func TestIsRobot() {
	flag, err := robotclient.IsRobotPlayer(275)
	fmt.Println("++++++++++++++++++++++++++++")
	fmt.Printf("err(%v)\n", err)
	fmt.Printf("rsp(%v)\n", flag)
}

func TestUpdataWinRate() {
	flag, err := robotclient.UpdataRobotPlayerWinRate(283, 2, 80)
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
		GameID:      3,
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
	NewState := false
	flag, err := robotclient.SetRobotPlayerState(uint64(playerID), NewState)
	fmt.Printf("err(%v)\n", err)
	fmt.Printf("是否更改成功(%v)\n", flag)
}
