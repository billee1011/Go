package handle

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
)

// 所有游戏的选项配置
var optionInfos []optionInfo

// 单个游戏的选项
type optionInfo struct {
	Key string // 游戏名字
	Hsz bool   // 是否开启换三张
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

//Option 处理选项请求
func Option(resp http.ResponseWriter, req *http.Request) {
	opt := optionInfo{}

	// 游戏名字
	gameName := req.FormValue(game)
	if len(gameName) == 0 {
		respMSG(resp, fmt.Sprintf("缺少游戏ID"), 404)
		return
	}

	opt.Key = gameName

	// 换三张选项
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

	// 添加进来
	addOptionInfo(opt)

	logrus.WithFields(logrus.Fields{
		"游戏":    opt.Key,
		"换三张开关": opt.Hsz,
	}).Info("选项配置成功")
}
