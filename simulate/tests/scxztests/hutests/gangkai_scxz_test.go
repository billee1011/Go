package hutests

import (
	"steve/client_pb/common"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Gangkai 杠上开花测试， 杠牌后自摸，且墙牌还有
// 开始游戏后，seat0庄家出1W，seat1玩家杠，sea1玩家摸 8W，胡牌
// 期望：
// 1. 1号玩家可以收到 ROOM_CHUPAIWENXUN_NTF，并且 明杠 enable
// 2. 1号玩家发送杠请求后，收到 ROOM_MOPAI_NTF，
// 3. 1号玩家收到 MsgID_ROOM_ZIXUN_NTF， 并且 胡 enable
// 4. 1号玩家发送 ROOM_XINGPAI_ACTION_REQ 胡
// 5. 所有玩家收到胡牌通知 hutype=HuType_HT_GANGKAI
func Test_SCXZ_Gangkai(t *testing.T) {
	var Int1W uint32 = 11
	var Int8W uint32 = 18
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	gangSeat := 1
	huSeat := 1
	bankerSeat := params.BankerSeat
	params.DingqueColor[huSeat] = room.CardColor_CC_TONG

	params.WallCards = []uint32{18, 18}

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	// 0 号玩家收到出牌问询通知，直接出牌1W
	bankerPlayer := utils.GetDeskPlayerBySeat(bankerSeat, deskData)
	expector, _ := bankerPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, Int1W))

	// 1 号玩家收到出牌问询通知，并且可杠
	gangPlayer := utils.GetDeskPlayerBySeat(gangSeat, deskData)
	expector, _ = gangPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf := room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.True(t, ntf.GetEnableMinggang())

	// 1号玩家 发送明杠请求
	utils.SendGangReq(deskData, gangSeat, Int1W, room.GangType_MingGang)

	// 1号玩家 摸牌出牌
	expector, _ = gangPlayer.Expectors[msgid.MsgID_ROOM_MOPAI_NTF]
	mopaiNtf := room.RoomMopaiNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &mopaiNtf))
	assert.Equal(t, gangPlayer.Player.GetID(), mopaiNtf.GetPlayer())
	assert.Equal(t, Int8W, mopaiNtf.GetCard())

	// 1号玩家 收到 自询问通知
	huPlayer := utils.GetDeskPlayerBySeat(huSeat, deskData)
	expector, _ = huPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	zxNtf := room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &zxNtf))
	assert.True(t, zxNtf.GetEnableZimo())
	assert.True(t, zxNtf.GetEnableQi())
	// 1号玩家发送 行牌动作请求 胡
	utils.SendHuReq(deskData, huSeat)
	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{huSeat}, gangSeat, Int8W, room.HuType_HT_GANGKAI)
}
