package utils

import (
	"fmt"
	"net/http"
	"steve/simulate/config"

	"github.com/Sirupsen/logrus"
)

func hszSwitch(open bool) error {
	url := fmt.Sprintf("%s?hszswitch=%v", config.SwitchURL, open)
	return requestOpen(url)
}

func requestOpen(url string) error {
	logrus.WithField("url", url).Info("")
	_, err := http.DefaultClient.Get(url)
	return err
}
