package core

import (
	"steve/structs"
	"steve/structs/service"
)

type backCore struct {
	e *structs.Exposer
}

// NewService 创建服务
func NewService() service.Service {
	return new(backCore)
}

func (c *backCore) Init(e *structs.Exposer, param ...string) error {

	return nil
}

func (c *backCore) Start() error {
	back := make(chan bool)
	<-back
	return nil
}
