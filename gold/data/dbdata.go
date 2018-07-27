package data

import (
	"fmt"
	"steve/structs"
	"strconv"
)

/*
	功能： 服务数据保存到Mysql.
	作者： SkyWang
	日期： 2018-7-25

*/

var mapID2Name = map[int16]string{}
var mapName2ID = map[string]int16{}
// 如果玩家账号不存在，向DB中加入此玩家初始金币值
var bInitGold = true

// 设置货币类型列表
func SetGoldTypeList(list map[int16]string) {
	mapID2Name = list

	for k, v := range mapID2Name {
		mapName2ID[v] = k
	}
}

// 从DB加载玩家金币
func LoadGoldFromDB(uid uint64) (map[int16]int64, error) {

	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine("steve")
	if err != nil {
		return nil, fmt.Errorf("connect db error")
	}

	strCol := ""
	for _, col := range mapID2Name {
		if len(strCol) > 0 {
			strCol += ","
		}
		strCol += col
	}

	sql := fmt.Sprintf("select %s from t_player_currency  where playerID='%d';", strCol,  uid)
	res, err := engine.QueryString(sql)
	if err != nil {
		if bInitGold {
			return InitGoldToDB(uid)
		}
		return nil, err
	}

	if len(res) != 1 {
		return nil, fmt.Errorf("db result num != 1")
	}
	row := res[0]

	m := make(map[int16]int64)
	for k, v := range row {
		id := mapName2ID[k]
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {

		}
		m[id] = value
	}

	return m, nil
}

// 将玩家金币同步到DB
func SaveGoldToDB(uid uint64, goldList map[int16]int64) error {

	if len(goldList) == 0 {
		return fmt.Errorf("gold list = 0")
	}

	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine("steve")
	if err != nil {
		return fmt.Errorf("connect db error")
	}

	strCol := ""
	for k, v := range goldList {
		if len(strCol) > 0 {
			strCol += ","
		}
		c, ok := mapID2Name[k]
		if !ok {
			return fmt.Errorf("gold type no db col")
		}
		strCol += c
		strCol += "="
		strCol += fmt.Sprintf("'%d'", v)
	}

	sql := fmt.Sprintf("update t_player_currency set %s  where playerID=?;", strCol)
	res, err := engine.Exec(sql, uid)
	if err != nil {
		return err
	}
	if aff, err := res.RowsAffected(); aff == 0 {
		return err
	}
	return nil
}

func InitGoldToDB(uid uint64) (map[int16]int64, error) {

	goldList := make(map[int16]int64)
	goldList[1] = 100000
	goldList[2] = 100000
	goldList[3] = 100000

	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine("steve")
	if err != nil {
		return nil, fmt.Errorf("connect db error")
	}

	strCol := "playerID"
	for k := range goldList {
		if len(strCol) > 0 {
			strCol += ","
		}
		c, ok := mapID2Name[k]
		if !ok {
			return nil,  fmt.Errorf("gold type no db col")
		}
		strCol += c
	}

	strValue := fmt.Sprintf("%d", uid)
	for _, v := range goldList {
		if len(strValue) > 0 {
			strValue += ","
		}
		strValue += fmt.Sprintf("'%d'", v)
	}

	sql := fmt.Sprintf("insert into t_player_currency (%s) values(%s);", strCol, strValue)
	res, err := engine.Exec(sql)
	if err != nil {
		return nil, err
	}
	if aff, err := res.RowsAffected(); aff == 0 {
		return nil, err
	}
	return goldList, nil
}
