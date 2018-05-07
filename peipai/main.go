package peipai

import (
	"fmt"
	"net/http"
	"strconv"
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
	fmt.Println("lens:", len(lenValue))
	if len(lenValue) != 0 {
		lenNum, err := strconv.Atoi(lenValue)
		fmt.Println("lens:", lens)
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
	fmt.Println("game: ", game)
	fmt.Println("value: ", value)
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
	peipai_value[game] = value
	pp := peipaiInfo{
		key:   game,
		cards: value,
		len:   lens,
		fx:    fx,
	}
	addPeiPaiInfo(pp)
	fmt.Println(pp)
	okStr := "ok"
	w.WriteHeader(200)
	w.Write([]byte(okStr))
}

var peipai_value map[string]string

var peiPaiInfos []peipaiInfo

// var PP_PORT string
type peipaiInfo struct {
	//配牌关键字
	key string
	//配牌
	cards string
	//墙牌长度
	len int
	//换三张方向
	fx int
}

//GetPeiPai 通过配牌关键字拿到配牌
func GetPeiPai(key string) (string, error) {
	if value, ok := peipai_value[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("不存在这个配牌key")
}

//GetLensOfWallCards 牌墙长度
func GetLensOfWallCards(key string) int {
	for _, pp := range peiPaiInfos {
		if pp.key == key {
			return pp.len
		}
	}
	return 0
}

//GetHSZFangXiang 换三张方向
func GetHSZFangXiang(key string) int {
	for _, pp := range peiPaiInfos {
		if pp.key == key {
			return pp.fx
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
		if peipaiInfo.key == pp.key {
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
func init() {
	peipai_value = make(map[string]string)
	peiPaiInfos = make([]peipaiInfo, 0)
}

// Run 启动配牌
func Run(addr string) {
	// PP_PORT = addr
	http.HandleFunc("/", Demo)
	http.ListenAndServe(addr, nil)
}
