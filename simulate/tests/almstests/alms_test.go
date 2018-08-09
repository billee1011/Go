package almstests

import (
	"fmt"
	"steve/client_pb/alms"
	"steve/client_pb/common"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//玩家登陆接受到，救济金配合通知
func Test_Alms_Login(t *testing.T) {
	player, _ := utils.LoginNewPlayer()
	assert.NotNil(t, player)

	player.AddExpectors(msgid.MsgID_ALMS_LOGIN_GOLD_CONFIG_NTF)

	expector := player.GetExpector(msgid.MsgID_ALMS_LOGIN_GOLD_CONFIG_NTF)
	ntf := &alms.AlmsConfigNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf))
	fmt.Println(ntf)
}

// Test_Apply_Alms 测试申请救济金

func Test_Apply_Alms(t *testing.T) {
	player, _ := utils.LoginNewPlayer()
	assert.NotNil(t, player)

	player.AddExpectors(msgid.MsgID_ALMS_GET_GOLD_RSP)
	client := player.GetClient()
	aat := alms.AlmsApplyType_AAT_SELECTIONS
	req := &alms.AlmsGetGoldReq{}
	gameID, levelID, totalGold, version := common.GameId_GAMEID_XUELIU, int32(1), int64(101), int32(1)
	req.AlmsApplyType = &aat
	req.GameId = &gameID
	req.LevelId = &levelID
	req.TotalGold = &totalGold
	req.Version = &version
	client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ALMS_GET_GOLD_REQ), req)

	expector := player.GetExpector(msgid.MsgID_ALMS_GET_GOLD_RSP)
	rsq := &alms.AlmsGetGoldRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, rsq))
	fmt.Println(rsq)
}
