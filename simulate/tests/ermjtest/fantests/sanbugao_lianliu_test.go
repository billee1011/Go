package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test_sanBuGao_lianLiu_ziMo 三步高 连六
//步奏:庄家天胡自摸
//期望:所有玩家收到结算通知

func Test_sanBuGao_ziMo(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.PlayerSeatGold = map[int]uint64{0: 100000, 1: 100000}
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{13, 14, 15, 15, 16, 17, 17, 18, 19, 41, 41, 43, 43},
		{17, 18, 19, 16, 17, 18, 43, 18, 19, 45, 44, 19, 44},
	}
	params.WallCards = []uint32{43, 42, 42, 42}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	// 庄家
	banker := deskData.BankerSeat
	//开局 庄家 自询 能自摸
	utils.CheckZixunNtf(t, deskData, banker, false, false, true)
	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, banker))
	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{banker}, banker, 43, room.HuType_HT_TIANHU)
	// 检查番结算 天胡 88,无花 4,三步高 16,连六 1,一般高 1，混一色 6 = 116
	utils.CheckFanSettle(t, deskData, params.GameID, banker, 116, room.FanType_FT_SANBUGAO)
}
