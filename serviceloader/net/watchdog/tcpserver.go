package watchdog

import (
	"fmt"
	"io"
	"steve/base/socket"
)

type tcpUnpacker struct {
}

var _ socket.Unpacker = new(tcpUnpacker)

func (unpacker *tcpUnpacker) Unpack(r io.Reader) ([]byte, error) {
	byteSz := make([]byte, 2)
	_, err := r.Read(byteSz)
	if err != nil {
		return nil, err
	}
	sz := (int(byteSz[0]) << 8) | int(byteSz[1])

	if sz < 2 {
		return nil, fmt.Errorf("包长错误， %v", sz)
	}
	data := make([]byte, sz-2)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}

type tcpExchanger struct {
	sock socket.Socket
}

var _ exchanger = new(tcpExchanger)

func (e *tcpExchanger) Recv() ([]byte, error) {
	return e.sock.RecvPackage()
}

func (e *tcpExchanger) Send(data []byte) error {
	dataSz := len(data) + 2

	header := make([]byte, 2)
	header[0] = byte((dataSz >> 8) & 0xff)
	header[1] = byte(dataSz & 0xff)

	wholeData := make([]byte, dataSz)
	copy(wholeData[:2], header)
	copy(wholeData[2:], data)

	return e.sock.SendPackage(wholeData)
}

type tcpServer struct {
	worker       workerFunc
	socketServer socket.Server
	addr         string
}

var _ server = new(tcpServer)

func (s *tcpServer) workOnSocket(sock socket.Socket) {
	s.worker(&tcpExchanger{
		sock: sock,
	})
	sock.Close()
}

func (s *tcpServer) Serve(addr string, worker workerFunc) error {
	s.worker = worker
	socketServer, err := socket.NewServer(&tcpUnpacker{}, nil)
	if err != nil {
		return err
	}
	if err := socketServer.Listen(addr); err != nil {
		return err
	}
	s.socketServer = socketServer
	for {
		sock, err := socketServer.Accept()
		if err != nil {
			return err
		}
		go s.workOnSocket(sock)
	}
}

func (s *tcpServer) Close() {
	s.socketServer.Close()
}
