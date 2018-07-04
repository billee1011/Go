package utils

import (
	"fmt"
	"net/http"
	"steve/simulate/config"

	"github.com/Sirupsen/logrus"
)

// 通知服务器：是否开启换三张
func majongOption(gameName string, open bool) error {
	url := fmt.Sprintf("%s/option/?game=%v&hszswitch=%v", config.MaJongConfigURL, gameName, open)
	return requestOpen(url)
}

// 通知服务器：所有玩家的金币数
// 参数seatGold : 座位ID 与 金币 的map
// 参数seatID 	: 座位ID 与 playerID 的map
func majongPlayerGold(seatGold, seatID map[int]uint64) error {
	for seat, playerID := range seatID {
		if gold, isExist := seatGold[seat]; isExist {
			url := fmt.Sprintf("%s/setgold/?player_id=%v&gold=%v", config.MaJongConfigURL, playerID, gold)
			if err := requestOpen(url); err != nil {
				return err
			}
		}
	}
	return nil
}

// 发出url的get请求
func requestOpen(url string) error {
	logrus.WithField("url", url).Info("")
	_, err := http.DefaultClient.Get(url)
	return err
}
