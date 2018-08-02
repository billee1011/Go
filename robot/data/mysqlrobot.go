package data

import (
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"

	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/xorm"
)

// 是否需要验证操作玩家是否有权限
const (
	// MysqldbName 数据库名
	MysqldbName             = "steve"
	playerTableName         = "t_player"          // 玩家表
	playerCurrencyTableName = "t_player_currency" // 玩家货币表
	playerGameTableName     = "t_player_game"     // 玩家游戏表
	// 玩家类型机器人
	playerType = 2
)

//MysqlEnginefunc 单元测试需要
var MysqlEnginefunc = getMysqlEngineByName

func getMysqlEngineByName(mysqlName string) (*xorm.Engine, error) {
	engine, err := Exposer.MysqlEngineMgr.GetEngine(mysqlName)
	if err != nil {
		return nil, fmt.Errorf("获取 mysql 引擎失败：%v", err)
	}
	if err := engine.Ping(); err != nil {
		return nil, fmt.Errorf("engine ping 失败：%v", err)
	}
	return engine, nil
}

//根据玩家ID获取机器人的各个属性
func getMysqlRobotPropByPlayerID(playerID uint64) *cache.RobotPlayer {
	robotPlayer := &cache.RobotPlayer{}
	playerCurrency, err := getMysqlPlayerCurrencyByPlayerID(playerID, "coins") // 金币
	if err != nil {
		logrus.Errorf("msql获取金币失败 err(%v)", playerID, err)
	}
	player, err := getMysqlPlayerByPlayerID(playerID, "nickname") // 昵称
	if err != nil {
		logrus.Errorf("msql获取昵称失败 err(%v)", playerID, err)
	}
	playerGame, err := getMysqlPlayerGameByPlayerID(playerID, "gameID,winningRate") //游戏ID和胜率
	if err != nil {
		logrus.Errorf("msql获取游戏ID和胜率失败 err(%v)", playerID, err)
	}
	if robotPlayer.GameIDWinRate == nil || len(robotPlayer.GameIDWinRate) == 0 {
		robotPlayer.GameIDWinRate = map[uint64]uint64{uint64(playerGame.Gameid): uint64(playerGame.Winningrate)}
	} else {
		robotPlayer.GameIDWinRate[uint64(playerGame.Gameid)] = uint64(playerGame.Winningrate)
	}
	robotPlayer.Coin = uint64(playerCurrency.Coins)
	robotPlayer.NickName = player.Nickname
	return robotPlayer
}

// 获取所有机器人PlayerID
func getRobotIDAll() ([]uint64, error) {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return nil, err
	}
	robots := make([]*db.TPlayer, 0)
	err = engine.Table(playerTableName).Where(fmt.Sprintf("type=%v", playerType)).Find(&robots)
	if err != nil {
		return nil, err
	}
	robotsIDAll := make([]uint64, 0, len(robots))
	for _, robot := range robots {
		robotsIDAll = append(robotsIDAll, uint64(robot.Playerid))
	}
	return robotsIDAll, nil
}

//根据玩家ID获取获取玩家货币表上的数据
func getMysqlPlayerCurrencyByPlayerID(playerID uint64, result string) (*db.TPlayerCurrency, error) {
	pct := &db.TPlayerCurrency{}
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return pct, err
	}
	where := fmt.Sprintf("playerID=%v", playerID)
	exist, err := engine.Table(playerCurrencyTableName).Where(where).Select(result).Get(pct)
	if err != nil {
		return pct, err
	}
	if !exist {
		return pct, fmt.Errorf("TPlayerCurrency获取失败 : %v", playerID)
	}
	return pct, nil
}

// 根据玩家ID获取玩家表上的数据
func getMysqlPlayerByPlayerID(playerID uint64, result string) (*db.TPlayer, error) {
	pt := &db.TPlayer{}
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return pt, err
	}
	where := fmt.Sprintf("playerID=%v", playerID)
	exist, err := engine.Table(playerTableName).Where(where).Select(result).Get(pt)
	if err != nil {
		return pt, err
	}
	if !exist {
		return pt, fmt.Errorf("TPlayer获取失败 : %v", playerID)
	}
	return pt, nil
}

//根据玩家ID获取玩家游戏表上的数据
func getMysqlPlayerGameByPlayerID(playerID uint64, result string) (*db.TPlayerGame, error) {
	pt := &db.TPlayerGame{}
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return pt, err
	}
	where := fmt.Sprintf("playerID=%v", playerID)
	exist, err := engine.Table(playerGameTableName).Where(where).Select(result).Get(pt)
	if err != nil {
		return pt, err
	}
	if !exist {
		return pt, fmt.Errorf("TPlayerGame获取失败 : %v", playerID)
	}
	return pt, nil
}

//获取所有机器人的游戏ID和对应的胜率
func getMysqlRobotGameWinRateAll() ([]*db.TPlayerGame, error) {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return nil, err
	}
	// 胜率
	robotsPGs := make([]*db.TPlayerGame, 0)
	idEqu := fmt.Sprintf("%v.playerID = %v.playerID", playerTableName, playerGameTableName)
	where := fmt.Sprintf("type=%v", playerType)
	Select := fmt.Sprintf("%v.playerID,%v.winningRate,%v.gameID", playerGameTableName, playerGameTableName, playerGameTableName)
	if err := engine.Table(playerGameTableName).Join("INNER", playerTableName, idEqu).Where(where).Select(Select).Find(&robotsPGs); err != nil {
		return nil, err
	}
	return robotsPGs, nil
}

//获取所有机器人的金币
func getMysqlRobotCoinAll() ([]*db.TPlayerCurrency, error) {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return nil, err
	}
	// 金币
	robotsTPCs := make([]*db.TPlayerCurrency, 0)
	idEqu := fmt.Sprintf("%v.playerID = %v.playerID", playerTableName, playerCurrencyTableName)
	where := fmt.Sprintf("type=%v", playerType)
	Select := fmt.Sprintf("%v.playerID,%v.coins", playerTableName, playerCurrencyTableName)
	if err := engine.Table(playerCurrencyTableName).Join("INNER", playerTableName, idEqu).Where(where).Select(Select).Find(&robotsTPCs); err != nil {
		return nil, err
	}
	return robotsTPCs, nil
}

// 获取所有机器人的昵称
func getMysqlRobotNicknameAll() ([]*db.TPlayer, error) {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return nil, err
	}
	// 昵称
	where := fmt.Sprintf("type=%v", playerType)
	robotsTPs := make([]*db.TPlayer, 0)
	if err := engine.Table(playerTableName).Select("playerID,nickname").Where(where).Find(&robotsTPs); err != nil {
		return nil, err
	}
	return robotsTPs, nil
}
