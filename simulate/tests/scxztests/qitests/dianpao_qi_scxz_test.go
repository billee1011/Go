package qitests

import (
	"steve/client_pb/common"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Dianpao_qi 点炮弃测试
// 开始游戏后，庄家出9W，1号玩家可以胡，其他玩家都不可以胡
// 期望：
// 1号玩家点弃，由于1号玩家是庄家的下家，1号玩家将收到自询通知
func Test_SCXZ_Dianpao_qi(t *testing.T) {
	var Int9W uint32 = 19
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
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
	//bankerSeat玩家收到可自询通知
	utils.WaitZixunNtf(deskData, bankerSeat)
	// 庄家出 9W
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, Int9W))

	// 1 号玩家收到出牌问询通知， 可以胡
	huPlayer := utils.GetDeskPlayerBySeat(huSeat, deskData)
	expector, _ := huPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.True(t, ntf.GetEnableDianpao())
	assert.True(t, ntf.GetEnableQi())

	// 发送弃请求
	assert.Nil(t, utils.SendQiReq(deskData, huSeat))

	// huSeat玩家收到可自询通知
	utils.WaitZixunNtf(deskData, huSeat)
}
