package jiaodizhu

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/utils"
	"testing"

	"github.com/Sirupsen/logrus"
)

//TestJiabei 加倍测试
//游戏过程中0号玩家发起加倍
//期望：
//     1. 所有玩家都收到，0号玩家的加倍广播
func TestJiabei(t *testing.T) {

	/* 	logEntry := logrus.WithFields(logrus.Fields{
	   		"func_name": "TestJiabei()",
	   	})

	   	// 正常开始游戏
	   	params := global.NewStartDDZGameParams()
	   	deskData, err := utils.StartDDZGame(params)
	   	assert.NotNil(t, deskData)
	   	assert.Nil(t, err)

	   	// 加倍广播的消息期待
	   	// 所有人都要收到那个人的加倍广播
	   	expectors := []interfaces.MessageExpector{}
	   	for _, player := range deskData.Players {
	   		expector, _ := player.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_DOUBLE_NTF)
	   		expectors = append(expectors, expector)
	   	}

	   	// 0号玩家发起加倍
	   	player := utils.GetDeskPlayerBySeat(0, deskData)
	   	//player := deskData.Players[deskData.DDZData.AssignLordID]

	   	// 当前状态的时间间隔
	   	logEntry.Infof("当前状态 = %v, 进入下一状态等待时间 = %d", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())

	   	// 发出加倍请求
	   	assert.Nil(t, sendJiabeiReq(player))

	   	// 期望收到加倍广播通知，且加倍的是前面请求的玩家
	   	for i := 0; i < len(expectors); i++ {

	   		// 加倍广播消息
	   		ntf := room.DDZDoubleNtf{}

	   		if err := expectors[i].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
	   			logEntry.Error("接收加倍广播消息超时")
	   			assert.NotNil(t, nil)
	   			return
	   		}

	   		logEntry.Info("收到了斗地主的加倍广播，加倍的玩家是%v", ntf.GetPlayerId())

	   		// 两者必须相同
	   		assert.Equal(t, player.Player.GetID(), ntf.GetPlayerId())
	   	} */
}

// sendJiabeiReq 发送加倍请求
func sendJiabeiReq(player *utils.DeskPlayer) error {
	logrus.WithFields(logrus.Fields{
		"func_name": "sendJiabeiReq()",
	}).Info("发出加倍请求，玩家 = ", player.Player.GetID())

	// 加倍为true
	double := true

	client := player.Player.GetClient()
	_, err := client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_DDZ_DOUBLE_REQ), &room.DDZDoubleReq{
		IsDouble: &double,
	})

	return err
}
