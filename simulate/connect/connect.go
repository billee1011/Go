package connect

import (
	"fmt"
	"io"
	"steve/base/socket"
)

type tcpUnpacker struct {
}

func (unpacker *tcpUnpacker) Unpack(r io.Reader) ([]byte, error) {
	byteSz := make([]byte, 2)
	_, err := r.Read(byteSz)
	if err != nil {
		return nil, err
	}
	sz := 0
	sz += int(byteSz[0]) << 8
	sz += int(byteSz[1])

	if sz < 2 {
		return nil, fmt.Errorf("包长错误， %v", sz)
	}

	data := make([]byte, sz-2)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}

func connectServer(addr string) (socket.Socket, error) {
	client, err := socket.NewClient(&tcpUnpacker{}, nil)
	if err != nil {
		return nil, fmt.Errorf("创建客户端失败 : %v", err)
	}
	sock, err := client.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("连接服务器失败： %v", err)
	}
	return sock, nil
}
