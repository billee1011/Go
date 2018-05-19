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

// Test_Angang 暗杠测试用例1
// 游戏流程:  1.定缺完成后,庄家自询可以暗杠2w和3w,庄家选择出牌,出3w(期待:收到庄家的出牌通知)
// 		    2.下家摸牌,进入自询问状态,可暗杠6w和7w,下家选择杠6w(期待:收到下家的杠牌通知)
func Test_Angang(t *testing.T) {
	// utils.StartGameParams
	thisParams := global.NewCommonStartGameParams()
	thisParams.WallCards = append(thisParams.WallCards, &global.Card9B)
	deskData, err := utils.StartGame(thisParams)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	zjPlayer := utils.GetDeskPlayerBySeat(deskData.BankerSeat, deskData)
	expector, _ := zjPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	zixunNtf := room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(3*time.Second, &zixunNtf))
	assert.Subset(t, zixunNtf.GetEnableAngangCards(), []uint32{12, 13})
	//庄家出牌
	assert.Nil(t, utils.SendChupaiReq(deskData, deskData.BankerSeat, uint32(13)))
	//所有客户端接受出牌通知
	utils.CheckChuPaiNotify(t, deskData, uint32(13), deskData.BankerSeat)
	//下家这时候摸到牌后，进入自询状态，自询状态下可以暗杠
	xjPlayer := utils.CheckMoPaiNotify(t, deskData, (deskData.BankerSeat+1)%len(deskData.Players), 31)
	//检查自询通知
	xjexpector, _ := xjPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	xjzixunNtf := room.RoomZixunNtf{}
	assert.Nil(t, xjexpector.Recv(3*time.Second, &xjzixunNtf))
	assert.Subset(t, xjzixunNtf.GetEnableAngangCards(), []uint32{16, 17})
	//下家请求暗杠
	utils.SendGangReq(deskData, xjPlayer.Seat, uint32(16), room.GangType_AnGang)
	//检查下家暗杠的通知
	utils.CheckGangNotify(t, deskData, xjPlayer.Player.GetID(), xjPlayer.Player.GetID(), uint32(16), room.GangType_AnGang)

}

// Test_Angang1 暗杠测试2,测试庄家暗杠
// 游戏流程 : 1.定缺完成后,庄家收到自询通知,期待:可以暗杠 1万,2万,3万
//			2.庄家请求暗杠1w,期待:收到暗杠1w的杠通知
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
	//庄家请求暗杠
	utils.SendGangReq(deskData, gangSeat, uint32(11), room.GangType_AnGang)
	//检查庄家暗杠的通知
	utils.CheckGangNotify(t, deskData, gangPlayer.Player.GetID(), gangPlayer.Player.GetID(), uint32(11), room.GangType_AnGang)
}
