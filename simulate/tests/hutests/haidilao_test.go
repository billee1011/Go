package hutests

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

// Test_Zimo_Haidilao 海底捞测试
// 开始游戏后，庄家出1筒，没有人可以碰杠胡。1 号玩家摸 9W 且可以自摸，和是最后一张
// 期望：
//1号玩家摸牌后收到自询通知，且可以自摸胡
//1号玩家发送胡请求后，所有玩家收到胡通知， 胡牌者为1号玩家，胡类型为自摸，结算类型为海底捞，胡的牌为9W
func Test_Zimo_Haidilao(t *testing.T) {
	var Int1B uint32 = 31
	var Int9W uint32 = 19
	params := global.NewCommonStartGameParams()
	params.BankerSeat = 0
	zimoSeat := 1
	bankerSeat := params.BankerSeat

	// 庄家的最后一张牌改为 1B
	params.Cards[bankerSeat][13] = &global.Card1B
	// 1 号玩家最后1张牌改为 9W
	params.Cards[zimoSeat][12] = &global.Card9W
	// 牌墙大小设置为1
	params.WallCards = make([]*room.Card, 1)
	// 墙牌改成 9W
	params.WallCards[0] = &global.Card9W

	// 传入参数开始游戏
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	// 发送出牌请求，庄家出1筒
	assert.Nil(t, utils.WaitZixunNtf(deskData, bankerSeat))
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, Int1B))

	// 根据座号获取自摸玩家
	zimoPlayer := utils.GetDeskPlayerBySeat(zimoSeat, deskData)
	// 1 号玩家期望收到自询通知
	zixunexpector, _ := zimoPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	ntf := room.RoomZixunNtf{}
	// 1秒内接收到自询通知，并赋值到ntf
	assert.Nil(t, zixunexpector.Recv(global.DefaultWaitMessageTime, &ntf))
	// 判断自询通知中是否有自摸
	assert.True(t, ntf.GetEnableZimo())

	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, zimoSeat))

	// 检测所有玩家收到（海底捞算自摸）自摸结算通知,自摸-清一色-2根 = 2 * 4 *4 = 32
	winScro := 32 * (len(deskData.Players) - 1)
	checkZiMoSettleScoreNotify(t, deskData, zimoSeat, int64(winScro))

	// 检测所有玩家收到海底捞胡类型通知
	utils.CheckHuNotify(t, deskData, []int{zimoSeat}, zimoSeat, Int9W, room.HuType_HT_HAIDILAO)
}

// checkZiMoSettleScoreNotify 检查自摸分数结算通知
func checkZiMoSettleScoreNotify(t *testing.T, deskData *utils.DeskData, zimoSeat int, winScore int64) {
	zimoplayer := utils.GetDeskPlayerBySeat(zimoSeat, deskData)
	zimoID := zimoplayer.Player.GetID()
	expector, _ := zimoplayer.Expectors[msgid.MsgID_ROOM_INSTANT_SETTLE]
	ntf := room.RoomSettleInstantRsp{}
	expector.Recv(global.DefaultWaitMessageTime, &ntf)
	assert.Equal(t, len(deskData.Players), len(ntf.BillPlayersInfo))
	for _, billInfo := range ntf.BillPlayersInfo {
		if billInfo.GetPid() == zimoID {
			assert.Equal(t, billInfo.GetScore(), winScore)
		} else {
			assert.Equal(t, billInfo.GetScore(), -(winScore / 3))
		}
		fmt.Println(billInfo)
	}
}
