package data

import (
	"fmt"
	"steve/entity/db"
	"steve/structs"

	"github.com/go-xorm/xorm"
)

const (
	// MysqlConfigdbName 数据库名
	MysqlConfigdbName = "config"
	// MysqlPlayerdbName 数据库名
	MysqlPlayerdbName        = "player"
	hallInfoTableName        = "t_hall_info"         // 大厅信息表
	almsConfigTableName      = "t_alms_config"       // 救济金配置表
	playerTableName          = "t_player"            // 玩家表
	gameLevelConfigTableName = "t_game_level_config" //游戏场次配置表
)

//MysqlEnginefunc 单元测试需要
var MysqlEnginefunc = getMysqlEngineByName

func getMysqlEngineByName(mysqlName string) (*xorm.Engine, error) {
	e := structs.GetGlobalExposer()
	engine, err := e.MysqlEngineMgr.GetEngine(mysqlName)
	if err != nil {
		return nil, fmt.Errorf("获取 mysql 引擎失败：%v", err)
	}
	if err := engine.Ping(); err != nil {
		return nil, fmt.Errorf("engine ping 失败：%v", err)
	}
	return engine, nil
}

// 获取救济金配置数据
func getMysqAlmsConfigData() (*db.TAlmsConfig, error) {
	engine, err := MysqlEnginefunc(MysqlConfigdbName)
	if err != nil {
		return nil, err
	}
	almsConfigs := []*db.TAlmsConfig{}
	session := engine.Table(almsConfigTableName).Select("almsCountDonw,depositCountDonw,getNorm,getTimes,getNumber,version")
	if err := session.Find(&almsConfigs); err != nil {
		sql, _ := session.LastSQL()
		return nil, fmt.Errorf("从数据库获取数据失败：(%v), sql:(%s)", err, sql)
	}
	if len(almsConfigs) == 0 {
		return nil, fmt.Errorf("获取不到救济金配置")
	}
	return almsConfigs[0], nil
}

// 更新救济金配置版本号
func updataMysqAlmsConfigVersion(version int) error {
	engine, err := MysqlEnginefunc(MysqlConfigdbName)
	if err != nil {
		return err
	}
	ac := &db.TAlmsConfig{
		Version: version,
	}
	session := engine.Table(almsConfigTableName).Select("version")
	num, err := session.Update(ac)
	if err != nil {
		sql, _ := session.LastSQL()
		return fmt.Errorf("从数据库获取数据失败：(%v), sql:(%s)", err, sql)
	}
	if num == 0 {
		return fmt.Errorf("更新救济金配置版本号失败 version:(%v)", version)
	}
	return nil
}

// 获取救济金配置版本号
func getMysqAlmsConfigVersion() (int, error) {
	engine, err := MysqlEnginefunc(MysqlConfigdbName)
	if err != nil {
		return 0, err
	}
	var version int
	session := engine.Table(almsConfigTableName).Select("version")
	exist, err := session.Get(&version)
	if err != nil {
		sql, _ := session.LastSQL()
		return 0, fmt.Errorf("从数据库获取数据失败：(%v), sql:(%s)", err, sql)
	}
	if !exist { // 表数据不存在插入新的
		ac := db.TAlmsConfig{
			Version: version,
		}
		_, err := engine.Table(almsConfigTableName).Insert(&ac)
		return 1, err
	}
	return version, nil
}

// 获取游戏场次配置数据
func getMysqlGameLevelConfigData() ([]*db.TGameLevelConfig, error) {
	engine, err := MysqlEnginefunc(MysqlConfigdbName)
	if err != nil {
		return nil, err
	}
	glcs := []*db.TGameLevelConfig{}
	session := engine.Table(gameLevelConfigTableName).Select("gameID,levelID,lowScores,isAlms")
	if err := session.Find(&glcs); err != nil {
		sql, _ := session.LastSQL()
		return nil, fmt.Errorf("从数据库获取数据失败：(%v), sql:(%s)", err, sql)
	}
	return glcs, nil
}

// 获取玩家救济金已领取次数
func getMysqlPlayerGotTimesByPlayerID(playerID uint64) (int, error) {
	engine, err := MysqlEnginefunc(MysqlPlayerdbName)
	if err != nil {
		return 0, err
	}
	var almsGotTimes int
	session := engine.Table(hallInfoTableName).Select("almsGotTimes").Where(fmt.Sprintf("playerID=%v", playerID))
	exist, err := session.Get(&almsGotTimes)
	if err != nil {
		sql, _ := session.LastSQL()
		return 0, fmt.Errorf("从数据库获取数据失败：(%v), sql:(%s)", err, sql)
	}
	if !exist { // 不存在插入新的
		hi := &db.THallInfo{}
		hi.Playerid = int64(playerID)
		_, err := engine.Table(hallInfoTableName).Insert(hi)
		return 0, err
	}
	return almsGotTimes, nil
}

// 更改玩家救济金已领取次数
func updateMysqlPlayerGotTimesByPlayerID(playerID uint64, gotTimes int) error {
	engine, err := MysqlEnginefunc(MysqlPlayerdbName)
	if err != nil {
		return err
	}
	hi := &db.THallInfo{
		Almsgottimes: gotTimes,
	}
	session := engine.Table(hallInfoTableName).Select("almsGotTimes").Where(fmt.Sprintf("playerID=%v", playerID))
	num, err := session.Update(hi)
	if err != nil {
		sql, _ := session.LastSQL()
		return fmt.Errorf("从数据库获取数据失败：(%v), sql:(%s)", err, sql)
	}
	if num == 0 {
		return fmt.Errorf("修改玩家救济金已领取次数失败 : %v - %v", playerID, gotTimes)
	}
	return nil
}
