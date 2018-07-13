package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestDaSanYuan_ZiMo 大三元
//步奏:庄家天胡自摸
//期望:所有玩家收到结算通知

func TestDaSanYuan_ZiMo(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.PlayerSeatGold = map[int]uint64{0: 100000, 1: 100000}
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{11, 12, 13, 19, 19, 45, 45, 45, 46, 46, 46, 47, 47},
		{14, 14, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18},
	}
	params.WallCards = []uint32{47, 41, 41, 41}
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
	utils.CheckHuNotify(t, deskData, []int{banker}, banker, 47, room.HuType_HT_TIANHU)
	// 检查番结算 天胡 88,无花 4,大三元 88,混一色 6,全带幺 4,三暗刻 16,小于五 88 = 206
	utils.CheckFanSettle(t, deskData, params.GameID, banker, 206, room.FanType_FT_DASANYUAN)
}
