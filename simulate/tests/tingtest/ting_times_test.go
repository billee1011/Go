package tingtest

import (
	msgid "steve/client_pb/msgId"
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
