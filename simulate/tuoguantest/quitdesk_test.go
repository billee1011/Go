package tuoguantest

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_ChupauiwenxunTuoguan 测试出牌问询时，退出房间托管
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行换三张,定缺
//  2. 用户0-2在收到换三张完成通知后，请求定缺，花色为万
//  3. 用户1请求退出游戏，用户1执行摸牌
// 期望：
// 	1. 最迟1秒后，用户0-2收到用户1摸牌通知， 用户3不会收到定缺完成通知
func Test_QuitDesk(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.WallCards = []uint32{31}
	params.Cards[0] = []uint32{11, 11, 11, 31, 12, 12, 12, 32, 13, 13, 13, 33, 14, 18}
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	quitSeat := deskData.BankerSeat

	assert.Nil(t, utils.SendChupaiReq(deskData, quitSeat, 18))

	player := utils.GetDeskPlayerBySeat(quitSeat, deskData)
	assert.Nil(t, utils.SendQuitReq(deskData, quitSeat))
	expector := player.Expectors[msgId.MsgID_ROOM_DESK_QUIT_RSP]
	rsp := room.RoomDeskQuitRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &rsp))
	assert.Equal(t, room.RoomError_SUCCESS, rsp.GetErrCode())
}
