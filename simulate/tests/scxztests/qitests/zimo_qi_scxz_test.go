package qitests

import (
	"steve/client_pb/common"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Zimo_qi 自摸弃测试
// 期望：
// 庄家出9W后，1号玩家将收到出牌问询通知，可杠
// 1号玩家发出弃杠请求，保留原状态
func Test_SCXZ_Zimo_qi(t *testing.T) {
	var Int1B uint32 = 31
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	zimoSeat := 1
	bankerSeat := params.BankerSeat

	// 庄家的最后一张牌改为 1B
	params.Cards[bankerSeat][13] = 31
	// 1 号玩家最后1张牌改为 9W
	params.Cards[zimoSeat][12] = 19
	// 墙牌改成 9W 。 墙牌有两张，否则就是海底捞了
	params.WallCards = []uint32{19, 31}

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, Int1B))

	// 1 号玩家收到可自摸通知
	zimoPlayer := utils.GetDeskPlayerBySeat(zimoSeat, deskData)
	expector, _ := zimoPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	ntf := room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.True(t, ntf.GetEnableZimo())

	// 发送弃请求
	assert.Nil(t, utils.SendQiReq(deskData, zimoSeat))
}
