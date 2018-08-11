package data

import (
	"encoding/json"
	"strconv"

	"github.com/Sirupsen/logrus"
)

//InterToint64 接口转int64
func InterToint64(param interface{}) int64 {
	if param == nil {
		return 0
	}
	str := param.(string)
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func_name": "InterToint64",
			"param": param}).Infoln("InterToint64失败")
		return 0
	}
	return result
}

//检验redis 返回的 值
func checkMapStringInterface(m map[string]interface{}, checkString []string) bool {
	if len(m) != len(checkString) {
		return false
	}
	for _, str := range checkString {
		switch m[str].(type) {
		case string:
			if str == GameLeveConfigs && len(JSONToGameLeveConfig(m[str].(string))) <= 0 {
				return false
			}
		case int64:
			if InterToint64(m[str]) <= 0 {
				return false
			}
		default:
			logrus.WithFields(logrus.Fields{"func_name": "checkMapStringInterface",
				"m[str]": m[str]}).Infoln("检验redis 返回的")
			return false
		}
	}
	return true
}

// AlmsConfigToMap AlmsConfig to map[string]interface{}
func AlmsConfigToMap(ac *AlmsConfig) map[string]interface{} {
	almsConfigMap := make(map[string]interface{})
	if ac.GetNorm > 0 {
		almsConfigMap[AlmsGetNorm] = ac.GetNorm // 救济线
	}
	if ac.GetTimes > 0 {
		almsConfigMap[AlmsGetTimes] = ac.GetTimes // 最多领取次数
	}
	if ac.GetNumber > 0 {
		almsConfigMap[AlmsGetNumber] = ac.GetNumber // 领取数量
	}
	if ac.AlmsCountDonw > 0 {
		almsConfigMap[AlmsCountDonw] = ac.AlmsCountDonw // 救济倒计时
	}
	if ac.DepositCountDonw > 0 {
		almsConfigMap[DepositCountDonw] = ac.DepositCountDonw // 快冲倒计时
	}
	if len(ac.GameLeveConfigs) > 0 {
		almsConfigMap[GameLeveConfigs] = GameLeveConfigToJSON(ac.GameLeveConfigs) //游戏场次是否开启救济金
	}
	if ac.Version > 0 {
		almsConfigMap[AlmsVersion] = ac.Version // 版本号
	}
	return almsConfigMap
}

// GameLeveConfigToJSON 游戏场次配置 转 JSON
func GameLeveConfigToJSON(gemeLeveOK []*GameLeveConfig) string {
	if gemeLeveOK == nil {
		return ""
	}
	str, err := json.Marshal(gemeLeveOK)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func_name": "GameLeveConfigToJSON",
			"gemeLeveOK": gemeLeveOK}).Infoln("游戏场次配置 转 JSON失败")
	}
	return string(str)
}

// JSONToGameLeveConfig JSON 转 游戏场次配置
func JSONToGameLeveConfig(gemeLeveOKJSON string) []*GameLeveConfig {
	gemeLeveOK := []*GameLeveConfig{}
	if gemeLeveOKJSON == "" {
		return gemeLeveOK
	}
	globyte := []byte(gemeLeveOKJSON)
	if err := json.Unmarshal(globyte, &gemeLeveOK); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"func_name": "JSONToGameLeveConfig",
			"gemeLeveOKJSON": gemeLeveOKJSON}).Infoln("JSON 转 游戏场次配置失败")
	}
	return gemeLeveOK
}
