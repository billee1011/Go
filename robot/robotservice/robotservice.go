package robotservice

import (
	"context"
	"fmt"
	"steve/external/hallclient"
	"steve/robot/data"
	"steve/server_pb/robot"
	"steve/server_pb/user"

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
		ErrCode:       robot.ErrCode_EC_FAIL,
	}
	gameID := int(request.GetGame().GetGameId()) // 游戏
	coinsRange := request.GetCoinsRange()        // 金币范围
	winRateRange := request.GetWinRateRange()    // 胜率范围
	// 检验请求是否合法
	if !checkGetLeisureRobotArgs(coinsRange, winRateRange) {
		rsp.ErrCode = robot.ErrCode_EC_Args
		return rsp, fmt.Errorf("参数错误")
	}
	checkFunc := func(playerID uint64, robotPlayer *data.RobotInfo) bool {
		if robotPlayer == nil {
			logrus.Errorf("robotPlayer eq nil", playerID)
			return false
		}
		winRate, isExist := robotPlayer.GameWinRates[gameID]
		if isExist {
			if winRate > float64(winRateRange.High) || winRate < float64(winRateRange.Low) {
				return false
			}
		} else {
			robotPlayer.GameWinRates[gameID] = 50
			if 50 > winRateRange.High || 50 < winRateRange.Low {
				return false
			}
			logrus.Debugf("playerID(%d) gameID(%d) 不存在 ，默认胜率50", playerID, gameID)
		}
		gold := robotPlayer.Gold
		// 找到适合的
		if gold <= coinsRange.High && gold >= coinsRange.Low {
			return true
		}
		return false
	}
	initRobotsMapFalse := data.GetLeisureRobot()
	if len(initRobotsMapFalse) > 0 {
		for playerID, robotPlayer := range initRobotsMapFalse {
			if checkFunc(playerID, robotPlayer) {
				rsp.RobotPlayerId = playerID
				rsp.Coin = robotPlayer.Gold
				rsp.WinRate = robotPlayer.GameWinRates[gameID]
				rsp.ErrCode = robot.ErrCode_EC_SUCCESS
				return rsp, data.UpdataRobotState(playerID, true)
			}
		}
	}
	logrus.Debugln("从未初始化中，查找适合的机器人")
	notInitRobotMap := data.GetNoInitRobot() //未初始化
	if l := len(notInitRobotMap); l > 0 {
		i := 0 //防止死循环
		for {
			if len(notInitRobotMap) == 0 || i >= l {
				logrus.Debugf("从未初始化中，找不到适合的机器人 notInitRobotMaplen(%d)", len(notInitRobotMap))
				break
			}
			i++
			suitRobot := make([]uint64, 0, 10)
			for playerID, robotPlayer := range notInitRobotMap {
				if checkFunc(playerID, robotPlayer) {
					suitRobot = append(suitRobot, uint64(playerID))
				}
				if len(suitRobot) > 10 {
					break
				}
			}
			if len(suitRobot) == 0 { // 没有适合的机器人
				continue
			}
			hallrsp, err := hallclient.InitRobotPlayerState(suitRobot)
			if err != nil || hallrsp.GetErrCode() != int32(user.ErrCode_EC_SUCCESS) {
				logrus.WithError(err).Errorf("hall-初始化机器人失败 %d", hallrsp.GetErrCode())
				continue
			}
			if len(hallrsp.GetRobotState()) == 0 {
				logrus.Warningf("hall-初始化机器人失败 hall get robotSate len %d", len(hallrsp.GetRobotState()))
				continue
			}
			playerID, robotInfo := data.ToInitRobotMapReturnLeisure(hallrsp.GetRobotState()) // 初始化
			if playerID > 0 && robotInfo != nil {
				rsp.RobotPlayerId = playerID
				rsp.Coin = robotInfo.Gold
				rsp.WinRate = robotInfo.GameWinRates[gameID]
				rsp.ErrCode = robot.ErrCode_EC_SUCCESS
				return rsp, data.UpdataRobotState(playerID, true)
			}
			notInitRobotMap = data.GetNoInitRobot()
		}
	} else {
		logrus.Debugln("未初始化 notInitRobotMap 为0")
	}
	defer func() {
		if rsp.ErrCode == robot.ErrCode_EC_SUCCESS {
			logrus.WithFields(logrus.Fields{
				"RobotPlayerId": rsp.GetRobotPlayerId(),
				"coin":          rsp.GetCoin(),
				"winRate":       rsp.GetWinRate(),
			}).Infoln("获取空闲机器人成功")
		} else {
			logrus.Debugln("获取空闲机器人失败")
		}
	}()
	return rsp, fmt.Errorf("找不到适合的机器人")
}

