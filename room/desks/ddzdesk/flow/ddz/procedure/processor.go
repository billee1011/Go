package procedure

import (
	"steve/room/desks/ddzdesk/flow/ddz/ddzmachine"
	"steve/room/desks/ddzdesk/flow/ddz/states"
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// Result 处理牌局事件的结果
type Result struct {
	Context           ddz.DDZContext // 最新现场
	HasAutoEvent      bool
	AutoEventID       int
	AutoEventContext  []byte
	AutoEventDuration time.Duration
	Succeed           bool // 是否成功
}

// Params 处理牌局事件的参数
type Params struct {
	Context      ddz.DDZContext           // 牌局现场
	Sender       ddzmachine.MessageSender // 消息发送器， 拆分后要修改
	EventID      int                      // 事件 ID
	EventContext []byte                   // 事件现场
}

// HandleEvent 处理牌局事件
func HandleEvent(params Params) (result Result) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleEvent",
		"params":    params,
	})

	cloneContext := *proto.Clone(&params.Context).(*ddz.DDZContext)

	result = Result{
		Context:      cloneContext,
		Succeed:      false,
		HasAutoEvent: false,
	}
	m := ddzmachine.CreateDDZMachine(&cloneContext, states.NewFactory(), params.Sender)
	err := m.ProcessEvent(machine.Event{
		EventID:   params.EventID,
		EventData: params.EventContext,
	})
	if err != nil {
		logEntry.WithError(err).Errorln("处理事件失败")
		return
	}
	result.Context = *m.GetDDZContext()
	e, d := m.GetAutoEvent()
	if e != nil {
		result.HasAutoEvent = true
		result.AutoEventID = e.EventID
		result.AutoEventContext = e.EventData
		result.AutoEventDuration = d
	} else {
		result.HasAutoEvent = false
	}
	result.Succeed = true
	return
}
