package fantests

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tianhu(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_ERRENMJ // 二人
	params.PeiPaiGame = "ermj"
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		{41, 41, 46, 45, 45, 44, 44, 43, 43, 47, 47, 42, 42},
		{41, 41, 46, 15, 15, 15, 15, 43, 43, 47, 47, 42, 42},
	}
	params.WallCards = []uint32{46, 11, 16}
	params.IsHsz = false
	params.IsDq = false
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	return deskData
}

//TestFan_TianHu_Zimo_ERM 天胡立即结算自摸测试
// 庄摸牌46,自摸
//期望赢分：180 = [88（大七星）  +4（无花）+ 88（天胡）]* 1
func TestFan_TianHu_Zimo_ERM(t *testing.T) {
	deskData := tianhu(t)
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家摸到46,自摸
	assert.Nil(t, utils.SendHuReq(deskData, 0))

	// 检测分数
	winScro := 180 * (len(deskData.Players) - 1)

	utils.CheckFanSettle(t, deskData, 4, 0, int64(winScro), deskData.DiFen, room.FanType_FT_DAQIXING)
}
