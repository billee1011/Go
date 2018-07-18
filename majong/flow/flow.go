package flow

import (
	"errors"
	"steve/majong/interfaces"
	"steve/majong/transition"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"

	_ "steve/majong/fantype"        // init fantype
	_ "steve/majong/states/factory" // init state facotry

	"github.com/Sirupsen/logrus"
	"steve/majong/bus"
)

type flow struct {
	context             majongpb.MajongContext
	autoEvent           *majongpb.AutoEvent
	stateFactory        interfaces.MajongStateFactory
	transitionValidator interfaces.TransitionValidator
	msgs                []majongpb.ReplyClientMessage
	timeCheckInfos      []majongpb.TimeCheckInfo
}

// NewFlow 创建 Flow
func NewFlow(mjContext majongpb.MajongContext) interfaces.MajongFlow {
	transitionFactory := transition.NewFactory()

	return &flow{
		context:             mjContext,
		stateFactory:        bus.GetMajongStateFactory(),
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

var errCreateState = errors.New("创建状态对象失败")
var errStateProcess = errors.New("当前状态处理事件失败")
var errTransitionNotExist = errors.New("不存在转换关系")

// stateProcess 派发到状态处理事件
func (f *flow) stateProcess(entry *logrus.Entry, eventID majongpb.EventID, eventContext []byte) (majongpb.StateID, error) {
	curStateID := f.context.CurState
	oldState := f.stateFactory.CreateState(f.context.GameId, curStateID)
	if oldState == nil {
		entry.Error(errCreateState)
		return curStateID, errCreateState
	}
	newStateID, err := oldState.ProcessEvent(eventID, eventContext, f)
	if err != nil {
		entry.WithError(err).Error(errStateProcess)
		return curStateID, errStateProcess
	}
	return newStateID, nil
}

// switchState 状态切换
func (f *flow) switchState(entry *logrus.Entry, eventID majongpb.EventID, newStateID majongpb.StateID) error {
	oldStateID := f.context.CurState
	entry = entry.WithFields(logrus.Fields{
		"old_state": f.context.CurState,
		"new_state": newStateID,
	})
	entry.Debugln("状态切换")
	if newStateID == oldStateID {
		return nil
	}
	if err := f.transitionValidator.Valid(oldStateID, newStateID, eventID, f); err != nil {
		entry.WithError(err).Error(errTransitionNotExist)
		return errTransitionNotExist
	}
	oldState := f.stateFactory.CreateState(f.context.GameId, oldStateID)
	newState := f.stateFactory.CreateState(f.context.GameId, newStateID)
	if oldState == nil || newState == nil {
		entry.Error(errCreateState)
		return errCreateState
	}
	oldState.OnExit(f)
	f.context.CurState = newStateID
	newState.OnEntry(f)
	return nil
}

// processAutoEvent 处理自动事件
func (f *flow) processAutoEvent(entry *logrus.Entry) error {
	if f.autoEvent == nil {
		return nil
	}
	ae := f.autoEvent
	f.autoEvent = nil

	entry.WithFields(logrus.Fields{
		"event_id": ae.EventId,
	}).Debugln("处理自动事件")

	return f.ProcessEvent(ae.EventId, ae.EventContext)
}

// ProcessEvent 处理外部事件
func (f *flow) ProcessEvent(eventID majongpb.EventID, eventContext []byte) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":        "flow.ProcessEvent",
		"event_id":         eventID,
		"current_state_id": f.context.CurState,
		"game_id":          f.context.GameId,
	})

	var err error
	var newStateID majongpb.StateID
	if newStateID, err = f.stateProcess(entry, eventID, eventContext); err != nil {
		return err
	}
	// 自动事件交给 room 处理，确保消息时序正确
	return f.switchState(entry, eventID, newStateID)
	// if err = f.switchState(entry, eventID, newStateID); err != nil {
	// 	return err
	// }
	// return f.processAutoEvent(entry)
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
