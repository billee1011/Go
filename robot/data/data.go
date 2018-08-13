package data

import (
	"github.com/Sirupsen/logrus"
)

var robotsMap map[int64]*RobotPlayer

// RobotPlayer 机器人玩家
type RobotPlayer struct {
	State         int `protobuf:"varint,5,opt,name=state" json:"state,omitempty"`
	GameIDWinRate []*PlayerGameGW
}

//PlayerGameGW 游戏ID对应的胜率
type PlayerGameGW struct {
	GameID  int
	WinRate float64
}

//InitRobotRedis 初始化机器人redis
func InitRobotRedis() error {
	Startlimit = 0
	robotsMap = make(map[int64]*RobotPlayer) //先清空
	if err := GetMysqlRobotFieldValuedAll(robotsMap); err != nil {
		logrus.WithError(err).Errorln("初始化从mysql获取机器人失败")
		return err
	}
	logrus.Debugf("初始化从mysql获取机器人完成 robotsMaplen(%d)\n", len(robotsMap))
	return nil
}

//GetRobotsMap 获取机器人信息
func GetRobotsMap() map[int64]*RobotPlayer {
	return robotsMap
}

// UpdataRobotState 更新机器人状态
func UpdataRobotState(playerID int64, state int) bool {
	if rp, isExist := robotsMap[playerID]; isExist {
		rp.State = state
		robotsMap[playerID] = rp
		return true
	}
	return false
}

// UpdataRobotWinRate 更新机器人胜率
func UpdataRobotWinRate(playerID int64, gameID int, winrate float64, flag bool) {
	rp := robotsMap[playerID]
	if flag {
		for _, gw := range rp.GameIDWinRate {
			if gw.GameID == gameID {
				gw.WinRate = winrate
			}
		}
	} else {
		// 找不到游戏ID的情况
		pg := &PlayerGameGW{
			GameID:  gameID,
			WinRate: winrate,
		}
		if len(rp.GameIDWinRate) == 0 {
			rp.GameIDWinRate = []*PlayerGameGW{pg}
		} else {
			rp.GameIDWinRate = append(rp.GameIDWinRate, pg)
		}
	}
}

//GetMysqlRobotFieldValuedAll 获取机器人需要的值
func GetMysqlRobotFieldValuedAll(currRobotMap map[int64]*RobotPlayer) error {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return err
	}
	//gameid-winrate 游戏id对应的胜率
	robotsPGs, err := getRobotInfo(engine)
	if err != nil {
		return err
	}
	for _, robot := range robotsPGs {
		playerID := robot.Playerid
		if playerID == 0 {
			continue
		}
		nrp := &RobotPlayer{
			State:         0,
			GameIDWinRate: []*PlayerGameGW{},
		}
		if robot.Gameid > 0 {
			pg := &PlayerGameGW{
				GameID:  robot.Gameid,
				WinRate: robot.Winningrate,
			}
			rp, isExist := currRobotMap[robot.Playerid]
			if isExist && len(rp.GameIDWinRate) != 0 {
				nrp.GameIDWinRate = append(rp.GameIDWinRate, pg)
			} else {
				nrp.GameIDWinRate = []*PlayerGameGW{pg}
			}
		}
		currRobotMap[robot.Playerid] = nrp
	}
	return nil
}
