package tests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/connect"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Xipai(t *testing.T) {
	xipaiExpectors := map[int]interfaces.MessageExpector{}

	for i := 0; i < 4; i++ {
		client := connect.NewTestClient(ServerAddr, ClientVersion)
		assert.NotNil(t, client)
		player, err := utils.LoginUser(client, "test_user")
		assert.Nil(t, err)
		assert.NotNil(t, player)

		expector, err := client.ExpectMessage(msgid.MsgID_ROOM_XIPAI_NTF)
		assert.Nil(t, err)
		xipaiExpectors[i] = expector

		assert.Nil(t, utils.ApplyJoinDesk(player))
	}

	for _, e := range xipaiExpectors {
		xipaiNtf := room.RoomXipaiNtf{}
		assert.Nil(t, e.Recv(time.Second*1, &xipaiNtf))
		assert.True(t, xipaiNtf.GetTotalCard() > 0)
		assert.Equal(t, 2, len(xipaiNtf.GetDices()))
		zjIdx := xipaiNtf.GetBankerSeat()
		assert.True(t, zjIdx >= 0 && zjIdx < 4)
	}
}
