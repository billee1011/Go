package process

import (
	"steve/majong/flow"
	server_pb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// HandleMajongEventResult 处理牌局事件的结果
type HandleMajongEventResult struct {
	NewContext     server_pb.MajongContext        // 处理后的牌局现场
	TimeCheckInfos []server_pb.TimeCheckInfo      // 需要取消或者设置的时间检测信息
	ReplyMsgs      []server_pb.ReplyClientMessage // 回复给客户端的消息
	Succeed        bool                           // 是否成功
}

// HandleMajongEventParams 处理牌局事件的参数
type HandleMajongEventParams struct {
	MajongContext server_pb.MajongContext // 牌局现场
	EventID       server_pb.EventID       // 事件 ID
	EventContext  []byte                  // 事件现场
}

// HandleMajongEvent 处理牌局事件
func HandleMajongEvent(params HandleMajongEventParams) (result HandleMajongEventResult) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleMajongEvent",
		"params":    params,
	})
	result = HandleMajongEventResult{
		NewContext:     params.MajongContext,
		TimeCheckInfos: make([]server_pb.TimeCheckInfo, 0),
		ReplyMsgs:      make([]server_pb.ReplyClientMessage, 0),
		Succeed:        false,
	}

	flow := flow.NewFlow(params.MajongContext)
	err := flow.ProcessEvent(params.EventID, params.EventContext)
	if err != nil {
		logEntry.WithError(err).Errorln("处理事件失败")
		return
	}
	result.NewContext = *flow.GetMajongContext()
	result.TimeCheckInfos = flow.GetTimeCheckInfos()
	result.ReplyMsgs = flow.GetMessages()
	result.Succeed = true
	return
}
