package data

import (
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"

	"github.com/go-xorm/xorm"
)

// 是否需要验证操作玩家是否有权限
const (
	// MysqldbName 数据库名
	MysqldbName             = "robot"
	playerTableName         = "t_player"
	playerCurrencyTableName = "t_player_currency"
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

//获取所有机器人的各个属性
func getMysqlRobotProp(robotMap map[int64]*cache.RobotPlayer) error {
	err := getMysqlRobotAllCoins(robotMap) // 获取金币
	return err
}

//根据玩家ID获取机器人的各个属性
func getMysqlRobotPropByPlayerID(playerID uint64) (*cache.RobotPlayer, error) {
	robotPlayer := &cache.RobotPlayer{}
	coin, err := getMysqlRobotCoinsByPlayerID(playerID) // 金币
	robotPlayer.Coin = uint64(coin)
	return robotPlayer, err
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
		return nil, fmt.Errorf("获取所有机器人失败:%v", err)
	}
	robotsIDAll := make([]uint64, 0, len(robots))
	for _, robot := range robots {
		robotsIDAll = append(robotsIDAll, uint64(robot.Playerid))
	}
	return robotsIDAll, nil
}

//根据玩家ID获取机器人的金币
func getMysqlRobotCoinsByPlayerID(playerID uint64) (int, error) {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return 0, err
	}
	pct := &db.TPlayerCurrency{}
	where := fmt.Sprintf("playerID=%v", playerID)
	exist, err := engine.Table(playerCurrencyTableName).Where(where).Select("coins").Get(pct)
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, fmt.Errorf("获取机器人金币失败 : %v", playerID)
	}
	return pct.Coins, nil
}

//Join 查询机器人对应在playerCurrencyTable 的金币
func getMysqlRobotAllCoins(robotMap map[int64]*cache.RobotPlayer) error {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return err
	}
	robots := make([]*db.TPlayerCurrency, 0)
	idEqu := fmt.Sprintf("%v.playerID = %v.playerID", playerTableName, playerCurrencyTableName)
	where := fmt.Sprintf("type=%v", playerType)
	Select := fmt.Sprintf("%v.playerID,%v.coins", playerCurrencyTableName, playerCurrencyTableName)
	err = engine.Table(playerCurrencyTableName).Join("INNER", playerTableName, idEqu).Where(where).Select(Select).Find(&robots)
	for _, robot := range robots {
		if rp := robotMap[robot.Playerid]; rp != nil {
			rp.PlayerID = uint64(robot.Playerid)
			rp.Coin = uint64(robot.Coins)
			robotMap[robot.Playerid] = rp
		} else {
			robotMap[robot.Playerid] = &cache.RobotPlayer{
				PlayerID: uint64(robot.Playerid),
				Coin:     uint64(robot.Coins),
			}
		}
	}
	return err
}
