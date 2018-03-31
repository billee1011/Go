package main

import (
	"fmt"
	"steve/structs"
	"steve/structs/configuration"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
)

type configExampleService struct {
	e *structs.Exposer
}

func (res *configExampleService) Init(e *structs.Exposer, param ...string) error {
	res.e = e

	return nil
}

func (res *configExampleService) Start() error {
	c := res.e.Configuration

	logrus.WithField("configuration", c).Debug("config")

	version, err := c.GetConfigVer(configuration.Product)
	if err != nil {
		return err
	}
	gameConfig, err := c.GetConfig(&configuration.ConfigGetParam{
		Env:     configuration.Product,
		Version: version,
		Key:     "game",
		Prefix:  true,
	})
	if err != nil {
		return err
	}
	for k, v := range gameConfig {
		fmt.Println(k, v)
	}
	helloConfig, err := c.GetConfig(&configuration.ConfigGetParam{
		Env:     configuration.Product,
		Version: version,
		Key:     "hello",
		Prefix:  false,
	})
	if err != nil {
		return err
	}
	fmt.Println("hello", helloConfig["hello"])
	return nil
}

// GetService 提供给 serviceloader 的接口
func GetService() service.Service {
	return &configExampleService{}
}

func main() {

}
