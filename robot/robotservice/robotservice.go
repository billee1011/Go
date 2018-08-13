package robotservice

import (
	"context"
	"fmt"
	"steve/external/goldclient"
	"steve/robot/data"
	"steve/server_pb/gold"
	"steve/server_pb/robot"

	"github.com/Sirupsen/logrus"
)

//Robotservice 机器人服务
type Robotservice struct{}

var defaultObject = new(Robotservice)
var _ robot.RobotServiceServer = DefaultRobot()

// DefaultRobot 默认对象
func DefaultRobot() *Robotservice {
	return defaultObject
}

//GetLeisureRobotInfoByInfo 获取空闲机器人信息
func (r *Robotservice) GetLeisureRobotInfoByInfo(ctx context.Context, request *robot.GetLeisureRobotInfoReq) (*robot.GetLeisureRobotInfoRsp, error) {
	logrus.Debugln("GetLeisureRobotInfoByInfo req", *request)
	rsp := &robot.GetLeisureRobotInfoRsp{
		RobotPlayerId: 0,
		Coin:          0,
		WinRate:       0,
		ErrCode:       int32(robot.ErrCode_EC_FAIL),
	}
	gameID := int(request.GetGame().GetGameId()) // 游戏
	coinsRange := request.GetCoinsRange()        // 金币范围
	winRateRange := request.GetWinRateRange()    // 胜率范围
	newState := int(request.GetNewState())       // 获取成功时设置的状态
	// 检验请求是否合法
	if !checkGetLeisureRobotArgs(coinsRange, winRateRange, newState) {
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		return rsp, fmt.Errorf("参数错误")
	}
	checkFunc := func(currRobotMap map[int64]*data.RobotPlayer, robotMap2 map[int64]*data.RobotPlayer) bool {
	next:
		for id, robotPlayer := range currRobotMap {
			if robotMap2 != nil {
				rp, isExist := robotMap2[id]
				if isExist { // 存在说明已经找过了
					continue
				} else { // 不存在，加入到新的
					robotMap2[id] = rp
				}
			}
			if robotPlayer == nil {
				logrus.Errorf("robotPlayer eq nil", id)
				continue
			}
			// 是空闲的情况
			if robotPlayer.State == int(robot.RobotPlayerState_RPS_IDIE) {
				gws := robotPlayer.GameIDWinRate
				currWinRate := float64(0)
				flag := false
				//找到相同GAMEid
				for _, gw := range gws {
					currGameID := gw.GameID
					if currGameID == gameID {
						flag = true //gameID存在
						currWinRate = gw.WinRate
						if currWinRate > float64(winRateRange.High) || currWinRate < float64(winRateRange.Low) {
							continue next // 胜率不符合
						}
						break // 胜率符合
					}
				}
				// 找不到对应ID
				if !flag {
					currWinRate = float64(50) // 默认50，情况下胜率是否符合
					if int32(currWinRate) > winRateRange.High || int32(currWinRate) < winRateRange.Low {
						continue
					}
					logrus.Debugf("gameID(%d) 不存在 ，默认胜率50", gameID)
				}
				// 从金币服获取
				gold, err := goldclient.GetGold(uint64(id), int16(gold.GoldType_GOLD_COIN))
				if err != nil {
					logrus.WithError(err).Errorf("获取金币失败 playerID(%v)", id)
					continue
				}
				// 找到适合的
				if gold <= coinsRange.High && gold >= coinsRange.Low {
					rsp.RobotPlayerId = uint64(id)
					rsp.Coin = gold
					rsp.WinRate = currWinRate
					rsp.ErrCode = int32(robot.ErrCode_EC_SUCCESS)
					return true
				}
			}
		}
		return false
	}
	robotsMap := data.GetRobotsMap()
	// 更新状态
	defer func() {
		if rsp.ErrCode == int32(robot.ErrCode_EC_SUCCESS) {
			data.UpdataRobotState(int64(rsp.RobotPlayerId), newState)
			logrus.WithFields(logrus.Fields{
				"RobotPlayerId": rsp.GetRobotPlayerId(),
				"coin":          rsp.GetCoin(),
				"winRate":       rsp.GetWinRate(),
				"newState":      newState,
				"Startlimit":    data.Startlimit,
			}).Infoln("获取空闲机器人成功")
		} else {
			logrus.WithFields(logrus.Fields{
				"Startlimit": data.Startlimit,
			}).Infoln("获取空闲机器人失败")
		}
	}()
	if checkFunc(robotsMap, nil) {
		return rsp, nil
	}
	// 已存的找不到的情况,再从数据库获取
	logrus.Debugln("已存的找不到的情况,再从数据库获取")
	for {
		newRobotsMap := make(map[int64]*data.RobotPlayer) //先清空
		if err := data.GetMysqlRobotFieldValuedAll(newRobotsMap); err != nil {
			rsp.ErrCode = int32(robot.ErrCode_EC_FAIL)
			return rsp, err
		}
		if checkFunc(newRobotsMap, robotsMap) {
			return rsp, nil
		}
		if data.Startlimit == 0 {
			break
		}
	}
	rsp.ErrCode = int32(robot.ErrCode_EC_FAIL)
	return rsp, fmt.Errorf("找不到适合的机器人")
}

