package tests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_StartGame_Ting_QiDuiHu 测试游戏开始定缺后，闲家能不能听7对牌
// 游戏开始流程包括： 登录，加入房间，配牌，洗牌，发牌，换三张，定缺
func Test_StartGame_Ting_QiDuiHu(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.Cards = [][]uint32{
		{11, 11, 11, 25, 22, 22, 22, 22, 13, 13, 13, 13, 14, 14},
		{15, 15, 15, 11, 26, 26, 26, 26, 17, 17, 17, 17, 28},
		{21, 21, 21, 15, 12, 12, 12, 12, 23, 23, 23, 23, 24},
		{25, 25, 25, 21, 16, 16, 16, 16, 27, 27, 27, 27, 18},
	}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TIAO, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG}
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	//等待听牌信息
	assert.Nil(t, utils.WaitTingInfoNtf(t, deskData, 1, []uint32{28}...))
	assert.Nil(t, utils.WaitTingInfoNtf(t, deskData, 2, []uint32{24}...))
	assert.Nil(t, utils.WaitTingInfoNtf(t, deskData, 3, []uint32{18}...))
}

// Test_StartGame_Ting_TuiDaoHu 测试游戏开始定缺后，闲家能不能听推到胡的牌
// 游戏开始流程包括： 登录，加入房间，配牌，洗牌，发牌，换三张，定缺
func Test_StartGame_Ting_TuiDaoHu(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.Cards = [][]uint32{
		{11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13, 14, 14},
		{15, 15, 15, 15, 16, 17, 31, 32, 33, 34, 35, 36, 18},
		{21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24},
		{25, 25, 25, 25, 26, 26, 26, 26, 27, 27, 27, 27, 28},
	}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	//等待听牌信息
	assert.Nil(t, utils.WaitTingInfoNtf(t, deskData, 1, []uint32{18, 15}...))
}
