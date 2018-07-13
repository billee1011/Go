package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test_laoShaoFu_ZiMo 老少副
//步奏:庄家天胡自摸
//期望:所有玩家收到结算通知

func Test_laoShaoFu_ZiMo(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.PlayerSeatGold = map[int]uint64{0: 100000, 1: 100000}
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{11, 12, 13, 17, 18, 19, 42, 42, 42, 43, 43, 43, 44},
		{14, 14, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18},
	}
	params.WallCards = []uint32{44, 41, 41, 41}
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
	utils.CheckHuNotify(t, deskData, []int{banker}, banker, 44, room.HuType_HT_TIANHU)
	// 检查番结算 天胡 88,无花 4,老少副 1,全带幺 4,一般高 1,混一色 6,双暗刻 2,小三风 6 = 112
	utils.CheckFanSettle(t, deskData, params.GameID, banker, 112, room.FanType_FT_LAOSHAOFU)
}
