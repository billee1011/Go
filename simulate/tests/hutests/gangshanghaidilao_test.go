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

// Test_Zimo_Gangshanghaidilao 杠上海底捞测试
// 开始游戏后，庄家出1B，庄家的下家明杠1B，庄家的下家摸牌9w且自摸9w
// 期望：
//1号玩家将收到出牌问询通知, 通知中的可选行为包括杠
//1号玩家请求执行杠后， 所有玩家收到杠通知消息, 其中杠的牌为 1B， 杠的玩家为 1 号玩家， 杠类型为明杠
//1好玩家将将收到摸牌通知，所有玩家收到1号玩家摸牌通知,摸的牌为9W，且是最后一张
//1号玩家收到自询通知，且可以自摸胡
//1号玩家发送胡请求后，所有玩家收到胡通知， 胡牌者为1号玩家，胡类型为自摸，结算类型为杠上海底捞，胡的牌为9W
func Test_Zimo_Gangshanghaidilao(t *testing.T) {
	var Int1B uint32 = 31

	params := global.CommonStartGameParams

	params.BankerSeat = 0
	zimoSeat := 1
	bankerSeat := params.BankerSeat

	// 庄家的最后一张牌改为 1B
	params.Cards[bankerSeat][13] = &global.Card1B
	// 1 号玩家最后1张牌改为 9W
	params.Cards[zimoSeat][12] = &global.Card9W
	// 1 号玩家修改手牌，5w,6w,7w各一个改1B
	params.Cards[zimoSeat][3] = &global.Card1B
	params.Cards[zimoSeat][7] = &global.Card1B
	params.Cards[zimoSeat][11] = &global.Card1B

	// 牌墙大小设置为1
	params.WallCards = make([]*room.Card, 1)
	// 墙牌改成 9W
	params.WallCards[0] = &global.Card9W

	// 传入参数开始游戏
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	// 发送出牌请求，庄家出1筒
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, Int1B))
	// 根据座号获取自摸玩家
	zimoPlayer := utils.GetDeskPlayerBySeat(zimoSeat, deskData)

	// 自摸玩家期望接收到出牌问询通知
	chupaiexpector, _ := zimoPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, chupaiexpector.Recv(time.Second*2, &ntf))
	assert.Equal(t, Int1B, ntf.GetCard())
	assert.True(t, ntf.GetEnableMinggang())
	assert.True(t, ntf.GetEnableQi())

	// 发送杠请求
	assert.Nil(t, utils.SendGangReq(deskData, zimoSeat, Int1B, room.GangType_MingGang))

	// 检测所有玩家收到杠通知
	bankerPlayer := utils.GetDeskPlayerBySeat(bankerSeat, deskData)
	utils.CheckGangNotify(t, deskData, zimoPlayer.Player.GetID(), bankerPlayer.Player.GetID(), Int1B, room.GangType_MingGang)

	var Int9W uint32 = 19
	// 所有玩家收到1号玩家摸牌通知
	for _, deskPlayer := range deskData.Players {
		mopaiexpector, _ := deskPlayer.Expectors[msgid.MsgID_ROOM_MOPAI_NTF]
		ntf := room.RoomMopaiNtf{}
		assert.Nil(t, mopaiexpector.Recv(time.Second*2, &ntf))
		assert.Equal(t, zimoPlayer.Player.GetID(), ntf.GetPlayer())
		assert.Equal(t, false, ntf.GetBack())
		if deskPlayer.Seat == zimoSeat {
			assert.Equal(t, Int9W, ntf.GetCard()) // 摸了一张 9W
		}
	}

	// 1 号玩家期望收到自询通知
	zimoexpector, _ := zimoPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	zixunNtf := room.RoomZixunNtf{}
	// 2秒内接收到自询通知，并赋值到ntf
	assert.Nil(t, zimoexpector.Recv(time.Second*2, &zixunNtf))
	// 判断自询通知中是否有自摸
	assert.True(t, zixunNtf.GetEnableZimo())

	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, zimoSeat))

	// 检测所有玩家收到自摸海底捞通知
	utils.CheckHuNotify(t, deskData, []int{zimoSeat}, zimoSeat, Int9W, room.HuType_GangShangHaiDiLao)
}
