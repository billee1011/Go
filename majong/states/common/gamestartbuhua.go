package common

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// GameStartBuhuaState 开局补花
type GameStartBuhuaState struct{}

// ProcessEvent 处理事件,目前二人是自动补花，如果存在其他麻将有手动补花，补花需要请求
func (bh *GameStartBuhuaState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_gamestart_buhua_finish:
		{
			return majongpb.StateID_state_mopai, nil
		}
	}
	return majongpb.StateID_state_gamestart_buhua, nil
}

// doBuhua 执行
func (bh *GameStartBuhuaState) doBuhua(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	players := mjContext.GetPlayers()
	//初始化第一个实际补花的玩家
	buhuaPlayerID := players[int(mjContext.GetZhuangjiaIndex())].GetPalyerId()
	//循环补花开关
	continueBuhua := true
	//补花次数
	buhuaNum := 0
	//当所有人都没有花牌的时候，结束补花循环
	for continueBuhua {
		nextBuhuaPlayerID := bh.ntfBuhua(flow, buhuaPlayerID, buhuaNum)
		buhuaNum++
		if nextBuhuaPlayerID == 0 {
			continueBuhua = false
		} else {
			buhuaPlayerID = nextBuhuaPlayerID
		}
	}
	//所有人补花完成后，设置庄家为摸牌玩家
	mjContext.MopaiPlayer = mjContext.Players[int(mjContext.GetZhuangjiaIndex())].GetPalyerId()
}

func (bh *GameStartBuhuaState) devlopNextBuhuaPlayer(players []*majongpb.Player, curBuhuaPlayerID uint64, mjContext *majongpb.MajongContext) uint64 {
	player := utils.GetPlayerByID(players, curBuhuaPlayerID)
	//先判断当前补完花的玩家时候还有花牌，如果有花牌，则继续补花，没有的话，下家补花
	if len(player.GetHandCards()) < 13 || len(bh.getHuaCards(player)) != 0 {
		return curBuhuaPlayerID
	}
	for i := 0; i < len(players); i++ {
		nextplayer := utils.GetNextXpPlayerByID(curBuhuaPlayerID, players, mjContext)
		if len(nextplayer.GetHandCards()) < 13 || len(bh.getHuaCards(nextplayer)) != 0 {
			return nextplayer.GetPalyerId()
		}
		curBuhuaPlayerID = nextplayer.GetPalyerId()
	}
	return 0
}

func (bh *GameStartBuhuaState) ntfBuhua(flow interfaces.MajongFlow, buhuaPlayerID uint64, buhuaNum int) uint64 {
	var nextBuhuaPlayerID uint64
	if buhuaNum == 0 {
		nextBuhuaPlayerID = bh.ntfFirstBuhua(flow, buhuaPlayerID)
	} else {
		nextBuhuaPlayerID = bh.ntfOtherBuhua(flow, buhuaPlayerID)
	}
	return nextBuhuaPlayerID
}

func (bh *GameStartBuhuaState) ntfFirstBuhua(flow interfaces.MajongFlow, buhuaPlayerID uint64) uint64 {
	// 第一次补花需要广播所有人的花牌，并且移除所有人的花牌,且庄家需要摸补牌
	mjContext := flow.GetMajongContext()
	players := mjContext.GetPlayers()
	infos := []*room.RoomBuHuaInfo{}
	//将所有人的花牌拿出来,移除手牌中的花
	for _, player := range players {
		huaCards := bh.getHuaCards(player)
		for _, card := range huaCards {
			var ok bool
			player.HandCards, ok = utils.RemoveCards(player.HandCards, card, 1)
			player.HuaCards = append(player.GetHuaCards(), card)
			if !ok {
				logrus.WithFields(logrus.Fields{
					"func_name":       "GameStartBuhuaState.getHuaCards",
					"hand_cards":      player.HandCards,
					"buhua_player_id": player.GetPalyerId(),
				}).Errorln("移除补花者的花牌失败")
			}
		}
		info := room.RoomBuHuaInfo{
			PlayerId:    proto.Uint64(player.GetPalyerId()),
			OutHuaCards: utils.ServerCards2Uint32(huaCards),
		}
		infos = append(infos, &info)
	}
	//广播补花
	for _, player := range players {
		if buhuaPlayerID == player.GetPalyerId() {
			huaCards := bh.getHuaCards(player)
			//从墙牌中摸牌
			wallCards := mjContext.GetWallCards()
			l := len(huaCards)
			for _, info := range infos {
				if info.GetPlayerId() == player.GetPalyerId() {
					player.HandCards = append(player.HandCards, wallCards[0:l]...)
					info.BuCards = utils.ServerCards2Uint32(wallCards[0:l])
					wallCards = wallCards[l:]
				}
			}
		}
		toClientMessage := interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_ROOM_BUHUA_NTF),
			Msg: &room.RoomBuHuaNtf{
				BuhuaInfo: infos,
			},
		}
		logrus.WithFields(logrus.Fields{
			"func_name":     "gamestartbuhua_ntfFirstBuhua",
			"buhua_infos":   infos,
			"ntf_to_player": player.GetPalyerId(),
			"hand_cards":    utils.ServerCards2Uint32(player.GetHandCards()),
		}).Info("开局首次补花，全员亮花牌")
		flow.PushMessages([]uint64{player.GetPalyerId()}, toClientMessage)
	}
	return bh.devlopNextBuhuaPlayer(players, buhuaPlayerID, mjContext)
}

