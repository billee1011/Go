package procedure

import (
	"steve/client_pb/room"
	"steve/entity/poker/ddz"
	"steve/room/flows/ddzflow/ddz/ddzmachine"
	"steve/room/flows/ddzflow/ddz/states"
	"steve/room/flows/ddzflow/machine"
	playerpkg "steve/room/player"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// Result 处理牌局事件的结果
type Result struct {
	Context           ddz.DDZContext // 最新现场
	HasAutoEvent      bool
	AutoEventID       int
	AutoEventContext  interface{}
	AutoEventDuration time.Duration
	Succeed           bool // 是否成功
}

// Params 处理牌局事件的参数
type Params struct {
	// PlayerMgr    interfaces.DeskPlayerMgr // 是否托管
	Context      ddz.DDZContext           // 牌局现场
	Sender       ddzmachine.MessageSender // 消息发送器， 拆分后要修改
	EventID      int                      // 事件 ID
	EventContext interface{}              // 事件现场
}

// HandleEvent 处理牌局事件
func HandleEvent(params Params) (result Result) {
	start := time.Now()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleEvent",
		"params":    params,
	})

	result = Result{
		Context:      params.Context,
		Succeed:      false,
		HasAutoEvent: false,
	}
	m := ddzmachine.CreateDDZMachine(&params.Context, states.NewFactory(), params.Sender)

	// 处理恢复对局的请求
	// if params.EventID == int(ddz.EventID_event_resume_request) {
	// 	resumeErr := dealResumeRequest(&params, m, &cloneContext)
	// 	if resumeErr != nil {
	// 		logEntry.WithError(resumeErr).Errorln("处理恢复对局失败")
	// 	}
	// 	return
	// }

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

	end := time.Now()
	logrus.WithField("duration", end.Sub(start)).Debug("状态机从创建到退出")
	return
}

// TranslateDDZPlayerToRoomPlayer 将 ddzPlayer 转换成 RoomPlayerInfo
func TranslateDDZPlayerToRoomPlayer(ddzPlayer ddz.Player, seat uint32) room.RoomPlayerInfo {
	playerMgr := playerpkg.GetPlayerMgr()
	playerID := ddzPlayer.GetPlayerId()
	roomPlayer := playerMgr.GetPlayer(playerID)
	var coin uint64
	if roomPlayer != nil {
		coin = roomPlayer.GetCoin()
	}

	return room.RoomPlayerInfo{
		PlayerId: proto.Uint64(playerID),
		Name:     proto.String("player" + string(playerID)),
		Coin:     proto.Uint64(coin),
		Seat:     proto.Uint32(seat),
		// Location: TODO 没地方拿
	}
}
