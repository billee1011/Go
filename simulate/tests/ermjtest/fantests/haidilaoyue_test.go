package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_haiDiLaoYue_Zimo_ErRen 海底捞月测试
// 牌墙设置为30 张，开始游戏后，庄家出41，没有人可以碰杠胡。1 号玩家摸 42 出42,庄玩家摸 19出 19,1 号玩家点炮19
// 期望：
//1号玩家点炮胡
//1号玩家发送胡请求后，所有玩家收到胡通知， 胡牌者为1号玩家，来源玩家庄家
func Test_haiDiLaoYue_Zimo_ErRen(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.PlayerSeatGold = map[int]uint64{0: 100000, 1: 100000}
	params.IsDq = false
	params.IsHsz = false
	huSeat := 1
	bankerSeat := params.BankerSeat
	params.Cards = [][]uint32{
		{11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13, 14},
		{15, 15, 15, 15, 16, 16, 16, 16, 17, 17, 17, 17, 19},
	}
	// 牌墙大小设置为1
	params.WallCards = []uint32{41, 42, 19}

	// 传入参数开始游戏
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	// 发送出牌请求，庄家出41
	assert.Nil(t, utils.WaitZixunNtf(deskData, bankerSeat))
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, 41))

	// 发送出牌请求，对家出42
	assert.Nil(t, utils.WaitZixunNtf(deskData, huSeat))
	assert.Nil(t, utils.SendChupaiReq(deskData, huSeat, 42))

	//发送出牌请求，庄家出19
	assert.Nil(t, utils.WaitZixunNtf(deskData, bankerSeat))
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, 19))

	// 对家收到自询通知
	utils.WaitChupaiWenxunNtf(deskData, huSeat, false, true, false)
	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, huSeat))
	// 检测所有玩家收到点炮通知
	utils.CheckHuNotify(t, deskData, []int{huSeat}, bankerSeat, 19, room.HuType_HT_DIANPAO)
	// 检查番结算 无花 4,门前清 2,海底捞月 8,四同顺 48,三暗刻 16,平和 1,清一色 16 = 95
	utils.CheckFanSettle(t, deskData, params.GameID, huSeat, 95, room.FanType_FT_HAIDILAOYUE)

}
