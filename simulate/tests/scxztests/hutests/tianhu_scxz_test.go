package hutests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_SCXZ_TianHu 测试天胡
// 步骤：庄家起手自摸
// 期望：庄家收到自摸，类型为天胡
func Test_SCXZ_TianHu(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.Cards = [][]uint32{
		{11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13, 14, 14},
		{15, 15, 15, 15, 16, 16, 16, 16, 17, 17, 17, 17, 18},
		{21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24},
		{25, 25, 25, 25, 26, 26, 26, 26, 27, 27, 27, 27, 28},
	}
	params.HszCards = [][]uint32{}
	params.PeiPaiGame = "scxz"
	params.IsHsz = false // 不换三张
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	banker := params.BankerSeat
	// 庄家自摸
	assert.Nil(t, utils.WaitZixunNtf(deskData, params.BankerSeat))
	assert.Nil(t, utils.SendHuReq(deskData, banker))
	var Int1W uint32 = 11
	// 检测所有玩家收到天胡通知
	utils.CheckHuNotify(t, deskData, []int{banker}, banker, Int1W, room.HuType_HT_TIANHU)

}
