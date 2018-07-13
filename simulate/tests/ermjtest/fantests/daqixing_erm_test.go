package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func daqixing(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_ERRENMJ // 二人
	params.PeiPaiGame = "ermj"
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		{41, 41, 46, 45, 45, 44, 44, 43, 43, 47, 47, 42, 42},
		{41, 41, 46, 15, 15, 15, 15, 43, 43, 47, 47, 42, 42},
	}
	params.WallCards = []uint32{16, 11, 46}
	params.IsHsz = false
	params.IsDq = false
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	return deskData
}

//TestFan_DAQIXING_Zimo_ERM 大七星立即结算自摸测试
// 庄摸牌46,自摸
//期望赢分：96 = [88（大七星） + 4（不求人） +4（无花）]* 1
func TestFan_DAQIXING_Zimo_ERM(t *testing.T) {
	deskData := daqixing(t)
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家出16
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 16))
	//开局 1 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1玩家出11
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 11))
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家摸到46,自摸
	assert.Nil(t, utils.SendHuReq(deskData, 0))

	// 检测分数
	winScro := 96 * (len(deskData.Players) - 1)

	utils.CheckFanSettle(t, deskData, 4, 0, int64(winScro), room.FanType_FT_DAQIXING)
}
