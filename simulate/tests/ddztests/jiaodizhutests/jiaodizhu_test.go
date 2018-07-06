package jiaodizhu

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
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

	// 当前状态的时间间隔
	logEntry.Infof("当前状态 = %v, 进入下一状态等待时间 = %v", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())

	// ---------------------------------------------------------	叫地主状态	-----------------------------------------------------------

	// 叫地主的玩家
	jiaoPlayer := deskData.Players[deskData.DDZData.AssignLordID]

	// 执行四次，A叫一次，B抢一次，C抢一次，A再抢一次
	for i := 1; i <= 4; i++ {

		// 发出第i次叫地主请求
		assert.Nil(t, sendJiaodizhuReq(&jiaoPlayer))

		// 期望收到叫地主广播通知，且叫地主的是前面请求的玩家
		for playerID, player := range deskData.Players {

			// 叫地主的广播
			ntf := room.DDZGrabLordNtf{}

			if err := player.Expectors[msgid.MsgID_ROOM_DDZ_GRAB_LORD_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				logEntry.Errorf("接收第%v次的叫地主广播消息超时", i)
				assert.NotNil(t, nil)
				return
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
			logEntry.Error("接收最终地主广播消息超时")
			assert.NotNil(t, nil)
			return
		}

		logEntry.Infof("%v收到了最终的地主广播，最终地主的玩家是%v", playerID, ntf.GetPlayerId())

		// 记录状态
		deskData.DDZData.NextState = ntf.GetNextStage()

		// 记录最终的地主
		deskData.DDZData.ResultLordID = ntf.GetPlayerId()

		logEntry.Infof("最新状态 = %v, 进入下一状态等待时间 = %d", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())
	}

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
		assert.Nil(t, sendDoubleReq(operaPlayer, double))

		// 期望收到加倍广播通知，且加倍的是前面请求的玩家
		for playerID, player := range deskData.Players {

			// 加倍的广播
			ntf := room.DDZDoubleNtf{}

			if err := player.Expectors[msgid.MsgID_ROOM_DDZ_DOUBLE_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				logEntry.Errorf("接收第%d次的加倍广播消息超时", i)
				assert.NotNil(t, nil)
				return
			}

			logEntry.Infof("%v收到了第%v次加倍操作的广播，操作玩家是%v，加倍：%v， 当前加倍总倍数%v，下一状态：%v， 进入下一状态时间：%v",
				playerID, i, ntf.GetPlayerId(), ntf.GetIsDouble(), ntf.GetTotalDouble(), ntf.GetNextStage().GetStage(), ntf.GetNextStage().GetTime())

			// 记录状态
			deskData.DDZData.NextState = ntf.GetNextStage()
		}
	}

	// ---------------------------------------------------------	行牌状态	-----------------------------------------------------------
	// 最终地主的deskPlayer
	lordPlayer := deskData.Players[deskData.DDZData.ResultLordID]

	// 最终地主的座位号
	lordseatID := lordPlayer.Seat

	// 最终地主手里的牌
	lordCards := deskData.DDZData.Params.Cards[lordseatID]

	// 农民1
	farmer1 := utils.DeskPlayer{}
	farmer1.Seat = -1

	// 农民2
	farmer2 := utils.DeskPlayer{}
	farmer2.Seat = -1

	// 第几次出牌
	i := 0

	// 给两个农民赋值
	for playerID, player := range deskData.Players {
		// 不是地主
		if playerID != deskData.DDZData.ResultLordID {

			if farmer1.Seat == -1 {
				farmer1 = player
				continue
			}

			if farmer2.Seat == -1 {
				farmer2 = player
				continue
			}
		}
	}

	for j := 0; j < 3; /* len(lordCards) */ j++ {

		i++

		// 地主出牌，一次出一张牌
		assert.Nil(t, sendPlayCardReq(&lordPlayer, []uint32{lordCards[j]}, room.CardType_CT_SINGLE))

		// 检测第i次地主出牌成功
		ntf := room.DDZPlayCardRsp{}
		if err := lordPlayer.Expectors[msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
			logEntry.Errorf("地主第%d次出牌回应超时", i)
			assert.NotNil(t, nil)
			return
		}

		// 农民1检测第i次出牌广播
		nextPlayerID, err := listenPlayCardNtf(&farmer1, i)
		assert.Nil(t, err)

		// 农民2检测第1次出牌广播
		nextPlayerID, err = listenPlayCardNtf(&farmer2, i)
		assert.Nil(t, err)

		logEntry.Info("确定下次出牌玩家ID为", nextPlayerID)

		// 农民1
		if nextPlayerID == farmer1.Player.GetID() {

			// 农民1放弃出牌
			assert.Nil(t, sendPlayCardReq(&farmer1, []uint32{}, room.CardType_CT_NONE))

			// 农民2放弃出牌
			assert.Nil(t, sendPlayCardReq(&farmer2, []uint32{}, room.CardType_CT_NONE))
		} else {
			// 农民2放弃出牌
			assert.Nil(t, sendPlayCardReq(&farmer2, []uint32{}, room.CardType_CT_NONE))

			// 农民1放弃出牌
			assert.Nil(t, sendPlayCardReq(&farmer1, []uint32{}, room.CardType_CT_NONE))
		}

		// 重新建立监听

		assert.NotNil(t, nil)
	}

	/* 	// ---------------------------------------------------------	第一次出牌	-----------------------------------------------------------
	   	i++

	   	// 三个5
	   	cards1 := []uint32{
	   		uint32(room.PokerSuit_PS_CLUB) + uint32(room.PokerValue_PV_5),  // 梅花5
	   		uint32(room.PokerSuit_PS_HEART) + uint32(room.PokerValue_PV_5), // 红桃5
	   		uint32(room.PokerSuit_PS_SPADE) + uint32(room.PokerValue_PV_5), // 黑桃5}
	   	}

	   	assert.Nil(t, sendPlayCardReq(&lordPlayer, cards1, room.CardType_CT_PAIRS))

	   	// 检测第1次地主出牌成功
	   	ntf1 := room.DDZPlayCardRsp{}
	   	if err := lordPlayer.Expectors[msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP].Recv(global.DefaultWaitMessageTime, &ntf1); err != nil {
	   		logEntry.Errorf("地主第%d次出牌回应超时", i)
	   		assert.NotNil(t, nil)
	   		return
	   	}

	   	// 农民1检测第1次出牌广播
	   	assert.Nil(t, listenPlayCardNtf(&farmer1, i))

	   	// 农民2检测第1次出牌广播
	   	assert.Nil(t, listenPlayCardNtf(&farmer2, i))


	   	// ---------------------------------------------------------	第二次出牌	-----------------------------------------------------------
	   	i++ */
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

