package peipai

import (
	"net/http"
	"steve/room/peipai/handle"
	// "github.com/go-redis/redis"
	// "github.com/spf13/viper"
)

//checkPeiPaiKey 检测配牌key
func checkPeiPaiKey(gameName string) {
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

// Run 启动配牌
func Run(addr string) error {
	// PP_PORT = addr
	http.HandleFunc("/", handle.Peipai)
	http.HandleFunc("/option/", handle.Option)
	http.HandleFunc("/setgold/", handle.SetGoldHandle)
	return http.ListenAndServe(addr, nil)
}
