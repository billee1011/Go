package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test_xiaoYuWu_ziMo 小于五
//步奏:庄家天胡自摸
//期望:所有玩家收到结算通知

func Test_xiaoYuWu_ziMo(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = room.GameId_GAMEID_ERRENMJ
	params.PlayerSeatGold = map[int]uint64{0: 100000, 1: 100000}
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{11, 11, 11, 12, 12, 12, 13, 13, 13, 12, 13, 14, 14},
		{17, 18, 19, 16, 17, 18, 17, 18, 19, 17, 18, 19, 19},
	}
	params.WallCards = []uint32{14, 41, 41, 41, 42}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	// 庄家
	banker := deskData.BankerSeat
	//开局 庄家 自询 能自摸
	utils.CheckZixunNtf(t, deskData, banker, false, true, true)
	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, banker))
	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{banker}, banker, 14, room.HuType_HT_TIANHU)
	// 检查番结算 天胡 88,无花 4,小于五 88,清一色 16,四归一 2,三同顺 24,三暗刻 16，平胡 1 = 239
	utils.CheckFanSettle(t, deskData, params.GameID, banker, 239, room.FanType_FT_XIAOYUWU)
}
