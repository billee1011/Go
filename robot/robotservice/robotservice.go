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
		ErrCode:       int32(robot.ErrCode_EC_FAIL),
	}
	gameID := request.GetGame().GetGameId()   // 游戏ID
	coinsRange := request.GetCoinsRange()     // 金币范围
	winRateRange := request.GetWinRateRange() // 胜率范围
	newState := int(request.GetNewState())    // 获取成功时设置的状态
	// 检验请求是否合法
	if !checkGetLeisureRobotArgs(coinsRange, winRateRange, newState) {
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
	if RobotPlayerID == 0 {
		return rsp, fmt.Errorf("没有适合的机器人")
	}
	//获取到机器人ID,并将redis该ID的状态为匹配状态
	if err := data.SetRobotWatch(RobotPlayerID, cache.GameState, newState, data.RedisTimeOut); err != nil {
		return rsp, err
	}
	rsp.ErrCode = int32(robot.ErrCode_EC_SUCCESS)
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
	newState := int(request.GetNewState())
	oldState := int(request.GetOldState())
	severType := int(request.GetServerType())
	serverAddr := request.GetServerAddr()

	// 检验请求是否合法
	if !checkSateArgs(playerID, newState, oldState, severType, serverAddr) {
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		return rsp, nil
	}

	//比较请求旧状态是否是当前状态
	val, _ := data.GetRobotStringFiled(playerID, cache.GameState)
	state, _ := strconv.Atoi(val)
	if oldState != state {
		rsp.ErrCode = int32(robot.ErrCode_EC_Args)
		return rsp, nil
	}

	//修改状态和服务地址
	serverField := map[robot.ServerType]string{
		robot.ServerType_ST_GATE:  cache.GateAddr,
		robot.ServerType_ST_MATCH: cache.MatchAddr,
		robot.ServerType_ST_ROOM:  cache.RoomAddr,
	}[robot.ServerType(severType)]
	rfields := map[string]interface{}{
		cache.GameState: fmt.Sprintf("%d", newState),
		serverField:     serverAddr,
	}
	if err := data.SetRobotPlayerWatchs(playerID, rfields, data.RedisTimeOut); err != nil {
		rsp.ErrCode = int32(robot.ErrCode_EC_FAIL)
		rsp.Result = false
		return rsp, err
	}
	return rsp, nil
}
