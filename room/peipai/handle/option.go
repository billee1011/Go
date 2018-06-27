package handle

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

var optionInfos []optionInfo

type optionInfo struct {
	Key   string
	Hsz   bool
	Gold  uint64
	Golds map[uint64]uint64
}

//addOptionInfo 添加新的选项请求
func addOptionInfo(opt optionInfo) {
	index, exist := checkOptionInfo(opt)
	if exist {
		optionInfos = append(optionInfos[:index], optionInfos[index+1:]...)
	}
	optionInfos = append(optionInfos, opt)
}

//checkOptionInfo 检查新的选项请求是否存在
func checkOptionInfo(opt optionInfo) (int, bool) {
	for index, optionInfo := range optionInfos {
		if optionInfo.Key == opt.Key {
			return index, true
		}
	}
	return 0, false
}

//GetHsz 获取换三张开关
func GetHsz(gameID int) bool {
	key := idIntToStr(gameID)
	for _, opt := range optionInfos {
		if opt.Key == key {
			return opt.Hsz
		}
	}
	return true
}

//GetGolds 获取玩家的金币数
func GetGolds(gameID int) map[uint64]uint64 {
	key := idIntToStr(gameID)
	for _, opt := range optionInfos {
		if opt.Key == key {
			return opt.Golds
		}
	}
	return nil
}

//GetGold 获取全局金币数
func GetGold(gameID int) uint64 {
	key := idIntToStr(gameID)
	for _, opt := range optionInfos {
		if opt.Key == key {
			return opt.Gold
		}
	}
	return 0
}

//Option 处理选项请求
func Option(resp http.ResponseWriter, req *http.Request) {
	opt := optionInfo{}
	gameName := req.FormValue(game)
	if len(gameName) == 0 {
		respMSG(resp, fmt.Sprintf("缺少游戏ID"), 404)
		return
	}
	opt.Key = gameName
	value := req.FormValue(HszSwitch)
	if len(value) != 0 {
		open, err := strconv.ParseBool(value)
		if err != nil {
			respMSG(resp, fmt.Sprintf("switch对应的值有误:%v", err), 404)
			return
		}
		opt.Hsz = open
		respMSG(resp, fmt.Sprintf("配置换三张开关成功,当前为:%v", opt.Hsz), 200)
	}
	goldValue := req.FormValue(gold)
	if len(goldValue) != 0 {
		goldSum, err := strconv.ParseUint(goldValue, 10, 0)
		if err != nil {
			respMSG(resp, fmt.Sprintf("gold对应的值有误:%v", err), 404)
			return
		}
		opt.Gold = goldSum
		respMSG(resp, fmt.Sprintf("配置金币成功,当前为:%v", opt.Gold), 200)
	}
	goldsValue := req.FormValue(golds)
	if len(goldsValue) != 0 {
		strs := strings.Split(goldsValue, ",")
		mp := make(map[uint64]uint64, 0)
		for _, str := range strs {
			kv := strings.Split(str, "-")
			if len(kv) != 2 {
				respMSG(resp, fmt.Sprintf("玩家与金币的键值对有误:%v", kv), 404)
				return
			}
			id, err := strconv.ParseUint(kv[0], 10, 0)
			if err != nil {
				respMSG(resp, fmt.Sprintf("玩家id应该为纯数字:%v", id), 404)
				return
			}
			gold, err := strconv.ParseUint(kv[1], 10, 0)
			if err != nil {
				respMSG(resp, fmt.Sprintf("玩家金币应该为纯数字:%v", gold), 404)
				return
			}
			mp[id] = gold
		}
		opt.Golds = mp
	}
	addOptionInfo(opt)
	logrus.WithFields(logrus.Fields{
		"游戏":    opt.Key,
		"全局金币":  opt.Gold,
		"玩家金币":  opt.Golds,
		"换三张开关": opt.Hsz,
	}).Info("选项配置成功")
}
