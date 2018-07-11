package tests

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// Test_Chupai 测试开局后的出牌
func Test_Chupai(t *testing.T) {
	params := global.NewCommonStartGameParams()
	deskData, err := utils.StartGame(params)

	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 庄家出一万
	zjPlayer := utils.GetDeskPlayerBySeat(params.BankerSeat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	zjClient.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_CHUPAI_REQ), &room.RoomChupaiReq{
		Card: proto.Uint32(11),
	})

	// 其他玩家收到庄家出牌通知
	for _, deskPlayer := range deskData.Players {
		expector, _ := deskPlayer.Expectors[msgid.MsgID_ROOM_CHUPAI_NTF]
		ntf := room.RoomChupaiNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, ntf.GetPlayer(), zjPlayer.Player.GetID())
		assert.Equal(t, uint32(11), ntf.GetCard())
	}
}
