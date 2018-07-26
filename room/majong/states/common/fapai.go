package common

import (
	"errors"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/gutils"
	"steve/room/majong/interfaces"
	"steve/room/majong/utils"

	majongpb "steve/entity/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// FapaiState 发牌状态
type FapaiState struct{}

const (
	initHandCardCount int = 13
)

// ProcessEvent 处理事件
func (f *FapaiState) ProcessEvent(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_fapai_finish:
		{
			return f.nextState(flow.GetMajongContext()), nil
		}
	case majongpb.EventID_event_cartoon_finish_request:
		{
			return f.onCartoonFinish(flow, eventContext)
		}
	}
	return f.curState(), nil
}

// OnEntry 进入状态给每位玩家发手牌
func (f *FapaiState) OnEntry(flow interfaces.MajongFlow) {
	f.fapai(flow)
	f.notifyPlayer(flow)

	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_fapai_finish,
		EventContext: nil,
		WaitTime:     flow.GetMajongContext().GetOption().GetMaxFapaiCartoonTime(),
	})
}

// OnExit 退出状态
func (f *FapaiState) OnExit(flow interfaces.MajongFlow) {
	flow.GetMajongContext().TempData = new(majongpb.TempDatas) //清除临时数据
}

// nextState 下个状态
func (f *FapaiState) nextState(mjcontext *majongpb.MajongContext) majongpb.StateID {
	nextState := f.getNextState(mjcontext)
	logrus.WithFields(logrus.Fields{
		"func_name": "FapaiState.nextState",
		"nextState": nextState,
	}).Infoln("发牌下一状态")
	return nextState
}

// curState 当前状态
func (f *FapaiState) curState() majongpb.StateID {
	return majongpb.StateID_state_fapai
}

// onCartoonFinish 动画播放完毕
func (f *FapaiState) onCartoonFinish(flow interfaces.MajongFlow, eventContext interface{}) (newState majongpb.StateID, err error) {
	cartoonFinishData := CartoonFinishData{
		CurState:        f.curState(),
		NextState:       f.getNextState(flow.GetMajongContext()),
		NeedCartoonType: room.CartoonType_CTNT_FAPAI,
		EventContext:    eventContext,
	}
	return OnCartoonFinish(cartoonFinishData, flow.GetMajongContext())
}

var errCardsNotEnough = errors.New("墙牌不足")

func (f *FapaiState) fapaiToPlayer(flow interfaces.MajongFlow, p *majongpb.Player, count int) error {
	majongContext := flow.GetMajongContext()
	wallCards := majongContext.GetWallCards()
	if count > len(wallCards) {
		logrus.WithError(errCardsNotEnough).WithFields(logrus.Fields{
			"player_id":       p.GetPalyerId(),
			"count":           count,
			"wallcards_count": len(wallCards),
		})
		return errCardsNotEnough
	}
	p.HandCards = append(p.HandCards, wallCards[:count]...)
	majongContext.WallCards = wallCards[count:]
	return nil
}

func (f *FapaiState) fapai(flow interfaces.MajongFlow) {
	majongContext := flow.GetMajongContext()
	playerCount := len(majongContext.Players)

	zjIndex := int(majongContext.GetZhuangjiaIndex())
	if zjIndex >= playerCount {
		logrus.WithField("index", zjIndex).Panic("庄家索引越界")
	}
	xpOption := mjoption.GetXingpaiOption(int(majongContext.GetXingpaiOptionId()))
	zjPlayer := majongContext.Players[zjIndex]
	switch xpOption.FapaiType {
	case mjoption.NomarlFapai:
		{
			f.fapaiToPlayer(flow, zjPlayer, 1)
			f.fapaiToPlayers(flow, zjIndex, playerCount)
		}
	case mjoption.ErrenFapai:
		{
			f.fapaiToPlayers(flow, zjIndex, playerCount)
		}
	}
	majongContext.LastMopaiPlayer = zjPlayer.GetPalyerId()
	zjHandCards := zjPlayer.GetHandCards()
	majongContext.LastMopaiCard = zjHandCards[len(zjHandCards)-1]
}

func (f *FapaiState) fapaiToPlayers(flow interfaces.MajongFlow, zjIndex int, playerCount int) {
	mjContext := flow.GetMajongContext()
	for i := 0; i < playerCount; i++ {
		player := mjContext.Players[(i+zjIndex)%playerCount]
		f.fapaiToPlayer(flow, player, initHandCardCount)
	}
}

// notifyPlayer 通知玩家发牌消息
func (f *FapaiState) notifyPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	playerCardCount := []*room.PlayerCardCount{}

	for _, player := range mjContext.Players {
		playerCardCount = append(playerCardCount, &room.PlayerCardCount{
			PlayerId:  proto.Uint64(player.GetPalyerId()),
			CardCount: proto.Uint32(uint32(len(player.GetHandCards()))),
		})
	}

	for _, player := range mjContext.Players {
		msg := &room.RoomFapaiNtf{
			Cards:            utils.ServerCards2Uint32(player.GetHandCards()),
			PlayerCardCounts: playerCardCount,
		}
		flow.PushMessages([]uint64{player.GetPalyerId()}, interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_ROOM_FAPAI_NTF),
			Msg:   msg,
		})
	}
}

// 下一状态获取
func (f *FapaiState) getNextState(mjContext *majongpb.MajongContext) majongpb.StateID {
	//先要判断游戏有没有换三张的玩法，有换三张的玩法，再判断需不需要配置换三张
	xpOption := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId()))
	if gutils.GameHasHszState(mjContext) {
		return majongpb.StateID_state_huansanzhang
	}
	if xpOption.EnableDingque {
		return majongpb.StateID_state_dingque
	}
	if xpOption.EnableKaijuAddflower {
		return majongpb.StateID_state_gamestart_buhua
	}
	return majongpb.StateID_state_zixun
}
