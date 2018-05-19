package tests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// Test_Mopai 摸牌测试
func Test_Mopai(t *testing.T) {
	params := global.NewCommonStartGameParams()
	// 庄家最后一张牌改为1筒
	params.Cards[0][13] = &global.Card1B
	deskData, err := utils.StartGame(params)

	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 庄家出1筒
	assert.Nil(t, utils.WaitZixunNtf(deskData, params.BankerSeat))
	zjPlayer := utils.GetDeskPlayerBySeat(params.BankerSeat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	zjClient.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_CHUPAI_REQ), &room.RoomChupaiReq{
		Card: proto.Uint32(31),
	})

	mopaiSeat := (params.BankerSeat + 1) % len(deskData.Players)
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
