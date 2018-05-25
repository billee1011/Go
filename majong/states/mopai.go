package states

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/utils"
	"steve/peipai"
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

func (s *MoPaiState) notifyMopai(flow interfaces.MajongFlow, playerID uint64, back bool, card *majongpb.Card) {
	context := flow.GetMajongContext()
	for _, player := range context.Players {
		ntf := &room.RoomMopaiNtf{}
		if player.PalyerId == context.GetMopaiPlayer() {
			ntf.Card = proto.Uint32(utils.ServerCard2Uint32(card))
		}
		ntf.Player = &context.MopaiPlayer
		ntf.Back = proto.Bool(back)
		toClientMessage := interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_ROOM_MOPAI_NTF),
			Msg:   ntf,
		}
		flow.PushMessages([]uint64{player.GetPalyerId()}, toClientMessage)
	}
}

//mopai 摸牌处理
func (s *MoPaiState) mopai(flow interfaces.MajongFlow) (majongpb.StateID, error) {
	context := flow.GetMajongContext()
	logEntry := logrus.WithField("func_name", "MoPaiState.mopai")
	logEntry = utils.WithMajongContext(logEntry, context)

	players := context.GetPlayers()
	activePlayer := utils.GetPlayerByID(players, context.GetMopaiPlayer())
	context.ActivePlayer = activePlayer.GetPalyerId()
	if s.checkGameOver(flow) {
		// if len(context.WallCards) == 0 {
		logEntry.Infoln("没牌了")
		return majongpb.StateID_state_gameover, nil
	}
	//从墙牌中移除一张牌
	card := context.WallCards[0]
	logEntry.WithFields(logrus.Fields{
		"wall_card_count": len(context.GetWallCards()),
		"card":            card,
	}).Infoln("执行摸牌")

	context.WallCards = context.WallCards[1:]
	//将这张牌添加到手牌中
	activePlayer.HandCards = append(activePlayer.GetHandCards(), card)
	context.LastMopaiPlayer = context.MopaiPlayer
	context.LastMopaiCard = card
	context.ZixunType = majongpb.ZixunType_ZXT_NORMAL
	activePlayer.MopaiCount++

	s.notifyMopai(flow, context.GetMopaiPlayer(), false, card)
	return majongpb.StateID_state_zixun, nil
}

func (s *MoPaiState) checkGameOver(flow interfaces.MajongFlow) bool {
	context := flow.GetMajongContext()
	if len(context.WallCards) == 0 {
		return true
	}
	//TODO 由配牌控制是否gameover,配牌长度为0走正常gameover,配牌长度不为0走配牌长度流局
	length := peipai.GetLensOfWallCards(utils.GetGameName(flow))
	if utils.GetAllMopaiCount(context) == length-53 {
		return true
	}
	return false
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
