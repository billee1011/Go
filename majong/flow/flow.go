package flow

import (
	"errors"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type flow struct {
	context             majongpb.MajongContext
	autoEvent           *majongpb.AutoEvent
	stateFactory        interfaces.MajongStateFactory
	transitionValidator interfaces.TransitionValidator
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
	if err := f.transitionValidator.Valid(f.context.CurState, newStateID, eventID); err != nil {
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
	logrus.Warn("TODO")
}

func (f *flow) GetMessages() []proto.Message {
	logrus.Warn("TODO")
	return nil
}
