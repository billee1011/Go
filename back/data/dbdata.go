package data

import (
	"fmt"
	"steve/entity/db"
	"steve/structs"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/xorm"
)

const dbName = "steve"

// GetTotalBureau 获取玩家总场次
func GetTotalBureau(userID uint64, gameID int) (uint64, error) {
	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(dbName)
	if err != nil {
		return 0, err
	}
	sql := fmt.Sprintf("select totalBureau from t_player_game where userID = '%v' and gameID = '%v';", userID, gameID)
	result, err := engine.QueryString(sql)
	if err != nil {
		return 0, err
	}
	if len(result) != 1 {
		return 0, fmt.Errorf("not only result")
	}
	kv := result[0]
	for _, v := range kv {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {

		}
		return uint64(value), nil
	}
	return 0, fmt.Errorf("no value")
}

// UpdateTotalBureau 更新玩家总局数
func UpdateTotalBureau(userID uint64, gameID int, totalCount uint64) error {
	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(dbName)
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("update t_player_game set totalBureau=%v where userID = '%v' and gameID = '%v';", totalCount, userID, gameID)
	_, err = engine.Exec(sql)
	if err != nil {
		return err
	}
	// _, err := result.RowsAffected()
	return nil
}

// UpdateMaxWinningStreak 更新玩家最高连胜
func UpdateMaxWinningStreak() {

}

// UpdateMaxMultiple 更新玩家获胜最大倍数
func UpdateMaxMultiple() {

}

// GetWinningRate 更新胜率
func GetWinningRate() {
	//读mysql获取当前的胜率，通过胜率×总局数获取赢的局数
	//通过赢的局数+1/总局数+1得到最新的胜率
	//将最新的胜率储存到mysql
}

// GetMaxWinningStreak 更新最高连胜
func GetMaxWinningStreak() {
	//先读redis，redis有连胜记录，并且当前是赢家
	//redis连胜+1，并且储存,如果时输家则redis连胜记录职位置为0
	//读mysql，拿出表记录的最高连胜
	//对比两个胜率，如果redis储存的大于mysql的，update mysql的数据
}

// GetMaxMultiple 最大获胜倍数
func GetMaxMultiple() {
	//读mysql拿到最大倍数
	//将订阅到的最大倍数与数据库的最大倍数进行比对
	//如果大于，更新
}

// GetTPlayerGame 获取t_player_game的信息
func GetTPlayerGame(gameID uint64, playerID uint64) (db.TPlayerGame, error) {
	tpg := db.TPlayerGame{}
	engine, err := MysqlEngineGetter(dbName)
	if err != nil {
		logrus.Errorln(err)
		return tpg, err
	}
	sql := fmt.Sprintf("select * from t_player_game where userID='%v' and gameID='%v'", playerID, gameID)
	result, err := engine.QueryString(sql)
	if err != nil {
		return tpg, err
	}
	if len(result) != 1 {
		return tpg, fmt.Errorf("num of result is not only")
	}
	translationTGP(result[0], &tpg)
	return tpg, nil
}

func translationTGP(kv map[string]string, tpg *db.TPlayerGame) {
	for k, v := range kv {
		switch k {
		case "id":
			tpg.Id, _ = strconv.ParseInt(v, 10, 64)
		case "userID":
			tpg.Userid, _ = strconv.ParseInt(v, 10, 64)
		}
	}
}

// InsertSummary 向db添加Summary信息
func InsertSummary(summary *db.TGameSumary) error {
	engine, err := MysqlEngineGetter(dbName)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	_, err = engine.Insert(summary)
	if err != nil {
		logrus.Errorf("failed to Insert Sunmmary,err:%v", err)
		return err
	}
	return nil
}

// InsertDetail 向db库添加detail信息
func InsertDetail(detail *db.TGameDetail) error {
	engine, err := MysqlEngineGetter(dbName)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	_, err = engine.Insert(detail)
	if err != nil {
		logrus.Errorf("failed to Insert detail,err:%v", err)
		return err
	}
	return nil
}

func getMysqlEngine(mysqlName string) (*xorm.Engine, error) {
	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(mysqlName)
	if err != nil {
		return nil, fmt.Errorf("获取 mysql 引擎失败：%v", err)
	}
	return engine, nil
}

// MysqlEngineGetter 单元测试通过这两个值修改 mysql 引擎获取和 redis cli 获取
var MysqlEngineGetter = getMysqlEngine
