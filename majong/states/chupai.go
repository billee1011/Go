package states

import (
	"steve/clientpb"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// ChupaiState 初始化状态
type ChupaiState struct {
}

var _ interfaces.MajongState = new(ChupaiState)

// ProcessEvent 处理事件
func (s *ChupaiState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_chupai_finish {
		context := flow.GetMajongContext()
		players := context.GetPlayers()
		card := context.GetLastOutCard()
		var hasChupaiwenxun bool
		for _, player := range players {
			actionInfos := checkActions(context, player, card)
			if len(actionInfos) > 0 {
				//TODO暂时不广播
				hasChupaiwenxun = true
			}
		}
		if hasChupaiwenxun {
			return majongpb.StateID_state_chupaiwenxun, nil
		}
		s.mopai(flow)
		return majongpb.StateID_state_chupai, nil
	}
	return majongpb.StateID_state_init, errInvalidEvent
}

//checkActions 检查玩家可以有哪些操作
func checkActions(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) []*clientpb.ActionInfo {
	actionInfos := []*clientpb.ActionInfo{}
	canMingGang, mingGangInfo := checkMingGang(context, player, card)
	if canMingGang {
		actionInfos = append(actionInfos, mingGangInfo)
		player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_minggang)
	}
	canDianPao, dianPaoInfo := checkDianPao(context, player, card)
	if canDianPao {
		actionInfos = append(actionInfos, dianPaoInfo)
		player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_dianpao)
	}
	return actionInfos
}

//checkMingGang 查明杠
func checkMingGang(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) (bool, *clientpb.ActionInfo) {
	// 没有墙牌不能明杠
	if len(context.WallCards) == 0 {
		return false, nil
	}
	cpPlayerID := context.GetActivePlayer()
	cpPlayer := utils.GetPlayerByID(context.GetPlayers(), cpPlayerID)
	outCard := context.GetLastOutCard()
	color := player.GetDingqueColor()
	//定缺牌不查
	if outCard.Color == color {
		return false, nil
	}
	if cpPlayer.PalyerId != player.PalyerId {
		cards := player.HandCards
		num := 0
		for _, card := range cards {
			if utils.CardEqual(card, outCard) {
				num++
			}
		}
		if num == 3 {
			if len(player.GetHuCards()) > 0 {
				//创建副本，移除相应的杠牌进行查胡
				newcards := make([]*majongpb.Card, 0, len(cards))
				newcards = append(newcards, cards...)
				newcards, _ = utils.DeleteCardFromLast(newcards, outCard)
				newcards, _ = utils.DeleteCardFromLast(newcards, outCard)
				newcards, _ = utils.DeleteCardFromLast(newcards, outCard)
				newcardsI, _ := utils.CardsToInt(newcards)
				cardsI := utils.IntToUtilCard(newcardsI)
				laizi := make(map[utils.Card]bool)
				huCards := utils.FastCheckTingV2(cardsI, laizi)
				if utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(player.HuCards)) {
					//暂时不广播消息
					cardToClient, _ := utils.CardToInt(*outCard)
					actionInfo := &clientpb.ActionInfo{
						ActionID:    clientpb.ActionID_MingGang.Enum(),
						ActionCards: []uint32{uint32(*cardToClient)},
						FromPid:     proto.Uint64(cpPlayer.GetPalyerId()),
						Pid:         proto.Uint64(player.GetPalyerId()),
					}
					return true, actionInfo
				}
			} else {
				//暂时不广播消息
				cardToClient, _ := utils.CardToInt(*outCard)
				actionInfo := &clientpb.ActionInfo{
					ActionID:    clientpb.ActionID_MingGang.Enum(),
					ActionCards: []uint32{uint32(*cardToClient)},
					FromPid:     proto.Uint64(cpPlayer.GetPalyerId()),
					Pid:         proto.Uint64(player.GetPalyerId()),
				}
				return true, actionInfo
			}
		}
	}
	return false, &clientpb.ActionInfo{}
}

//checkDianPao 查点炮
func checkDianPao(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) (bool, *clientpb.ActionInfo) {
	cpPlayer := utils.GetPlayerByID(context.GetPlayers(), context.ActivePlayer)
	cpCard := context.GetLastOutCard()
	if cpPlayer.PalyerId != player.PalyerId {
		color := player.GetDingqueColor()
		hasDingQueCard := utils.CheckHasDingQueCard(player.HandCards, color)
		if hasDingQueCard {
			return false, nil
		}
		handCard := player.GetHandCards() // 当前点炮胡玩家手牌
		cardI, _ := utils.CardToInt(*cpCard)
		flag := utils.CheckHu(handCard, uint32(*cardI))
		if flag {
			actionInfo := &clientpb.ActionInfo{
				ActionID:    clientpb.ActionID_DianPao.Enum(),
				ActionCards: []uint32{uint32(*cardI)},
				FromPid:     proto.Uint64(cpPlayer.PalyerId),
				Pid:         proto.Uint64(player.PalyerId),
			}
			return true, actionInfo
		}
	}
	return false, nil
}

//mopai 摸牌处理
func (s *ChupaiState) mopai(flow interfaces.MajongFlow) (majongpb.StateID, error) {
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	activePlayer := utils.GetNextPlayerByID(players, context.ActivePlayer)
	//TODO：目前只在这个地方改变操作玩家（感觉碰，明杠，点炮这三种情况也需要改变activePlayer）
	context.ActivePlayer = activePlayer.GetPalyerId()
	//从墙牌中移除一张牌
	drowCard := context.WallCards[0]
	context.WallCards = context.WallCards[1:]
	//将这张牌添加到手牌中
	activePlayer.HandCards = append(activePlayer.HandCards, drowCard)
	return majongpb.StateID_state_zixun, nil
}

// OnEntry 进入状态
func (s *ChupaiState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_chupai_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *ChupaiState) OnExit(flow interfaces.MajongFlow) {

}
