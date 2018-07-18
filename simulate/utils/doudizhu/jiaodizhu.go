package doudizhu

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"

	"github.com/Sirupsen/logrus"
)

// sendChatReq 发送叫地主请求
func sendJiaodizhuReq(player *utils.DeskPlayer, bJiao bool) error {
	logrus.WithFields(logrus.Fields{
		"func_name": "TestJiaodizhu()",
	}).Info("发出叫地主请求，玩家 = ", player.Player.GetID())

	// 叫地主为true，不叫为false
	jiao := bJiao

	client := player.Player.GetClient()
	_, err := client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_DDZ_GRAB_LORD_REQ), &room.DDZGrabLordReq{
		Grab: &jiao,
	})

	return err
}

// JiaodizhuTest1 叫地主测试1
func JiaodizhuTest1(deskData *utils.DeskData) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.go::JiaodizhuTest1",
	})

	// ---------------------------------------------------------	叫地主状态	-----------------------------------------------------------
	// 叫地主的玩家
	jiaoPlayer := deskData.Players[deskData.DDZData.AssignLordID]

	// 执行四次，A叫一次，B抢一次，C抢一次，A再抢一次
	for i := 1; i <= 4; i++ {

		// 发出第i次叫地主请求
		sendJiaodizhuReq(&jiaoPlayer, true)

		// 期望收到叫地主广播通知，且叫地主的是前面请求的玩家
		for playerID, player := range deskData.Players {

			// 叫地主的广播
			ntf := room.DDZGrabLordNtf{}

			if err := player.Expectors[msgid.MsgID_ROOM_DDZ_GRAB_LORD_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				return fmt.Errorf("接收第%v次的叫地主广播消息超时", i)
			}

			logEntry.Infof("%v收到了第%v次叫地主的广播，叫地主的玩家是%v，下一个操作的玩家是%v", playerID, i, ntf.GetPlayerId(), ntf.GetNextPlayerId())

			// 下一次抢地主的玩家
			jiaoPlayer = deskData.Players[ntf.GetNextPlayerId()]
		}
	}

	// 期望收到最终的地主广播通知
	for playerID, player := range deskData.Players {
		// 最终的地主广播
		ntf := room.DDZLordNtf{}

		// 等久一些，一直等待到抢地主状态结束
		if err := player.Expectors[msgid.MsgID_ROOM_DDZ_LORD_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
			return fmt.Errorf("接收最终地主广播消息超时")
		}

		logEntry.Infof("%v收到了最终的地主广播，最终地主的玩家是%v", playerID, ntf.GetPlayerId())

		// 记录状态
		deskData.DDZData.NextState = ntf.GetNextStage()

		// 记录最终的地主
		deskData.DDZData.ResultLordID = ntf.GetPlayerId()

		logEntry.Infof("最新状态 = %v, 进入下一状态等待时间 = %d", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())
	}

	return nil
}