func (bh *GameStartBuhuaState) ntfOtherBuhua(flow interfaces.MajongFlow, buhuaPlayerID uint64) uint64 {
	// 由指定玩家摸补牌，如果补上的牌是花牌，继续补完，如果补完后则下家开始补
	mjContext := flow.GetMajongContext()
	players := mjContext.GetPlayers()
	buhuaPlayer := utils.GetPlayerByID(players, buhuaPlayerID)
	mopaiNum := 13 - len(buhuaPlayer.GetHandCards())
	if mopaiNum > 0 {
		bh.ntf(flow, players, buhuaPlayerID, make([]*majongpb.Card, 0), mopaiNum)
	}
	//补完牌后需要检查是否有花要补
	huaCards := bh.getHuaCards(buhuaPlayer)
	num := len(huaCards)
	if num > 0 {
		bh.ntf(flow, players, buhuaPlayerID, huaCards, num)
	}
	return bh.devlopNextBuhuaPlayer(players, buhuaPlayerID, mjContext)
}

func (bh *GameStartBuhuaState) getHuaCards(player *majongpb.Player) []*majongpb.Card {
	handCards := player.GetHandCards()
	huaCards := []*majongpb.Card{}
	for _, card := range handCards {
		if card.GetColor() == majongpb.CardColor_ColorHua {
			huaCards = append(huaCards, card)
		}
	}
	return huaCards
}

func (bh *GameStartBuhuaState) ntf(flow interfaces.MajongFlow, players []*majongpb.Player, curPlayerID uint64, huaCards []*majongpb.Card, buCardNum int) {
	mjContext := flow.GetMajongContext()
	if len(huaCards) > 0 {
		buCardNum = len(huaCards)
	}
	for _, player := range players {
		info := &room.RoomBuHuaInfo{
			PlayerId:    proto.Uint64(curPlayerID),
			BuCards:     make([]uint32, buCardNum),
			OutHuaCards: utils.ServerCards2Uint32(huaCards),
		}
		if player.GetPalyerId() == curPlayerID {
			info.BuCards = utils.ServerCards2Uint32(mjContext.WallCards[0:buCardNum])
			player.HandCards = append(player.HandCards, mjContext.WallCards[0:buCardNum]...)
			for _, card := range huaCards {
				var ok bool
				player.HandCards, ok = utils.RemoveCards(player.HandCards, card, 1)
				player.HuaCards = append(player.GetHuaCards(), card)
				if !ok {
					logrus.WithFields(logrus.Fields{
						"func_name":       "GameStartBuhuaState.getHuaCards",
						"hand_cards":      player.HandCards,
						"buhua_player_id": player.GetPalyerId(),
					}).Errorln("移除补花者的花牌失败")
				}
			}
			mjContext.WallCards = mjContext.WallCards[buCardNum:]
		}
		toClientMessage := interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_ROOM_BUHUA_NTF),
			Msg: &room.RoomBuHuaNtf{
				BuhuaInfo: []*room.RoomBuHuaInfo{info},
			},
		}
		logrus.WithFields(logrus.Fields{
			"buhua_player":  curPlayerID,
			"ntf_to_player": player.GetPalyerId(),
			"hua_cards":     info.GetOutHuaCards(),
			"bu_cards":      info.GetBuCards(),
		}).Info("补花通知")
		flow.PushMessages([]uint64{player.GetPalyerId()}, toClientMessage)
	}
}

// OnEntry 进入状态前需要一个补花完成的自动事件
func (bh *GameStartBuhuaState) OnEntry(flow interfaces.MajongFlow) {
	bh.doBuhua(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_gamestart_buhua_finish,
		EventContext: nil,
		WaitTime:     0,
	})
}

// OnExit 离开状态前需要做什么
func (bh *GameStartBuhuaState) OnExit(flow interfaces.MajongFlow) {

}
