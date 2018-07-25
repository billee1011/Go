package transition

import (
	"errors"
	majongpb "steve/entity/majong"

	"github.com/Sirupsen/logrus"
)

func stringToStateID(s string) (majongpb.StateID, bool) {
	id, ok := majongpb.StateID_value[s]
	return majongpb.StateID(id), ok
}
func stringToEventID(s string) (majongpb.EventID, bool) {
	id, ok := majongpb.EventID_value[s]
	return majongpb.EventID(id), ok
}

var errInvalidStateName = errors.New("状态名无效")
var errInvalidEventName = errors.New("事件名无效")

func originTransitionToMap(t *transition) (transitionMap, error) {
	result := make(transitionMap)
	for _, state := range t.States {
		curState, ok := stringToStateID(state.CurState)
		if !ok {
			logrus.WithField("state", state.CurState).Error("状态不存在")
			return nil, errInvalidStateName
		}
		if _, ok = result[curState]; !ok {
			result[curState] = make(map[majongpb.EventID][]majongpb.StateID)
		}

		for _, tran := range state.Trans {
			if err := convertStateTrans(&tran, result[curState]); err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

func convertStateTrans(t *stateTran, destination map[majongpb.EventID][]majongpb.StateID) error {
	stateID, ok := stringToStateID(t.NextState)
	if !ok {
		logrus.WithField("state", t.NextState).Error(errInvalidStateName)
		return errInvalidStateName
	}

	for _, event := range t.Events {
		eventID, ok := stringToEventID(event)
		if !ok {
			logrus.WithField("event", event).Error(errInvalidEventName)
			return errInvalidEventName
		}
		if _, ok := destination[eventID]; !ok {
			destination[eventID] = []majongpb.StateID{}
		}
		destination[eventID] = append(destination[eventID], stateID)
	}
	return nil
}
