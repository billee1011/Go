package process

import (
	"steve/majong/flow"
	server_pb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// HandleMajongEvent 处理麻将事件
func HandleMajongEvent(mjContext server_pb.MajongContext,
	eventID server_pb.EventID, eventContext []byte) (newContext server_pb.MajongContext,
	autoEvent *server_pb.AutoEvent, replyMsgs []server_pb.ReplyClientMessage, succeed bool) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleMajongEvent",
		"event_id":  eventID,
		"cur_state": mjContext.GetCurState(),
	})
	succeed = false
	replyMsgs = []server_pb.ReplyClientMessage{}
	newContext = mjContext
	autoEvent = nil

	flow := flow.NewFlow(mjContext)
	err := flow.ProcessEvent(eventID, eventContext)
	if err != nil {
		logEntry.WithError(err).Errorln("处理事件失败")
		return
	}
	newContext = *flow.GetMajongContext()
	replyMsgs = flow.GetMessages()
	autoEvent = flow.GetAutoEvent()
	succeed = true

	return
}
