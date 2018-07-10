package ddzdesk

import (
	msgid "steve/client_pb/msgId"
	"steve/room/desks/ddzdesk/flow/ddz/ddzmachine"
	"steve/room/desks/ddzdesk/flow/ddz/procedure"
	"steve/room/desks/deskbase"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/server_pb/ddz"
	"steve/structs/proto/gate_rpc"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"runtime/debug"
	"context"
)

// deskEvent 牌桌事件
type deskEvent struct {
	eventID      int
	eventContext []byte
	eventType interfaces.EventType
	playerID uint64
}

// desk 斗地主牌桌
type desk struct {
	*deskbase.DeskBase
	eventChannel   chan deskEvent
	closingChannel chan struct{}
	ddzContext     *ddz.DDZContext
	cancel  context.CancelFunc     // 取消事件处理
}

// initDDZContext 初始化斗地主现场
func (d *desk) initDDZContext() {
	d.ddzContext = procedure.CreateInitDDZContext(facade.GetDeskPlayerIDs(d))
}

// Start 启动牌桌逻辑
// finish : 当牌桌逻辑完成时调用
func (d *desk) Start(finish func()) error {
	d.eventChannel = make(chan deskEvent, 4)
	d.closingChannel = make(chan struct{})

	d.initDDZContext()
	go func() {
		d.run()
		finish()
	}()

	var ctx context.Context
	ctx, d.cancel = context.WithCancel(context.Background())
	go func() {
		d.timerTask(ctx)
	}()
	d.pushEvent(&deskEvent{
		eventID: int(ddz.EventID_event_start_game),
	})
	return nil
}

// timerTask 定时任务，产生自动事件
func (d *desk) timerTask(ctx context.Context) {
	defer func() {
		if x := recover(); x != nil {
			debug.PrintStack()
		}
	}()

	t := time.NewTicker(time.Millisecond * 200)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			{
				d.genTimerEvent()
			}
		case <-ctx.Done():
			{
				return
			}
		}
	}
}

// genTimerEvent 生成计时事件
func (d *desk) genTimerEvent() {
	g := global.GetDeskAutoEventGenerator()
	// 先将 context 指针读出来拷贝， 后面的 context 修改都会分配一块新的内存
	dContext := d.ddzContext
	tuoGuanPlayers := facade.GetTuoguanPlayers(d)
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":       "desk.genTimerEvent",
		"tuoguan_players": tuoGuanPlayers,
	})
	startTime := time.Time{}
	startTime.UnmarshalBinary(dContext.StartTime)
	result := g.GenerateV2(&interfaces.AutoEventGenerateParams{
		Desk:       d,
		DDZContext: dContext,
		PlayerIds:  dContext.CountDownPlayers,
		StartTime:  startTime,
		Duration:   dContext.Duration,
		RobotLv:    map[uint64]int{},
	})
	for _, event := range result.Events {
		logEntry.WithFields(logrus.Fields{
			"event_id":     event.ID,
			"event_player": event.PlayerID,
			"event_type":   event.EventType,
		}).Debugln("注入计时事件")
		d.eventChannel <- deskEvent{
			eventID:       int(event.ID),
			eventContext: event.Context,
			eventType: event.EventType,
			playerID: event.PlayerID,
		}
	}
}

// Stop 停止牌桌
func (d *desk) Stop() error {
	d.cancel()
	d.closingChannel <- struct{}{}
	return nil
}

// PushRequest 压入玩家请求
func (d *desk) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.PushRequest",
		"player_id": playerID,
		"msg_id":    head.GetMsgId(),
	})

	translator := global.GetReqEventTranslator()
	eventID, eventData, err := translator.Translate(playerID, head, bodyData)
	if err != nil {
		entry.WithError(err).Errorln("事件转换失败")
		return
	}
	if eventID == 0 {
		entry.Warningln("没有对应事件")
		return
	}
	eventMessage, ok := eventData.(proto.Message)
	if !ok {
		entry.Errorln("事件数据不是 proto.Message 类型")
		return
	}
	eventContext, err := proto.Marshal(eventMessage)
	if err != nil {
		entry.WithError(err).Errorln("事件消息序列化失败")
		return
	}
	d.pushEvent(&deskEvent{
		eventID:      eventID,
		eventContext: eventContext,
	})
}

func (d *desk) pushEvent(e *deskEvent) {
	d.eventChannel <- *e
}

// PushEvent 压入事件
func (d *desk) PushEvent(event interfaces.Event) {
	return
}

