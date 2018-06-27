package utils

import (
	"fmt"
	"net/http"
	"steve/simulate/config"

	"github.com/Sirupsen/logrus"
)

func mjconfig(gameName string, open bool, gold uint64) error {
	url := fmt.Sprintf("%s/option/?game=%v&hszswitch=%v&gold=%v", config.PeipaiURL, gameName, open, gold)
	return requestOpen(url)
}

func requestOpen(url string) error {
	logrus.WithField("url", url).Info("")
	_, err := http.DefaultClient.Get(url)
	return err
}
