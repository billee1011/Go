package ddz

import (
	"fmt"
	"steve/entity/poker/ddz"
	. "steve/room/desks/ddzdesk/flow/ddz/states"
	"steve/room/interfaces"

	"github.com/Sirupsen/logrus"
	"steve/entity/poker"
)

type playStateAI struct {
}

// GenerateAIEvent 生成 出牌AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (playAI *playStateAI) GenerateAIEvent(params interfaces.AIEventGenerateParams) (result interfaces.AIEventGenerateResult, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "play.go:GenerateAIEvent()"})

	// 产生的事件结果
	result, err = interfaces.AIEventGenerateResult{
		Events: []interfaces.AIEvent{},
	}, nil

	ddzContext := params.DDZContext

	// 没到自己打牌
	if ddzContext.GetCurrentPlayerId() != params.PlayerID {
		return result, nil
	}

	// 当前玩家
	var curPlayer *ddz.Player
	for _, player := range ddzContext.GetPlayers() {
		if player.GetPlayerId() == params.PlayerID {
			curPlayer = player
		}
	}

	// 无效玩家
	if curPlayer == nil {
		logEntry.Errorf("无效玩家%d", params.PlayerID)
		return result, fmt.Errorf("无效玩家%d", params.PlayerID)
	}

	// 没有牌型时说明是主动打牌
	if ddzContext.GetCurCardType() == poker.CardType_CT_NONE {

		// 主动产生
		if event := playAI.getActivePlayCardEvent(ddzContext, curPlayer); event != nil {
			result.Events = append(result.Events, *event)
		}
	} else {
		// 被动产生
		if event := playAI.getPassivePlayCardEvent(ddzContext, curPlayer); event != nil {
			result.Events = append(result.Events, *event)
		}
	}

	return
}