// JiaodizhuTest2 叫地主测试2
func JiaodizhuTest2(deskData *utils.DeskData) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.go::JiaodizhuTest2",
	})

	// ---------------------------------------------------------	叫地主状态	-----------------------------------------------------------
	// 叫地主的玩家
	jiaoPlayer := deskData.Players[deskData.DDZData.AssignLordID]

	// 执行四次，A叫一次，B抢一次，C不抢，A不抢
	for i := 1; i <= 4; i++ {

		switch i {
		case 1:
			sendJiaodizhuReq(&jiaoPlayer, true)
			break
		case 2:
			sendJiaodizhuReq(&jiaoPlayer, true)
			break
		case 3:
			sendJiaodizhuReq(&jiaoPlayer, false)
			break
		case 4:
			sendJiaodizhuReq(&jiaoPlayer, false)
			break
		}

		// 期望收到叫地主广播通知，且叫地主的是前面请求的玩家
		for playerID, player := range deskData.Players {

			// 叫地主的广播
			ntf := room.DDZGrabLordNtf{}

			if err := player.Expectors[msgid.MsgID_ROOM_DDZ_GRAB_LORD_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				return fmt.Errorf("接收第%v次的叫地主广播消息超时", i)
			}

			logEntry.Infof("%v收到了第%v次叫地主的广播，操作叫地主的玩家是:%v，操作是:%v，下一个操作的玩家是%v", playerID, i, ntf.GetPlayerId(), ntf.GetGrab(), ntf.GetNextPlayerId())

			// 下一次抢地主的玩家
			jiaoPlayer = deskData.Players[ntf.GetNextPlayerId()]
		}
	}

	// 期望收到最终的地主广播通知
	for playerID, player := range deskData.Players {
		// 最终的地主广播
		ntf := room.DDZLordNtf{}

		// 等久一些，一直等待到抢地主状态结束
		if err := player.Expectors[msgid.MsgID_ROOM_DDZ_LORD_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
			return fmt.Errorf("接收最终地主广播消息超时")
		}

		logEntry.Infof("%v收到了最终的地主广播，最终地主的玩家是%v", playerID, ntf.GetPlayerId())

		// 记录状态
		deskData.DDZData.NextState = ntf.GetNextStage()

		// 记录最终的地主
		deskData.DDZData.ResultLordID = ntf.GetPlayerId()

		logEntry.Infof("最新状态 = %v, 进入下一状态等待时间 = %d", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())
	}

	//return fmt.Errorf("叫地主测试用例2正常结束")
	return nil
}

// JiaodizhuTest3 叫地主测试3
func JiaodizhuTest3(deskData *utils.DeskData) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.go::JiaodizhuTest3",
	})

	// ---------------------------------------------------------	叫地主状态	-----------------------------------------------------------
	// 叫地主的玩家
	jiaoPlayer := deskData.Players[deskData.DDZData.AssignLordID]

	// 执行四次，A叫一次，B不抢，C不抢
	for i := 1; i <= 3; i++ {

		switch i {
		case 1:
			sendJiaodizhuReq(&jiaoPlayer, true)
			break
		case 2:
			sendJiaodizhuReq(&jiaoPlayer, false)
			break
		case 3:
			sendJiaodizhuReq(&jiaoPlayer, false)
			break
		}

		// 期望收到叫地主广播通知，且叫地主的是前面请求的玩家
		for playerID, player := range deskData.Players {

			// 叫地主的广播
			ntf := room.DDZGrabLordNtf{}

			if err := player.Expectors[msgid.MsgID_ROOM_DDZ_GRAB_LORD_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				return fmt.Errorf("接收第%v次的叫地主广播消息超时", i)
			}

			logEntry.Infof("%v收到了第%v次叫地主的广播，操作叫地主的玩家是:%v，操作是:%v，下一个操作的玩家是%v", playerID, i, ntf.GetPlayerId(), ntf.GetGrab(), ntf.GetNextPlayerId())

			// 下一次抢地主的玩家
			jiaoPlayer = deskData.Players[ntf.GetNextPlayerId()]
		}
	}

	// 期望收到最终的地主广播通知
	for playerID, player := range deskData.Players {
		// 最终的地主广播
		ntf := room.DDZLordNtf{}

		// 等久一些，一直等待到抢地主状态结束
		if err := player.Expectors[msgid.MsgID_ROOM_DDZ_LORD_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
			return fmt.Errorf("接收最终地主广播消息超时")
		}

		logEntry.Infof("%v收到了最终的地主广播，最终地主的玩家是%v", playerID, ntf.GetPlayerId())

		// 记录状态
		deskData.DDZData.NextState = ntf.GetNextStage()

		// 记录最终的地主
		deskData.DDZData.ResultLordID = ntf.GetPlayerId()

		logEntry.Infof("最新状态 = %v, 进入下一状态等待时间 = %d", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())
	}

	//return fmt.Errorf("叫地主测试用例2正常结束")
	return nil
}
