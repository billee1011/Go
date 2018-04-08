package socket

import (
	"net"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Server(t *testing.T) {
	ctrl := gomock.NewController(t)
	unpacker := NewMockUnpacker(ctrl)
	socketFactory := NewMockSocketFactory(ctrl)
	svrSocket := NewMockSocket(ctrl)

	svr, err := NewServer(unpacker, socketFactory)
	assert.Nil(t, err)
	assert.Nil(t, svr.Listen("127.0.0.1:36001"))
	defer func() {
		assert.Nil(t, svr.Close())
	}()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		socketFactory.EXPECT().NewSocket(gomock.Any(), Unpacker(unpacker)).Return(svrSocket, nil)
		socket, err := svr.Accept()
		assert.Nil(t, err)
		assert.Equal(t, svrSocket, socket)
	}()

	cc, err := net.Dial("tcp", "127.0.0.1:36001")
	assert.Nil(t, err)
	assert.NotNil(t, cc)
	wg.Wait()
}
