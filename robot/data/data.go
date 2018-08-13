package data

import (
	"fmt"
	"steve/external/goldclient"
	"steve/server_pb/gold"
	"steve/server_pb/user"

	"github.com/Sirupsen/logrus"
)

//RobotInfo 机器人信息
type RobotInfo struct {
	Gold         int64           // 金币
	GameWinRates map[int]float64 // key gameid - value winRate
}

// 未初始化
var robotsMap map[uint64]*RobotInfo

// 已经初始  false 为使用，true,使用
var initRobotsMap map[bool]map[uint64]*RobotInfo

//GetNoInitRobot 获取未初始化map
func GetNoInitRobot() map[uint64]*RobotInfo {
	return robotsMap
}

//GetRobotMapByState 获取初始机器人map
func GetRobotMapByState(state user.PlayerState) map[uint64]*RobotInfo {
	if state == user.PlayerState_PS_IDIE {
		return GetLeisureRobot()
	}
	return getNotLeisureRobot()
}

//GetLeisureRobot 获取空闲机器人map
func GetLeisureRobot() map[uint64]*RobotInfo {
	return initRobotsMap[false]
}

// 获取非空闲机器人map
func getNotLeisureRobot() map[uint64]*RobotInfo {
	return initRobotsMap[true]
}

//InitRobotRedis 初始化机器人redis
func InitRobotRedis() error {
	robotsMap = make(map[uint64]*RobotInfo) //先清空
	if err := GetMysqlRobotFieldValuedAll(robotsMap); err != nil {
		logrus.WithError(err).Errorln("初始化从mysql获取机器人失败")
		return err
	}
	logrus.Debugf("初始化从mysql获取机器人完成 robotsMaplen(%d)\n", len(robotsMap))
	return nil
}

//GetRobotInfoByPlayerID get robotinfo
func GetRobotInfoByPlayerID(playerID uint64) (*RobotInfo, error) {
	if nr, isExist := initRobotsMap[false][playerID]; isExist {
		return nr, nil
	} else if nr2, isExist := initRobotsMap[true][playerID]; isExist {
		return nr2, nil
	} else if nr3, isExist := robotsMap[playerID]; isExist {
		return nr3, nil
	}
	return nil, fmt.Errorf("get robotinfo failed playerid(%v)", playerID)
}

//ToInitRobotMapReturnLeisure 初始化RobotMap
func ToInitRobotMapReturnLeisure(playerStates map[uint64]user.PlayerState) (rplayerID uint64, robotInfo *RobotInfo) {
	for playerID, state := range playerStates {
		if robot, isExist := robotsMap[playerID]; isExist {
			delete(robotsMap, playerID)
			if state == user.PlayerState_PS_IDIE {
				initRobotsMap[false][playerID] = robot
				if rplayerID == 0 {
					rplayerID = playerID
					robotInfo = robot
				}
			} else {
				initRobotsMap[true][playerID] = robot
			}
		} else {
			logrus.Debugf("init robotMap failed playerID(%d) not Exist", playerID)
		}
	}
	return rplayerID, robotInfo
}

// UpdataRobotState 更新机器人状态 true 空闲到使用，false 使用到空闲
func UpdataRobotState(playerID uint64, state bool) (err error) {
	curr, isExist2 := initRobotsMap[!state]
	if isExist2 && len(curr) != 0 {
		if robot, isExist := curr[playerID]; isExist {
			delete(initRobotsMap[!state], playerID)
			initRobotsMap[state][playerID] = robot
		} else {
			err = fmt.Errorf("playerID(%d) state(%v) initRobotsMap not Exist", playerID, state)
		}
	} else {
		err = fmt.Errorf("state(%v) initRobotsMap len eq 0 or not Exist", state)
	}
	return err
}

// UpdataRobotWinRate 更新机器人胜率
func UpdataRobotWinRate(playerID uint64, gameID int, winrate float64) error {
	robotInfo, err := GetRobotInfoByPlayerID(playerID)
	if err != nil {
		return err
	}
	robotInfo.GameWinRates[gameID] = winrate
	return nil
}

//GetMysqlRobotFieldValuedAll 获取机器人需要的值
func GetMysqlRobotFieldValuedAll(currRobotMap map[uint64]*RobotInfo) error {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return err
	}
	//gameid-winrate 游戏id对应的胜率
	robotsPGs, err := getMysqlRobotInfo(engine)
	if err != nil {
		return err
	}
	for _, robot := range robotsPGs {
		if robot.Gameid == 0 {
			continue
		}
		playerID := robot.Playerid
		if playerID == 0 {
			continue
		}
		rp, isExist := currRobotMap[uint64(playerID)]
		if isExist {
			rp.GameWinRates[robot.Gameid] = robot.Winningrate
		} else {
			rp.GameWinRates = map[int]float64{robot.Gameid: robot.Winningrate}
			// 从金币服获取
			gold, err := goldclient.GetGold(uint64(playerID), int16(gold.GoldType_GOLD_COIN))
			if err != nil {
				logrus.WithError(err).Errorf("获取金币失败 playerID(%v)", playerID)
				continue
			}
			rp.Gold = gold
		}
		currRobotMap[uint64(playerID)] = rp
	}
	return nil
}