//SetRobotPlayerState 设置机器人玩家状态
func (r *Robotservice) SetRobotPlayerState(ctx context.Context, request *robot.SetRobotPlayerStateReq) (*robot.SetRobotPlayerStateRsp, error) {
	logrus.Debugln("SetRobotPlayerState req", *request)
	rsp := &robot.SetRobotPlayerStateRsp{
		Result:  true,
		ErrCode: int32(robot.ErrCode_EC_SUCCESS),
	}
	playerID := int64(request.GetRobotPlayerId())
	newState := int(request.GetNewState())
	oldState := int(request.GetOldState())

	// 检验请求是否合法
	if !checkSateArgs(playerID, newState, oldState) {
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		return rsp, nil
	}
	robotsMap := data.GetRobotsMap()
	// 判断玩家ID是否存在
	rp, isExist := robotsMap[playerID]
	if !isExist { //要先判断是否时机器人
		flag, err := isRobot(robotsMap, playerID)
		if err != nil {
			rsp.ErrCode = int32(robot.ErrCode_EC_FAIL)
			rsp.Result = false
			return rsp, err
		}
		if !flag {
			logrus.Warningln("不是机器人 playerID(%d)", playerID)
			rsp.ErrCode = int32(robot.ErrCode_EC_NOTROBOT)
			rsp.Result = false
			return rsp, fmt.Errorf("不是机器人 playerID(%d)", playerID)
		}
		rp = robotsMap[playerID]
	}
	// 判断旧的状态是否符合
	if rp.State != int(oldState) {
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		rsp.Result = false
		return rsp, fmt.Errorf("currState(%d) 旧状态不符合 oldState(%d)", rp.State, oldState)
	}
	// 更新状态
	data.UpdataRobotState(playerID, int(newState))
	logrus.WithFields(logrus.Fields{
		"RobotPlayerId": playerID,
		"oldState":      oldState,
		"newState":      robotsMap[playerID].State,
	}).Infoln("更新空闲机器人状态成功")
	return rsp, nil
}

