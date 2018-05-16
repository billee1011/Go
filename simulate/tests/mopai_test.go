package tests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// Test_Mopai 摸牌测试
func Test_Mopai(t *testing.T) {
	deskData, err := utils.StartGame(commonStartGameParams)

	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 庄家出一万
	zjPlayer := utils.GetDeskPlayerBySeat(commonStartGameParams.BankerSeat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	zjClient.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_CHUPAI_REQ), &room.RoomChupaiReq{
		Card: proto.Uint32(11),
	})

	mopaiSeat := (commonStartGameParams.BankerSeat + 1) % len(deskData.Players)
	mopaiPlayer := utils.GetDeskPlayerBySeat(mopaiSeat, deskData)
	mopaiPlayerID := mopaiPlayer.Player.GetID()
	// 所有玩家收到庄家下家摸牌通知
	for _, deskPlayer := range deskData.Players {
		expector, _ := deskPlayer.Expectors[msgid.MsgID_ROOM_MOPAI_NTF]
		ntf := room.RoomMopaiNtf{}
		assert.Nil(t, expector.Recv(time.Second*1, &ntf))
		assert.Equal(t, mopaiPlayerID, ntf.GetPlayer())
		assert.Equal(t, false, ntf.GetBack())
		if deskPlayer.Seat == mopaiSeat {
			assert.Equal(t, uint32(31), ntf.GetCard()) // 摸了一张 1 筒
		}
	}
}
