package socket

import (
	"io"
	"net"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_socketImpl_SendReceive(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:36001")
	assert.Nil(t, err)
	assert.NotNil(t, lis)

	var clientSocket, serverSocket *socketImpl
	ctl := gomock.NewController(t)
	clientUnpakcer, serverUnpakcer := NewMockUnpacker(ctl), NewMockUnpacker(ctl)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := lis.Accept()
		assert.Nil(t, err)
		assert.NotNil(t, conn)
		serverSocket = &socketImpl{
			conn:     conn,
			unpacker: serverUnpakcer,
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:36001")
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	clientSocket = &socketImpl{
		conn:     conn,
		unpacker: clientUnpakcer,
	}
	wg.Wait()

	defer func() {
		assert.Nil(t, clientSocket.Close())
		assert.Nil(t, serverSocket.Close())
	}()

	sendData := "hello, world!"

	wg = new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := clientSocket.SendPackage([]byte(sendData))
		assert.Nilf(t, err, "%v", err)
	}()

	go func() {
		defer wg.Done()
		serverUnpakcer.EXPECT().Unpack(io.Reader(serverSocket.conn)).DoAndReturn(func(r io.Reader) ([]byte, error) {
			data := make([]byte, len(sendData))
			_, err := io.ReadFull(r, data)
			if err != nil {
				return nil, err
			}
			return data, nil
		})

		data, err := serverSocket.RecvPackage()
		assert.Nilf(t, err, "%v", err)
		assert.Equal(t, sendData, string(data))
	}()

	wg.Wait()
}
