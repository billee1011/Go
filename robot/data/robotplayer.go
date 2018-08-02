package data

import (
	"steve/entity/cache"
	"steve/server_pb/robot"

	"github.com/Sirupsen/logrus"
)

//getRedisLeisureRobotPlayer 从redis 获取 空闲的RobotPlayer
func getRedisLeisureRobotPlayer(robotPlayerIDAll []uint64) ([]*cache.RobotPlayer, []uint64) {
	robotsIDCoins := make([]*cache.RobotPlayer, 0)
	lackRobotsID := make([]uint64, 0) // 没有存入redis的机器人
	for _, robotPlayerID := range robotPlayerIDAll {
		robotPlayerInfo, err := GetRobotFields(robotPlayerID, RobotPlayerCoinField, cache.PlayerStateField, RobotPlayerGameIDWinRate)
		if err != nil || len(robotPlayerInfo) == 0 {
			lackRobotsID = append(lackRobotsID, robotPlayerID)
			continue
		}
		robotPlayer := &cache.RobotPlayer{}
		robotPlayer.State = InterToUint64(robotPlayerInfo[cache.PlayerStateField]) // 玩家状态
		if robotPlayer.State != uint64(robot.RobotPlayerState_RPS_IDIE) {          //是空闲状态
			continue
		}
		robotPlayer.PlayerID = robotPlayerID                                                                // 玩家ID
		robotPlayer.Coin = InterToUint64(robotPlayerInfo[RobotPlayerCoinField])                             // 金币
		robotPlayer.GameIDWinRate = JSONToGameIDWinRate(robotPlayerInfo[RobotPlayerGameIDWinRate].(string)) // 游戏对应的胜率
		robotsIDCoins = append(robotsIDCoins, robotPlayer)
	}
	return robotsIDCoins, lackRobotsID
}

//getMysqlLeisureRobotPlayer 从mysql中获取空闲的玩家,并存入redis
func getMysqlLeisureRobotPlayer(robotsIDCoins []*cache.RobotPlayer, lackRobotsID []uint64) []*cache.RobotPlayer {
	log := logrus.WithFields(logrus.Fields{"func_name": "getMysqlLeisureRobotPlayer"})
	failedIDErrMpa := make(map[uint64]error) //存入redis 失败 playerID
	for _, playerID := range lackRobotsID {
		robotPlayer, err := getMysqlRobotPropByPlayerID(playerID) // 从mysql获取 的一定是空闲的
		if err != nil {
			failedIDErrMpa[playerID] = err
			continue
		}
		err = AddRobotWatch(playerID, FmtRobotPlayer(robotPlayer), RedisTimeOut) // 存入redis
		if err != nil {
			failedIDErrMpa[playerID] = err
		}
		robotPlayer.PlayerID = playerID
		robotsIDCoins = append(robotsIDCoins, robotPlayer)
	}
	if len(failedIDErrMpa) > 0 {
		log = log.WithFields(logrus.Fields{"failedIDErrMpa": failedIDErrMpa})
	}
	log.Info("从mysql获取机器人,并存入redis")
	return robotsIDCoins
}
