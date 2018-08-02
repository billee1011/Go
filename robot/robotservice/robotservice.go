package robotservice

import (
	"context"
	"fmt"
	"steve/entity/cache"
	"steve/robot/data"
	"steve/server_pb/robot"
	"strconv"

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

//GetRobotPlayerIDByInfo 根据请求信息获取机器人玩家ID
func (r *Robotservice) GetRobotPlayerIDByInfo(ctx context.Context, request *robot.GetRobotPlayerIDReq) (*robot.GetRobotPlayerIDRsp, error) {
	logrus.Debugln("GetRobotPlayerIDByInfo req", *request)
	rsp := &robot.GetRobotPlayerIDRsp{
		RobotPlayerId: 0,
		ErrCode:       int32(robot.ErrCode_EC_SUCCESS),
	}
	gameID := request.GetGame().GetGameId()   // 游戏ID
	coinsRange := request.GetCoinsRange()     // 金币范围
	winRateRange := request.GetWinRateRange() // 胜率范围
	// 检验请求是否合法
	if !checkCoinsWinRtaeRange(coinsRange, winRateRange) {
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		return rsp, nil
	}
	robotsPlayers, err := data.GetLeisureRobot() // 获取机空闲的器人
	if err != nil {
		rsp.ErrCode = int32(robot.ErrCode_EC_FAIL)
		return rsp, err
	}
	var RobotPlayerID uint64
	// 符合的指定金币数和胜率的机器人
	for _, robotPlayer := range robotsPlayers {
		// 游戏ID对应的胜率
		winRate, exist := robotPlayer.GameIDWinRate[uint64(gameID)]
		if !exist || winRate > uint64(winRateRange.High) || winRate < uint64(winRateRange.Low) {
			continue
		}
		// 金币
		currCoins := int64(robotPlayer.Coin)
		if currCoins <= coinsRange.High && currCoins >= coinsRange.Low {
			RobotPlayerID = robotPlayer.PlayerID
			break
		}
	}
	rsp.RobotPlayerId = RobotPlayerID
	return rsp, err
}

//SetRobotPlayerState 设置机器人玩家状态
func (r *Robotservice) SetRobotPlayerState(ctx context.Context, request *robot.SetRobotPlayerStateReq) (*robot.SetRobotPlayerStateRsp, error) {
	logrus.Debugln("SetRobotPlayerState req", *request)
	rsp := &robot.SetRobotPlayerStateRsp{
		Result:  true,
		ErrCode: int32(robot.ErrCode_EC_SUCCESS),
	}
	playerID := request.GetRobotPlayerId()
	newState := int(request.GetNewstate())
	oldState := int(request.GetOldstate())
	severType := int(request.GetServerType())
	serverAddr := request.GetServerAddr()

	// 检验请求是否合法
	if !checkSateArgs(playerID, newState, oldState, severType, serverAddr) {
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		return rsp, nil
	}

	//比较请求旧状态是否是当前状态
	val, _ := data.GetRobotStringFiled(playerID, cache.PlayerStateField)
	state, _ := strconv.Atoi(val)
	if oldState != state {
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		return rsp, nil
	}

	//修改状态和服务地址
	serverField := map[robot.ServerType]string{
		robot.ServerType_ST_GATE:  cache.GateAddrField,
		robot.ServerType_ST_MATCH: cache.MatchAddrField,
		robot.ServerType_ST_ROOM:  cache.RoomAddrField,
	}[robot.ServerType(severType)]
	rfields := map[string]interface{}{
		cache.PlayerStateField: fmt.Sprintf("%d", newState),
		serverField:            serverAddr,
	}
	if err := data.SetRobotPlayerWatchs(playerID, rfields, data.RedisTimeOut); err != nil {
		rsp.ErrCode = int32(robot.ErrCode_EC_FAIL)
		rsp.Result = false
		return rsp, err
	}
	return rsp, nil
}
