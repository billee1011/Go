package qitests

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_MingGang_Qi 测试明杠弃
// 期望：
// 庄家出9W后，1号玩家将收到出牌问询通知，可杠9W
// 1号玩家发出弃杠请求后，1号玩家摸牌，1号玩家将收到自询通知
func Test_MingGang_Qi(t *testing.T) {
	var Int9W uint32 = 19
	params := global.NewCommonStartGameParams()

	params.BankerSeat = 0
	gangSeat := 1
	bankerSeat := params.BankerSeat

	// 庄家的最后一张牌改为 9W
	params.Cards[bankerSeat][13] = 19
	// 1 号玩家最后3张牌改为 9W
	params.Cards[gangSeat][10] = 19
	params.Cards[gangSeat][11] = 19
	params.Cards[gangSeat][12] = 19

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, Int9W))

	// 1 号玩家收到可杠通知
	gangPlayer := utils.GetDeskPlayerBySeat(gangSeat, deskData)
	expector, _ := gangPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.Equal(t, Int9W, ntf.GetCard())
	assert.True(t, ntf.GetEnableMinggang())
	assert.True(t, ntf.GetEnableQi())

	// 发送弃请求
	assert.Nil(t, utils.SendQiReq(deskData, gangSeat))

	// 1 号玩家收到可自询通知
	utils.WaitZixunNtf(deskData, gangSeat)
}
