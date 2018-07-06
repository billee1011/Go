package utils

import (
	"fmt"
	"net/http"
	"steve/simulate/config"

	"github.com/Sirupsen/logrus"
)

func majongOption(gameName string, open bool) error {
	url := fmt.Sprintf("%s/option/?game=%v&hszswitch=%v", config.GetPeipaiURL(), gameName, open)
	return requestOpen(url)
}

func majongPlayerGold(seatGold, seatID map[int]uint64) error {
	for seat, playerID := range seatID {
		if gold, isExist := seatGold[seat]; isExist {
			url := fmt.Sprintf("%s/setgold/?player_id=%v&gold=%v", config.GetPeipaiURL(), playerID, gold)
			if err := requestOpen(url); err != nil {
				return err
			}
		}
	}
	return nil
}

func requestOpen(url string) error {
	logrus.WithField("url", url).Info("")
	_, err := http.DefaultClient.Get(url)
	return err
}
