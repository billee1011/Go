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

// Test_Tianhu 天胡测试
// 开始游戏后，庄家定缺筒，庄家胡牌
// 期望：
// 1. 庄家，即0号玩家，收到自询通知，且可以胡牌
// 2. 0号玩家发送胡请求后，所有玩家收到胡通知， 胡牌者为0号玩家，胡类型为自摸，胡的牌为9W
func Test_Tianhu(t *testing.T) {
	var Int1W uint32 = 11
	params := global.NewCommonStartGameParams()
	params.BankerSeat = 0
	zimoSeat := params.BankerSeat
	bankerSeat := params.BankerSeat
	params.DingqueColor[0] = room.CardColor_CC_TONG
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	// 0 号玩家收到可自摸通知
	bankerPlayer := utils.GetDeskPlayerBySeat(bankerSeat, deskData)
	expector, _ := bankerPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	ntf := room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(time.Second*1, &ntf))
	assert.True(t, ntf.GetEnableZimo())

	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, bankerSeat))

	// 检测所有玩家收到天胡通知
	utils.CheckHuNotify(t, deskData, []int{zimoSeat}, zimoSeat, Int1W, room.HuType_HT_TIANHU)

	// 检测所有玩家收到天胡结算通知
	utils.CheckZiMoSettleNotify(t, deskData, []int{zimoSeat}, zimoSeat, Int1W, room.HuType_HT_TIANHU)

}
