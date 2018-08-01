package data

import (
	"steve/entity/cache"
	"steve/structs"

	"github.com/Sirupsen/logrus"
)

const (
	//RobotPlayerIDField 玩家ID字段名
	RobotPlayerIDField string = "playerID"
	//RobotPlayerCoinField 玩家金币数字段名
	RobotPlayerCoinField string = "coin"
	//RobotPlayerStateField 玩家状态字段名
	RobotPlayerStateField string = "state"
	//RobotPlayerGameIDWinRate 玩家游戏ID对应的胜率字段名
	RobotPlayerGameIDWinRate string = "gameidwinrate"
	//RobotPlayerGameIDField 玩家游戏 ID 字段名
	RobotPlayerGameIDField string = "gameid"
	//RobotPlayerNickNameField 玩家昵称字段
	RobotPlayerNickNameField string = "nickname"
	//RobotPlayerHeadImageField 玩家头像字段
	RobotPlayerHeadImageField string = "head_image"
)

var Exposer *structs.Exposer

//InitRobotRedis 初始化机器人redis
func InitRobotRedis() {
	robotMap := make(map[int64]*cache.RobotPlayer)
	log := logrus.WithFields(logrus.Fields{"func_name": "initRobotRedis"})
	if err := getMysqlRobotProp(robotMap); err != nil {
		log.WithError(err).Errorln("初始化从mysql获取机器人失败")
		return
	}
	failedIDErrMpa := make(map[uint64]error) //存入redis 失败 playerID
	for playerID, robotPlayer := range robotMap {
		err := AddRobottFiled(uint64(playerID), FmtRobotPlayer(robotPlayer))
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
	robotsIDCoins, lackRobotsID := getRedisLeisureRobotPlayer(robotPlayerIDAll) // 从redis 获取 空闲的RobotPlayer
	if len(lackRobotsID) > 0 {                                                  // 存在redis 中 不存在 机器人ID
		robotsIDCoins = getMysqlLeisureRobotPlayer(robotsIDCoins, lackRobotsID) //从mysql中获取空闲的玩家,并存入redis
	}
	log.WithFields(logrus.Fields{"lackRobotsID": lackRobotsID, "robotsIDCoins": robotsIDCoins}).Infoln("获取空闲机器人")
	return robotsIDCoins, nil
}
