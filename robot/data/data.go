package data

import (
	"steve/entity/cache"

	"github.com/Sirupsen/logrus"
)

const (
	//RobotPlayerIDField 玩家ID字段名
	RobotPlayerIDField string = "playerID"
	//RobotPlayerGameIDWinRate 玩家游戏ID对应的胜率字段名
	RobotPlayerGameIDWinRate string = "gameidwinrate"
	//RobotPlayerGameIDField 玩家游戏 ID 字段名
	RobotPlayerGameIDField string = "gameid"
	//RobotPlayerNickNameField 玩家昵称字段
	RobotPlayerNickNameField string = "nickname"
	// RobotPlayerAvatarField 玩家头像
	RobotPlayerAvatarField = "avatar"
	// RobotPlayerStateField ...玩家状态
	RobotPlayerStateField = "game_state"
	// RobotPlayerGateAddrField ...网关服地址
	RobotPlayerGateAddrField = "gate_addr"
	// RobotPlayerMatchAddrField ...匹配服地址
	RobotPlayerMatchAddrField = "match_addr"
	// RobotPlayerRoomAddrField ...房间服地址
	RobotPlayerRoomAddrField = "room_addr"
)

//InitRobotRedis 初始化机器人redis
func InitRobotRedis() {
	robotMap := make(map[int64]*cache.RobotPlayer)
	log := logrus.WithFields(logrus.Fields{"func_name": "initRobotRedis"})
	if err := getMysqlRobotFieldValuedAll(robotMap); err != nil {
		log.WithError(err).Errorln("初始化从mysql获取机器人失败")
		return
	}
	failedIDErrMpa := make(map[uint64]error) //存入redis 失败 playerID
	for playerID, robotPlayer := range robotMap {
		err := AddRobotWatch(uint64(playerID), FmtRobotPlayer(robotPlayer), RedisTimeOut)
		if err != nil {
			failedIDErrMpa[uint64(playerID)] = err
			continue
		}
	}
	if len(failedIDErrMpa) > 0 {
		log.WithFields(logrus.Fields{"failedIDErrMpa": failedIDErrMpa}).Info("失败的playerID")
	}
}

//GetLeisureRobot 获取空闲机器人
func GetLeisureRobot() ([]*cache.RobotPlayer, error) {
	log := logrus.WithFields(logrus.Fields{"func_name": "GetLeisureRobot"})
	robotPlayerIDAll, err := getRobotIDAll() // 获取所有机器人的玩家ID
	if err != nil {
		return nil, err
	}
	if len(robotPlayerIDAll) == 0 {
		log.Info("数据库中不存在机器人")
		return []*cache.RobotPlayer{}, nil
	}
	robots, lackRobotsID := getRedisLeisureRobotPlayer(robotPlayerIDAll) // 从redis 获取 空闲的RobotPlayer
	if len(lackRobotsID) > 0 {                                           // 存在redis 中 不存在 机器人ID
		robots = getMysqlLeisureRobotPlayer(robots, lackRobotsID) //从mysql中获取空闲的玩家,并存入redis
	}
	log.WithFields(logrus.Fields{"lackRobotsID": lackRobotsID}).Infoln("获取空闲机器人")
	return robots, nil
}
