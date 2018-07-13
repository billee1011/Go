package tests

import (
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/connect"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Xipai(t *testing.T) {
	xipaiExpectors := map[int]interfaces.MessageExpector{}

	for i := 0; i < 4; i++ {
		client := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
		assert.NotNil(t, client)
		player, err := utils.LoginUser(client, global.AllocUserName())
		assert.Nil(t, err)
		assert.NotNil(t, player)

		expector, err := client.ExpectMessage(msgid.MsgID_ROOM_XIPAI_NTF)
		assert.Nil(t, err)
		xipaiExpectors[i] = expector
		gameID := room.GameId(1)
		_, err = utils.ApplyJoinDesk(player, gameID)
		assert.Nil(t, err)
	}

	for _, e := range xipaiExpectors {
		xipaiNtf := room.RoomXipaiNtf{}
		assert.Nil(t, e.Recv(time.Second*1, &xipaiNtf))
		assert.True(t, xipaiNtf.GetTotalCard() > 0)
		// assert.Equal(t, uint32(108), xipaiNtf.GetTotalCard())
		assert.Equal(t, 2, len(xipaiNtf.GetDices()))
		zjIdx := xipaiNtf.GetBankerSeat()
		assert.True(t, zjIdx >= 0 && zjIdx < 4)
	}
}