// Play 生成出牌请求事件(被动出牌)
func (playAI *playStateAI) getPassivePlayCardEvent(ddzContext *ddz.DDZContext, player *ddz.Player) *interfaces.AIEvent {
	// 最终打出去的牌
	resultCards := []uint32{}

	// 最终打出去的牌型
	resultCardType := poker.CardType_CT_NONE

	// 玩家手中的牌
	handCards := player.GetHandCards()

	// 转换为poke
	handPokes := ToDDZCards(handCards)

	// 按照排序权重进行排序
	//DDZPokerSort(handPokes)

	// 当前牌型
	curCardType := ddzContext.GetCurCardType()

	// 上家出的牌，转换为poke
	curOutPokes := ToDDZCards(ddzContext.GetCurOutCards())

	// 是否有压制的牌
	bSuc := false

	// 压制牌的数组
	sendPukes := []Poker{}

	switch curCardType {
	case poker.CardType_CT_SINGLE: // 单牌
		bSuc, sendPukes = GetMinBiggerSingle(handPokes, curOutPokes)
		break
	case poker.CardType_CT_PAIR: // 对子
		bSuc, sendPukes = GetMinBiggerPair(handPokes, curOutPokes)
		break
	case poker.CardType_CT_SHUNZI: // 顺子
		bSuc, sendPukes = GetMinBiggerShunzi(handPokes, curOutPokes)
		break
	case poker.CardType_CT_PAIRS: // 连对
		bSuc, sendPukes = GetMinBiggerPairs(handPokes, curOutPokes)
		break
	case poker.CardType_CT_TRIPLE: // 三张
		bSuc, sendPukes = GetMinBiggerTriple(handPokes, curOutPokes)
		break
	case poker.CardType_CT_3AND1: // 三带一
		bSuc, sendPukes = GetMinBigger3And1(handPokes, curOutPokes)
		break
	case poker.CardType_CT_3AND2: // 三带二
		bSuc, sendPukes = GetMinBigger3And2(handPokes, curOutPokes)
		break
	case poker.CardType_CT_TRIPLES: // 飞机
		bSuc, sendPukes = GetMinBiggerTriples(handPokes, curOutPokes)
		break
	case poker.CardType_CT_3SAND1S: // 飞机带翅膀1，例：JJJQQQKKK + 856
		bSuc, sendPukes = GetMinBigger3sAnd1s(handPokes, curOutPokes)
		break
	case poker.CardType_CT_3SAND2S: // 飞机带翅膀2，例：JJJQQQKKK + 885566
		bSuc, sendPukes = GetMinBigger3sAnd2s(handPokes, curOutPokes)
		break
	case poker.CardType_CT_4SAND1S: // 四带两个单张
		bSuc, sendPukes = GetMinBigger4sAnd1s(handPokes, curOutPokes)
		break
	case poker.CardType_CT_4SAND2S: // 四带两个对子
		bSuc, sendPukes = GetMinBigger4sAnd2s(handPokes, curOutPokes)
		break
	case poker.CardType_CT_BOMB: // 炸弹
		bSuc, sendPukes = GetMinBiggerBoom(handPokes, curOutPokes)
		break
	case poker.CardType_CT_KINGBOMB: // 火箭
		bSuc, sendPukes = GetMinBiggerKingBoom(handPokes, curOutPokes)
		break
	}

	// 有压制的牌，则出的牌和上家牌型一致
	if bSuc {
		resultCardType = curCardType
	} else {

		// 无压制的牌，且当前牌型是炸弹，则判断自己有无火箭
		if !bSuc && curCardType == poker.CardType_CT_BOMB {
			bSuc, sendPukes = GetKingBoom(handPokes)

			if bSuc {
				resultCardType = poker.CardType_CT_KINGBOMB
			}
		}

		// 无压制的牌，且当前牌型不是炸弹，也不是火箭，则判断自己有无炸弹，无炸弹时再检测火箭
		if !bSuc && curCardType != poker.CardType_CT_BOMB && curCardType != poker.CardType_CT_KINGBOMB {

			// 优先检测炸弹
			bSuc, sendPukes = GetBoom(handPokes)
			if bSuc {
				resultCardType = poker.CardType_CT_BOMB // 用炸弹来压
			} else {
				// 无炸弹时检测有无火箭
				bSuc, sendPukes = GetKingBoom(handPokes)
				if bSuc {
					resultCardType = poker.CardType_CT_KINGBOMB // 用火箭来压
				}
			}
		}
	}

	// 下面是回复消息

	// 有压制的牌，转换数组
	if bSuc {
		resultCards = ToInts(sendPukes)
	}

	logrus.Info("托管被动出牌：%v", resultCards)

	request := &ddz.PlayCardRequestEvent{
		Head: &ddz.RequestEventHead{
			PlayerId: player.GetPlayerId()},
		Cards:    resultCards,    // 打出去的牌
		CardType: resultCardType, // 打出去的牌型
	}

	event := interfaces.AIEvent{
		ID:      int32(ddz.EventID_event_chupai_request),
		Context: request,
	}

	return &event
}

// Play 生成出牌请求事件(主动出牌)
func (playAI *playStateAI) getActivePlayCardEvent(ddzContext *ddz.DDZContext, player *ddz.Player) *interfaces.AIEvent {

	// 玩家手中的牌
	handCards := player.GetHandCards()

	// 转换为poke
	handPokes := ToDDZCards(handCards)

	// 按照排序权重进行排序
	DDZPokerSort(handPokes)

	// 最终打出去的牌（打最小的那个牌）
	resultCards := []uint32{handPokes[0].ToInt()}

	logrus.Info("托管主动出牌：%v", resultCards)

	// 最终打出去的牌型（单张）
	resultCardType := poker.CardType_CT_SINGLE

	// 下面是回复消息
	request := &ddz.PlayCardRequestEvent{
		Head: &ddz.RequestEventHead{
			PlayerId: player.GetPlayerId()},
		Cards:    resultCards,    // 打出去的牌
		CardType: resultCardType, // 打出去的牌型
	}

	event := interfaces.AIEvent{
		ID:      int32(ddz.EventID_event_chupai_request),
		Context: request,
	}

	return &event
}
