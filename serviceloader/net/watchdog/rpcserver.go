package watchdog

import (
	"net"
	"steve/structs/proto/base"

	"google.golang.org/grpc"
)

type rpcServer struct {
	work       workerFunc
	grpcServer *grpc.Server
}

var _ server = new(rpcServer)
var _ base.ExchangerServer = new(rpcServer)

type rpcExchanger struct {
	e base.Exchanger_ExchangeServer
}

var _ exchanger = new(rpcExchanger)

func (e *rpcExchanger) Recv() ([]byte, error) {
	var ctx *base.ExchangeContext
	var err error
	if ctx, err = e.e.Recv(); err != nil {
		return []byte{}, err
	}
	return ctx.Data, err
}

func (e *rpcExchanger) Send(data []byte) error {
	return e.e.Send(&base.ExchangeContext{
		Data: data,
	})
}

func (s *rpcServer) Exchange(e base.Exchanger_ExchangeServer) error {
	return s.work(&rpcExchanger{
		e: e,
	})
}

func (s *rpcServer) Serve(addr string, worker workerFunc) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.work = worker

	server := grpc.NewServer()
	base.RegisterExchangerServer(server, s)

	s.grpcServer = server
	return server.Serve(lis)
}

func (s *rpcServer) Close() {
	if s.grpcServer == nil {
		return
	}
	s.grpcServer.Stop()
}
