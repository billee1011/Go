package process

import (
	"bytes"
	"encoding/gob"
	server_pb "steve/entity/majong"
	"steve/room/majong/flow"

	"github.com/Sirupsen/logrus"
)

// HandleMajongEventResult 处理牌局事件的结果
type HandleMajongEventResult struct {
	NewContext server_pb.MajongContext        // 处理后的牌局现场
	AutoEvent  *server_pb.AutoEvent           // 自动事件
	ReplyMsgs  []server_pb.ReplyClientMessage // 回复给客户端的消息
	Succeed    bool                           // 是否成功
}

// HandleMajongEventParams 处理牌局事件的参数
type HandleMajongEventParams struct {
	MajongContext server_pb.MajongContext // 牌局现场
	EventID       server_pb.EventID       // 事件 ID
	EventContext  interface{}             // 事件现场
}

// HandleMajongEvent 处理牌局事件
func HandleMajongEvent(params HandleMajongEventParams) (result HandleMajongEventResult) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleMajongEvent",
		"params":    params,
	})

	cloneContext := &params.MajongContext

	result = HandleMajongEventResult{
		NewContext: *cloneContext,
		ReplyMsgs:  make([]server_pb.ReplyClientMessage, 0),
		Succeed:    false,
	}
	flow := flow.NewFlow(*cloneContext)
	err := flow.ProcessEvent(params.EventID, params.EventContext)
	if err != nil {
		logEntry.WithError(err).Errorln("处理事件失败")
		return
	}
	result.NewContext = *flow.GetMajongContext()
	result.ReplyMsgs = flow.GetMessages()
	result.AutoEvent = flow.GetAutoEvent()
	result.Succeed = true
	return
}

func deepCopyMjongContext(src server_pb.MajongContext) (*server_pb.MajongContext, error) {
	var buf bytes.Buffer
	var err error
	if err = gob.NewEncoder(&buf).Encode(src); err != nil {
		return nil, err
	}
	var dst server_pb.MajongContext
	err = gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&dst)
	return &dst, err
}
