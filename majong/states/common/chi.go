package common

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/gutils"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// ChiState 吃状态
type ChiState struct {
}

var _ interfaces.MajongState = new(XipaiState)

// ProcessEvent 处理事件
func (s *ChiState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_chi_finish:
		mjContext := flow.GetMajongContext()
		mjContext.ZixunType = majongpb.ZixunType_ZXT_CHI
		return majongpb.StateID_state_zixun, nil
	}
	return majongpb.StateID_state_chi, nil
}

// OnEntry 进入状态
func (s *ChiState) OnEntry(flow interfaces.MajongFlow) {
	s.doChi(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_chi_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *ChiState) OnExit(flow interfaces.MajongFlow) {
	s.clearInfo(flow)
}

func (s *ChiState) clearInfo(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	playerID := mjContext.GetLastChiPlayer()
	player := utils.GetPlayerByID(mjContext.GetPlayers(), playerID)
	player.DesignChiCards = make([]*majongpb.Card, 0)
}

// doChi 执行碰操作
func (s *ChiState) doChi(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	card := mjContext.GetLastOutCard()
	chiPlayerID := mjContext.GetLastChiPlayer()
	chiPlayer := utils.GetMajongPlayer(chiPlayerID, mjContext)
	srcPlayerID := mjContext.GetLastChupaiPlayer()
	srcPlayer := utils.GetMajongPlayer(srcPlayerID, mjContext)

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":       "ChiState.doChi",
		"chi_player_id":   chiPlayerID,
		"chied_player_id": srcPlayerID,
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	gutils.SortCards(chiPlayer.GetDesignChiCards())
	checkCards := utils.ServerCards2Numbers(chiPlayer.GetDesignChiCards())
	if len(checkCards) != 3 {
		return
	}
	if checkCards[0]+1 != checkCards[1] || checkCards[0]+2 != checkCards[2] {
		return
	}
	// 从被吃玩家的outCards移除被吃牌
	srcOutCards := srcPlayer.GetOutCards()
	srcPlayer.OutCards = removeLastCard(logEntry, srcOutCards, card)
	// 从吃牌玩家的handCards移除吃牌
	newCards := make([]*majongpb.Card, 0)
	newCards = append(newCards, chiPlayer.GetHandCards()...)
	for _, designCard := range chiPlayer.DesignChiCards {
		if utils.CardEqual(card, designCard) {
			continue
		}
		var ok bool
		newCards, ok = utils.RemoveCards(newCards, designCard, 1)
		if !ok {
			logEntry.Errorln("移除玩家手牌失败")
			return
		}
	}
	chiPlayer.HandCards = newCards
	s.notifyChi(flow, chiPlayer.DesignChiCards, card, srcPlayerID, chiPlayerID)
	s.addChiCard(chiPlayer.DesignChiCards[0], card, chiPlayer, srcPlayerID)
	logEntry = logEntry.WithFields(logrus.Fields{
		"chi_cards": checkCards,
		"out_card":  utils.ServerCard2Number(card),
		"chiCard":   chiPlayer.GetChiCards(),
	})
	logEntry.Infoln("吃牌成功")
	return
}

// addChiCard 添加碰的牌
func (s *ChiState) addChiCard(card *majongpb.Card, oprCard *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	player.ChiCards = append(player.GetChiCards(), &majongpb.ChiCard{
		Card:      card,
		OprCard:   oprCard,
		SrcPlayer: srcPlayerID,
	})
}

func (s *ChiState) notifyChi(flow interfaces.MajongFlow, cards []*majongpb.Card, chiCard *majongpb.Card, from uint64, to uint64) {
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_CHI_NTF, &room.RoomChiNtf{
		Cards:        utils.ServerCards2Uint32(cards),
		ChiCard:      proto.Uint32(utils.ServerCard2Uint32(chiCard)),
		FromPlayerId: proto.Uint64(from),
		ToPlayerId:   proto.Uint64(to),
	})
}
