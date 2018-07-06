package handle

import (
	"net/http"
	"steve/room/peipai/utils"
	"strconv"

	"github.com/Sirupsen/logrus"
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
	gameName := r.FormValue(game)
	if len(gameName) == 0 {
		respMSG(w, "配牌失败，需要配牌关键字", 404)
		return
	}
	value := r.FormValue(cards)
	lenValue := r.FormValue(num)
	lens := 0
	if len(lenValue) != 0 {
		lenNum, err := strconv.Atoi(lenValue)
		if err != nil {
			respMSG(w, "墙牌长度应为纯数字", 404)
			return
		}
		if lenNum < 54 || lenNum > 108 {
			respMSG(w, "墙牌长度不能少于54且不能大于108", 404)
			return
		}
		lens = lenNum
	}
	fxValue := r.FormValue("hszfx")
	fx := utils.GetHszFx(fxValue)
	zhuangValue := r.FormValue("zhuang")
	zhuangIndex := -1
	if len(zhuangValue) != 0 {
		index, err := strconv.Atoi(zhuangValue)
		if err != nil {
			respMSG(w, "庄家index应该为纯数字", 404)
			return
		}
		zhuangIndex = index
	}

	var cards []int
	for i := 0; i < len(value); i = i + 3 {
		in := value[i : i+2]
		ca, err := strconv.Atoi(in)
		if err != nil {
			respMSG(w, "配牌失败", 404)
			return
		}
		cards = append(cards, ca)
	}
	if len(cards) > 108 {
		respMSG(w, "配牌越界，您的配牌超过了108张", 404)
		return
	}
	cardsNum := make(map[int]int)
	for _, c1 := range cards {
		num := 0
		if _, ok := cardsNum[c1]; ok {
			continue
		}
		for _, c2 := range cards {
			if c2 == c1 {
				num++
			}
			cardsNum[c1] = num
		}
	}

	for c, num := range cardsNum {
		if c%10 == 0 || c/10 > 3 {
			data := "牌：" + strconv.Itoa(c) + "不存在墙牌中，请检查配牌"
			respMSG(w, data, 404)
			return
		}
		if num > 4 {
			data := "牌：" + strconv.Itoa(c) + "的配牌数量为：" + strconv.Itoa(num) + "，超过了配牌值,请重新配牌"
			respMSG(w, data, 404)
			return
		}
	}
	pp := peipaiInfo{
		Key:    gameName,
		Cards:  value,
		Len:    lens,
		Fx:     fx,
		Zhuang: zhuangIndex,
	}
	addPeiPaiInfo(pp)
	okStr := "ok"
	w.WriteHeader(200)
	w.Write([]byte(okStr))
	logrus.WithFields(logrus.Fields{
		"游戏":    pp.Key,
		"墙牌长度":  pp.Len,
		"配牌":    pp.Cards,
		"换三张方向": pp.Fx,
		"庄家座位号": pp.Zhuang,
	}).Info("配牌成功")
}
