package jiaodizhu

import (
	"fmt"
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"
	"time"

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

	// 加上三张底牌
	lordCards = append(lordCards, uint32(room.PokerSuit_PS_SPADE)+uint32(room.PokerValue_PV_K))          // 红桃K)
	lordCards = append(lordCards, uint32(room.PokerSuit_PS_NONE)+uint32(room.PokerValue_PV_BLACK_JOKER)) // 小王)
	lordCards = append(lordCards, uint32(room.PokerSuit_PS_NONE)+uint32(room.PokerValue_PV_RED_JOKER))   // 大王)

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

	// 建立GameOver的消息期望

	// 地主收到GameOver的期望
	lordPlayer.Player.GetClient().RemoveMsgExpect(msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF)
	lordPlayerOverExpect, _ := lordPlayer.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF)

	// 农民1收到GameOver的期望
	farmer1.Player.GetClient().RemoveMsgExpect(msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF)
	farmer1OverExpect, _ := farmer1.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF)

	// 地主收到GameOver的期望
	farmer2.Player.GetClient().RemoveMsgExpect(msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF)
	farm2OverExpect, _ := farmer2.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF)

	for j := 0; j < len(lordCards); j++ {

		i++

		// 地主出牌回应的期望
		lordPlayer.Player.GetClient().RemoveMsgExpect(msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP)
		lordPlayerRspExpect, _ := lordPlayer.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP)

		// 地主监听其他人出牌广播的期望
		//lordPlayer.Player.GetClient().RemoveMsgExpect(msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF)
		//lordPlayerNtfExpect, _ := lordPlayer.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF)

		// 农民1出牌回应的期望
		farmer1.Player.GetClient().RemoveMsgExpect(msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP)
		farmer1RspExpect, _ := farmer1.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP)

		// 农民1监听其他人出牌广播的期望
		farmer1.Player.GetClient().RemoveMsgExpect(msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF)
		farmer1NtfExpect, _ := farmer1.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF)

		// 农民2出牌回应的期望
		farmer2.Player.GetClient().RemoveMsgExpect(msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP)
		farmer2RspExpect, _ := farmer2.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP)

		// 农民2监听其他人出牌广播的期望
		farmer2.Player.GetClient().RemoveMsgExpect(msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF)
		farmer2NtfExpect, _ := farmer2.Player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF)

		// -----------------------------  一轮出牌	-------------------------------

		// 地主出牌，一次出一张牌
		assert.Nil(t, sendPlayCardReq(&lordPlayer, []uint32{lordCards[j]}, room.CardType_CT_SINGLE))

		// 若是最后一张牌，则不再关心回应，因为游戏应该结束了
		if j == len(lordCards)-1 {
			break
		}

		// 检测第i次地主出牌成功
		ntf := room.DDZPlayCardRsp{}
		if err := lordPlayerRspExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
			logEntry.Errorf("地主第%d次出牌回应超时", i)
			assert.NotNil(t, nil)
			return
		}

		// 农民1检测第i次出牌广播
		nextPlayerID, err := listenPlayCardNtf(farmer1NtfExpect, &farmer1, i)
		assert.Nil(t, err)

		// 农民2检测第i次出牌广播
		nextPlayerID, err = listenPlayCardNtf(farmer2NtfExpect, &farmer2, i)
		assert.Nil(t, err)

		logEntry.Info("确定下次出牌玩家ID为", nextPlayerID)

		// 农民1
		if nextPlayerID == farmer1.Player.GetID() {

			// 暂停
			time.Sleep(200 * time.Millisecond)

			// 农民1放弃出牌
			assert.Nil(t, sendPlayCardReq(&farmer1, []uint32{}, room.CardType_CT_NONE))

			// 暂停
			time.Sleep(200 * time.Millisecond)

			// 农民2放弃出牌
			assert.Nil(t, sendPlayCardReq(&farmer2, []uint32{}, room.CardType_CT_NONE))

			// 检测农民2出牌结果
			ntf := room.DDZPlayCardRsp{}
			if err := farmer2RspExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				logEntry.Errorf("玩家%d 第%d次出牌回应超时", farmer2.Player.GetID(), i)
				assert.NotNil(t, nil)
				return
			}

			logEntry.Errorf("玩家%d 第%d次出牌回应结果为%s", farmer2.Player.GetID(), i, ntf.GetResult().GetErrDesc())
		} else {
			// 暂停
			time.Sleep(200 * time.Millisecond)

			// 农民2放弃出牌
			assert.Nil(t, sendPlayCardReq(&farmer2, []uint32{}, room.CardType_CT_NONE))

			// 暂停
			time.Sleep(200 * time.Millisecond)

			// 农民1放弃出牌
			assert.Nil(t, sendPlayCardReq(&farmer1, []uint32{}, room.CardType_CT_NONE))

			// 检测农民1出牌结果
			ntf := room.DDZPlayCardRsp{}
			if err := farmer1RspExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				logEntry.Errorf("玩家%d 第%d次出牌回应超时", farmer1.Player.GetID(), i)
				assert.NotNil(t, nil)
				return
			}

			logEntry.Errorf("玩家%d 第%d次出牌回应结果为%s", farmer1.Player.GetID(), i, ntf.GetResult().GetErrDesc())
		}

		// 暂停2秒
		time.Sleep(200 * time.Millisecond)
	}

	// 牌已出完，期待游戏结束通知

	// 地主应收到游戏结束的通知
	ntf := room.DDZGameOverNtf{}
	if err := lordPlayerOverExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		logEntry.Errorf("地主玩家%d没有收到游戏结束的通知", lordPlayer.Player.GetID())
		assert.NotNil(t, nil)
		return
	}

	// 农民1应收到游戏结束的通知
	ntf = room.DDZGameOverNtf{}
	if err := farmer1OverExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		logEntry.Errorf("农民玩家%d没有收到游戏结束的通知", farmer1.Player.GetID())
		assert.NotNil(t, nil)
		return
	}

	// 农民2应收到游戏结束的通知
	ntf = room.DDZGameOverNtf{}
	if err := farm2OverExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		logEntry.Errorf("农民玩家%d没有收到游戏结束的通知", farmer2.Player.GetID())
		assert.NotNil(t, nil)
		return
	}

	// 胜利者是地主
	if ntf.GetWinnerId() != lordPlayer.Player.GetID() {
		logEntry.Errorf("游戏结束时，胜利者竟然不是地主！胜利者ID = ", ntf.GetWinnerId())
		assert.NotNil(t, nil)
		return
	}

	// 打印游戏结束信息
	logEntry.Infof("游戏结束，胜利者ID = %d，摊牌时间 = ", ntf.GetWinnerId(), ntf.GetShowHandTime())

	for i := 0; i < len(ntf.GetBills()); i++ {
		playrInfo := ntf.GetBills()[i]
		logEntry.Infof("玩家:%d，名字:%s，底分:%d，输赢倍数:%d，输赢分数:%d，当前分数:%d，是否为地主:%v，已出的牌:%v，手中的牌:%v",
			playrInfo.GetPlayerId(), playrInfo.GetPlayerName(), playrInfo.GetBase(), playrInfo.GetMultiple(),
			playrInfo.GetScore(), playrInfo.GetCurrentScore(), playrInfo.GetLord(), playrInfo.GetOutCards(), playrInfo.GetHandCards())
	}

	// ------------------------------------------------ 	恢复对局	  ---------------------------------------
	/* 	// 最终地主的deskPlayer
	   	lordPlayer = deskData.Players[deskData.DDZData.ResultLordID]

	   	// 地主断开连接
	   	assert.Nil(t, lordPlayer.Player.GetClient().Stop())
	   	time.Sleep(time.Millisecond * 200) // 等200毫秒，确保连接断开

	   	// 再新建一个连接
	   	client := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
	   	assert.NotNil(t, client)
	   	player, err := utils.LoginUser(client, lordPlayer.Player.GetUsrName())
	   	assert.Nil(t, err)
	   	assert.NotNil(t, player)
	   	assert.Equal(t, lordPlayer.Player.GetID(), player.GetID())

	   	// 步骤4
	   	utils.UpdatePlayerClientInfo(client, player, deskData)
	   	// 监听恢复对局的回复消息
	   	resumeRspExpect, _ := player.GetClient().ExpectMessage(msgid.MsgID_ROOM_DDZ_RESUME_RSP)

	   	// 发出恢复对局请求
	   	assert.Nil(t, SendDDZRecoverGameReq(lordPlayer.GetSeat(), deskData))

	   	resumeRsp := room.DDZResumeGameRsp{}
	   	if err := resumeRspExpect.Recv(global.DefaultWaitMessageTime, &resumeRsp); err != nil {
	   		logEntry.Errorf("玩家%d没有收到恢复对局的回复", lordPlayer.Player.GetID())
	   		assert.NotNil(t, nil)
	   		return
	   	}

	   	// 打印恢复对局回复的信息
	   	logEntry.Infof("玩家%d收到恢复对局的回复,resultCode = %d， resultStr = %s", resumeRsp.GetResult().GetErrCode(), resumeRsp.GetResult().GetErrDesc())

	   	// 成功时打印游戏信息
	   	if resumeRsp.GetResult().GetErrCode() == 0 {
	   		ddzDeskInfo := resumeRsp.GetGameInfo()
	   		// 桌子里面的每一个玩家
	   		for i := 0; i < len(ddzDeskInfo.GetPlayers()); i++ {
	   			ddzPlayrInfo := ddzDeskInfo.GetPlayers()[i]
	   			roomPlayerInfo := ddzPlayrInfo.GetPlayerInfo()
	   			logEntry.Infof("玩家:%d，名字:%s，金币数:%d，座位号:%d，已打出的牌:%v，手中的牌:%v，是否为地主:%v，是否托管:%d，是否加倍:%d",
	   				roomPlayerInfo.GetPlayerId(), roomPlayerInfo.GetName(), roomPlayerInfo.GetCoin(), roomPlayerInfo.GetSeat(),
	   				ddzPlayrInfo.GetOutCards(), ddzPlayrInfo.GetHandCards(), ddzPlayrInfo.GetLord(), ddzPlayrInfo.GetTuoguan(), ddzPlayrInfo.GetIsDouble())
	   		}

	   		// 当前状态
	   		logEntry.Infof("当前状态：%v，进入下一状态的等待时间:%d", ddzDeskInfo.GetStage().GetStage(), ddzDeskInfo.GetStage().GetTime())
	   	} */
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
	}).Infof("玩家%d发出出牌请求", player.Player.GetID())

	client := player.Player.GetClient()
	_, err := client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_DDZ_PLAY_CARD_REQ), &room.DDZPlayCardReq{
		Cards:    cards,     // 牌数据
		CardType: &cardType, // 牌类型
	})

	return err
}

// 指定的deskPlayer监听出牌的消息
func listenPlayCardNtf(expect interfaces.MessageExpector, player *utils.DeskPlayer, i int) (nextPlayerID uint64, err error) {
	logrus.WithFields(logrus.Fields{
		"func_name": "listenPlayCardNtf()",
	})

	ntf := room.DDZPlayCardNtf{}
	if err := expect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return 0, fmt.Errorf("%d监听第%d次出牌广播超时", player.Player.GetID(), i)
	}

	logrus.Infof("玩家%d监听到玩家%d第%d次出牌为%v，下一个出牌玩家为%v", player.Player.GetID(), ntf.GetPlayerId(), i, ntf.GetCards(), ntf.GetNextPlayerId())

	return ntf.GetNextPlayerId(), nil
}

// sendResumeGameReq 发送恢复对局请求
func sendResumeGameReq(player *utils.DeskPlayer) error {
	logrus.WithFields(logrus.Fields{
		"func_name": "sendResumeGameReq()",
	}).Infof("玩家%d发出恢复请求", player.Player.GetID())

	client := player.Player.GetClient()
	_, err := client.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_DDZ_RESUME_REQ), &room.DDZResumeGameReq{})

	return err
}
