package util

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"strconv"
	"steve/room/peipai/utils"
	"steve/gutils"
	"fmt"
)

const (
	// HszSwitch 换三张开关关键字
	HszSwitch = "hszswitch"
	// PlayerIDKey 玩家ID
	PlayerIDKey = "player_id"
	// GoldKey 金币
	GoldKey = "gold"
	// game 游戏关键字
	game = "game"
	// cards 墙牌
	cards = "cards"
	// len 墙牌数量
	num = "len"

	//SCXL 四川血流
	SCXL = "scxl"
	//SCXZ 四川血战
	SCXZ = "scxz"
	//DDZ 斗地主
	DDZ = "doudizhu"
	//ERMJ 二人麻将
	ERMJ = "ermj"
)

var peiPaiInfos []peipaiInfo

// var PP_PORT string
type peipaiInfo struct {
	//配牌关键字
	Key string
	//配牌
	Cards string
	//墙牌长度
	Len int
	//换三张方向
	Fx int
	//庄的index
	Zhuang int
}

//addPeiPaiInfo 添加新的配牌请求
func addPeiPaiInfo(pp peipaiInfo) {
	index, exist := checkPeiPaiInfo(pp)

	// 若已经存在，则先删除，再替换
	if exist {
		peiPaiInfos = append(peiPaiInfos[:index], peiPaiInfos[index+1:]...)
	}

	peiPaiInfos = append(peiPaiInfos, pp)
}

//checkPeiPaiInfo 检查新的配牌请求是否存在
func checkPeiPaiInfo(pp peipaiInfo) (int, bool) {
	for index, peipaiInfo := range peiPaiInfos {
		if peipaiInfo.Key == pp.Key {
			return index, true
		}
	}
	return 0, false
}

//GetPeiPai 通过配牌关键字拿到配牌
func GetPeiPai(gameID int) string {
	// 根据游戏ID获取游戏名字
	key := idIntToStr(gameID)

	for _, pp := range peiPaiInfos {
		if pp.Key == key {
			logrus.WithField("peipai", pp.Cards).Debugln("获取配牌")
			return pp.Cards
		}
	}
	logrus.WithField("key", key).Debugln("获取配牌为空")
	return ""
}

//GetLensOfWallCards 牌墙长度
func GetLensOfWallCards(gameID int) int {
	key := idIntToStr(gameID)
	for _, pp := range peiPaiInfos {
		if pp.Key == key {
			return pp.Len
		}
	}
	return 0
}

//GetHSZFangXiang 换三张方向
func GetHSZFangXiang(gameID int) int {
	key := idIntToStr(gameID)
	for _, pp := range peiPaiInfos {
		if pp.Key == key {
			return pp.Fx
		}
	}
	return -1
}

//GetZhuangIndex 定义庄家的index
func GetZhuangIndex(gameID int) int {
	key := idIntToStr(gameID)
	for _, pp := range peiPaiInfos {
		if pp.Key == key {
			return pp.Zhuang
		}
	}
	return -1
}

