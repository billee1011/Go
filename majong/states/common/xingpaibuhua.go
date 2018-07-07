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

// XingPaiBuhuaState 行牌补花状态
type XingPaiBuhuaState struct {
}

var _ interfaces.MajongState = new(XipaiState)

// ProcessEvent 处理事件
func (s *XingPaiBuhuaState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_xingpai_buhua_finish:
		return majongpb.StateID_state_zixun, nil
	}
	return majongpb.StateID_state_xingpai_buhua, nil
}

// OnEntry 进入状态
func (s *XingPaiBuhuaState) OnEntry(flow interfaces.MajongFlow) {
	s.doBuhua(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_xingpai_buhua_finish,
		EventContext: nil,
		WaitTime:     0,
	})
}

// OnExit 退出状态
func (s *XingPaiBuhuaState) OnExit(flow interfaces.MajongFlow) {
}

func (s *XingPaiBuhuaState) doBuhua(flow interfaces.MajongFlow) {
	//补到没花可补，进入自询
	stop := false
	mjContext := flow.GetMajongContext()
	players := mjContext.GetPlayers()
	activePlayer := utils.GetPlayerByID(players, mjContext.GetLastMopaiPlayer())
	for !stop {
		huaCards := s.getHuaCards(activePlayer)
		if len(huaCards) > 0 {
			s.ntf(flow, players, mjContext.GetLastMopaiPlayer(), huaCards, len(huaCards))
		} else {
			stop = true
		}
	}
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
			info.BuCards = utils.ServerCards2Uint32(mjContext.WallCards[0:buCardNum])
			player.HandCards = append(player.HandCards, mjContext.WallCards[0:buCardNum]...)
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
