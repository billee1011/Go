package qitests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_qi_peng 测试胡过的玩家不能摸牌
// 庄家自摸
// 下家对家摸牌打牌，尾家摸牌打出18
// 此时下家可以胡
// 下家弃胡，因为庄家胡过，所以期待摸牌玩家是下家
func Test_qi_peng(t *testing.T) {
	param := global.NewCommonStartGameParams()
	param.GameID = room.GameId_GAMEID_XUEZHAN
	param.PeiPaiGame = "scxz"
	param.IsHsz = false
	param.Cards[3][12] = 18
	param.WallCards = []uint32{31, 31, 31, 31}
	data, err := utils.StartGame(param)
	assert.Nil(t, err)
	utils.CheckZixunNotify(t, data, 0)
	assert.Nil(t, utils.SendHuReq(data, 0))
	utils.CheckHuNotify(t, data, []int{0}, 0, 14, room.HuType_HT_TIANHU)
	utils.CheckMoPaiNotify(t, data, 1, 31)
	utils.SendChupaiReq(data, 1, 31)
	utils.CheckMoPaiNotify(t, data, 2, 31)
	utils.SendChupaiReq(data, 2, 31)
	utils.CheckMoPaiNotify(t, data, 3, 31)
	utils.SendChupaiReq(data, 3, 18)
	assert.Nil(t, utils.WaitChupaiWenxunNtf(data, 1, false, true, false))
	assert.Nil(t, utils.SendQiReq(data, 1))
	utils.CheckMoPaiNotify(t, data, 1, 31)

}
