package peipai

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	// "github.com/go-redis/redis"
	// "github.com/spf13/viper"
)

//Demo 接受http请求并处理
func Demo(w http.ResponseWriter, r *http.Request) {
	game := r.FormValue("game")
	if len(game) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("配牌失败，需要配牌关键字"))
		return
	}
	value := r.FormValue("cards")
	lenValue := r.FormValue("len")
	lens := 0
	if len(lenValue) != 0 {
		lenNum, err := strconv.Atoi(lenValue)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("墙牌长度应为纯数字"))
			return
		}
		if lenNum < 54 || lenNum > 108 {
			w.WriteHeader(404)
			w.Write([]byte("墙牌长度不能少于54且不能大于108"))
			return
		}
		lens = lenNum
	}
	fxValue := r.FormValue("hszfx")
	fx := -1
	if fxValue == "dui" {
		fx = 1
	} else if fxValue == "shun" {
		fx = 0
	} else if fxValue == "ni" {
		fx = 2
	}

	zhuangValue := r.FormValue("zhuang")
	zhuangIndex := -1
	if len(zhuangValue) != 0 {
		index, err := strconv.Atoi(zhuangValue)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("庄家index应该为纯数字"))
			return
		}
		zhuangIndex = index
	}

	var cards []int
	for i := 0; i < len(value); i = i + 3 {
		in := value[i : i+2]
		ca, err := strconv.Atoi(in)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("配牌失败"))
			return
		}
		cards = append(cards, ca)
	}
	if len(cards) > 108 {
		w.WriteHeader(404)
		w.Write([]byte("配牌越界，您的配牌超过了108张"))
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
			w.WriteHeader(404)
			data := "牌：" + strconv.Itoa(c) + "不存在墙牌中，请检查配牌"
			w.Write([]byte(data))
			return
		}
		if num > 4 {
			w.WriteHeader(404)
			data := "牌：" + strconv.Itoa(c) + "的配牌数量为：" + strconv.Itoa(num) + "，超过了配牌值,请重新配牌"
			w.Write([]byte(data))
		}
	}
	// peipaiValue[game] = value
	pp := peipaiInfo{
		Key:    game,
		Cards:  value,
		Len:    lens,
		Fx:     fx,
		Zhuang: zhuangIndex,
	}
	addPeiPaiInfo(pp)
	// for k, info := range peiPaiInfos {
	// 	logrus.WithFields(logrus.Fields{
	// 		strconv.Itoa(k): logrus.Fields{
	// 			"game":   info.Key,
	// 			"cards":  info.Cards,
	// 			"len":    info.Len,
	// 			"fx":     info.Fx,
	// 			"zhuang": info.Zhuang,
	// 		},
	// 	}).Info("所有的配牌信息")
	// }

	okStr := "ok"
	w.WriteHeader(200)
	w.Write([]byte(okStr))
}

var peipaiValue map[string]string

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

//GetPeiPai 通过配牌关键字拿到配牌
func GetPeiPai(key string) (string, error) {
	for _, pp := range peiPaiInfos {
		if pp.Key == key {
			return pp.Cards, nil
		}
	}
	return "", fmt.Errorf("不存在这个配牌key")
}

//GetLensOfWallCards 牌墙长度
func GetLensOfWallCards(key string) int {
	for _, pp := range peiPaiInfos {
		if pp.Key == key {
			return pp.Len
		}
	}
	return 0
}

//GetHSZFangXiang 换三张方向
func GetHSZFangXiang(key string) int {
	for _, pp := range peiPaiInfos {
		if pp.Key == key {
			return pp.Fx
		}
	}
	return -1
}

//GetZhuangIndex 定义庄家的index
func GetZhuangIndex(key string) int {
	for _, pp := range peiPaiInfos {
		if pp.Key == key {
			return pp.Zhuang
		}
	}
	return -1
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

//checkPeiPaiKey 检测配牌key
func checkPeiPaiKey(key string) {

}

//checkPeiPaiValue 检测配牌value
func checkPeiPaiValue(value string) {

}

//checkPeiPaiLen 检测配牌的墙牌长度
func checkPeiPaiLen(lenValue string) {

}

//checkPeiPaiHSZFX 检测换三张方向
func checkPeiPaiHSZFX(fxValue string) {

}

//LogPeiPaiInfos 打印配牌信息
func LogPeiPaiInfos() {
	fmt.Println(peiPaiInfos)
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

// func init() {
// 	peipaiValue = make(map[string]string)
// 	peiPaiInfos = []*peipaiInfo{}
// }

// Run 启动配牌
func Run(addr string) error {
	// PP_PORT = addr
	http.HandleFunc("/", Demo)
	return http.ListenAndServe(addr, nil)
}
