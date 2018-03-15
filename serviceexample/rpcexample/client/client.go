package main

import (
	"context"
	"fmt"
	"steve/serviceexample/rpcexample/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:7878", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := hw.NewHelloWorldClient(conn)
	rsp, err := client.HelloWorld(context.Background(), &hw.HelloWorldRequest{
		Name: "Adam",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.GetEcho())
}
