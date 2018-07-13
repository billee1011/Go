package qitests

import (
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Peng_Qi 测试自询碰弃
// 期望：
// 庄家出5W后，1号玩家将收到出牌问询通知，可碰5W
// 1号玩家发出弃碰请求后1号玩家摸牌，1号玩家将收到自询通知
func Test_SCXZ_Peng_Qi(t *testing.T) {
	var Int5w uint32 = 15
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	// 0 号玩家的最后一张牌改成 5W， 打出后 1 号玩家可碰
	params.Cards[0][13] = 15
	params.Cards[1][0] = 14
	// 修改换三张的牌
	params.HszCards = [][]uint32{
		{13, 13, 13},
		{17, 17, 17},
		{23, 23, 23},
		{27, 27, 27},
	}

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	//bankerSeat玩家收到可自询通知
	utils.WaitZixunNtf(deskData, params.BankerSeat)
	// 庄家出 5W
	assert.Nil(t, utils.SendChupaiReq(deskData, params.BankerSeat, Int5w))

	// 1 号玩家收到出牌问询通知，可以碰
	pengSeat := (params.BankerSeat + 1) % len(deskData.Players)
	pengPlayer := utils.GetDeskPlayerBySeat(pengSeat, deskData)

	expector, _ := pengPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.Equal(t, Int5w, ntf.GetCard())
	assert.True(t, ntf.GetEnablePeng())
	assert.True(t, ntf.GetEnableQi())

	// 发送弃请求
	assert.Nil(t, utils.SendQiReq(deskData, pengSeat))

	// 1 号玩家收到可自询通知
	utils.WaitZixunNtf(deskData, 1)
}
