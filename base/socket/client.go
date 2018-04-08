package socket

import (
	"fmt"
	"net"
)

// Client 客户端接口
type Client interface {

	// Dial 连接服务器， addr 为服务器地址， 格式： ip:port， 如 127.0.0.1:8001
	Dial(addr string) (Socket, error)
}

type clientImpl struct {
	Socket
	unpacker      Unpacker
	socketFactory SocketFactory
}

var _ Client = new(clientImpl)

// NewClient 创建客户端接口
// unpacker 为消息解包器， 不能为空，
func NewClient(unpacker Unpacker, socketFactory SocketFactory) (Client, error) {
	if unpacker == nil {
		return nil, fmt.Errorf("unpakcer should have a non-nil value")
	}
	if socketFactory == nil {
		socketFactory = &socketFactoryImpl{}
	}
	return &clientImpl{
		unpacker:      unpacker,
		socketFactory: socketFactory,
	}, nil
}

func (client *clientImpl) Dial(addr string) (Socket, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	socket, err := client.socketFactory.NewSocket(conn, client.unpacker)
	return socket, err
}
