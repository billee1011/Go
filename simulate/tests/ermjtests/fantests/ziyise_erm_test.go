package fantests

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ziyise(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_ERRENMJ // 二人
	params.PeiPaiGame = "ermj"
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		//{41, 41, 46, 45, 45, 44, 44, 43, 43, 47, 47, 42, 42},
		{42, 42, 42, 45, 45, 45, 44, 44, 44, 47, 47, 47, 46},
		{11, 11, 16, 15, 15, 15, 15, 13, 13, 17, 17, 12, 12},
	}
	params.WallCards = []uint32{16, 11, 46}
	params.IsHsz = false
	params.IsDq = false
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	return deskData
}

//TestFan_ZiYiSe_Zimo_ERM 子一色立即结算自摸测试
// 庄摸牌46,自摸
//期望赢分：201 = [64（子一色） + 64（四暗刻） + 64（小三元）+ 4（不求人） +4（无花）+1（单钓将）]* 1
func TestFan_ZiYiSe_Zimo_ERM(t *testing.T) {
	deskData := ziyise(t)
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家出16
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 16))
	//1玩家能碰,能胡13
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, false, true, false))
	//开局 1 弃
	assert.Nil(t, utils.SendQiReq(deskData, 1))
	//开局 1 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1玩家出11
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 11))
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家摸到46,自摸
	assert.Nil(t, utils.SendHuReq(deskData, 0))

	// 检测分数
	utils.CheckFanSettle(t, deskData, 4, 0, 209, deskData.DiFen, room.FanType_FT_ZIYISE)
}
