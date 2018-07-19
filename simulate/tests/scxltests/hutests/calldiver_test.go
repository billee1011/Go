package hutests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_Dianpao 呼叫转移测试
func Test_CallDIve(t *testing.T) {
	var Int9W uint32 = 19
	params := global.NewCommonStartGameParams()
	params.Cards = [][]uint32{
		{12, 13, 14, 15, 16, 17, 19, 19, 19, 19, 39, 27, 27, 26},
		{21, 21, 21, 24, 22, 22, 22, 22, 23, 23, 33, 33, 34},
		{29, 29, 29, 29, 28, 28, 28, 28, 27, 27, 11, 11, 11},
		{31, 31, 31, 13, 32, 32, 32, 32, 33, 33, 23, 23, 24},
	}
	params.PlayerSeatGold = map[int]uint64{
		0: 10000, 1: 10000, 2: 1, 3: 10000,
	}
	params.HszCards = [][]uint32{
		{27, 27, 26},
		{33, 33, 34},
		{11, 11, 11},
		{23, 23, 24},
	}
	params.WallCards = []uint32{11, 24}
	params.HszDir = room.Direction_Opposite
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TIAO}
	params.BankerSeat = 0
	bankerSeat := params.BankerSeat

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	// 庄家暗杠9W
	assert.Nil(t, utils.SendGangReq(deskData, bankerSeat, Int9W, room.GangType_AnGang))
	assert.Nil(t, utils.SendGangReq(deskData, bankerSeat, 11, room.GangType_AnGang))

	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, 24))

	assert.Nil(t, utils.SendHuReq(deskData, 1))
}
