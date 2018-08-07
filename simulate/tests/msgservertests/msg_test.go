package msgservertests

import (
	"testing"
	"steve/simulate/utils"
	"github.com/stretchr/testify/assert"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/client_pb/msgserver"
	"github.com/Sirupsen/logrus"
)

// Test_GetPlayerInfo 测试获取玩家数据
func Test_GetHorseRate(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)

	player.AddExpectors(msgid.MsgID_MSGSVR_GET_HORSE_RACE_RSP)
	player.GetClient().SendPackage(utils.CreateMsgHead(msgid.MsgID_MSGSVR_GET_HORSE_RACE_REQ), &msgserver.MsgSvrGetHorseRaceReq{})
	expector := player.GetExpector(msgid.MsgID_MSGSVR_GET_HORSE_RACE_RSP)

	response := msgserver.MsgSvrGetHorseRaceRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &response))
	assert.Zero(t, response.GetErrCode())

	logrus.Infoln("GetHorseRate成功:", response)
	// assert.NotEmpty(t, response.GetNickName())
}