//Peipai 接受http请求并处理
func Peipai(w http.ResponseWriter, r *http.Request) {

	// 游戏名字
	gameName := r.FormValue(game)
	if len(gameName) == 0 {
		respMSG(w, "配牌失败，需要配牌关键字", 404)
		return
	}

	// 配牌
	value := r.FormValue(cards)

	// 配牌的长度-字符串
	lenValue := r.FormValue(num)
	lens := 0
	if len(lenValue) != 0 {
		//配牌的长度 -数字
		lenNum, err := strconv.Atoi(lenValue)
		if err != nil {
			respMSG(w, "墙牌长度应为纯数字", 404)
			return
		}
		// if lenNum < 54 {
		// 	respMSG(w, "墙牌长度不能少于54", 404)
		// 	return
		// }
		lens = lenNum
	}

	// 换三张方向-字符串
	fxValue := r.FormValue("hszfx")

	// 换三张方向-数字
	fx := utils.GetHszFx(fxValue)

	// 庄家玩家的座位号-字符串
	zhuangValue := r.FormValue("zhuang")
	zhuangIndex := -1
	if len(zhuangValue) != 0 {
		// 庄家玩家的座位号-数字
		index, err := strconv.Atoi(zhuangValue)
		if err != nil {
			respMSG(w, "庄家index应该为纯数字", 404)
			return
		}
		zhuangIndex = index
	}

	var cards []int

	// i + 3 是因为数字占2位，逗号占1位
	for i := 0; i < len(value); i = i + 3 {

		// 每两位一个数字
		in := value[i : i+2]

		// 转换成数字
		ca, err := strconv.Atoi(in)
		if err != nil {
			respMSG(w, "配牌失败", 404)
			return
		}

		cards = append(cards, ca)
	}
	// if len(cards) > 108 {
	// 	respMSG(w, "配牌越界，您的配牌超过了108张", 404)
	// 	return
	// }
	cardsNum := make(map[int]int)

	for _, c1 := range cards {
		num := 0

		// c1牌已存在的，跳过
		if _, ok := cardsNum[c1]; ok {
			continue
		}

		// 统计c1牌的个数
		for _, c2 := range cards {
			if c2 == c1 {
				num++
			}
			cardsNum[c1] = num
		}
	}

	// 检测牌的正确性
	for c, num := range cardsNum {

		// 斗地主
		if gameName == "doudizhu" {
			// 牌的范围在14,15,17-29,33-45,49-61,65-77之间
			if !((c == 14) || (c == 15) || (c >= 17 && c <= 29) || (c >= 33 && c <= 45) || (c >= 49 && c <= 61) || (c >= 65 && c <= 77)) {
				data := "检测到不应该存在的牌：" + strconv.Itoa(c) + "，请检查配牌"
				respMSG(w, data, 404)
				return
			}
		} else {
			// 不能是10的整数倍
			// 不能超过40
			if c%10 == 0 || c/10 > 5 {
				data := "牌：" + strconv.Itoa(c) + "不存在墙牌中，请检查配牌"
				respMSG(w, data, 404)
				return
			}
		}

		// 同一类型的牌不能超过4个
		if num > 4 {
			data := "牌：" + strconv.Itoa(c) + "的配牌数量为：" + strconv.Itoa(num) + "，超过了配牌值,请重新配牌"
			respMSG(w, data, 404)
			return
		}
	}

	// 记录该配牌信息
	pp := peipaiInfo{
		Key:    gameName,
		Cards:  value,
		Len:    lens,
		Fx:     fx,
		Zhuang: zhuangIndex,
	}

	// 添加进来
	addPeiPaiInfo(pp)

	// 回复客户端
	okStr := "ok"
	w.WriteHeader(200)
	w.Write([]byte(okStr))

	// 日志
	logrus.WithFields(logrus.Fields{
		"游戏":    pp.Key,
		"墙牌长度":  pp.Len,
		"配牌":    pp.Cards,
		"换三张方向": pp.Fx,
		"庄家座位号": pp.Zhuang,
	}).Info("配牌成功")
}


func respMSG(resp http.ResponseWriter, message string, code int) {
	resp.WriteHeader(code)
	resp.Write([]byte(message))
	switch code {
	case 200:
		logrus.Infoln(message)
	default:
		logrus.Debugln(message)
	}
}

//LogPeiPaiInfos 打印配牌信息
func LogPeiPaiInfos() {
	for k, info := range peiPaiInfos {
		logrus.WithFields(logrus.Fields{
			"game":   info.Key,
			"cards":  info.Cards,
			"len":    info.Len,
			"fx":     info.Fx,
			"zhuang": info.Zhuang,
		}).Info(k)
	}
}

//LogOptionInfos 打印选项信息
func LogOptionInfos() {
	for k, info := range optionInfos {
		logrus.WithFields(logrus.Fields{
			"game": info.Key,
			"fx":   info.Hsz,
		}).Info(k)
	}
}

// idIntToStr 根据游戏ID得到游戏名字
func idIntToStr(gameID int) string {
	switch gameID {
	case gutils.SCXLGameID:
		return SCXL
	case gutils.SCXZGameID:
		return SCXZ
	case gutils.DDZGameID:
		return DDZ
	case gutils.ERMJGameID:
		return ERMJ
	}
	return ""
}


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