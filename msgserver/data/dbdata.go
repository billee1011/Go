package data

import (
	"encoding/json"
	"fmt"
	"steve/msgserver/define"
	"steve/structs"
	"strconv"

	"github.com/Sirupsen/logrus"
	"steve/external/configclient"
)

/*
	功能： 服务数据保存到Mysql.
	作者： SkyWang
	日期： 2018-7-25

*/

const dbName = "config"

// 从Config加载跑马灯
func LoadHorseFromConfig() (map[int64]*define.HorseRace, error) {
	strJson, err := configclient.GetConfig("horse", "config")
	if err != nil {
		logrus.Errorf("LoadHorseFromConfig  err:", err)
		return nil,  err
	}

	if len(strJson) == 0 {
		logrus.Errorf("parseHorseConfig config is empty err")
		return nil,  fmt.Errorf("parseHorseConfig config is empty err")
	}

	hList := parseHorseConfig(strJson)
	if hList == nil {
		logrus.Errorf("parseHorseConfig parse err:")
		return nil,  fmt.Errorf("parseHorseConfig parse err")
	}

	return hList, nil
}

// 解析跑马灯 config json
func parseHorseConfig(strJson string) map[int64]*define.HorseRace {

	jsonList := make([]*define.HorseRaceJson, 0, 5)
	err := json.Unmarshal([]byte(strJson), &jsonList)
	if err != nil {
		return nil
	}
	list := make(map[int64]*define.HorseRace)

	for _, jsonObject := range jsonList {

		horse := new(define.HorseRace)
		horse.Id = jsonObject.Id
		horse.IsUse = jsonObject.IsUse
		horse.Prov = jsonObject.Prov
		horse.Channel = jsonObject.Channel
		horse.City = jsonObject.City
		horse.IsUseParent = jsonObject.IsUseParent

		horse.TickTime = jsonObject.TickTime
		horse.SleepTime = jsonObject.SleepTime
		horse.LastUpdateTime = jsonObject.LastUpdateTime

		for _, v := range jsonObject.Horse {
			if v.IsUse == 0 {
				continue
			}
			hc := new(define.HorseContent)
			hc.PlayType = v.PlayType
			hc.BeginDate = v.BeginDate
			hc.EndDate = v.EndDate
			hc.BeginTime = v.BeginTime
			hc.EndTime = v.EndTime
			hc.Content = v.Content
			hc.IsUse = v.IsUse


			hc.WeekDate = make(map[int8]bool, len(v.WeekDate))
			for _, t := range v.WeekDate {
				hc.WeekDate[t] = true
			}
			horse.Content = append(horse.Content, hc)

			logrus.Debugf("LoadHorseFromDB add content: %v", *hc)
		}
		list[horse.Id] = horse
	}

	logrus.Debugf("parseHorseConfig win:sum=%d", len(list))

	return list
}


// 从DB加载跑马灯
func LoadHorseFromDB() (map[int64]*define.HorseRace, error) {

	return LoadHorseFromConfig()

	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(dbName)
	if err != nil {
		logrus.Errorf("LoadHorseFromDB err1:%v", err)
		return nil, fmt.Errorf("connect db error")
	}

	sql := fmt.Sprintf("select n_id, n_channel, n_prov, n_city, n_bUse, n_bUseParent, n_horseData from t_horse_race ;")
	res, err := engine.QueryString(sql)
	if err != nil {
		logrus.Errorf("LoadHorseFromDB err2:%v", err)
		return nil, err
	}
	list := make(map[int64]*define.HorseRace)
	for _, row := range res {

		id, _ := strconv.ParseInt(row["n_id"], 10, 64)
		if id == 0 {
			continue
		}

		horse := new(define.HorseRace)
		horse.Id = id

		horse.Channel, _ = strconv.ParseInt(row["n_channel"], 10, 64)
		horse.Prov, _ = strconv.ParseInt(row["n_prov"], 10, 64)
		horse.City, _ = strconv.ParseInt(row["n_city"], 10, 64)
		isUse, _ := strconv.ParseInt(row["n_bUse"], 10, 8)
		isUseParent, _ := strconv.ParseInt(row["n_bUseParent"], 10, 8)
		horse.IsUse = int8(isUse)
		horse.IsUseParent = int8(isUseParent)
		parseHorseJson(horse, row["n_horseData"])
		list[id] = horse
		logrus.Debugf("LoadHorseFromDB add one: %v", *horse)
	}
	logrus.Debugf("LoadHorseFromDB win:sum=%d", len(list))


	return list, nil
}

// 解析跑马灯json
func parseHorseJson(horse *define.HorseRace, strJson string) bool {

	jsonObject := &define.HorseRaceJson{}
	err := json.Unmarshal([]byte(strJson), jsonObject)
	if err != nil {
		return false
	}

	horse.TickTime = jsonObject.TickTime
	horse.SleepTime = jsonObject.SleepTime
	horse.LastUpdateTime = jsonObject.LastUpdateTime

	for _, v := range jsonObject.Horse {
		hc := new(define.HorseContent)
		hc.PlayType = v.PlayType
		hc.BeginDate = v.BeginDate
		hc.EndDate = v.EndDate
		hc.BeginTime = v.BeginTime
		hc.EndTime = v.EndTime
		hc.Content = v.Content

		hc.WeekDate = make(map[int8]bool, len(v.WeekDate))
		for _, t := range v.WeekDate {
			hc.WeekDate[t] = true
		}
		horse.Content = append(horse.Content, hc)
		logrus.Debugf("LoadHorseFromDB add content: %v", *hc)
	}

	return true
}

func MarshalHorseJson(horse *define.HorseRaceJson) (string, error) {
	data, err := json.Marshal(horse)
	return string(data), err
}
