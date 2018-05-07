package states

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// MoPaiState 摸牌状态
type MoPaiState struct {
}

var _ interfaces.MajongState = new(MoPaiState)

// ProcessEvent 处理事件
func (s *MoPaiState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_mopai_finish {
		mjContext := flow.GetMajongContext()
		wallCards := mjContext.GetWallCards()
		if len(wallCards) == 0 {
			return majongpb.StateID_state_gameover, nil
		}
		//进入自询状态，需要查当前玩家可以有的特殊操作
		s.checkActions(flow)
		return majongpb.StateID_state_zixun, nil
	}
	return majongpb.StateID_state_mopai, errInvalidEvent
}

//checkActions 检测进入自询状态下，玩家有哪些可以可行的事件
func (s *MoPaiState) checkActions(flow interfaces.MajongFlow) {
	context := flow.GetMajongContext()
	zixunNtf := &room.RoomZixunNtf{}
	canZiMo := s.checkZiMo(context)
	zixunNtf.EnableZimo = proto.Bool(canZiMo)
	canAnGang, enablieAngangCards := s.checkAnGang(context)
	if canAnGang {
		zixunNtf.EnableAngangCards = enablieAngangCards
	}
	canBuGang, enablieBugangCards := s.checkBuGang(context)
	if canBuGang {
		zixunNtf.EnableBugangCards = enablieBugangCards
	}
	if canZiMo {
		//TODO:可以出的牌，在胡牌后可能需要
		// enableChupaiCards :=
	}
	playerIDs := make([]uint64, 0, 0)
	playerIDs = append(playerIDs, context.ActivePlayer)
	toClient := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_room_zixun_ntf),
		Msg:   zixunNtf,
	}
	if canAnGang || canBuGang || canZiMo {
		flow.PushMessages(playerIDs, toClient)
	}
}

//checkZiMo 查自摸
func (s *MoPaiState) checkZiMo(context *majongpb.MajongContext) bool {
	activePlayerID := context.GetActivePlayer()
	activePlayer := utils.GetPlayerByID(context.Players, activePlayerID)
	handCard := activePlayer.GetHandCards()
	if utils.CheckHasDingQueCard(handCard, activePlayer.GetDingqueColor()) {
		return false
	}
	l := len(handCard)
	if l%3 != 2 {
		return false
	}
	flag := utils.CheckHu(handCard, 0)
	if flag {
		activePlayer.PossibleActions = append(activePlayer.PossibleActions, majongpb.Action_action_zimo)
	}
	return flag
}

//checkAnGang 查暗杠
func (s *MoPaiState) checkAnGang(context *majongpb.MajongContext) (bool, []*room.Card) {
	if len(context.WallCards) == 0 {
		return false, nil
	}
	activePlayerID := context.GetActivePlayer()
	activePlayer := utils.GetPlayerByID(context.Players, activePlayerID)
	//分两种情况查暗杠，一种是胡牌前，一种胡牌后
	hasHu := len(activePlayer.GetHuCards()) > 0
	handCard := activePlayer.GetHandCards()
	enableAngangCards := make([]*room.Card, 0, 0)
	cardsI, _ := utils.CardsToInt(handCard)
	cardNum := make(map[int32]int)
	for i := 0; i < len(cardsI); i++ {
		num := cardNum[cardsI[i]]
		num++
		cardNum[cardsI[i]] = num
	}
	color := activePlayer.GetDingqueColor()
	for k, num := range cardNum {
		if k/10 != int32(color) && num == 4 {
			if hasHu {
				//创建副本，移除相应的杠牌进行查胡
				newcardsI := make([]int32, 0, len(cardsI))
				newcardsI = append(newcardsI, cardsI...)
				newcardsI, _ = utils.DeleteIntCardFromLast(newcardsI, k)
				newcardsI, _ = utils.DeleteIntCardFromLast(newcardsI, k)
				newcardsI, _ = utils.DeleteIntCardFromLast(newcardsI, k)
				newcardsI, _ = utils.DeleteIntCardFromLast(newcardsI, k)
				cardsI := utils.IntToUtilCard(newcardsI)
				laizi := make(map[utils.Card]bool)
				huCards := utils.FastCheckTingV2(cardsI, laizi)
				if utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
					roomCard, _ := utils.IntToRoomCard(k)
					enableAngangCards = append(enableAngangCards, roomCard)
				}
			} else {
				roomCard, _ := utils.IntToRoomCard(k)
				enableAngangCards = append(enableAngangCards, roomCard)
			}
		}
	}
	if len(enableAngangCards) > 0 {
		activePlayer.PossibleActions = append(activePlayer.PossibleActions, majongpb.Action_action_angang)
	}
	return len(enableAngangCards) > 0, enableAngangCards
}

//checkBuGang 查补杠
func (s *MoPaiState) checkBuGang(context *majongpb.MajongContext) (bool, []*room.Card) {
	// 没有墙牌不能杠
	if len(context.WallCards) == 0 {
		return false, nil
	}
	activePlayerID := context.GetActivePlayer()
	activePlayer := utils.GetPlayerByID(context.Players, activePlayerID)
	//分两种情况查暗杠，一种是胡牌前，一种胡牌后
	hasHu := len(activePlayer.GetHuCards()) > 0
	pengCards := activePlayer.GetPengCards()
	enableBugangCards := make([]*room.Card, 0, 0)
	// actioninfos := []*clientpb.ActionInfo{}
	for _, touchCard := range activePlayer.HandCards {
		for _, pengCard := range pengCards {
			if *pengCard.Card == *touchCard {
				removeCard, _ := utils.CardToInt(*touchCard)
				if hasHu {
					//创建副本，移除相应的杠牌进行查胡
					cardsI, _ := utils.CardsToInt(activePlayer.HandCards)
					newcardsI := make([]int32, 0, len(cardsI))
					newcardsI = append(newcardsI, cardsI...)
					newcardsI, _ = utils.DeleteIntCardFromLast(newcardsI, *removeCard)
					utilCards := utils.IntToUtilCard(newcardsI)
					laizi := make(map[utils.Card]bool)
					huCards := utils.FastCheckTingV2(utilCards, laizi)
					if utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
						roomCard, _ := utils.IntToRoomCard(*removeCard)
						enableBugangCards = append(enableBugangCards, roomCard)
					}
				} else {
					roomCard, _ := utils.IntToRoomCard(*removeCard)
					enableBugangCards = append(enableBugangCards, roomCard)
				}
			}
		}
	}
	if len(enableBugangCards) > 0 {
		activePlayer.PossibleActions = append(activePlayer.PossibleActions, majongpb.Action_action_bugang)
	}
	return len(enableBugangCards) > 0, enableBugangCards
}

// OnEntry 进入状态
func (s *MoPaiState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_mopai_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *MoPaiState) OnExit(flow interfaces.MajongFlow) {

}
