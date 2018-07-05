package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"steve/majong/global"
	"github.com/gogo/protobuf/proto"
	"steve/client_pb/room"
	"steve/client_pb/msgId"
	"github.com/pkg/errors"
)

type playState struct{}

func (s *playState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入出牌状态")
}

func (s *playState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开出牌状态")
}

func (s *playState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID != int(ddz.EventID_event_chupai_request) {
		return int(ddz.StateID_state_playing), global.ErrInvalidEvent
	}

	message := &ddz.PlayCardRequestEvent{}
	err := proto.Unmarshal(event.EventData, message)
	if err != nil {
		return int(ddz.StateID_state_playing), global.ErrUnmarshalEvent
	}

	context := getDDZContext(m);
	playerId := message.GetHead().GetPlayerId()
	if context.CurrentPlayerId != playerId {
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: &room.Result{ErrCode:proto.Uint32(1), ErrDesc: proto.String("未轮到本玩家出牌")},
		})
		return int(ddz.StateID_state_playing), global.ErrInvalidRequestPlayer
	}

	player := GetPlayerByID(context.GetPlayers(), playerId)
	outCards := toDDZCards(message.GetCards())
	nextPlayerId := GetNextPlayerByID(context.GetPlayers(), playerId).PalyerId

	if len(outCards) == 0 {//pass
		if context.CurCardType == ddz.CardType_CT_NONE {//该你出牌时不出牌，报错
			sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
				Result: &room.Result{ErrCode:proto.Uint32(6), ErrDesc: proto.String("首轮出牌玩家不能过牌")},
			})
			return int(ddz.StateID_state_playing), errors.New("首轮出牌玩家不能过牌")
		}
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{//成功pass
			Result: &room.Result{ErrCode:proto.Uint32(0), ErrDesc: proto.String("")},
		})

		stage := room.DDZStage_DDZ_STAGE_PLAYING
		broadcastExcept(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF, &room.DDZPlayCardNtf{//广播pass
			PlayerId: &playerId,
			Cards: message.GetCards(),
			CardType: nil,
			TotalBomb: &context.TotalBomb,
			NextPlayerId: &nextPlayerId,
			NextStage: &room.NextStage{
				Stage: &stage,
				Time: proto.Uint32(15),
			},
		})
		context.PassCount++
		if context.PassCount >= 2 {//两个玩家都过，清空当前牌型
			context.CurCardType = ddz.CardType_CT_NONE
			context.CurOutCards = []uint32{}
			context.CardTypePivot = 0
		}
		return int(ddz.StateID_state_playing), nil
	}

	handCards := toDDZCards(player.HandCards)
	if !ContainsAll(handCards, outCards ){//检查所出的牌是否在手牌中
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: &room.Result{ErrCode:proto.Uint32(2), ErrDesc: proto.String("所出的牌不在手牌中")},
		})
		return int(ddz.StateID_state_playing), errors.New("所出的牌不在手牌中")
	}

	cardType, pivot := getCardType(outCards)
	if cardType == ddz.CardType_CT_NONE {//检查所出的牌能否组成牌型
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: &room.Result{ErrCode:proto.Uint32(3), ErrDesc: proto.String("无法组成牌型")},
		})
		return int(ddz.StateID_state_playing), errors.New("无法组成牌型")
	}

	if context.CurCardType != ddz.CardType_CT_NONE &&
		(!canBiggerThan(cardType, context.CurCardType) || //牌型与上家不符(炸弹不算不符)
			(context.CurCardType == ddz.CardType_CT_SHUNZI && cardType == ddz.CardType_CT_SHUNZI && len(outCards) != len(context.CurOutCards))) {//顺子牌数不足
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: &room.Result{ErrCode:proto.Uint32(4), ErrDesc: proto.String("牌型与上家不符")},
		})
		return int(ddz.StateID_state_playing), errors.New("牌型与上家不符")
	}

	lastPivot := toDDZCard(context.CardTypePivot)
	currPivot := *pivot
	if lastPivot.pointBiggerThan(currPivot) {
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: &room.Result{ErrCode:proto.Uint32(5), ErrDesc: proto.String("牌比上家小")},
		})
		return int(ddz.StateID_state_playing), errors.New("牌比上家小")
	}

	//更新玩家手牌和已出的牌
	handCards = RemoveAll(handCards, outCards)
	player.HandCards = toInts(handCards)
	lastOutCards := toDDZCards(player.OutCards)
	lastOutCards = AppendAll(lastOutCards, outCards)
	player.OutCards = toInts(lastOutCards)

	//更新context
	context.CurrentPlayerId = nextPlayerId
	context.LastPlayerId = playerId
	context.CurOutCards = message.GetCards()
	context.CurCardType = cardType
	context.CardTypePivot = (*pivot).toInt()
	if cardType == ddz.CardType_CT_BOMB || cardType == ddz.CardType_CT_KINGBOMB {
		context.TotalBomb = context.TotalBomb * 2
	}
	if playerId != context.LordPlayerId {
		context.Spring = false //农民出牌了，没有春天了
	}
	if context.Spring == false && playerId == context.LordPlayerId {
		context.AntiSpring = false // 地主第二次出牌了，没有反春天了
	}

	sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{//成功出牌
		Result: &room.Result{ErrCode:proto.Uint32(0), ErrDesc: proto.String("")},
	})

	var nextStage room.DDZStage
	if len(player.HandCards) == 0 {
		nextStage = room.DDZStage_DDZ_STAGE_OVER
	} else {
		nextStage = room.DDZStage_DDZ_STAGE_PLAYING
	}
	clientCardType := room.CardType(int32(cardType))
	broadcastExcept(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF, &room.DDZPlayCardNtf{//广播出牌
		PlayerId: &playerId,
		Cards: message.GetCards(),
		CardType: &clientCardType,
		TotalBomb: &context.TotalBomb,
		NextPlayerId: &nextPlayerId,
		NextStage: &room.NextStage{
			Stage: &nextStage,
			Time: proto.Uint32(15),
		},
	})

	if len(player.HandCards) == 0 {
		return int(ddz.StateID_state_over), nil
	} else {
		return int(ddz.StateID_state_playing), nil
	}
}
