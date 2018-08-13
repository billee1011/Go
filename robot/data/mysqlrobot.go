package data

import (
	"fmt"
	"steve/entity/db"
	"steve/structs"

	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/xorm"
)

const (
	// MysqldbName 数据库名
	MysqldbName             = "player"
	playerTableName         = "t_player"          // 玩家表
	playerCurrencyTableName = "t_player_currency" // 玩家货币表
	playerGameTableName     = "t_player_game"     // 玩家游戏表
	//PlayerType 玩家类型机器人
	PlayerType = 2
)

//MysqlEnginefunc 单元测试需要
var MysqlEnginefunc = getMysqlEngineByName

func getMysqlEngineByName(mysqlName string) (*xorm.Engine, error) {
	e := structs.GetGlobalExposer()
	engine, err := e.MysqlEngineMgr.GetEngine(mysqlName)
	if err != nil {
		return nil, fmt.Errorf("获取 mysql 引擎失败：%v", err)
	}
	return engine, nil
}

// Startlimit 分页
var Startlimit int

const limit int = 100 // 固定页数 100

//IsMysqlRobot 判断是否时机器人
func IsMysqlRobot(playerID int64) (bool, error) {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return false, err
	}
	p := &db.TPlayer{}
	session := engine.Table(playerTableName).Select("type").Where(fmt.Sprintf("playerID=%d", playerID))
	exist, err := session.Get(p)
	if err != nil {
		sql, _ := session.LastSQL()
		return false, fmt.Errorf("err(%v),sql(%v)", err, sql)
	}
	if !exist {
		logrus.Debugf("该玩家ID不存在 playerID (%d) ", playerID)
		return false, fmt.Errorf("该玩家ID不存在 playerID (%d)", playerID)
	}
	return p.Type == PlayerType, nil
}

//GetRobotInfoByPlayerID 根据玩家ID获取游戏信息
func GetRobotInfoByPlayerID(playerID int64) ([]*db.TPlayerGame, error) {
	engine, err := MysqlEnginefunc(MysqldbName)
	if err != nil {
		return nil, err
	}
	robotsPGs := make([]*db.TPlayerGame, 0)
	session := engine.Table(playerGameTableName).Select("playerID,winningRate,gameID").Where(fmt.Sprintf("playerID=%d", playerID))
	if err := session.Find(robotsPGs); err != nil {
		sql, _ := session.LastSQL()
		return nil, fmt.Errorf("err(%v),sql(%v)", err, sql)
	}
	return robotsPGs, nil
}

//获取机器人的的游戏ID和对应的胜率
func getRobotInfo(engine *xorm.Engine) ([]*db.TPlayerGame, error) {
	// 游戏胜率
	robotsPGs := make([]*db.TPlayerGame, 0)
	idEqu := fmt.Sprintf("%v.playerID = %v.playerID", playerTableName, playerGameTableName)
	where := fmt.Sprintf("type=%d", PlayerType)
	Select := fmt.Sprintf("%v.playerID,%v.winningRate,%v.gameID", playerTableName, playerGameTableName, playerGameTableName)
	session := engine.Table(playerTableName).Join("LEFT", playerGameTableName, idEqu).Where(where).Select(Select).Limit(limit, Startlimit)
	if err := session.Find(&robotsPGs); err != nil {
		sql, _ := session.LastSQL()
		return []*db.TPlayerGame{}, fmt.Errorf("err(%v),sql(%v)", err, sql)
	}
	Startlimit = Startlimit + len(robotsPGs)
	if len(robotsPGs) != limit {
		logrus.Debugf(" maxStartlimit (%d) \n", Startlimit)
		Startlimit = 0
	}
	return robotsPGs, nil
}
