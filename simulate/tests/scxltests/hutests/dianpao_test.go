package hutests

import (
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_Dianpao 点炮测试
// 开始游戏后，庄家出9W，1 号玩家可以胡，其他玩家都不可以胡
// 期望：
// 1. 1号玩家收到出牌问询通知，且可以胡
// 2. 1号玩家发送胡请求后，所有玩家收到胡通知， 胡牌者为1号玩家，胡类型为点炮，胡的牌为9W
func Test_Dianpao(t *testing.T) {
	var Int9W uint32 = 19
	params := global.NewCommonStartGameParams()

	params.BankerSeat = 0
	huSeat := 1
	bankerSeat := params.BankerSeat
	// 庄家的最后一张牌改为 9W
	params.Cards[bankerSeat][13] = 19
	// 1 号玩家最后1张牌改为 9W
	params.Cards[huSeat][12] = 19

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	// 庄家出 9W
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, Int9W))

	// 1 号玩家收到出牌问询通知， 可以胡
	huPlayer := utils.GetDeskPlayerBySeat(huSeat, deskData)
	expector, _ := huPlayer.Expectors[msgId.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.True(t, ntf.GetEnableDianpao())
	assert.True(t, ntf.GetEnableQi())

	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, huSeat))

	// 检测所有玩家收到点炮通知
	utils.CheckHuNotify(t, deskData, []int{huSeat}, bankerSeat, Int9W, room.HuType_HT_DIANPAO)

	// 检测所有玩家收到点炮结算通知
	utils.CheckDianPaoSettleNotify(t, deskData, []int{huSeat}, bankerSeat, Int9W, room.HuType_HT_DIANPAO)
}
