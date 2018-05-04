package tests

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
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
					MsgID: uint32(msgid.MsgID_room_login_req),
				},
			}, &room.RoomLoginReq{
				UserName: &userName,
			}, time.Second*5)

			assert.Nil(t, err)
			assert.NotNil(t, rsp)
		}()
	}
	wg.Wait()
}
