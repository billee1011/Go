package hutests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_Ganghoupao 杠后炮测试
// 0 号玩家暗杠，摸1筒，并且打出9万， 然后1号玩家可胡
// 期望：
// 1. 1号玩家收到出牌问询通知，且可胡
// 2. 1号玩家请求胡，所有玩家收到胡通知，2号玩家摸牌后收到自询通知
func Test_Ganghoupao(t *testing.T) {
	params := global.CommonStartGameParams

	params.BankerSeat = 0
	// 0 号玩家手牌改成 Card1W, Card1W, Card1W, Card1W, Card2W, Card2W, Card2W, Card2W, Card3W, Card3W, Card3W, Card3W, Card4W, Card9W
	// 换三张后手牌为 Card5T, Card5T, Card5T， Card1W, Card2W, Card2W, Card2W, Card2W, Card3W, Card3W, Card3W, Card3W, Card4W, Card9W
	params.Cards[0][13] = &global.Card9W
	// 1 号玩家手牌改成 Card5W, Card5W, Card5W, Card5W, Card6W, Card6W, Card6W, Card6W, Card7W, Card7W, Card7W, Card7W, Card9W
	params.Cards[1][12] = &global.Card9W
	// 1号玩家换三张后手牌为  Card1W, Card1W, Card1W, Card5W, Card5W, Card5W,  Card6W, Card6W, Card6W,  Card7W, Card7W, Card7W, Card9W
	params.HszCards[1] = []*room.Card{&global.Card5W, &global.Card6W, &global.Card7W}

	// 墙牌改为 1筒
	params.WallCards = []*room.Card{&global.Card1B}

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	// 庄家杠 2W
	// assert.Nil(t, utils.WaitZixunNtf(deskData, 0))	// 定缺后的自询通知还没完成
	assert.Nil(t, utils.SendGangReq(deskData, 0, 12, room.GangType_AnGang))
	// 庄家等到自询通知后出 9W
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 19))

	// 1 号玩家收到出牌问询通知， 可以胡
	huPlayer := utils.GetDeskPlayerBySeat(1, deskData)
	expector, _ := huPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(time.Second*1, &ntf))
	assert.True(t, ntf.GetEnableDianpao())
	assert.True(t, ntf.GetEnableQi())

	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 1))
	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{1}, 0, 19, room.HuType_HT_GANGHOUPAO)
}
