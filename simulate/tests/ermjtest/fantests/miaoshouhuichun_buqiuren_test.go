package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_miaoShouHuiChun_buQiuRen_Zimo_ErRen 妙手回春测试 不求人
// 牌墙设置为30 张，开始游戏后，庄家出41，没有人可以碰杠胡。1 号玩家摸 42 出42,庄玩家摸 43 出43,1 号玩家摸 19,自摸19
// 期望：
//1号玩家摸牌后收到自询通知，且可以自摸胡
//1号玩家发送胡请求后，所有玩家收到胡通知， 胡牌者为1号玩家，胡类型为自摸，胡的牌为9W
func Test_miaoShouHuiChun_buQiuRen_Zimo_ErRen(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.PlayerSeatGold = map[int]uint64{0: 100000, 1: 100000}
	params.IsDq = false
	params.IsHsz = false
	zimoSeat := 1
	bankerSeat := params.BankerSeat
	params.Cards = [][]uint32{
		{11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13, 14},
		{15, 15, 15, 15, 16, 16, 16, 16, 17, 17, 17, 17, 19},
	}
	// 牌墙大小设置为1
	params.WallCards = []uint32{41, 42, 43, 19}

	// 传入参数开始游戏
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	// 发送出牌请求，庄家出41
	assert.Nil(t, utils.WaitZixunNtf(deskData, bankerSeat))
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, 41))

	// 发送出牌请求，对家出42
	assert.Nil(t, utils.WaitZixunNtf(deskData, zimoSeat))
	assert.Nil(t, utils.SendChupaiReq(deskData, zimoSeat, 42))

	//发送出牌请求，庄家出43
	assert.Nil(t, utils.WaitZixunNtf(deskData, bankerSeat))
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, 43))

	// 对家收到自询通知
	utils.CheckZixunNtf(t, deskData, zimoSeat, false, false, true)
	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, zimoSeat))
	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{zimoSeat}, zimoSeat, 19, room.HuType_HT_HAIDILAO)
	// 检查番结算 无花 4,三暗刻 16,清一色 16,不求人 4，平胡 1，单吊将 1,四同顺 48,妙手回春 8 = 98
	utils.CheckFanSettle(t, deskData, params.GameID, zimoSeat, 98, room.FanType_FT_MIAOSHOUHUICHUN)

}
