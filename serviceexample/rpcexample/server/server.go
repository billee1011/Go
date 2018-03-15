package main

import (
	"context"
	"net"
	"steve/serviceexample/rpcexample/proto"
	"steve/structs"
	"steve/structs/service"
)

type RPCExampleService struct {
}

type HelloWorldService struct {
}

func (hws *HelloWorldService) HelloWorld(ctx context.Context, req *hw.HelloWorldRequest) (rsp *hw.HelloWorldResponse, err error) {
	rsp = &hw.HelloWorldResponse{}
	rsp.Echo = "Hello," + req.GetName()
	err = nil
	return
}

func (res *RPCExampleService) Start(e *structs.Exposer, param ...string) error {
	rpcServer := e.RPCMgr.NewRPCServer()
	err := rpcServer.RegisterService(hw.RegisterHelloWorldServer, &HelloWorldService{})
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", "0.0.0.0:7878")
	if err != nil {
		return err
	}
	err = rpcServer.StartServer(lis)
	return err
}

func GetService() service.Service {
	return &RPCExampleService{}
}

func main() {

}
