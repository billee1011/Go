package tingtest

import (
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_Ting_times TAPD缺陷 :id 1009888(问题描述:天听时胡牌提示倍数错误)
// 流程:　玩家起手听牌　，期待：打出1t,可以听胡的牌除了１ｗ和９ｗ时８倍，其他都是４倍
func Test_Ting_times(t *testing.T) {
	thisParams := global.NewCommonStartGameParams()
	thisParams.Cards[0] = []uint32{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 33, 33, 34, 21}
	thisParams.Cards[1] = []uint32{21, 21, 21, 22, 23, 24, 25, 26, 27, 28, 37, 37, 36}
	thisParams.Cards[2] = []uint32{31, 31, 31, 31, 32, 32, 32, 32, 33, 33, 19, 19, 19}
	thisParams.Cards[3] = []uint32{39, 39, 39, 39, 38, 38, 38, 38, 37, 37, 29, 29, 29}
	thisParams.WallCards = []uint32{35, 16, 15, 35}
	thisParams.DingqueColor[0] = room.CardColor_CC_TONG
	thisParams.HszDir = room.Direction_Opposite
	thisParams.HszCards = [][]uint32{
		{33, 33, 34},
		{37, 37, 36},
		{19, 19, 19},
		{29, 29, 29},
	}
	deskData, err := utils.StartGame(thisParams)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	//1w=8 9w=8 其他4
	zjDeskPlayer := utils.GetDeskPlayerBySeat(deskData.BankerSeat, deskData)
	zixunExpectors := zjDeskPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	zixunNtf := room.RoomZixunNtf{}
	zixunExpectors.Recv(2*time.Second, &zixunNtf)
	assert.Equal(t, 4, len(zixunNtf.GetCanTingCardInfo()))
	for _, cantingCardInfo := range zixunNtf.GetCanTingCardInfo() {
		if uint32(21) == cantingCardInfo.GetOutCard() {
			for _, tingCards := range cantingCardInfo.GetTingCardInfo() {
				if tingCards.GetTingCard() == uint32(11) || tingCards.GetTingCard() == uint32(19) {
					assert.Equal(t, uint32(8), tingCards.GetTimes())
				} else {
					assert.Equal(t, uint32(4), tingCards.GetTimes())
				}
			}
		}
	}
}

// Test_ChuPaiwenxun_Actions 玩家可以有的操作模拟测试
// 流程: 庄家天胡1w,下家摸5b,并且打出7w,检测玩家可以有的操作
// 期待: 庄家只能胡牌,对家可胡可弃
func Test_ChuPaiwenxun_Actions(t *testing.T) {
	thisParams := global.NewCommonStartGameParams()
	thisParams.Cards[0] = []uint32{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 36, 37, 38}
	thisParams.Cards[1] = []uint32{23, 23, 24, 12, 25, 25, 26, 26, 26, 27, 27, 27, 28}
	thisParams.Cards[2] = []uint32{16, 17, 12, 31, 32, 33, 34, 35, 36, 37, 38, 39, 17}
	thisParams.Cards[3] = []uint32{24, 31, 31, 32, 32, 32, 33, 33, 33, 34, 34, 34, 35}
	thisParams.WallCards = []uint32{35, 16, 15, 35}
	thisParams.DingqueColor[0] = room.CardColor_CC_TONG
	thisParams.DingqueColor[2] = room.CardColor_CC_TIAO
	thisParams.BankerSeat = 0
	thisParams.HszDir = room.Direction_Opposite
	thisParams.HszCards = [][]uint32{
		{36, 37, 38},
		{23, 23, 24},
		{16, 17, 17},
		{31, 32, 33},
	}
	deskData, err := utils.StartGame(thisParams)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	zjDeskPlayer := utils.GetDeskPlayerBySeat(deskData.BankerSeat, deskData)
	zixunExpectors := zjDeskPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	zixunNtf := room.RoomZixunNtf{}
	zixunExpectors.Recv(2*time.Second, &zixunNtf)
	utils.SendHuReq(deskData, deskData.BankerSeat)
	utils.CheckHuNotify(t, deskData, []int{deskData.BankerSeat}, deskData.BankerSeat, uint32(12), room.HuType_HT_TIANHU)
	xjSeat := (deskData.BankerSeat + 1) % len(deskData.Players)
	utils.CheckMoPaiNotify(t, deskData, xjSeat, uint32(35))
	utils.SendChupaiReq(deskData, xjSeat, uint32(12))
	utils.CheckChuPaiNotify(t, deskData, uint32(12), xjSeat)
	utils.WaitChupaiWenxunNtf0(deskData, 0, false, true, false, false)
	utils.WaitChupaiWenxunNtf0(deskData, 2, false, true, false, true)
}
