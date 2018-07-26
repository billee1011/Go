package ermjtest

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHu(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.PlayerNum = 2
	params.BankerSeat = 0
	params.PeiPaiGame = "ermj"
	params.GameID = common.GameId_GAMEID_ERRENMJ
	params.IsDq = false
	params.IsHsz = false
	params.Cards = [][]uint32{
		{11, 11, 11, 51, 52, 14, 12, 12, 13, 13, 13, 14, 14},
		{53, 54, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18},
	}
	params.WallCards = []uint32{11, 55, 12, 56, 13, 18, 57, 58, 12, 41}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	utils.CheckZixunNtfWithTing(t, deskData, 0, false, true, true, true)
	assert.Nil(t, utils.SendHuReq(deskData, 0))
	utils.CheckHuNotifyBySeats(t, deskData, []int{0}, 0, uint32(12), room.HuType_HT_TIANHU, []int{0, 1})
	// time.Sleep(time.Second * 2)
	// utils.SendChupaiReq(deskData, 1, uint32(18))
}
