package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test_siZiKe_shuanJianKe_ziMo 四字刻 双箭刻
//步奏:庄家天胡自摸
//期望:所有玩家收到结算通知

func Test_siZiKe_shuanJianKe_ziMo(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.PlayerSeatGold = map[int]uint64{0: 100000, 1: 100000}
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{41, 41, 41, 43, 43, 43, 45, 45, 45, 46, 46, 46, 11},
		{17, 18, 12, 16, 17, 18, 17, 18, 12, 17, 18, 12, 12},
	}
	params.WallCards = []uint32{11, 42, 42, 24}
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
	utils.CheckHuNotify(t, deskData, []int{banker}, banker, 11, room.HuType_HT_TIANHU)
	// 检查番结算 天胡 88,无花 4,四字刻 24,混一色 6,全带幺 4,四暗刻 64,双箭刻 6,圈风刻 2,门风刻 2= 200
	utils.CheckFanSettle(t, deskData, params.GameID, banker, 200, room.FanType_FT_SIZIKE)
}
