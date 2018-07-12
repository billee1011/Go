package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"time"
)

type playState struct{}

func (s *playState) OnEnter(m machine.Machine) {
	context := getDDZContext(m)
	context.CurStage = ddz.DDZStage_DDZ_STAGE_PLAYING
	//产生超时事件
	context.CountDownPlayers = []uint64{context.CurrentPlayerId}
	context.StartTime, _ = time.Now().MarshalBinary()
	context.Duration = StageTime[room.DDZStage_DDZ_STAGE_PLAYING]

	logrus.WithField("context", context).Debugln("进入出牌状态")
}

func (s *playState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开出牌状态")
}

func (s *playState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID != int(ddz.EventID_event_chupai_request) {
		logrus.Error("playState can only handle ddz.EventID_event_chupai_request, invalid event")
		return int(ddz.StateID_state_playing), global.ErrInvalidEvent
	}

	message := &ddz.PlayCardRequestEvent{}
	err := proto.Unmarshal(event.EventData, message)
	if err != nil {
		logrus.Error("playState unmarshal event error!")
		return int(ddz.StateID_state_playing), global.ErrUnmarshalEvent
	}

	context := getDDZContext(m)
	playerId := message.GetHead().GetPlayerId()
	outCards := toDDZCards(message.GetCards())
	logrus.WithField("playerId", playerId).WithField("outCards", outCards).Debug("玩家出牌")
	if context.CurrentPlayerId != playerId {
		logrus.WithField("expected player:", context.CurrentPlayerId).WithField("fact player", playerId).Error("未到本玩家出牌")
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: genResult(1, "未轮到本玩家出牌"),
		})
		return int(ddz.StateID_state_playing), global.ErrInvalidRequestPlayer
	}

	nextPlayerId := GetNextPlayerByID(context.GetPlayers(), playerId).PalyerId
	if len(outCards) == 0 { //pass
		if context.CurCardType == ddz.CardType_CT_NONE { //该你出牌时不出牌，报错
			sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
				Result: genResult(6, "首轮出牌玩家不能过牌"),
			})
			return int(ddz.StateID_state_playing), errors.New("首轮出牌玩家不能过牌")
		}
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{ //成功pass
			Result: genResult(0, ""),
		})

		stage := room.DDZStage_DDZ_STAGE_PLAYING
		broadcast(m, msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF, &room.DDZPlayCardNtf{ //广播pass
			PlayerId:     &playerId,
			Cards:        message.GetCards(),
			CardType:     nil,
			TotalBomb:    &context.TotalBomb,
			NextPlayerId: &nextPlayerId,
			NextStage: &room.NextStage{
				Stage: &stage,
				Time:  proto.Uint32(15),
			},
		})

		context.CurrentPlayerId = nextPlayerId
		//产生超时事件
		context.CountDownPlayers = []uint64{context.CurrentPlayerId}
		context.StartTime, _ = time.Now().MarshalBinary()
		context.Duration = StageTime[room.DDZStage_DDZ_STAGE_PLAYING]

		context.PassCount++
		if context.PassCount >= 2 { //两个玩家都过，清空当前牌型
			context.CurCardType = ddz.CardType_CT_NONE
			context.CurOutCards = []uint32{}
			context.CardTypePivot = 0
			context.PassCount = 0
		}
		return int(ddz.StateID_state_playing), nil
	}

	player := GetPlayerByID(context.GetPlayers(), playerId)
	handCards := toDDZCards(player.HandCards)
	if !ContainsAll(handCards, outCards) { //检查所出的牌是否在手牌中
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: genResult(2, "所出的牌不在手牌中"),
		})
		return int(ddz.StateID_state_playing), errors.New("所出的牌不在手牌中")
	}

	cardType, pivot := getCardType(outCards)
	if cardType == ddz.CardType_CT_NONE { //检查所出的牌能否组成牌型
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: genResult(3, "无法组成牌型"),
		})
		return int(ddz.StateID_state_playing), errors.New("无法组成牌型")
	}

	if context.CurCardType != ddz.CardType_CT_NONE &&
		(!canBiggerThan(cardType, context.CurCardType) || //牌型与上家不符(炸弹不算不符)
			(context.CurCardType == ddz.CardType_CT_SHUNZI && cardType == ddz.CardType_CT_SHUNZI && len(outCards) != len(context.CurOutCards))) { //顺子牌数不足
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: genResult(4, "牌型与上家不符"),
		})
		return int(ddz.StateID_state_playing), errors.New("牌型与上家不符")
	}

	lastPivot := toDDZCard(context.CardTypePivot)
	currPivot := *pivot
	bigger := false
	if cardType == ddz.CardType_CT_KINGBOMB {
		bigger = true
	} else if context.CurCardType == ddz.CardType_CT_KINGBOMB {
		bigger = false
	} else if cardType == ddz.CardType_CT_BOMB && context.CurCardType == ddz.CardType_CT_BOMB {
		bigger = currPivot.pointBiggerThan(lastPivot)
	} else if cardType == ddz.CardType_CT_BOMB && context.CurCardType != ddz.CardType_CT_BOMB {
		bigger = true
	} else if cardType != ddz.CardType_CT_BOMB && context.CurCardType == ddz.CardType_CT_BOMB {
		bigger = false
	} else if cardType != ddz.CardType_CT_BOMB && context.CurCardType != ddz.CardType_CT_BOMB {
		bigger = currPivot.pointBiggerThan(lastPivot)
	}

	if !bigger {
		sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{
			Result: genResult(5, "牌比上家小"),
		})
		return int(ddz.StateID_state_playing), errors.New("牌比上家小")
	}

	//更新玩家手牌和已出的牌
	handCards = RemoveAll(handCards, outCards)
	player.HandCards = toInts(handCards)
	player.OutCards = message.GetCards()

	//lastOutCards := toDDZCards(player.OutCards)
	//lastOutCards = AppendAll(lastOutCards, outCards)
	//player.AllOutCards = toInts(lastOutCards) // for 记牌器

	//更新context
	context.CurrentPlayerId = nextPlayerId
	//产生超时事件
	context.CountDownPlayers = []uint64{context.CurrentPlayerId}
	context.StartTime, _ = time.Now().MarshalBinary()
	context.Duration = StageTime[room.DDZStage_DDZ_STAGE_PLAYING]

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

	sendToPlayer(m, playerId, msgid.MsgID_ROOM_DDZ_PLAY_CARD_RSP, &room.DDZPlayCardRsp{ //成功出牌
		Result: genResult(0, ""),
	})

	var nextStage room.DDZStage
	if len(player.HandCards) == 0 {
		nextStage = room.DDZStage_DDZ_STAGE_OVER
	} else {
		nextStage = room.DDZStage_DDZ_STAGE_PLAYING
	}
	clientCardType := room.CardType(int32(cardType))
	broadcast(m, msgid.MsgID_ROOM_DDZ_PLAY_CARD_NTF, &room.DDZPlayCardNtf{ //广播出牌
		PlayerId:     &playerId,
		Cards:        message.GetCards(),
		CardType:     &clientCardType,
		TotalBomb:    &context.TotalBomb,
		NextPlayerId: &nextPlayerId,
		NextStage:    GenNextStage(nextStage),
	})

	if len(player.HandCards) == 0 {
		context.WinnerId = playerId
		return int(ddz.StateID_state_over), nil
	} else {
		return int(ddz.StateID_state_playing), nil
	}
}
