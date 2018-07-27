package doudizhu

import (
	"fmt"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"time"

	"github.com/Sirupsen/logrus"
)

// 获取一张牌的点数
func getCardPoint(card uint32) uint32 {
	return card % 16
}

// 获取一张牌的花色(1,2,3,4分别代表方块，梅花，红桃，黑桃，0则代表大小王)
func getCardColor(card uint32) uint32 {
	return card / 16
}

// 两张牌是否点数相同
func isSamePoint(card1 uint32, card2 uint32) bool {
	if card1%16 == card2%16 {
		return true
	}

	return false
}

// 两张牌是否是对子
func isPair(card1 uint32, card2 uint32) bool {
	if isSamePoint(card1, card2) {
		return true
	}

	return false
}

// IsTriple 三张牌是否是三张
func IsTriple(card1 uint32, card2 uint32, card3 uint32) bool {
	card1Point := getCardPoint(card1)
	card2Point := getCardPoint(card2)
	card3Point := getCardPoint(card3)

	if (card1Point == card2Point) &&
		(card1Point == card3Point) {
		return true
	}

	return false
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

// PlaycardTest1 出牌测试1
func PlaycardTest1(deskData *utils.DeskData) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.go::PlaycardTest1",
	})
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

		// 出的牌
		outCards := []uint32{}
		outCards = append(outCards, lordCards[j])

		// 出牌的类型
		outCardType := room.CardType_CT_SINGLE

		/* 		// ------------------------------  对子 ------------------------------------
		   		// 差距16的话，说明是对子
		   		bDoubleType := false
		   		if j+1 < len(lordCards) {
		   			if IsPair(lordCards[j], lordCards[j+1]) {
		   				outCards = append(outCards, lordCards[j+1])
		   				// 后移一位
		   				j++
		   				bDoubleType = true
		   				outCardType = room.CardType_CT_PAIR
		   			}
		   		} */

		// 地主出牌，一次出一张牌
		sendPlayCardReq(&lordPlayer, outCards, outCardType)

		// 若是最后一张牌，则不再关心回应，因为游戏应该结束了
		if j == len(lordCards)-1 {
			break
		}

		// 检测第i次地主出牌成功
		ntf := room.DDZPlayCardRsp{}
		if err := lordPlayerRspExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
			return fmt.Errorf("地主第%d次出牌回应超时", i)
		}

		// 农民1检测第i次出牌广播
		nextPlayerID, err := listenPlayCardNtf(farmer1NtfExpect, &farmer1, i)
		if err != nil {
			return err
		}

		// 农民2检测第i次出牌广播
		nextPlayerID, err = listenPlayCardNtf(farmer2NtfExpect, &farmer2, i)
		if err != nil {
			return err
		}

		logEntry.Info("确定下次出牌玩家ID为", nextPlayerID)

		// 农民1
		if nextPlayerID == farmer1.Player.GetID() {

			// 暂停
			time.Sleep(10 * time.Millisecond)

			// 农民1放弃出牌
			err := sendPlayCardReq(&farmer1, []uint32{}, room.CardType_CT_NONE)
			if err != nil {
				return err
			}

			// 暂停
			time.Sleep(10 * time.Millisecond)

			// 农民2放弃出牌
			err = sendPlayCardReq(&farmer2, []uint32{}, room.CardType_CT_NONE)
			if err != nil {
				return err
			}

			// 检测农民2出牌结果
			ntf := room.DDZPlayCardRsp{}
			if err := farmer2RspExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				return fmt.Errorf("玩家%d 第%d次出牌回应超时", farmer2.Player.GetID(), i)
			}

			logEntry.Errorf("玩家%d 第%d次出牌回应结果为%s", farmer2.Player.GetID(), i, ntf.GetResult().GetErrDesc())
		} else {
			// 暂停
			time.Sleep(10 * time.Millisecond)

			// 农民2放弃出牌
			err := sendPlayCardReq(&farmer2, []uint32{}, room.CardType_CT_NONE)
			if err != nil {
				return err
			}

			// 暂停
			time.Sleep(10 * time.Millisecond)

			// 农民1放弃出牌
			err = sendPlayCardReq(&farmer1, []uint32{}, room.CardType_CT_NONE)
			if err != nil {
				return err
			}

			// 检测农民1出牌结果
			ntf := room.DDZPlayCardRsp{}
			if err := farmer1RspExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
				return fmt.Errorf("玩家%d 第%d次出牌回应超时", farmer1.Player.GetID(), i)
			}

			logEntry.Errorf("玩家%d 第%d次出牌回应结果为%s", farmer1.Player.GetID(), i, ntf.GetResult().GetErrDesc())
		}

		// 对子
		/* 		if bDoubleType {
			logEntry.Errorf("地主玩家%d第%d次出牌为：对子", lordPlayer.Player.GetID(), i)
			assert.NotNil(t, nil)
			return
		} */

		// 暂停
		time.Sleep(10 * time.Millisecond)
	}

	// 牌已出完，期待游戏结束通知

	// 地主应收到游戏结束的通知
	ntf := room.DDZGameOverNtf{}
	if err := lordPlayerOverExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return fmt.Errorf("地主玩家%d没有收到游戏结束的通知", lordPlayer.Player.GetID())
	}

	// 农民1应收到游戏结束的通知
	ntf = room.DDZGameOverNtf{}
	if err := farmer1OverExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return fmt.Errorf("农民玩家%d没有收到游戏结束的通知", farmer1.Player.GetID())
	}

	// 农民2应收到游戏结束的通知
	ntf = room.DDZGameOverNtf{}
	if err := farm2OverExpect.Recv(global.DefaultWaitMessageTime, &ntf); err != nil {
		return fmt.Errorf("农民玩家%d没有收到游戏结束的通知", farmer2.Player.GetID())
	}

	// 胜利者是地主
	if ntf.GetWinnerId() != lordPlayer.Player.GetID() {
		return fmt.Errorf("游戏结束时，胜利者竟然不是地主！胜利者ID = %v", ntf.GetWinnerId())
	}

	// 打印游戏结束信息
	logEntry.Infof("游戏结束，胜利者ID = %d，摊牌时间 = %d", ntf.GetWinnerId(), ntf.GetShowHandTime())

	for i := 0; i < len(ntf.GetBills()); i++ {
		playrInfo := ntf.GetBills()[i]
		logEntry.Infof("玩家:%d，底分:%d，输赢倍数:%d，输赢分数:%d，当前分数:%d，是否为地主:%v，已出的牌:%v，手中的牌:%v",
			playrInfo.GetPlayerId(), playrInfo.GetBase(), playrInfo.GetMultiple(),
			playrInfo.GetScore(), playrInfo.GetCurrentScore(), playrInfo.GetLord(), playrInfo.GetOutCards(), playrInfo.GetHandCards())
	}

	return nil
}
