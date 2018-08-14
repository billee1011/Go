package data

import (
	"fmt"
	"steve/server_pb/user"
	"testing"
)

func Test_data(t *testing.T) {
	ri := &RobotInfo{
		Gold:         10,
		GameWinRates: map[int]float64{1: 12.2},
	}
	ri2 := &RobotInfo{
		Gold:         12,
		GameWinRates: map[int]float64{2: 50.2},
	}
	ri3 := &RobotInfo{
		Gold:         13,
		GameWinRates: map[int]float64{3: 99.2},
	}
	robotsMap = map[uint64]*RobotInfo{1: ri, 2: ri2, 3: ri3}

	u1 := &user.RobotState{
		RobotId:    1,
		RobotState: user.PlayerState_PS_IDIE,
	}
	u2 := &user.RobotState{
		RobotId:    2,
		RobotState: user.PlayerState_PS_MATCHING,
	}
	u3 := &user.RobotState{
		RobotId:    3,
		RobotState: user.PlayerState_PS_GAMEING,
	}
	usermap := []*user.RobotState{u1, u2, u3}

	playerID, rinfo := ToInitRobotMapReturnLeisure(usermap)
	fmt.Printf("playerID(%d)\n", playerID)
	fmt.Printf("rinfo(%v)\n", *rinfo)
	fmt.Printf("robotsMap(%v)\n", robotsMap)
	// fmt.Printf("initRobotsMapFalse(%v)\n", initRobotsMap[false])
	// fmt.Printf("initRobotsMapTrue(%v)\n", initRobotsMap[true])

	// fmt.Printf("%v\n", UpdataRobotState(1, true))
	// fmt.Printf("2 initRobotsMapFalse(%v) \n", initRobotsMap[false])
	// fmt.Printf("2 initRobotsMapTrue(%v) \n", initRobotsMap[true])

	// fmt.Printf("%v\n", UpdataRobotState(2, false))
	// fmt.Printf("%v\n", UpdataRobotState(3, false))
	// fmt.Printf("3 initRobotsMapFalse(%v) \n", initRobotsMap[false])
	// fmt.Printf("3 initRobotsMapTrue(%v) \n", initRobotsMap[true])

	// UpdataRobotWinRate(2, 2, 90)
	// fmt.Printf("4 initRobotsMapFalse robotInfo(%v)\n", *initRobotsMap[false][2])
	// UpdataRobotWinRate(2, 1, 50)
	// fmt.Printf("4 initRobotsMapFalse robotInfo(%v)\n", *initRobotsMap[false][2])
}
