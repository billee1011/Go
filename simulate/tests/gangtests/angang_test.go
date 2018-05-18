package gangtests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Angang(t *testing.T) {
	// utils.StartGameParams
	thisParams := global.NewCommonStartGameParams()
	thisParams.WallCards = append(thisParams.WallCards, &global.Card9B)
	deskData, err := utils.StartGame(thisParams)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	//庄家出牌
	assert.Nil(t, utils.SendChupaiReq(deskData, deskData.BankerSeat, uint32(13)))
	//所有客户端接受出牌通知
	utils.CheckChuPaiNotify(t, deskData, uint32(13), utils.GetDeskPlayerBySeat(deskData.BankerSeat, deskData).Player.GetID())
	//下家这时候摸到牌后，进入自询状态，自询状态下可以暗杠
	xjPlayer := utils.CheckMoPaiNotify(t, deskData, (deskData.BankerSeat+1)%len(deskData.Players))
	//下家请求暗杠
	utils.SendGangReq(deskData, xjPlayer.Seat, uint32(16), room.GangType_AnGang)
	//检查下家暗杠的通知
	utils.CheckGangNotify(t, deskData, xjPlayer.Player.GetID(), xjPlayer.Player.GetID(), uint32(16), room.GangType_AnGang)

}
