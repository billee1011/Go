package fantests

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func renhu(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_ERRENMJ // 二人
	params.PeiPaiGame = "ermj"
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		{41, 41, 46, 15, 15, 15, 15, 43, 43, 47, 47, 42, 42},
		{41, 41, 46, 45, 45, 44, 44, 43, 43, 47, 47, 42, 42},
	}
	params.WallCards = []uint32{46, 11, 16}
	params.IsHsz = false
	params.IsDq = false
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	return deskData
}

//TestFan_RenHu_dianpao_ERM 人胡立即结算自摸测试
//期望赢分：156 = [88（大七星） +4（无花）+ 64（人胡）]* 1
func TestFan_RenHu_dianpao_ERM(t *testing.T) {
	deskData := renhu(t)
	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0玩家出11
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 46))
	// 1玩家
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, false, true, false))
	// 1玩家
	assert.Nil(t, utils.SendHuReq(deskData, 1))
	// 检测分数
	winScro := 156 * (len(deskData.Players) - 1)

	utils.CheckFanSettle(t, deskData, 4, 1, int64(winScro), deskData.DiFen, room.FanType_FT_DAQIXING)
}
