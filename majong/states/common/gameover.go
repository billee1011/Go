package common

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// GameOverState 游戏结束状态
type GameOverState struct {
}

var _ interfaces.MajongState = new(GameOverState)

// ProcessEvent 处理事件
func (s *GameOverState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	return majongpb.StateID_state_gameover, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *GameOverState) OnEntry(flow interfaces.MajongFlow) {
	s.notifyGameOver(flow)
	s.doRoundSettle(flow)
}

// OnExit 退出状态
func (s *GameOverState) OnExit(flow interfaces.MajongFlow) {

}

// notifyGameOver 通知玩家游戏结束消息
func (s *GameOverState) notifyGameOver(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	cardsGroups := make([]*room.PlayerCardsGroup, 0)
	gameflow := true
	for _, player := range mjContext.Players {
		playerCardsGroup := &room.PlayerCardsGroup{
			PlayerId:   proto.Uint64(player.GetPalyerId()),
			CardsGroup: utils.GetCardsGroup(player),
		}
		cardsGroups = append(cardsGroups, playerCardsGroup)
		if len(player.HuCards) != 0 && gameflow {
			gameflow = false
		}
	}
	roomGameOverNtf := &room.RoomGameOverNtf{
		PlayerCardsGroup: cardsGroups,
		GameFlow:         proto.Bool(gameflow),
	}
	// 广播牌局结束消息
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_GAMEOVER_NTF, roomGameOverNtf)
	// 日志
	logrus.WithFields(logrus.Fields{
		"msgID":           msgid.MsgID_ROOM_GAMEOVER_NTF,
		"roomGameOverNtf": roomGameOverNtf,
	}).Info("-----牌局结束-推倒牌墙")
}

// roundSettle 处理查花猪，查大叫，退税 结算
func (s *GameOverState) doRoundSettle(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	// 花猪
	flowerPigPlayers := make([]uint64, 0)
	// 胡玩家
	huPlayers := make([]uint64, 0)
	// 未听玩家
	noTingPlayers := make([]uint64, 0)
	// 听牌未胡玩家信息
	tingPlayersInfo := make(map[uint64]int64)

	for _, player := range mjContext.Players {
		if len(player.HuCards) != 0 {
			huPlayers = append(huPlayers, player.GetPalyerId())

		}
		if isFlowerPig(player) {
			flowerPigPlayers = append(flowerPigPlayers, player.GetPalyerId())
		}
		if isNoTingPlayers(player) {
			noTingPlayers = append(noTingPlayers, player.GetPalyerId())
		}
	}
	tingPlayersInfo, _ = getTingPlayerInfo(mjContext)

	params := interfaces.RoundSettleParams{
		FlowerPigPlayers: flowerPigPlayers,
		HuPlayers:        huPlayers,
		TingPlayersInfo:  tingPlayersInfo,
		NotTingPlayers:   noTingPlayers,
		SettleInfos:      mjContext.SettleInfos,
		SettleID:         mjContext.CurrentSettleId,
	}
	settleInfos, raxbeatIds := facade.SettleRound(global.GetGameSettlerFactory(), int(mjContext.GetGameId()), params)
	for _, settleInfo := range settleInfos {
		mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
	}
	mjContext.RevertSettles = raxbeatIds
}

//isFlowerPig 修改为： 判断玩家是否是花猪,牌局结束结束后该玩家手上还有定缺牌，此时该玩家被查花猪
func isFlowerPig(player *majongpb.Player) bool {
	for _, card := range player.HandCards {
		if card.Color == player.DingqueColor {
			return true
		}
	}
	return false
}

// isNoTingPlayers 判断玩家是否未听，不包括花猪，因为查花猪包括了查大叫，所以未听玩家，中是花猪的，都不用再进行查大叫
func isNoTingPlayers(player *majongpb.Player) bool {
	// 胡过的不算
	if len(player.HuCards) > 0 {
		return false
	}

	if !hasDingQueCard(player.HandCards, player.DingqueColor) { // 手牌中没有定缺牌，检查该玩家是否可听，不可听返回true
		// 查听
		tingCards, _ := utils.GetTingCards(player.HandCards, nil)
		// 不能听
		if len(tingCards) == 0 {
			return true
		}
	} else { //  手牌中若有定缺牌，必是花猪
		// if !isFlowerPig(player) {
		// 	return true
		// }
		return false
	}
	return false
}

// getTingPlayerInfo 判断玩家是否能听,和返回能听玩家的最大倍数
// 未上听者需赔上听者最大可能番数（自摸、杠后炮、杠上开花、抢杠胡、海底捞、海底炮不参与）的牌型钱。注：查大叫时，若上听者牌型中有根，则根也要未上听者包给上听者。
func getTingPlayerInfo(context *majongpb.MajongContext) (map[uint64]int64, error) {
	players := context.Players
	tingPlayers := make(map[uint64]int64, 0)
	for i := 0; i < len(players); i++ {
		// 胡过的不算
		if len(players[i].HuCards) > 0 {
			continue
		}
		handCardSum := len(players[i].HandCards)
		var maxMulti int64
		//只差1张牌就能胡，并且玩家手牌不存在花牌
		if handCardSum%3 == 1 && !hasDingQueCard(players[i].HandCards, players[i].DingqueColor) {
			tingCards, err := utils.GetTingCards(players[i].HandCards, nil)
			if err != nil {
				return nil, err
			}
			for j := 0; j < len(tingCards); j++ {
				// 获取最大番型*根数
				cardParams := interfaces.CardCalcParams{
					HandCard: players[i].HandCards,
					PengCard: utils.TransPengCard(players[i].PengCards),
					GangCard: utils.TransGangCard(players[i].GangCards),
					HuCard:   tingCards[j],
					GameID:   int(context.GetGameId()),
				}
				calculator := global.GetCardTypeCalculator()
				total, _ := facade.CalculateCardValue(calculator, cardParams)
				if maxMulti < int64(total) {
					maxMulti = int64(total)
				}
			}
			if len(tingCards) != 0 {
				tingPlayers[players[i].GetPalyerId()] = maxMulti
			}
		}
	}
	return tingPlayers, nil
}

//hasDingQueCard 检查牌里面是否含有定缺的牌
func hasDingQueCard(cards []*majongpb.Card, color majongpb.CardColor) bool {
	for _, card := range cards {
		if card.Color == color {
			return true
		}
	}
	return false
}
