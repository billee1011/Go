package tests

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_StartGame_NoHsz 测试游戏开始
// 游戏开始流程包括： 登录，加入房间，配牌，洗牌，发牌,定缺
// 庄家只有3个1万,并且换牌3个1万，换牌换三张关闭，庄家3个1万，牌不变，出牌1万，其他玩家收到出牌通知
// 期望不出现换三张
func Test_SCXZ_StartGame_NoHsz(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN
	params.Cards = [][]uint32{
		{11, 11, 11, 37, 12, 12, 12, 12, 13, 13, 13, 13, 14, 14},
		{15, 15, 15, 15, 16, 16, 16, 16, 17, 17, 17, 17, 18},
		{21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24},
		{25, 25, 25, 25, 26, 26, 26, 26, 27, 27, 27, 27, 28},
	}
	params.HszCards = [][]uint32{
		{11, 11, 11},
		{15, 15, 15},
		{21, 21, 21},
		{25, 25, 25},
	}
	params.PeiPaiGame = "scxz"
	params.IsHsz = false // 不换三张
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// // 庄家出1万
	zjPlayer := utils.GetDeskPlayerBySeat(params.BankerSeat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	zjClient.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_CHUPAI_REQ), &room.RoomChupaiReq{
		Card: proto.Uint32(11),
	})

	// 其他玩家收到庄家出牌通知
	for _, deskPlayer := range deskData.Players {
		expector, _ := deskPlayer.Expectors[msgid.MsgID_ROOM_CHUPAI_NTF]
		ntf := room.RoomChupaiNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, ntf.GetPlayer(), zjPlayer.Player.GetID())
		assert.Equal(t, uint32(11), ntf.GetCard())
	}

}

// Test_SCXZ_StartGame_NoHsz 测试游戏开始
// 游戏开始流程包括： 登录，加入房间，配牌，洗牌，发牌,定缺
// 步骤:换三张成功,庄家拿到五条，出五条，其他玩家收到出牌消息五条通知
// 期望出现换三张
func Test_SCXZ_StartGame_Hsz(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN
	params.PeiPaiGame = "scxz"
	params.IsHsz = true // 换三张
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 庄家出五条
	zjPlayer := utils.GetDeskPlayerBySeat(params.BankerSeat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	zjClient.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_CHUPAI_REQ), &room.RoomChupaiReq{
		Card: proto.Uint32(25),
	})

	// 其他玩家收到庄家出牌通知
	for _, deskPlayer := range deskData.Players {
		expector, _ := deskPlayer.Expectors[msgid.MsgID_ROOM_CHUPAI_NTF]
		ntf := room.RoomChupaiNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, ntf.GetPlayer(), zjPlayer.Player.GetID())
		assert.Equal(t, uint32(25), ntf.GetCard())
	}

}
