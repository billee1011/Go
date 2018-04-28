package socket

import (
	"fmt"
	"net"
)

type Socket interface {
	SendPackage(pkg []byte) error
	RecvPackage() (pkg []byte, err error)
	Close() error
}

type SocketFactory interface {
	NewSocket(conn net.Conn, unpacker Unpacker) (Socket, error)
}

type socketImpl struct {
	conn     net.Conn
	unpacker Unpacker
}

var _ Socket = new(socketImpl)

func (socket *socketImpl) SendPackage(pkg []byte) error {
	sz, err := socket.conn.Write(pkg)
	if err != nil {
		return err
	}
	if sz != len(pkg) {
		return fmt.Errorf("package send failed. write size(%v) does not equal to expected size(%v)", sz, len(pkg))
	}
	return nil
}

func (socket *socketImpl) RecvPackage() (pkg []byte, err error) {
	return socket.unpacker.Unpack(socket.conn)
}

func (socket *socketImpl) Close() error {
	return socket.conn.Close()
}

type socketFactoryImpl struct{}

var _ SocketFactory = new(socketFactoryImpl)

func (factory *socketFactoryImpl) NewSocket(conn net.Conn, unpacker Unpacker) (Socket, error) {
	if conn == nil || unpacker == nil {
		return nil, fmt.Errorf("conn and unpacker can't be nil")
	}
	return &socketImpl{
		conn:     conn,
		unpacker: unpacker,
	}, nil
}
