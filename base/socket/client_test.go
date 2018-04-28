package socket

import (
	"net"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Client(t *testing.T) {
	ctrl := gomock.NewController(t)
	unpacker := NewMockUnpacker(ctrl)
	socketFactory := NewMockSocketFactory(ctrl)
	clientSocket := NewMockSocket(ctrl)

	client, err := NewClient(unpacker, socketFactory)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	lis, err := net.Listen("tcp", "127.0.0.1:36001")
	assert.Nil(t, err)
	assert.NotNil(t, lis)

	defer lis.Close()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		lis.Accept()
	}()

	socketFactory.EXPECT().NewSocket(gomock.Any(), Unpacker(unpacker)).Times(1).Return(clientSocket, nil)

	client.Dial("127.0.0.1:36001")
}
