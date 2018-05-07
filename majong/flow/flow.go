package flow

import (
	"errors"
	"steve/majong/interfaces"
	"steve/majong/states"
	"steve/majong/transition"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"

	"github.com/Sirupsen/logrus"
)

type flow struct {
	context             majongpb.MajongContext
	autoEvent           *majongpb.AutoEvent
	stateFactory        interfaces.MajongStateFactory
	transitionValidator interfaces.TransitionValidator
	msgs                []majongpb.ReplyClientMessage
}

// NewFlow 创建 Flow
func NewFlow(mjContext majongpb.MajongContext) interfaces.MajongFlow {
	transitionFactory := transition.NewFactory()

	return &flow{
		context:             mjContext,
		stateFactory:        states.NewFactory(),
		transitionValidator: transitionFactory.CreateTransitionValidator(int(mjContext.GetGameId())),
		msgs:                make([]majongpb.ReplyClientMessage, 0),
	}
}

func (f *flow) GetMajongContext() *majongpb.MajongContext {
	return &f.context
}

func (f *flow) SetAutoEvent(autoEvent majongpb.AutoEvent) {
	f.autoEvent = &autoEvent
}

var errCreateState = errors.New("创建当前状态对象失败")
var errStateProcess = errors.New("当前状态处理事件失败")
var errCreateNewState = errors.New("创建新状态对象失败")
var errTransitionNotExist = errors.New("不存在转换关系")

func (f *flow) ProcessEvent(eventID majongpb.EventID, eventContext []byte) error {
	entry := logrus.WithFields(logrus.Fields{
		"event_id":         eventID,
		"current_state_id": f.context.CurState,
		"game_id":          f.context.GameId,
	})

	oldState := f.stateFactory.CreateState(int(f.context.GameId), f.context.CurState)
	if oldState == nil {
		entry.Error(errCreateState)
		return errCreateState
	}
	newStateID, err := oldState.ProcessEvent(eventID, eventContext, f)
	if err != nil {
		entry.WithError(err).Error(errStateProcess)
		return errStateProcess
	}
	entry = entry.WithField("new_state_id", newStateID)

	if newStateID == f.context.CurState {
		return nil
	}
	if err := f.transitionValidator.Valid(f.context.CurState, newStateID, eventID, f); err != nil {
		entry.WithError(err).Error(errTransitionNotExist)
		return errTransitionNotExist
	}

	newState := f.stateFactory.CreateState(int(f.context.GameId), newStateID)
	if newState == nil {
		entry.Error(errCreateNewState)
		return errCreateNewState
	}

	oldState.OnExit(f)
	f.context.CurState = newStateID
	newState.OnEntry(f)
	return nil
}

func (f *flow) GetSettler(settlerType interfaces.SettlerType) interfaces.Settler {
	logrus.Warn("TODO")
	return nil
}

func (f *flow) PushMessages(playerIDs []uint64, msgs ...interfaces.ToClientMessage) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "flow.PushMessages",
		"players":   playerIDs,
	})

	for _, msg := range msgs {

		bodyData, err := proto.Marshal(msg.Msg)
		if err != nil {
			logEntry.WithField("msg_id", msg.MsgID).WithError(err).Errorln("消息序列化失败")
			continue
		}
		f.msgs = append(f.msgs, majongpb.ReplyClientMessage{
			Players: playerIDs,
			MsgId:   int32(msg.MsgID),
			Msg:     bodyData,
		})
	}
}

func (f *flow) GetMessages() []majongpb.ReplyClientMessage {
	return f.msgs
}

func (f *flow) GetAutoEvent() *majongpb.AutoEvent {
	return f.autoEvent
}
