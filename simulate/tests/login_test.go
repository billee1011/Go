package tests

import (
	"steve/structs/proto/msg"
	"sync"

	"steve/simulate/connect"

	"time"

	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {

	COUNT := 200
	var wg sync.WaitGroup
	wg.Add(COUNT)

	for i := 0; i < COUNT; i++ {
		go func() {
			defer wg.Done()
			c := connect.NewTestClient("127.0.0.1:36001", "1.0")
			assert.NotNil(t, c)
			rsp, err := c.Request(connect.SendHead{
				Head: connect.Head{
					MsgID: uint32(steve_proto_msg.MsgID_hall_login),
				},
			}, &steve_proto_msg.LoginReq{
				UserName: proto.String("Adam"),
			}, time.Second*5)
			assert.Nil(t, err)
			assert.NotNil(t, rsp)
			assert.Equal(t, steve_proto_msg.ErrorCode_err_OK, rsp.Body.(*steve_proto_msg.LoginRsp).GetResult())
		}()
	}
	wg.Wait()
}
