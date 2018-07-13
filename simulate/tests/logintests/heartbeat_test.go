package logintests

import (
	"steve/client_pb/gate"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_HeartBeat 测试心跳
//	1. 登录用户
//	2. 发送心跳请求
// 期望： 收到心跳回复
func Test_HeartBeat(t *testing.T) {
	player, _ := utils.LoginNewPlayer()
	assert.NotNil(t, player)

	player.AddExpectors(msgid.MsgID_GATE_HEART_BEAT_RSP)
	client := player.GetClient()
	client.SendPackage(utils.CreateMsgHead(msgid.MsgID_GATE_HEART_BEAT_REQ), &gate.GateHeartBeatReq{})

	expector := player.GetExpector(msgid.MsgID_GATE_HEART_BEAT_RSP)

	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, nil))
}
