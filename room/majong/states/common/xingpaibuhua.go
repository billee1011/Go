package common

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	majongpb "steve/entity/majong"
	"steve/room/majong/interfaces"
	"steve/room/majong/utils"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// XingPaiBuhuaState 行牌补花状态
type XingPaiBuhuaState struct {
}

var _ interfaces.MajongState = new(XipaiState)

// ProcessEvent 处理事件
func (s *XingPaiBuhuaState) ProcessEvent(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_xingpai_buhua_finish:
		return s.doBuhua(flow), nil
	}
	return majongpb.StateID_state_xingpai_buhua, nil
}

// OnEntry 进入状态
func (s *XingPaiBuhuaState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_xingpai_buhua_finish,
		EventContext: nil,
		WaitTime:     0,
	})
}

// OnExit 退出状态
func (s *XingPaiBuhuaState) OnExit(flow interfaces.MajongFlow) {
}

func (s *XingPaiBuhuaState) doBuhua(flow interfaces.MajongFlow) majongpb.StateID {
	//补到没花可补，进入自询
	stop := false
	mjContext := flow.GetMajongContext()
	players := mjContext.GetPlayers()
	activePlayer := utils.GetPlayerByID(players, mjContext.GetLastMopaiPlayer())
	for !stop {
		huaCards := s.getHuaCards(activePlayer)
		if len(huaCards) > 0 {
			if utils.HasAvailableWallCards(flow) {
				s.ntf(flow, players, mjContext.GetLastMopaiPlayer(), huaCards, len(huaCards))
			} else {
				return majongpb.StateID_state_gameover
			}
		} else {
			stop = true
		}
	}
	return majongpb.StateID_state_zixun
}
func (s *XingPaiBuhuaState) ntf(flow interfaces.MajongFlow, players []*majongpb.Player, curPlayerID uint64, huaCards []*majongpb.Card, buCardNum int) {
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
			buCard := mjContext.WallCards[0:buCardNum]
			info.BuCards = utils.ServerCards2Uint32(buCard)
			player.HandCards = append(player.HandCards, buCard...)
			mjContext.LastMopaiCard = buCard[0]
			for _, card := range huaCards {
				var ok bool
				player.HandCards, ok = utils.RemoveCards(player.HandCards, card, 1)
				player.HuaCards = append(player.GetHuaCards(), card)
				if !ok {
					logrus.WithFields(logrus.Fields{
						"func_name":       "XingPaiBuhuaState.getHuaCards",
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

func (s *XingPaiBuhuaState) getHuaCards(player *majongpb.Player) []*majongpb.Card {
	handCards := player.GetHandCards()
	huaCards := []*majongpb.Card{}
	for _, card := range handCards {
		if card.GetColor() == majongpb.CardColor_ColorHua {
			huaCards = append(huaCards, card)
		}
	}
	return huaCards
}
