package ermjtest

import (
	"steve/client_pb/common"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameOver(t *testing.T) {
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
	params.WallCards = []uint32{11, 55, 12, 56, 13, 14, 57, 58, 14, 41}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	utils.CheckZixunNtfWithTing(t, deskData, 0, false, true, true, true)
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, uint32(14)))
	utils.CheckChuPaiNotifyWithSeats(t, deskData, uint32(14), 0, []int{0, 1})
	assert.Nil(t, utils.SendQiReq(deskData, 1))
	utils.WaitZixunNtf(deskData, 1)
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, uint32(18)))
	utils.CheckChuPaiNotifyWithSeats(t, deskData, uint32(18), 1, []int{0, 1})
	utils.WaitGameOverNtf(t, deskData)
}
