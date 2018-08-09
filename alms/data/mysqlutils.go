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
	err = engine.Table(almsConfigTableName).Select("almsCountDonw,depositCountDonw,getNorm,getTimes,getNumber,version").Find(&almsConfigs)
	if err != nil {
		return nil, err
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
	num, err := engine.Table(almsConfigTableName).Select("version").Update(ac)
	if err != nil {
		return err
	}
	if num == 0 {
		return fmt.Errorf("更新救济金配置版本号失败 : %v", version)
	}
	return nil
}

// 获取救济金配置版本号
func getMysqAlmsConfigVersion() (int, error) {
	engine, err := MysqlEnginefunc(MysqlConfigdbName)
	if err != nil {
		return 0, err
	}
	ac := &db.TAlmsConfig{}
	exist, err := engine.Table(almsConfigTableName).Select("version").Get(ac)
	if err != nil {
		return 0, err
	}
	if !exist { // 表数据不存在插入新的
		ac.Version = 1
		_, err := engine.Table(almsConfigTableName).Insert(ac)
		return 1, err
	}
	return ac.Version, nil
}

// 获取游戏场次配置数据
func getMysqlGameLevelConfigData() ([]*db.TGameLevelConfig, error) {
	engine, err := MysqlEnginefunc(MysqlConfigdbName)
	if err != nil {
		return nil, err
	}
	glcs := []*db.TGameLevelConfig{}
	err = engine.Table(gameLevelConfigTableName).Select("gameID,levelID,isAlms").Find(&glcs)
	if err != nil {
		return nil, err
	}
	return glcs, nil
}

// 获取玩家救济金已领取次数
func getMysqlPlayerGotTimesByPlayerID(playerID uint64) (int, error) {
	engine, err := MysqlEnginefunc(MysqlPlayerdbName)
	if err != nil {
		return 0, err
	}
	hi := &db.THallInfo{}
	exist, err := engine.Table(hallInfoTableName).Select("almsGotTimes").Where(fmt.Sprintf("playerID=%v", playerID)).Get(hi)
	if err != nil {
		return 0, err
	}
	if !exist { // 不存在插入新的
		hi.Playerid = int64(playerID)
		_, err := engine.Table(hallInfoTableName).Insert(hi)
		return 0, err
	}
	return hi.Almsgottimes, nil
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
	num, err := engine.Table(hallInfoTableName).Select("almsGotTimes").Where(fmt.Sprintf("playerID=%v", playerID)).Update(hi)
	if err != nil {
		return err
	}
	if num == 0 {
		return fmt.Errorf("修改玩家救济金已领取次数失败 : %v - %v", playerID, gotTimes)
	}
	return nil
}
