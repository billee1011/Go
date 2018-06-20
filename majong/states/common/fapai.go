package common

import (
	"errors"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/utils"

	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// FapaiState 发牌状态
type FapaiState struct{}

const (
	initHandCardCount int = 13
)

// ProcessEvent 处理事件
func (f *FapaiState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_fapai_finish:
		{
			return f.nextState(), nil
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
}

// nextState 下个状态
func (f *FapaiState) nextState() majongpb.StateID {
	return majongpb.StateID_state_huansanzhang
}

// curState 当前状态
func (f *FapaiState) curState() majongpb.StateID {
	return majongpb.StateID_state_fapai
}

// onCartoonFinish 动画播放完毕
func (f *FapaiState) onCartoonFinish(flow interfaces.MajongFlow, eventContext []byte) (newState majongpb.StateID, err error) {
	return OnCartoonFinish(f.curState(), f.nextState(), room.CartoonType_CTNT_FAPAI, eventContext)
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
	zjPlayer := majongContext.Players[zjIndex]
	f.fapaiToPlayer(flow, zjPlayer, 1)

	for i := 0; i < playerCount; i++ {
		player := majongContext.Players[(i+zjIndex)%playerCount]
		f.fapaiToPlayer(flow, player, initHandCardCount)
	}
	majongContext.LastMopaiPlayer = zjPlayer.GetPalyerId()
	zjHandCards := zjPlayer.GetHandCards()
	majongContext.LastMopaiCard = zjHandCards[len(zjHandCards)-1]
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