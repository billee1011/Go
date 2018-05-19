package gangtests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

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
	utils.CheckChuPaiNotify(t, deskData, uint32(13), deskData.BankerSeat)
	//下家这时候摸到牌后，进入自询状态，自询状态下可以暗杠
	xjPlayer := utils.CheckMoPaiNotify(t, deskData, (deskData.BankerSeat+1)%len(deskData.Players), 31)
	//下家请求暗杠
	utils.SendGangReq(deskData, xjPlayer.Seat, uint32(16), room.GangType_AnGang)
	//检查下家暗杠的通知
	utils.CheckGangNotify(t, deskData, xjPlayer.Player.GetID(), xjPlayer.Player.GetID(), uint32(16), room.GangType_AnGang)

}

func Test_Angang1(t *testing.T) {
	// utils.StartGameParams
	thisParams := global.NewCommonStartGameParams()
	thisParams.Cards[0] = utils.MakeRoomCards(global.Card1W, global.Card1W, global.Card1W, global.Card1W, global.Card2W,
		global.Card2W, global.Card2W, global.Card2W, global.Card3W, global.Card3W, global.Card6T, global.Card7T, global.Card8T, global.Card9T)
	thisParams.Cards[1] = utils.MakeRoomCards(global.Card9T, global.Card9T, global.Card1B, global.Card1B, global.Card2B,
		global.Card2B, global.Card2B, global.Card2B, global.Card3B, global.Card3B, global.Card7B, global.Card7B, global.Card6B)
	thisParams.Cards[2] = utils.MakeRoomCards(global.Card1T, global.Card1T, global.Card2T, global.Card2T, global.Card3T,
		global.Card3T, global.Card3T, global.Card3T, global.Card4T, global.Card5T, global.Card3W, global.Card3W, global.Card4W)
	thisParams.Cards[3] = utils.MakeRoomCards(global.Card9B, global.Card9B, global.Card9B, global.Card9B, global.Card8B,
		global.Card8B, global.Card8B, global.Card8B, global.Card7B, global.Card7B, global.Card3B, global.Card3B, global.Card4B)
	thisParams.WallCards = append(thisParams.WallCards, &global.Card6T, &global.Card6B, &global.Card4W, &global.Card9T)
	thisParams.HszDir = room.Direction_Opposite
	thisParams.HszCards = [][]*room.Card{
		utils.MakeRoomCards(global.Card6T, global.Card7T, global.Card8T),
		utils.MakeRoomCards(global.Card2B, global.Card1B, global.Card1B),
		utils.MakeRoomCards(global.Card3W, global.Card3W, global.Card4W),
		utils.MakeRoomCards(global.Card9B, global.Card9B, global.Card9B),
	}
	deskData, err := utils.StartGame(thisParams)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	gangSeat := deskData.BankerSeat
	// 收到自询通知,可以暗杠 1万,2万,3万
	gangPlayer := utils.GetDeskPlayerBySeat(gangSeat, deskData)
	expector, _ := gangPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	zixunNtf := room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(3*time.Second, &zixunNtf))
	assert.Subset(t, zixunNtf.GetEnableAngangCards(), []uint32{11, 12, 13})
	//下家请求暗杠
	utils.SendGangReq(deskData, gangSeat, uint32(11), room.GangType_AnGang)
	//检查下家暗杠的通知
	utils.CheckGangNotify(t, deskData, gangPlayer.Player.GetID(), gangPlayer.Player.GetID(), uint32(11), room.GangType_AnGang)
}