// UpdataRobotGameWinRate 更新胜率
func (r *Robotservice) UpdataRobotGameWinRate(ctx context.Context, request *robot.UpdataRobotGameWinRateReq) (*robot.UpdataRobotGameWinRateRsp, error) {
	logrus.Debugln("UpdataRobotGameWinRate req", *request)
	rsp := &robot.UpdataRobotGameWinRateRsp{
		Result:  true,
		ErrCode: int32(robot.ErrCode_EC_SUCCESS),
	}
	playerID := int64(request.GetRobotPlayerId()) // TODO判断是否时机器人
	gameID := int(request.GetGameId())
	newWinRate := int(request.GetNewWinRate())
	oldWinRate := int(request.GetOldWinRate())
	if playerID < 0 {
		logrus.Warningln("Robot Player ID cannot be less than 0:%d", playerID)
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		return rsp, nil
	}
	if newWinRate == oldWinRate {
		logrus.Warningln("胜率未改变 newWinRate(%d) - oldWinRate(%d)", newWinRate, oldWinRate)
		return rsp, nil
	}
	robotsMap := data.GetRobotsMap()
	// 判断玩家ID是否存在
	rp, isExist := robotsMap[playerID]
	if !isExist {
		flag, err := isRobot(robotsMap, playerID)
		if err != nil {
			rsp.ErrCode = int32(robot.ErrCode_EC_FAIL)
			rsp.Result = false
			return rsp, err
		}
		if !flag {
			logrus.Warningln("不是机器人 playerID(%d)", playerID)
			rsp.ErrCode = int32(robot.ErrCode_EC_NOTROBOT)
			rsp.Result = false
			return rsp, fmt.Errorf("不是机器人 playerID(%d)", playerID)
		}
		rp = robotsMap[playerID]
	}
	// 判断旧胜率
	flag := false
	for _, gw := range rp.GameIDWinRate {
		if gameID == gw.GameID { //  找到
			if int(gw.WinRate) != oldWinRate {
				rsp.ErrCode = int32(robot.ErrCode_EC_Args)
				rsp.Result = false
				return rsp, fmt.Errorf("currGameID(%v) 旧胜率不符合 oldWinRate(%d)", gw.WinRate, oldWinRate)
			}
			flag = true
		}
	}
	// 游戏ID不存在
	if !flag && oldWinRate != 50 {
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		rsp.Result = false
		return rsp, fmt.Errorf("游戏id不存在！旧胜率不符合 oldWinRate(%d)", oldWinRate)
	}
	// 更新胜率
	data.UpdataRobotWinRate(playerID, gameID, float64(newWinRate), flag)
	logrus.WithFields(logrus.Fields{
		"RobotPlayerId": playerID,
		"gameID":        gameID,
		"oldWinRate":    oldWinRate,
		"newWinRate":    newWinRate,
	}).Infoln("更新空闲机器人胜率成功")
	return rsp, nil
}

//IsRobotPlayer 判断是否时机器人
func (r *Robotservice) IsRobotPlayer(ctx context.Context, request *robot.IsRobotPlayerReq) (*robot.IsRobotPlayerRsp, error) {
	logrus.Debugln("IsRobotPlayer req", *request)
	rsp := &robot.IsRobotPlayerRsp{
		Result: false,
	}
	playerID := int64(request.GetRobotPlayerId())
	if playerID <= 0 {
		return rsp, fmt.Errorf("参数错误")
	}
	robotsMap := data.GetRobotsMap()
	if _, isExist := robotsMap[playerID]; isExist {
		rsp.Result = true
		return rsp, nil
	}
	flag, err := isRobot(robotsMap, playerID)
	rsp.Result = flag
	return rsp, err
}

func isRobot(robotsMap map[int64]*data.RobotPlayer, playerID int64) (bool, error) {
	// 在从数据库获取
	falg, err := data.IsMysqlRobot(playerID)
	if err != nil {
		return false, err
	}
	if !falg { // 不是机器人
		return false, nil
	}
	// 是机器人
	pgs, err := data.GetRobotInfoByPlayerID(playerID)
	if err != nil {
		return true, err
	}
	nrp := &data.RobotPlayer{
		State:         0,
		GameIDWinRate: []*data.PlayerGameGW{},
	}
	for _, pg := range pgs {
		pggw := &data.PlayerGameGW{
			GameID:  pg.Gameid,
			WinRate: pg.Winningrate,
		}
		nrp.GameIDWinRate = append(nrp.GameIDWinRate, pggw)
	}
	logrus.WithFields(logrus.Fields{
		"RobotPlayerId": playerID,
		"nPlayerGameGW": fmtGameIDWinRate(nrp.GameIDWinRate),
	}).Infoln("判断是否时机器人 - 添加新的机器人")
	robotsMap[playerID] = nrp
	return true, nil
}

// 打印
func fmtGameIDWinRate(gws []*data.PlayerGameGW) []data.PlayerGameGW {
	new := make([]data.PlayerGameGW, 0, len(gws))
	for _, pg := range gws {
		new = append(new, *pg)
	}
	return new
}
