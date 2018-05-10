package states

import (
	"errors"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// WaitQiangganghuState 等待抢杠胡状态
type WaitQiangganghuState struct {
}

var _ interfaces.MajongState = new(WaitQiangganghuState)

// ProcessEvent 处理事件
func (s *WaitQiangganghuState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_hu_request:
		{
			return s.onHuRequest(eventContext, flow)
		}
	case majongpb.EventID_event_qi_request:
		{
			return s.onQiRequest(eventContext, flow)
		}
	}
	return majongpb.StateID_state_waitqiangganghu, errInvalidEvent
}

// OnEntry 进入状态
func (s *WaitQiangganghuState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *WaitQiangganghuState) OnExit(flow interfaces.MajongFlow) {

}

// onHuRequest 处理胡请求
func (s *WaitQiangganghuState) onHuRequest(eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "WaitQiangganghuState.onHuRequest",
	})
	newState, err = majongpb.StateID_state_waitqiangganghu, nil

	huRequest := majongpb.HuRequestEvent{}
	if err = proto.Unmarshal(eventContext, &huRequest); err != nil {
		err = errors.New("事件现场反序列化失败")
		logEntry.WithError(err).Errorln(err)
		return
	}
	playerID := huRequest.GetHead().GetPlayerId()
	logEntry = logEntry.WithField("request_player", playerID)

	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	player := utils.GetMajongPlayer(playerID, mjContext)
	if !utils.ExistPossibleAction(player, majongpb.Action_action_dianpao) {
		logEntry.Infoln("该玩家不能抢杠胡")
		return
	}
	if player.GetHasSelected() {
		logEntry.Infoln("该玩家已经做出过选择了")
		return
	}
	player.HasSelected, player.SelectedAction = true, majongpb.Action_action_dianpao
	return s.makeDecision(flow)
}

// onQiRequest 处理弃请求
func (s *WaitQiangganghuState) onQiRequest(eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "WaitQiangganghuState.onQiRequest",
	})
	newState, err = majongpb.StateID_state_waitqiangganghu, nil

	qiRequest := majongpb.QiRequestEvent{}
	if err = proto.Unmarshal(eventContext, &qiRequest); err != nil {
		err = errors.New("事件现场反序列化失败")
		logEntry.WithError(err).Errorln(err)
		return
	}
	playerID := qiRequest.GetHead().GetPlayerId()
	logEntry = logEntry.WithField("request_player", playerID)

	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	player := utils.GetMajongPlayer(playerID, mjContext)
	if !utils.ExistPossibleAction(player, majongpb.Action_action_dianpao) {
		logEntry.Infoln("该玩家不能抢杠胡")
		return
	}
	if player.GetHasSelected() {
		logEntry.Infoln("该玩家已经做出过选择了")
		return
	}
	player.HasSelected, player.SelectedAction = true, majongpb.Action_action_qi
	return s.makeDecision(flow)
}

// makeDecision 作决策
// step 1. 查找是否有玩家还没有做出选择， 如果有，保留原状态并结束
// step 2. 如果有玩家选择了胡操作，返回到抢杠胡状态。 否则返回补杠状态
func (s *WaitQiangganghuState) makeDecision(flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	huPlayers := []uint64{}

	mjContext := flow.GetMajongContext()
	for _, player := range mjContext.GetPlayers() {
		if len(player.GetPossibleActions()) < 0 {
			continue
		}
		if !player.GetHasSelected() {
			return majongpb.StateID_state_waitqiangganghu, nil
		}
		if player.SelectedAction == majongpb.Action_action_dianpao {
			huPlayers = append(huPlayers, player.GetPalyerId())
		}
	}
	if len(huPlayers) == 0 {
		return majongpb.StateID_state_bugang, nil
	}
	mjContext.LastHuPlayers = huPlayers
	return majongpb.StateID_state_qiangganghu, nil
}
