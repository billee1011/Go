package main

import (
	"context"
	"steve/serviceexample/rpcexample/proto"
	"steve/structs"
	"steve/structs/service"
)

type RPCExampleService struct {
	e *structs.Exposer
}

type HelloWorldService struct {
}

func (hws *HelloWorldService) HelloWorld(ctx context.Context, req *hw.HelloWorldRequest) (rsp *hw.HelloWorldResponse, err error) {
	rsp = &hw.HelloWorldResponse{}
	rsp.Echo = "Hello," + req.GetName()
	err = nil
	return
}

func (res *RPCExampleService) Init(e *structs.Exposer, param ...string) error {
	rpcServer := e.RPCServer
	err := rpcServer.RegisterService(hw.RegisterHelloWorldServer, &HelloWorldService{})
	if err != nil {
		return err
	}
	return nil
}

func (res *RPCExampleService) Start() error {
	return nil
}

func GetService() service.Service {
	return &RPCExampleService{}
}

func main() {

}
