package mailtests

import (
	"testing"
	"steve/simulate/utils"
	"github.com/stretchr/testify/assert"
	"steve/client_pb/msgid"
	"steve/client_pb/msgserver"
	"steve/simulate/global"
)

// Test_GetPlayerInfo 测试获取玩家数据
func Test_GetUnReadMailSum(t *testing.T) {

	reqCmd := msgid.MsgID_MAILSVR_GET_UNREAD_SUM_REQ
	rspCmd := msgid.MsgID_MAILSVR_GET_UNREAD_SUM_RSP
	req := &msgserver.MsgSvrGetHorseRaceReq{}
	rsp := &msgserver.MsgSvrGetHorseRaceRsp{}

	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)



	player.AddExpectors(rspCmd)
	player.GetClient().SendPackage(utils.CreateMsgHead(reqCmd), req)
	expector := player.GetExpector(rspCmd)


	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, rsp))
	assert.Zero(t, rsp.GetErrCode())

	t.Logf("Test_GetUnReadMailSum win:", rsp)
	// assert.NotEmpty(t, response.GetNickName())
}

