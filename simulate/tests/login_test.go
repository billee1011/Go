package tests

import (
	"fmt"
	"steve/structs/proto/msg"
	"sync"
	"unsafe"

	"steve/simulate/connect"

	"time"

	"testing"

	"github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

func TestLogin(t *testing.T) {

	COUNT := 1
	var wg sync.WaitGroup
	wg.Add(COUNT)

	for i := 0; i < COUNT; i++ {
		go func() {
			defer wg.Done()
			c := connect.NewTestClient("127.0.0.1:36001", "1.0")

			f := fuzz.New()
			var userName string
			f.NilChance(.0)
			f.Fuzz(&userName)

			assert.NotNil(t, c)
			rsp, err := c.Request(connect.SendHead{
				Head: connect.Head{
					MsgID: uint32(steve_proto_msg.MsgID_hall_login),
				},
			}, &steve_proto_msg.LoginReq{
				UserName: &userName,
			}, time.Second*5)
			assert.Nil(t, err)
			assert.NotNil(t, rsp)

			loginRsp := rsp.Body.(*steve_proto_msg.LoginRsp)
			assert.Equal(t, steve_proto_msg.ErrorCode_err_OK, loginRsp.GetResult())
			fmt.Println(loginRsp)
			assert.NotEqual(t, 0, loginRsp.GetUserId())
		}()
	}
	wg.Wait()
}
