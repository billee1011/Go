package ermjtest

import (
	"steve/client_pb/common"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStartBuhua 开局补花，补到没花牌位置
// 0、发牌，庄家闲家各13张牌
// 1、庄家闲家同时亮全部花牌，不补牌
// 2、庄家开始补牌，一次性补全
// 3、庄家补上的牌有一张花牌，继续补
// 4、庄家没有花牌后，闲家接着补，闲家步骤和庄家一样
// 5、所有人花牌补完后，庄家补摸一张牌
func TestStartBuhua(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = common.GameId_GAMEID_ERRENMJ
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{11, 11, 11, 51, 52, 12, 12, 12, 13, 13, 13, 14, 14},
		{53, 54, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18},
	}
	params.WallCards = []uint32{11, 55, 12, 56, 13, 14, 57, 58, 14, 19, 19, 19, 41, 41, 41}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	utils.CheckBuhuaNtf(t, []int{0, 1}, [][]uint32{{51, 52}, {53, 54}}, [][]uint32{nil, nil}, []int{0, 1}, deskData)
	utils.CheckBuhuaNtf(t, []int{0}, [][]uint32{nil}, [][]uint32{{11, 55}}, []int{0, 1}, deskData)
	utils.CheckBuhuaNtf(t, []int{0}, [][]uint32{{55}}, [][]uint32{{12}}, []int{0, 1}, deskData)
	utils.CheckBuhuaNtf(t, []int{1}, [][]uint32{nil}, [][]uint32{{56, 13}}, []int{0, 1}, deskData)
	utils.CheckBuhuaNtf(t, []int{1}, [][]uint32{{56}}, [][]uint32{{14}}, []int{0, 1}, deskData)
	utils.CheckMoPaiNotify(t, deskData, 0, 57)
	utils.CheckBuhuaNtf(t, []int{0}, [][]uint32{{57}}, [][]uint32{{58}}, []int{0, 1}, deskData)
	utils.CheckBuhuaNtf(t, []int{0}, [][]uint32{{58}}, [][]uint32{{14}}, []int{0, 1}, deskData)
	utils.CheckZixunNotify(t, deskData, 0)
}
