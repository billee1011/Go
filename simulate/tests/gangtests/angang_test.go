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
	thisParams.WallCards = append(thisParams.WallCards, 39)
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
	thisParams.Cards[0] = []uint32{11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 26, 27, 28, 29}
	thisParams.Cards[1] = []uint32{29, 29, 31, 31, 32, 32, 32, 32, 33, 33, 37, 37, 36}
	thisParams.Cards[2] = []uint32{21, 21, 22, 22, 23, 23, 23, 23, 24, 25, 13, 13, 14}
	thisParams.Cards[3] = []uint32{39, 39, 39, 39, 38, 38, 38, 38, 37, 37, 33, 33, 34}
	thisParams.WallCards = []uint32{26, 36, 14, 29}

	thisParams.HszDir = room.Direction_Opposite
	thisParams.HszCards = [][]uint32{
		{26, 27, 28},
		{32, 31, 31},
		{13, 13, 14},
		{39, 39, 39},
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
