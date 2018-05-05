package tests

import (
	"sync"
	"unsafe"

	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/connect"

	"time"

	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

func TestLogin(t *testing.T) {
	COUNT := 100
	var wg sync.WaitGroup
	wg.Add(COUNT)

	for i := 0; i < COUNT; i++ {
		go func() {
			defer wg.Done()
			c := connect.NewTestClient("127.0.0.1:36001", "1.0")
			assert.NotNil(t, c)

			rsp := room.RoomLoginRsp{}

			err := c.Request(connect.SendHead{
				Head: connect.Head{
					MsgID: uint32(msgid.MsgID_room_login_req),
				},
			}, &room.RoomLoginReq{
				UserName: proto.String("test_user"),
			}, time.Second*5, uint32(msgid.MsgID_room_login_rsp), &rsp)

			assert.Nil(t, err)
			assert.NotEqual(t, rsp.GetPlayerId(), 0)
		}()
	}
	wg.Wait()
}

/*
func TestLogin2(t *testing.T) {
	COUNT := 100

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

			result, err := c.SendPackage(connect.SendHead{
				Head: connect.Head{
					MsgID: uint32(steve_proto_msg.MsgID_hall_login),
				}},
				&steve_proto_msg.LoginReq{
					UserName: &userName,
				})

			assert.Nil(t, err)
			assert.NotNil(t, result)
			rsp, err := c.WaitMessage(context.Background(), uint32(steve_proto_msg.MsgID_hall_login), result.SendTimestamp)
			assert.Nil(t, err)

			loginRsp, ok := rsp.Body.(*steve_proto_msg.LoginRsp)
			assert.True(t, ok)
			assert.Equal(t, steve_proto_msg.ErrorCode_err_OK, loginRsp.GetResult())
			assert.NotEqual(t, 0, loginRsp.GetUserId())
		}()
	}
	wg.Wait()
}
*/