// sendPlayCardReq 发送出牌请求
// cards	：	要出的牌
func sendPlayCardReq(player *utils.DeskPlayer, cards []uint32, cardType room.CardType) error {
	logrus.WithFields(logrus.Fields{
		"func_name": "sendPlayCardReq()",
		"cards":     cards,
		"cardType":  cardType,
	}).Info("发出出牌请求，玩家 = ", player.Player.GetID())

	client := player.Player.GetClient()
	_, err := client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_DDZ_PLAY_CARD_REQ), &room.DDZPlayCardReq{
		Cards:    cards,     // 牌数据
		CardType: &cardType, // 牌类型
	})

	return err
}

// 指定的deskPlayer监听出牌的消息
func listenPlayCardNtf(player *utils.DeskPlayer, i int) (nextPlayerID uint64, err error) {
	logrus.WithFields(logrus.Fields{
		"func_name": "listenPlayCardNtf()",
	})

	ntf := room.DDZPlayCardNtf{}
	if err := player.Expectors[msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF].Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return 0, fmt.Errorf("%d监听第%d次出牌广播超时", player.Player.GetID(), i)
	}

	logrus.Infof("玩家%d监听到第%d次玩家%d出牌为%v，下一个出牌玩家为%v", player.Player.GetID(), i, ntf.GetPlayerId(), ntf.GetCards(), ntf.GetNextPlayerId())

	return ntf.GetNextPlayerId(), nil
}
