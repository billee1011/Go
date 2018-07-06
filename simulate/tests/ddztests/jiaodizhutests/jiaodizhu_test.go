package jiaodizhu

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//TestJiaodizhu 叫地主测试
//游戏过程中0号玩家发起叫地主
//期望：
//     1. 所有玩家都收到，0号玩家的叫地主广播
func TestJiaodizhu(t *testing.T) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "TestJiaodizhu()",
	})

	// 正常开始游戏
	params := global.NewStartDDZGameParams()
	deskData, err := utils.StartDDZGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 指定的玩家

	// 叫地主的消息期待
	// 所有人都要收到那个人的叫地主广播
	expectors := []interfaces.MessageExpector{}
	for _, player := range deskData.Players {
		expector, _ := player.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_GRAB_LORD_NTF)
		expectors = append(expectors, expector)
	}

	// 地主广播的期待
	// 所有人都要收到地主广播消息
	expectorsNtf := []interfaces.MessageExpector{}
	for _, player := range deskData.Players {
		expector, _ := player.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_LORD_NTF)
		expectorsNtf = append(expectorsNtf, expector)
	}

	// 叫地主的玩家
	player := deskData.Players[deskData.DDZData.AssignLordID]

	// 当前状态的时间间隔
	logEntry.Infof("当前状态 = %v, 进入下一状态等待时间 = %d", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())

	// 发出叫地主请求
	assert.Nil(t, sendJiaodizhuReq(&player))

	// 期望收到叫地主广播通知，且叫地主的是前面请求的玩家
	for i := 0; i < len(expectors); i++ {

		// 叫地主的广播
		ntf := room.DDZGrabLordNtf{}

		if err := expectors[i].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
			logEntry.Error("接收叫地主广播消息超时")
			assert.NotNil(t, nil)
			return
		}

		// 两者必须相同
		assert.Equal(t, player.Player.GetID(), ntf.GetPlayerId())

		logEntry.Infof("收到了叫地主的广播，抢地主的玩家是%v", ntf.GetPlayerId())
	}

	// 期望收到最终的地主广播通知
	for i := 0; i < len(expectorsNtf); i++ {

		// 最终的地主广播
		ntf := room.DDZLordNtf{}

		if err := expectorsNtf[i].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
			logEntry.Error("接收最终地主广播消息超时")
			assert.NotNil(t, nil)
			return
		}

		logEntry.Infof("收到了最终的地主广播，最终地主的玩家是%v", ntf.GetPlayerId())

		// 两者必须相同
		assert.Equal(t, player.Player.GetID(), ntf.GetPlayerId())

		// 记录状态
		deskData.DDZData.NextState = ntf.GetNextStage()

		logEntry.Infof("下一状态 = %v, 进入下一状态等待时间 = %d", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())
	}
}

// sendChatReq 发送叫地主请求
func sendJiaodizhuReq(player *utils.DeskPlayer) error {
	logrus.WithFields(logrus.Fields{
		"func_name": "TestJiaodizhu()",
	}).Info("发出叫地主请求，玩家 = ", player.Player.GetID())

	// 叫地主为true
	jiao := true

	client := player.Player.GetClient()
	_, err := client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_DDZ_GRAB_LORD_REQ), &room.DDZGrabLordReq{
		Grab: &jiao,
	})

	return err
}
