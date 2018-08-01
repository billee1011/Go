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
	MysqldbName             = "steve"
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

//根据玩家ID获取机器人的各个属性
func getMysqlRobotPropByPlayerID(playerID uint64) (*cache.RobotPlayer, error) {
	robotPlayer := &cache.RobotPlayer{}
	playerCurrency, err := getMysqlPlayerCurrencyByPlayerID(playerID) // 金币
	if err != nil {
		return robotPlayer, err
	}
	player, err := getMysqlPlayerByPlayerID(playerID) // 金币
	if err != nil {
		return robotPlayer, err
	}
	robotPlayer.Coin = uint64(playerCurrency.Coins)
	robotPlayer.NickName = player.Nickname
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
		return nil, err
	}
	robotsIDAll := make([]uint64, 0, len(robots))
	for _, robot := range robots {
		robotsIDAll = append(robotsIDAll, uint64(robot.Playerid))
	}
	return robotsIDAll, nil
}

//根据玩家ID获取获取玩家货币表上的数据
func getMysqlPlayerCurrencyByPlayerID(playerID uint64) (*db.TPlayerCurrency, error) {
	pct := &db.TPlayerCurrency{}
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return pct, err
	}
	where := fmt.Sprintf("playerID=%v", playerID)
	exist, err := engine.Table(playerCurrencyTableName).Where(where).Select("coins").Get(pct)
	if err != nil {
		return pct, err
	}
	if !exist {
		return pct, fmt.Errorf("TPlayerCurrency获取失败 : %v", playerID)
	}
	return pct, nil
}

// 根据玩家ID获取获取玩家表上的数据
func getMysqlPlayerByPlayerID(playerID uint64) (*db.TPlayer, error) {
	pt := &db.TPlayer{}
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return pt, err
	}
	where := fmt.Sprintf("playerID=%v", playerID)
	exist, err := engine.Table(playerTableName).Where(where).Select("nickname").Get(pt)
	if err != nil {
		return pt, err
	}
	if !exist {
		return pt, fmt.Errorf("TPlayer获取失败 : %v", playerID)
	}
	return pt, nil
}

//Join 获取机器人需要的值
func getMysqlRobotFieldValuedAll(robotMap map[int64]*cache.RobotPlayer) error {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return err
	}
	// 金币
	robotsTPCs := make([]*db.TPlayerCurrency, 0)
	idEqu := fmt.Sprintf("%v.playerID = %v.playerID", playerTableName, playerCurrencyTableName)
	where := fmt.Sprintf("type=%v", playerType)
	Select := fmt.Sprintf("%v.playerID,%v.coins", playerTableName, playerCurrencyTableName)
	if err := engine.Table(playerCurrencyTableName).Join("INNER", playerTableName, idEqu).Where(where).Select(Select).Find(&robotsTPCs); err != nil {
		return err
	}
	fmt.Println(robotsTPCs)
	for _, robot := range robotsTPCs {
		if rp := robotMap[robot.Playerid]; rp != nil {
			if rp.PlayerID == 0 {
				rp.PlayerID = uint64(robot.Playerid)
			}
			rp.Coin = uint64(robot.Coins)
			robotMap[robot.Playerid] = rp
		} else {
			robotMap[robot.Playerid] = &cache.RobotPlayer{
				PlayerID: uint64(robot.Playerid),
				Coin:     uint64(robot.Coins),
			}
		}
	}
	// 昵称
	robotsTPs := make([]*db.TPlayer, 0)
	if err := engine.Table(playerTableName).Select("playerID,nickname").Where(where).Find(&robotsTPs); err != nil {
		return err
	}
	for _, robot := range robotsTPs {
		if rp := robotMap[robot.Playerid]; rp != nil {
			if rp.PlayerID == 0 {
				rp.PlayerID = uint64(robot.Playerid)
			}
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
