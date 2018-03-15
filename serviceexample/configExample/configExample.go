package main

import (
	"fmt"
	"steve/structs"
	"steve/structs/configuration"
	"steve/structs/service"
)

type configExampleService struct{}

func (res *configExampleService) Start(e *structs.Exposer, param ...string) error {
	c := e.Configuration
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

func GetService() service.Service {
	return &configExampleService{}
}

func main() {

}
