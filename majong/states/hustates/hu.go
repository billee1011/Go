package hustates

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HuState 胡状态
// 进入胡状态时， 执行胡操作。设置胡完成事件
// 收到胡完成事件时，设置摸牌玩家，返回摸牌状态
type HuState struct {
}

var _ interfaces.MajongState = new(HuState)

// ProcessEvent 处理事件
func (s *HuState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_hu_finish {
		s.setMoPaiPlayer(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_hu, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *HuState) OnEntry(flow interfaces.MajongFlow) {
	s.doHu(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_hu_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *HuState) OnExit(flow interfaces.MajongFlow) {

}

// addHuCard 添加胡的牌
func (s *HuState) addHuCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	addHuCard(card, player, srcPlayerID, majongpb.HuType_hu_dianpao)
}

// doHu 执行胡操作
func (s *HuState) doHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HuState.doHu",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	players := mjContext.GetLastHuPlayers()

	for _, playerID := range players {
		player := utils.GetMajongPlayer(playerID, mjContext)
		card := mjContext.GetLastOutCard()
		s.addHuCard(card, player, playerID)
	}
	s.notifyHu(flow)
	return
}

// HuState 广播胡
func (s *HuState) notifyHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	huCard, _ := utils.CardToRoomCard(mjContext.GetLastOutCard())
	body := room.RoomHuNtf{
		Players:      mjContext.GetLastHuPlayers(),
		FromPlayerId: proto.Uint64(mjContext.GetLastChupaiPlayer()),
		Card:         huCard,
		HuType:       room.HuType_DianPao.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HU_NTF, &body)
}

// stepCount 计算从 srcPlayer 到 destPlayer 需要经过的距离
// TODO 未考虑性能， 但影响不大
func (s *HuState) stepCount(mjContext *majongpb.MajongContext, srcPlayer uint64, destPlayer uint64) (int, error) {
	players := mjContext.GetPlayers()
	srcIndex, err := utils.GetPlayerIndex(srcPlayer, players)
	if err != nil {
		return 0, err
	}
	destIndex, err := utils.GetPlayerIndex(destPlayer, players)
	if err != nil {
		return 0, err
	}
	return (srcIndex + len(players) - destIndex) % len(players), nil
}

// setMoPaiPlayer 设置摸牌玩家
// 从出牌玩家的位置算起，找到最后一个胡牌的玩家，他的下家就是摸牌的玩家
// TODO 未考虑性能， 但影响不大
func (s *HuState) setMoPaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HuState.setMoPaiPlayer",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	chupaiPlayer := mjContext.GetLastChupaiPlayer()
	hupaiPlayers := mjContext.GetLastHuPlayers()

	maxStepCount := -1
	var maxStepPlayer uint64
	for _, hupaiPlayer := range hupaiPlayers {
		c, err := s.stepCount(mjContext, chupaiPlayer, hupaiPlayer)
		if err != nil {
			logEntry.Errorln(err)
			return
		}
		if c > maxStepCount {
			maxStepCount = c
			maxStepPlayer = hupaiPlayer
		}
	}
	if maxStepPlayer == 0 {
		logEntry.Errorln("没有找到最后一个胡牌玩家")
		return
	}
	players := mjContext.GetPlayers()
	srcIndex, _ := utils.GetPlayerIndex(maxStepPlayer, players)
	moPaiIndex := (srcIndex + 1) % len(players)
	mjContext.MopaiPlayer = players[moPaiIndex].GetPalyerId()
}
