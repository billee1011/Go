package main

import (
	"steve/structs"
	"steve/structs/logger"
	"steve/structs/service"
)

type loggerExampleService struct{}

func (les *loggerExampleService) Start(e *structs.Exposer, param ...string) error {
	log := e.Logger

	{
		log := log.WithField("hello", "world")
		log.Debug("debug")
		log.Info("info")
		log.Warn("warn")
		log.Error("error")
		// logger.Fatal("fatal")  // will panic
	}
	{
		log := log.WithFields(logger.Fields{
			"hello": "world",
			"你好":    "世界",
		})
		log.Debug("debug")
		log.Info("info")
		log.Warn("warn")
		log.Error("error")
		// logger.Fatal("fatal")  // will panic
	}
	{
		log.Debug("debug")
		log.Info("info")
		log.Warn("warn")
		log.Error("error")
		// logger.Fatal("fatal")  // will panic
	}
	return nil
}

func GetService() service.Service {
	return &loggerExampleService{}
}

func main() {
}
