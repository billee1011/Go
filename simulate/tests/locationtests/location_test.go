package locationtests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocation(t *testing.T) {
	params := global.NewCommonStartGameParams()
	deskData, err := utils.StartGame(params)

	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	//开局后,庄家请求地理位置信息
	zjPlayer := utils.GetDeskPlayerBySeat(params.BankerSeat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	zjClient.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_PLAYER_LOCATION_REQ), &room.RoomPlayerLocationReq{})
	expecter, err := zjClient.ExpectMessage(msgid.MsgID_ROOM_PLAYER_LOCATION_RSP)
	assert.Nil(t, err)
	ntf := &room.RoomPlayerLocationRsp{}
	assert.Nil(t, expecter.Recv(time.Second*1, ntf))
	locations := ntf.GetLocations()
	for _, location := range locations {
		infos := location.GetLocation()
		for _, info := range infos {
			assert.Equal(t, 101.101, info.GetLongitude())
			assert.Equal(t, 202.202, info.GetLatitude())
		}
	}
}
