package data

import (
	"github.com/Sirupsen/logrus"
)

//UpdatePlayerGotTimesByPlayerID 根据玩家ID修改玩家已经领取数量(db,redis)
func UpdatePlayerGotTimesByPlayerID(playerID uint64, changeTimes int) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":   "getAlmsConfigByPlayerID",
		"changeTimes": changeTimes,
	})
	//db
	if err := updateMysqlPlayerGotTimesByPlayerID(playerID, changeTimes); err != nil {
		entry.WithError(err).Errorln("修改玩家已经领取数量 DB 失败 playerID(%v)", playerID)
		return err
	}
	// redis
	if err := UpdateAlmsPlayerGotTimes(playerID, changeTimes, RedisTimeOut); err != nil {
		entry.WithError(err).Errorln("修改玩家已经领取数量 redis 失败 playerID(%v)", playerID)
		return err
	}
	entry.Infoln("修改玩家已经领取数量")
	return nil
}

//UpdataAlmsConfigVersion 更新版本号
func UpdataAlmsConfigVersion() error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "UpdataAlmsConfigVersion",
	})
	var newVersion int
	versionStr, err := GetAlmsConfigFiled(AlmsVersion)
	if err != nil {
		entry.Errorln("redis 版本号获取失败")
		dbVersion, err := getMysqAlmsConfigVersion()
		if err != nil {
			return err
		}
		newVersion = dbVersion + 1
	} else {
		newVersion = int(InterToint64(versionStr)) + 1
	}
	// 修改redis
	if err := SetAlmsConfigWatch(AlmsVersion, newVersion, RedisTimeOut); err != nil {
		entry.WithError(err).Errorln("redis 救济金配置 Version 改变失败")
		return err
	}
	// 修改DB
	if err := updataMysqAlmsConfigVersion(newVersion); err != nil {
		entry.WithError(err).Errorln("DB 救济金配置 Version 改变失败")
		return err
	}
	return nil
}

//GetAlmsConfigByPlayerID 根据玩家ID获取救济金配置
func GetAlmsConfigByPlayerID(playerID uint64) (*AlmsConfig, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "getAlmsConfigByPlayerID",
	})
	//从 redis 获取救济金配置
	newConfig := []string{AlmsGetNorm, AlmsGetTimes, AlmsGetNumber, AlmsCountDonw, DepositCountDonw, AlmsGameLeveIsOK, AlmsVersion}
	acm, err := GetAlmsConfigFileds(newConfig...)
	ac := &AlmsConfig{}
	if err != nil || !checkMapStringInterface(acm, newConfig) {
		entry.Warnln("警告: 从redis获取救济金配置失败")
		//重新从数据库获取数据，存入redis
		ac, err = GetDBAlmsConfigData()
		if err != nil {
			return nil, err
		}
		// 存储到redis
		if err = SetAlmsConfigWatchs(AlmsConfigToMap(ac), RedisTimeOut); err != nil {
			entry.WithError(err).Errorln("存储救济金配置数据redis失败")
		}

	} else {
		ac = &AlmsConfig{
			GetNorm:             InterToint64(acm[AlmsGetNorm]),
			GetTimes:            int(InterToint64(acm[AlmsGetTimes])),
			GetNumber:           InterToint64(acm[AlmsGetNumber]),
			AlmsCountDonw:       int(InterToint64(acm[AlmsCountDonw])),
			DepositCountDonw:    int(InterToint64(acm[DepositCountDonw])),
			GemeLeveIsOpentAlms: JSONToGameLeveConfig(acm[AlmsGameLeveIsOK].(string)),
			Version:             int(InterToint64(acm[AlmsVersion])),
		}
	}
	// 获取当前玩家救济已领取次数t_hall_info
	times, err := GetPlayerGotTimesByPlayerID(playerID)
	if err != nil {
		return nil, err
	}
	ac.PlayerGotTimes = times
	return ac, nil
}

//GetPlayerGotTimesByPlayerID 根据玩家ID获取救济金已领取数量，先从redis,不存在从db，再存入redis
func GetPlayerGotTimesByPlayerID(playerID uint64) (int, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetPlayerGotTimesByPlayerID",
	})
	// redis获取当前玩家救济领取数量
	times, err := GetAlmsPlayerGotTimes(playerID)
	if err != nil {
		entry.WithError(err).Warnf("警告:从redis获取失败,重新从db获取数据 playerID(%v)", playerID)
		// 从t_hall_info数据库取数据
		times, err := getMysqlPlayerGotTimesByPlayerID(playerID) //t_hall_info 可能不存在该玩家id,
		if err != nil {
			entry.WithError(err).Errorln("获取救济金已领取数量失败")
			return 0, err
		}
		// 存入redis
		if err = UpdateAlmsPlayerGotTimes(playerID, times, RedisTimeOut); err != nil {
			entry.WithError(err).Errorln("存储玩家救济金领取次数失败")
		}
	}
	return times, nil
}

//GetDBAlmsConfigData 获取数据库救济金配置
func GetDBAlmsConfigData() (*AlmsConfig, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetDBAlmsConfigData",
	})
	// 救济倒计时，快充倒计时，救济线，救济领取次数，领取数量
	almsConfigData, err := getMysqAlmsConfigData()
	if err != nil {
		entry.WithError(err).Errorln("获取救济金配置数据失败")
		return nil, err
	}
	tgcs, err := getMysqlGameLevelConfigData()
	if err != nil {
		entry.WithError(err).Errorln("获取游戏场次配置数据失败")
		return nil, err
	}
	glos := make([]*GameLeveConfig, 0, len(tgcs))
	for _, tgc := range tgcs {
		glo := &GameLeveConfig{
			GameID:  int32(tgc.Gameid),
			LevelID: int32(tgc.Levelid),
			IsOpen:  tgc.Isalms,
		}
		glos = append(glos, glo)
	}
	ac := &AlmsConfig{
		GetNorm:             int64(almsConfigData.Getnorm),
		GetTimes:            almsConfigData.Gettimes,
		GetNumber:           int64(almsConfigData.Getnumber),
		AlmsCountDonw:       almsConfigData.Almscountdonw,
		DepositCountDonw:    almsConfigData.Depositcountdonw,
		GemeLeveIsOpentAlms: glos,
		Version:             almsConfigData.Version,
	}
	return ac, nil
}
