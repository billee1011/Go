package fantests

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func gangkai(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_ERRENMJ // 二人
	params.PeiPaiGame = "ermj"
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		{41, 41, 41, 45, 45, 45, 43, 43, 43, 47, 47, 42, 42},
		{11, 11, 46, 15, 15, 15, 15, 13, 13, 17, 17, 12, 12},
	}
	params.WallCards = []uint32{41, 42, 42}
	params.IsHsz = false
	params.IsDq = false
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	return deskData
}

//TestFan_TianHu_Zimo_ERM 天胡立即结算自摸测试
// 庄摸牌46,自摸
//期望赢分：280 = [64（子一色） +64（四暗刻）+2（暗杠）+不求人（4）+ 箭刻（2） +4（无花）]* 2(杠开)* 1
func TestFan_GangKai_Zimo_ERM(t *testing.T) {
	deskData := gangkai(t)
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家杠41
	assert.Nil(t, utils.SendGangReq(deskData, 0, 41, room.GangType_AnGang))
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家自摸
	assert.Nil(t, utils.SendHuReq(deskData, 0))
	// 检测分数
	winScro := 280 * (len(deskData.Players) - 1)

	utils.CheckFanSettle(t, deskData, 4, 0, int64(winScro), room.FanType_FT_ZIYISE)
}