//SetRobotPlayerState 设置机器人玩家状态  先判断是否是机器人，是机器人，在判断是否是空闲状态
func (r *Robotservice) SetRobotPlayerState(ctx context.Context, request *robot.SetRobotPlayerStateReq) (*robot.SetRobotPlayerStateRsp, error) {
	logrus.Debugln("SetRobotPlayerState req", *request)
	rsp := &robot.SetRobotPlayerStateRsp{
		Result:  true,
		ErrCode: robot.ErrCode_EC_SUCCESS,
	}
	playerID := request.GetRobotPlayerId()
	if playerID < 0 {
		logrus.Warningln("Robot Player ID cannot be less than 0:%v", playerID)
		rsp.Result = false
		rsp.ErrCode = robot.ErrCode_EC_Args
		return rsp, nil
	}
	newState := request.GetNewState()
	err := data.UpdataRobotState(playerID, newState)
	if err != nil {
		rsp.ErrCode = robot.ErrCode_EC_FAIL
		logrus.WithError(err).Debugln("更新空闲机器人状态失败")
		return rsp, err
	}
	logrus.WithFields(logrus.Fields{
		"RobotPlayerId": playerID,
		"newState":      newState,
	}).Debugln("更新空闲机器人状态成功")
	return rsp, nil
}

// UpdataRobotGameWinRate 更新胜率
func (r *Robotservice) UpdataRobotGameWinRate(ctx context.Context, request *robot.UpdataRobotGameWinRateReq) (*robot.UpdataRobotGameWinRateRsp, error) {
	logrus.Debugln("UpdataRobotGameWinRate req", *request)
	rsp := &robot.UpdataRobotGameWinRateRsp{
		Result:  true,
		ErrCode: robot.ErrCode_EC_SUCCESS,
	}
	playerID := request.GetRobotPlayerId()
	gameID := int(request.GetGameId())
	newWinRate := request.GetNewWinRate()

	if playerID < 0 {
		logrus.Warningln("Robot Player ID cannot be less than 0:%d", playerID)
		rsp.ErrCode = robot.ErrCode_EC_Args
		return rsp, nil
	}
	err := data.UpdataRobotWinRate(playerID, gameID, newWinRate)
	if err != nil {
		rsp.ErrCode = robot.ErrCode_EC_FAIL
		logrus.WithError(err).Debugln("更新空闲机器人胜率失败")
		return rsp, err
	}
	logrus.WithFields(logrus.Fields{
		"RobotPlayerId": playerID,
		"newWinRate":    newWinRate,
	}).Debugln("更新空闲机器人胜率成功")
	return rsp, nil
}

//IsRobotPlayer 判断是否时机器人
func (r *Robotservice) IsRobotPlayer(ctx context.Context, request *robot.IsRobotPlayerReq) (*robot.IsRobotPlayerRsp, error) {
	logrus.Debugln("IsRobotPlayer req", *request)
	rsp := &robot.IsRobotPlayerRsp{
		Result: false,
	}
	playerID := request.GetRobotPlayerId()
	if playerID <= 0 {
		return rsp, fmt.Errorf("参数错误")
	}
	robotInfo, err := data.GetRobotInfoByPlayerID(playerID)
	if robotInfo != nil {
		rsp.Result = true
		return rsp, err
	}
	return rsp, err
}
