package main

import (
	"io"
	"log"
	"net"
	"strconv"
	"google.golang.org/grpc"
	proto "steve/serviceexample/grpcStreamExample/proto"
	"context"
)


type Streamer struct{}

func (s *Streamer) BidStream(stream proto.Chat_BidStreamServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			log.Println("收到客户端通过context发出的终止信号")
			return ctx.Err()
		default:
			输入, err := stream.Recv()
			if err == io.EOF {
				log.Println("客户端发送的数据流结束")
				return nil
			}
			if err != nil {
				log.Println("接收数据出错:", err)
				return err
			}
			switch 输入.Input {
			case "结束对话\n":
				log.Println("收到'结束对话'指令")
				if err := stream.Send(&proto.Response{Output: "收到结束指令"}); err != nil {
					return err
				}
				return nil
			case "返回数据流\n":
				log.Println("收到'返回数据流'指令")
				for i := 0; i < 10; i++ {
					if err := stream.Send(&proto.Response{Output: "数据流 #" + strconv.Itoa(i)}); err != nil {
						return err
					}
				}
			default:
				log.Printf("[收到消息]: %s", 输入.Input)
				if err := stream.Send(&proto.Response{Output: "服务端返回: " + 输入.Input}); err != nil {
					return err
				}
			}
		}
	}
}

func (s *Streamer) BidNormal(context context.Context, request *proto.Request) (*proto.Response, error) {
	return &proto.Response{Output: request.Input + "World"}, nil
}

func main() {
	log.Println("启动服务端...")
	server := grpc.NewServer()

	proto.RegisterChatServer(server, &Streamer{})
	address, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	if err := server.Serve(address); err != nil {
		panic(err)
	}
}