// run 执行牌桌逻辑
func (d *desk) run() {

forstart:
	for {
		select {
		case event := <-d.eventChannel:
			{
				d.processEvent(&event)
				d.recordTuoguanOverTimeCount(event)
			}
		case <-d.closingChannel:
			{
				break forstart
			}
		}
	}
}

// recordTuoguanOverTimeCount 记录托管超时计数
func (d *desk) recordTuoguanOverTimeCount(event deskEvent) {
	if event.eventType != interfaces.OverTimeEvent {
		return
	}
	playerID := event.playerID
	if playerID == 0 {
		return
	}
	deskPlayer := facade.GetDeskPlayerByID(d, playerID)
	if deskPlayer != nil {
		deskPlayer.OnPlayerOverTime()
	}
}

func (d *desk) processEvent(e *deskEvent) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.processEvent",
		"event_id":  e.eventID,
	})

	params := procedure.Params{
		Context:      *d.ddzContext,
		Sender:       d.getMessageSender(),
		EventID:      e.eventID,
		EventContext: e.eventContext,
	}

	/* 	// 处理恢复对局的请求
	   	if e.eventID == int(ddz.EventID_event_resume_request) {
	   		message := &ddz.ResumeRequestEvent{}
	   		err := proto.Unmarshal(e.eventContext, message)
	   		if err != nil {
	   			//logEntry.WithError(err).Errorln("处理恢复对局事件失败")
	   			return
	   		}

	   		// 请求的玩家ID
	   		reqPlayerID := message.GetHead().GetPlayerId()

	   		bExist := false

	   		// 是否有这个玩家
	   		for _, player := range d.ddzContext.GetPlayers() {
	   			if player.GetPalyerId() == reqPlayerID {
	   				bExist = true
	   			}
	   		}

	   		// 存在的话，则发送回复消息
	   		if bExist {
	   			playersInfo := []*room.DDZPlayerInfo{}

	   			for _, player := range d.ddzContext.GetPlayers() {
	   				// Player转为RoomPlayer
	   				roomPlayerInfo := TranslateDDZPlayerToRoomPlayer(*player)
	   				lord := player.GetLord()
	   				double := player.GetIsDouble()
	   				tuoguan := false // TODO

	   				ddzPlayerInfo := room.DDZPlayerInfo{}
	   				ddzPlayerInfo.PlayerInfo = &roomPlayerInfo
	   				ddzPlayerInfo.OutCards = player.GetOutCards()
	   				ddzPlayerInfo.HandCards = player.GetHandCards()
	   				ddzPlayerInfo.Lord = &lord
	   				ddzPlayerInfo.IsDouble = &double
	   				ddzPlayerInfo.Tuoguan = &tuoguan

	   				playersInfo = append(playersInfo, &ddzPlayerInfo)
	   			}

	   			// 发送游戏信息
	   			d.getMessageSender().([]uint64{reqPlayerID}, msgid.MsgID_ROOM_DDZ_RESUME_REQ, &room.DDZResumeGameRsp{
	   				Result: genResult(0, ""),
	   				GameInfo: &room.DDZDeskInfo{
	   					Players: playersInfo,
	   					Stage:d.get
	   				},
	   			})
	   		}
	   	} */

	result := procedure.HandleEvent(params)
	if !result.Succeed {
		entry.Errorln("处理事件失败")
		return
	}

	d.ddzContext = &result.Context
	// 游戏结束
	if d.ddzContext.GetCurState() == ddz.StateID_state_over {
		go func() { d.Stop() }()
		return
	}
	if result.HasAutoEvent {
		if result.AutoEventDuration == time.Duration(0) {
			d.pushEvent(&deskEvent{
				eventID:      result.AutoEventID,
				eventContext: result.AutoEventContext,
			})
		} else {
			go func() {
				timer := time.NewTimer(result.AutoEventDuration)
				<-timer.C
				d.pushEvent(&deskEvent{
					eventID:      result.AutoEventID,
					eventContext: result.AutoEventContext,
				})
			}()
		}
	}
}

func (d *desk) getMessageSender() ddzmachine.MessageSender {
	return func(players []uint64, msgID msgid.MsgID, body proto.Message) error {
		//logEntry := logrus.WithFields(logrus.Fields{
		//	"func_name":       "deskPlayerMgr.BroadcastMessage",
		//	"dest_player_ids": players,
		//	"msg_id":          msgID,
		//	"msg":             body,
		//})
		//logEntry.Debug("斗地主广播")
		return facade.BroadCastDeskMessage(d, players, msgID, body, true)
	}
}
