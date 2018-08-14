package data

import (
	"fmt"
	"steve/structs"
	"strconv"

	"time"
	"github.com/Sirupsen/logrus"
)

/*
	功能： 服务数据保存到Mysql.
	作者： SkyWang
	日期： 2018-7-25

CREATE TABLE `t_player_props` (
  `playerID` bigint(20) NOT NULL COMMENT '玩家ID',
  `propID` bigint(20) NOT NULL COMMENT '道具ID',
  `count` bigint(20) NOT NULL COMMENT '道具数量',
  `createTime` datetime DEFAULT NULL COMMENT '创建时间',
  `createBy` varchar(100) DEFAULT NULL COMMENT '创建人',
  `updateTime` datetime DEFAULT NULL COMMENT '更新时间',
  `updateBy` varchar(100) DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`playerID`,`propID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='玩家道具表'


*/

const dbName = "player"

const dbLogName = "log"

// 从DB加载玩家道具
func LoadPropsFromDB(uid uint64) (map[uint64]int64, error) {

	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(dbName)
	if err != nil {
		return nil, fmt.Errorf("connect db error")
	}

	sql := fmt.Sprintf("select propID,count  from t_player_props  where playerID='%d';", uid)
	res, err := engine.QueryString(sql)
	if err != nil {
		return nil, err
	}

	m := make(map[uint64]int64)
	for _, row := range res {

		id, err := strconv.ParseUint(row["propID"], 10, 64)
		if err != nil {
			continue
		}
		value, err := strconv.ParseInt(row["count"], 10, 64)
		if err != nil {

		}
		m[id] = value
	}
	logrus.Debugf("LoadPropsFromDB win: uid=%d, propslist=%v", uid, m)
	return m, nil
}

// 将玩家道具同步到DB
func SavePropsToDB(uid uint64, propId uint64, propNum int64) error {

	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(dbName)
	if err != nil {
		return fmt.Errorf("connect db error")
	}
	sql := ""

	curDate := time.Now().Format("2006-01-02 15:04:05")


	sql = fmt.Sprintf("update  t_player_props set count=%d  where playerID='%d' and propID ='%d' ;", propNum,  uid, propId)

	res, err := engine.Exec(sql)
	if err != nil {
		return err
	}

	if aff, _ := res.RowsAffected(); aff == 0 {
		// 修改失败，再进行插入
		sql = fmt.Sprintf("insert into t_player_props (playerID, propID, count, createTime, updateTime) values('%d','%d','%d','%s','%s');",
			uid, propId, propNum, curDate, curDate)
		res, err = engine.Exec(sql)
		if err != nil {
			return err
		}

		if aff, err := res.RowsAffected(); aff == 0 {
			// 如果插入行=0，表明插入失败
			return err
		}
	}

	return nil
}
