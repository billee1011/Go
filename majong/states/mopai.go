package states

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

// MoPaiState 摸牌状态
type MoPaiState struct {
}

var _ interfaces.MajongState = new(MoPaiState)

// ProcessEvent 处理事件
func (s *MoPaiState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_mopai_finish {
		return s.mopai(flow)
	}
	return majongpb.StateID_state_mopai, global.ErrInvalidEvent
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
	playerIDs := make([]uint64, 0, 0)
	playerIDs = append(playerIDs, context.MopaiPlayer)
	toClient := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_ZIXUN_NTF),
		Msg:   zixunNtf,
	}
	flow.PushMessages(playerIDs, toClient)
	logrus.WithFields(logrus.Fields{
		"canZimo":   canZiMo,
		"canAnGang": canAnGang,
		"canBuGang": canBuGang,
	}).Infof("玩家%v可以有的操作", context.MopaiPlayer)
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
	return flag
}

//checkAnGang 查暗杠
func (s *MoPaiState) checkAnGang(context *majongpb.MajongContext) (bool, []uint32) {
	if len(context.WallCards) == 0 {
		return false, nil
	}
	activePlayerID := context.GetMopaiPlayer()
	activePlayer := utils.GetPlayerByID(context.Players, activePlayerID)
	//分两种情况查暗杠，一种是胡牌前，一种胡牌后
	hasHu := len(activePlayer.GetHuCards()) > 0
	handCard := activePlayer.GetHandCards()
	enableAngangCards := make([]uint32, 0, 0)
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
				huCards := utils.FastCheckTingV2(cardsI, map[utils.Card]bool{})
				if utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
					enableAngangCards = append(enableAngangCards, uint32(k))
				}
			} else {
				enableAngangCards = append(enableAngangCards, uint32(k))
			}
		}
	}
	return len(enableAngangCards) > 0, enableAngangCards
}

//checkBuGang 查补杠
func (s *MoPaiState) checkBuGang(context *majongpb.MajongContext) (bool, []uint32) {
	// 没有墙牌不能杠
	if len(context.WallCards) == 0 {
		return false, nil
	}
	activePlayerID := context.GetActivePlayer()
	activePlayer := utils.GetPlayerByID(context.Players, activePlayerID)
	//分两种情况查暗杠，一种是胡牌前，一种胡牌后
	hasHu := len(activePlayer.GetHuCards()) > 0
	pengCards := activePlayer.GetPengCards()
	enableBugangCards := make([]uint32, 0, 0)
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
						enableBugangCards = append(enableBugangCards, uint32(*removeCard))
					}
				} else {
					enableBugangCards = append(enableBugangCards, uint32(*removeCard))
				}
			}
		}
	}
	return len(enableBugangCards) > 0, enableBugangCards
}

//mopai 摸牌处理
func (s *MoPaiState) mopai(flow interfaces.MajongFlow) (majongpb.StateID, error) {
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	activePlayer := utils.GetPlayerByID(players, context.MopaiPlayer)
	//TODO：目前只在这个地方改变操作玩家（感觉碰，明杠，点炮这三种情况也需要改变activePlayer）
	context.ActivePlayer = activePlayer.GetPalyerId()
	if len(context.WallCards) == 0 {
		return majongpb.StateID_state_gameover, nil
	}
	//从墙牌中移除一张牌
	drowCard := context.WallCards[0]
	context.WallCards = context.WallCards[1:]
	//将这张牌添加到手牌中
	activePlayer.HandCards = append(activePlayer.HandCards, drowCard)
	context.LastMopaiPlayer = context.MopaiPlayer
	context.LastMopaiCard = drowCard
	s.checkActions(flow)
	// 清空其他玩家杠的标识
	for _, player := range players {
		if player.PalyerId != activePlayer.PalyerId {
			player.Properties["gang"] = []byte("false")
		}
	}
	return majongpb.StateID_state_zixun, nil
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
