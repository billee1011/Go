package gangtests

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_Minggang 测试明杠
// 开始游戏后， 庄家出 9W， 1 号玩家可杠， 没有其他玩家可胡
// 期望：
//		①.1 号玩家将收到出牌问询通知, 通知中的可选行为包括杠
//      ②.1 号玩家请求执行杠后， 所有玩家收到杠通知消息, 其中杠的牌为 9W， 杠的玩家为 1 号玩家， 杠类型为明杠
func Test_Minggang(t *testing.T) {
	var Int9W uint32 = 19
	params := global.NewCommonStartGameParams()

	params.BankerSeat = 0
	gangSeat := 1
	bankerSeat := params.BankerSeat

	// 庄家的最后一张牌改为 9W
	params.Cards[bankerSeat][13] = &global.Card9W
	// 1 号玩家最后3张牌改为 9W
	params.Cards[gangSeat][10] = &global.Card9W
	params.Cards[gangSeat][11] = &global.Card9W
	params.Cards[gangSeat][12] = &global.Card9W

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

	// 发送杠请求
	assert.Nil(t, utils.SendGangReq(deskData, gangSeat, Int9W, room.GangType_MingGang))

	// 检测所有玩家收到杠通知
	bankerPlayer := utils.GetDeskPlayerBySeat(bankerSeat, deskData)
	utils.CheckGangNotify(t, deskData, gangPlayer.Player.GetID(), bankerPlayer.Player.GetID(), Int9W, room.GangType_MingGang)
}
