package main

import (
	"context"
	"errors"
	"fmt"
	"steve/serviceexample/rpcexample/proto"
	"steve/structs"
	"steve/structs/service"
)

type RPCExampleClient struct {
}

func (rec *RPCExampleClient) Start(e *structs.Exposer, param ...string) error {
	cc, err := e.RPCClient.GetClientConnByServerName("example_service")
	if err != nil {
		return fmt.Errorf("Get client connection failed:%v", err)
	}
	if cc == nil {
		return errors.New("no service named example_service. ensure your consul agent is running and configed example_service")
	}
	client := hw.NewHelloWorldClient(cc)
	resp, err := client.HelloWorld(context.Background(), &hw.HelloWorldRequest{
		Name: "world",
	})
	if err != nil {
		return fmt.Errorf("call HelloWorld failed: %v", err)
	}
	fmt.Println("receive response from server:", resp.GetEcho())
	return nil
}

func GetService() service.Service {
	return &RPCExampleClient{}
}

func main() {

}
