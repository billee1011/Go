package robotservice

import (
	"errors"
	"fmt"
	"steve/robot/data"
	"steve/server_pb/robot"
	"strconv"
)

var errReqState = errors.New("已经是当前状态")
var errReqRange = errors.New("High 必须大于等于 Low")
var errReqRangeSize = errors.New("High 和 Low 都必须大于等于0")

// 检验金币和胜率
func checkCoinsWinRtaeRange(coinsRange *robot.CoinsRange, winRateRange *robot.WinRateRange) error {
	switch {
	case coinsRange.High < coinsRange.Low:
		return fmt.Errorf("coinsRange:%v", errReqRange)
	case coinsRange.High < 0 || coinsRange.Low < 0:
		return fmt.Errorf("coinsRange:%v", errReqRangeSize)
	case winRateRange.High < winRateRange.Low:
		return fmt.Errorf("winRateRange:%v", errReqRange)
	case winRateRange.High < 0 || winRateRange.Low < 0:
		return fmt.Errorf("winRateRange:%v", errReqRangeSize)
	}
	return nil
}

// 检验状态和playerID
func checkSetRobotPlayerStateReq(request *robot.SetRobotPlayerStateReq) error {
	playerID := request.GetRobotPlayerId()
	if playerID < 0 {
		return fmt.Errorf("playerID 小于0 :%v", playerID)
	}
	state := request.GetState()
	if state == robot.RobotPlayerState_RPS_IDIE || state == robot.RobotPlayerState_RPS_GAMEING || state == robot.RobotPlayerState_RPS_MATCHING {
		return nil
	}
	return fmt.Errorf("state 不在存 :%v", state)
}

// 先检验请求是否合法, 获取所有空闲的机器人,获取否符合金币PlayerID
func getRobotPlayerIDByInfo(request *robot.GetRobotPlayerIDReq) (uint64, error) {
	coinsRange := request.GetCoinsRange()
	winRateRange := request.GetWinRateRange()
	// 检验请求是否合法
	if err := checkCoinsWinRtaeRange(coinsRange, winRateRange); err != nil {
		return 0, err
	}
	robotsIDCoins, err := data.GetLeisureRobot() // 获取机空闲的器人
	if err != nil {
		return 0, err
	}
	var RobotPlayerID uint64
	// 符合的指定金币数的机器人
	for _, robotPlayer := range robotsIDCoins {
		currCoins := int32(robotPlayer.Coin)
		if currCoins <= coinsRange.High && currCoins >= coinsRange.Low {
			RobotPlayerID = robotPlayer.PlayerID
			break
		}
	}
	return RobotPlayerID, nil
}

// 先检验请求是否合法, 在设置机器人玩家状态
func setRobotPlayerState(request *robot.SetRobotPlayerStateReq) error {
	if err := checkSetRobotPlayerStateReq(request); err != nil {
		return err
	}
	playerID := request.GetRobotPlayerId()
	reqState := request.GetState()
	stateStr, err := data.GetRobotStringFiled(playerID, data.RobotPlayerStateField)
	if err != nil {
		return err
	}
	state, err := strconv.ParseUint(stateStr, 10, 64)
	if err != nil {
		return err
	}
	if state == uint64(reqState) {
		return errReqState
	}
	return data.SetRobotFiled(playerID, data.RobotPlayerStateField, int(reqState))
}
