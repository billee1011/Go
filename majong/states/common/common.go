package common

import (
	"steve/client_pb/room"
	"steve/majong/global"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// OnCartoonFinish 在某个状态上， 动画播放完成
func OnCartoonFinish(curState majongpb.StateID, nextState majongpb.StateID, needCartoonType room.CartoonType, eventContext []byte) (newState majongpb.StateID, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":         "OnCartoonFinish",
		"cur_state":         curState,
		"next_state":        nextState,
		"need_cartoon_type": needCartoonType,
	})

	req := new(majongpb.CartoonFinishRequestEvent)
	if marshalErr := proto.Unmarshal(eventContext, req); marshalErr != nil {
		logEntry.WithError(marshalErr).Errorln(global.ErrUnmarshalEvent)
		return curState, global.ErrUnmarshalEvent
	}
	reqCartoonType := req.GetCartoonType()
	logEntry.WithField("req_cartoon_type", reqCartoonType).Debugln("收到动画完成请求")
	if reqCartoonType != int32(needCartoonType) {
		return curState, nil
	}
	return nextState, nil
}
