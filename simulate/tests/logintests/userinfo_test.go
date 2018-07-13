package logintests

import (
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_GetPlayerInfo 测试获取玩家数据
func Test_GetPlayerInfo(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)

	player.AddExpectors(msgid.MsgID_HALL_GET_PLAYER_INFO_RSP)
	player.GetClient().SendPackage(utils.CreateMsgHead(msgid.MsgID_HALL_GET_PLAYER_INFO_REQ), &hall.HallGetPlayerInfoReq{})
	expector := player.GetExpector(msgid.MsgID_HALL_GET_PLAYER_INFO_RSP)

	response := hall.HallGetPlayerInfoRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &response))
	assert.Zero(t, response.GetErrCode())
	assert.NotEmpty(t, response.GetNickName())
}
