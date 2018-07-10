package ermjtest

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestChi 测试吃牌，庄家出14，闲家可吃13,14,15或者14,15,16
//闲家选择吃14,15,16
func TestChi(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{11, 11, 11, 51, 52, 12, 12, 12, 13, 13, 13, 14, 14},
		{53, 54, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18},
	}
	params.WallCards = []uint32{11, 55, 12, 56, 13, 14, 57, 58, 14, 19, 19, 19, 41, 41, 41}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	utils.CheckZixunNtfWithTing(t, deskData, 0, false, true, true, true)
	//等補花結束
	// time.Sleep(time.Second * 2)
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, uint32(14)))
	utils.CheckChuPaiNotifyWithSeats(t, deskData, uint32(14), 0, []int{0, 1})
	//出牌问询检查
	chiPlayer := utils.GetDeskPlayerBySeat(1, deskData)
	expector, _ := chiPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.Equal(t, []uint32{13, 14, 15, 14, 15, 16}, ntf.GetChiInfo().GetCards())
	assert.False(t, ntf.GetEnableDianpao())
	assert.True(t, ntf.GetEnableQi())
	//请求吃牌
	assert.Nil(t, utils.SendChiReq(t, deskData, []uint32{14, 15, 16}, 1))
	//检查吃牌广播
	utils.CheckChiNtfWithSeats(t, deskData, []uint32{14, 15, 16}, 1, 0, []int{0, 1})
}
