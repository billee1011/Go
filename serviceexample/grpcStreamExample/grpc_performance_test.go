package grpcStreamExample

import (
	"testing"
	"google.golang.org/grpc"
	"log"
	"steve/serviceexample/grpcStreamExample/proto"
	"context"
	"time"
)

func TestStreamGO(t *testing.T) {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	defer conn.Close()

	client := proto.NewChatClient(conn)
	ctx := context.Background()
	stream, err := client.BidStream(ctx)
	if err != nil {
		log.Printf("Create stream failed! error : [%v]\n", err)
	}

	if err := stream.Send(&proto.Request{"Hello, "}); err != nil {
		return
	}
	result, _ := stream.Recv()
	log.Println(result.Output)
}

func TestNormal(t *testing.T) {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		log.Printf("Connect failed! error : [%v]\n", err)
		return
	}
	defer conn.Close()

	client := proto.NewChatClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		_, err := client.BidNormal(ctx, &proto.Request{Input: "Hello, "})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	}
}

func Benchmark_Stream(b *testing.B) {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		log.Printf("Connect failed! error : [%v]\n", err)
		return
	}
	defer conn.Close()

	client := proto.NewChatClient(conn)
	ctx := context.Background()
	defer ctx.Done()

	stream, err := client.BidStream(ctx)
	defer stream.CloseSend()

	for i := 0; i < b.N; i++ {
		if err := stream.Send(&proto.Request{"Hello, "}); err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		stream.Recv()
	}
}

func Benchmark_Normal(b *testing.B) {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		log.Printf("Connect failed! error : [%v]\n", err)
		return
	}
	defer conn.Close()

	client := proto.NewChatClient(conn)
	ctx := context.Background()
	defer ctx.Done()

	for i := 0; i < b.N; i++ {
		_, err := client.BidNormal(ctx, &proto.Request{Input: "Hello, "})

		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	}
}