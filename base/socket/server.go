package socket

import (
	"fmt"
	"net"
)

type Server interface {
	Listen(addr string) error
	Accept() (Socket, error)
	Close() error
}

type serverImpl struct {
	lis           net.Listener
	unpacker      Unpacker
	socketFactory SocketFactory
}

func NewServer(unpacker Unpacker, socketFactory SocketFactory) (Server, error) {
	if unpacker == nil {
		return nil, fmt.Errorf("unpacker should have a non-nil value")
	}
	if socketFactory == nil {
		socketFactory = &socketFactoryImpl{}
	}
	return &serverImpl{
		unpacker:      unpacker,
		socketFactory: socketFactory,
	}, nil
}

var _ Server = &serverImpl{}

func (server *serverImpl) Listen(addr string) error {
	var err error
	server.lis, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return nil
}

func (server *serverImpl) Accept() (Socket, error) {
	if server.lis == nil {
		return nil, fmt.Errorf("not listen yet")
	}
	conn, err := server.lis.Accept()
	if err != nil {
		return nil, err
	}
	return server.socketFactory.NewSocket(conn, server.unpacker)
}

func (server *serverImpl) Close() error {
	if server.lis == nil {
		return nil
	}
	return server.lis.Close()
}
