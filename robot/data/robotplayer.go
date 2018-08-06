package data

import (
	"steve/entity/cache"
	"steve/external/goldclient"
	"steve/server_pb/gold"
	"steve/server_pb/robot"

	"github.com/Sirupsen/logrus"
)

//getRedisLeisureRobotPlayer 从redis 获取 空闲的RobotPlayer
func getRedisLeisureRobotPlayer(robotPlayerIDAll []uint64) ([]*cache.RobotPlayer, []uint64) {
	robotsIDCoins := make([]*cache.RobotPlayer, 0)
	lackRobotsID := make([]uint64, 0) // 没有存入redis的机器人
	for _, robotPlayerID := range robotPlayerIDAll {
		robotPlayerInfo, err := GetRobotFields(robotPlayerID, cache.GameState, RobotPlayerGameIDWinRate)
		if err != nil || len(robotPlayerInfo) == 0 {
			lackRobotsID = append(lackRobotsID, robotPlayerID)
			continue
		}
		robotPlayer := &cache.RobotPlayer{}
		robotPlayer.State = InterToUint64(robotPlayerInfo[cache.GameState]) // 玩家状态
		if robotPlayer.State != uint64(robot.RobotPlayerState_RPS_IDIE) {   //是空闲状态
			continue
		}
		robotPlayer.PlayerID = robotPlayerID
		// 从金币服获取
		gold, err := goldclient.GetGold(robotPlayerID, int16(gold.GoldType_GOLD_COIN))
		if err != nil {
			logrus.WithError(err).Errorf("获取金币失败 robotPlayerID(%v)", robotPlayerIDAll)
			continue
		}
		// 玩家ID
		robotPlayer.Coin = uint64(gold)                                                                     // 金币
		robotPlayer.GameIDWinRate = JSONToGameIDWinRate(robotPlayerInfo[RobotPlayerGameIDWinRate].(string)) // 游戏对应的胜率
		robotsIDCoins = append(robotsIDCoins, robotPlayer)
	}
	return robotsIDCoins, lackRobotsID
}

//getMysqlLeisureRobotPlayer 从mysql中获取空闲的玩家,并存入redis
func getMysqlLeisureRobotPlayer(robotsPlayers []*cache.RobotPlayer, lackRobotsID []uint64) []*cache.RobotPlayer {
	log := logrus.WithFields(logrus.Fields{"func_name": "getMysqlLeisureRobotPlayer"})
	failedIDErrMpa := make(map[uint64]error) //存入redis 失败 playerID
	for _, playerID := range lackRobotsID {
		robotPlayer := getMysqlRobotPropByPlayerID(playerID) // 从mysql获取 的一定是空闲的
		// 存入redis
		if err := SetRobotPlayerWatchs(playerID, FmtRobotPlayer(robotPlayer), RedisTimeOut); err != nil {
			failedIDErrMpa[playerID] = err
		}
		robotPlayer.PlayerID = playerID
		robotsPlayers = append(robotsPlayers, robotPlayer)
	}
	if len(failedIDErrMpa) > 0 {
		log = log.WithFields(logrus.Fields{"failedIDErrMpa": failedIDErrMpa})
	}
	log.Info("从mysql获取机器人,并存入redis")
	return robotsPlayers
}

//获取机器人需要的值
func getMysqlRobotFieldValuedAll(robotMap map[int64]*cache.RobotPlayer) error {
	//gameid-winrate 游戏id对应的胜率
	robotsPGs, err := getMysqlRobotGameWinRateAll()
	if err != nil {
		return err
	}
	for _, robot := range robotsPGs {
		if rp := robotMap[robot.Playerid]; rp != nil {
			rp.GameIDWinRate[uint64(robot.Gameid)] = uint64(robot.Winningrate)
			robotMap[robot.Playerid] = rp
		} else {
			robotMap[robot.Playerid] = &cache.RobotPlayer{
				PlayerID:      uint64(robot.Playerid),
				GameIDWinRate: map[uint64]uint64{uint64(robot.Gameid): uint64(robot.Winningrate)},
			}
		}
	}
	// 昵称
	robotsTPs, err := getMysqlRobotNicknameAll()
	if err != nil {
		return err
	}
	for _, robot := range robotsTPs {
		if rp := robotMap[robot.Playerid]; rp != nil {
			rp.NickName = robot.Nickname
			robotMap[robot.Playerid] = rp
		} else {
			robotMap[robot.Playerid] = &cache.RobotPlayer{
				PlayerID: uint64(robot.Playerid),
				NickName: robot.Nickname,
			}
		}
	}
	return err
}
