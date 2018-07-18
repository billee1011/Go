package doudizhu

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"

	"github.com/Sirupsen/logrus"
)

// sendDoubleReq 发送加倍请求
// double : 是否加倍
func sendDoubleReq(player *utils.DeskPlayer, double bool) error {
	logrus.WithFields(logrus.Fields{
		"func_name": "sendJiabeiReq()",
	}).Info("发出加倍请求，玩家 = ", player.Player.GetID())

	client := player.Player.GetClient()
	_, err := client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_DDZ_DOUBLE_REQ), &room.DDZDoubleReq{
		IsDouble: &double, // 加倍为true，不加倍为false
	})

	return err
}

// JiabeiTest1 加倍测试1
func JiabeiTest1(deskData *utils.DeskData) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.go::JiabeiTest1",
	})

	// ---------------------------------------------------------	加倍状态	-----------------------------------------------------------

	// 执行三次加倍(0号玩家不加倍，1号，2号玩家加倍)
	for i := 1; i <= 3; i++ {

		// 玩家
		operaPlayer := utils.GetDeskPlayerBySeat(i-1, deskData)

		double := true

		// 第1次：不加倍
		if i == 1 {
			double = false
		}
		sendDoubleReq(operaPlayer, double)

		// 期望收到加倍广播通知，且加倍的是前面请求的玩家
		for playerID, player := range deskData.Players {

			// 加倍的广播
			ntf := room.DDZDoubleNtf{}

			if err := player.Expectors[msgid.MsgID_ROOM_DDZ_DOUBLE_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				return fmt.Errorf("接收第%d次的加倍广播消息超时", i)
			}

			logEntry.Infof("%v收到了第%v次加倍操作的广播，操作玩家是%v，加倍：%v， 当前加倍总倍数%v，下一状态：%v， 进入下一状态时间：%v",
				playerID, i, ntf.GetPlayerId(), ntf.GetIsDouble(), ntf.GetTotalDouble(), ntf.GetNextStage().GetStage(), ntf.GetNextStage().GetTime())

			// 记录状态
			deskData.DDZData.NextState = ntf.GetNextStage()
		}
	}

	return nil
}
