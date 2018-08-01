package models

import (
	"context"
	"runtime/debug"
	"steve/client_pb/msgid"
	"steve/room/ai"
	context2 "steve/room/contexts"
	"steve/room/desk"
	"steve/room/fixed"
	"steve/room/flows/ddzflow/ddz/ddzmachine"
	"steve/room/flows/ddzflow/ddz/procedure"
	"steve/room/player"
	"steve/server_pb/ddz"
	"steve/structs/proto/gate_rpc"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// DDZEventModel 斗地主事件 model
type DDZEventModel struct {
	event chan desk.DeskEvent // 牌桌事件通道
	BaseModel
	closed bool
	mu     sync.Mutex
}

// NewDDZEventModel 创建斗地主事件 model
func NewDDZEventModel(desk *desk.Desk) DeskModel {
	result := &DDZEventModel{}
	result.SetDesk(desk)
	return result
}

// GetName 获取 model 名称
func (model *DDZEventModel) GetName() string {
	return fixed.EventModelName
}

// Active 激活 model
func (model *DDZEventModel) Active() {}

// Start 启动 model
func (model *DDZEventModel) Start() {
	model.event = make(chan desk.DeskEvent, 16)

	go func() {
		model.processEvents(context.Background())
		GetModelManager().StopDeskModel(model.GetDesk().GetUid())
	}()
	params := desk.CreateEventParams([]byte{}, 0)
	event := desk.NewDeskEvent(int(ddz.EventID_event_start_game), fixed.NormalEvent, model.GetDesk(), params)
	model.PushEvent(event)
}

// Stop 停止
func (model *DDZEventModel) Stop() {
	model.mu.Lock()
	if model.closed {
		model.mu.Unlock()
		return
	}
	model.closed = true
	close(model.event)
	model.mu.Unlock()
}

// PushEvent 压入事件
func (model *DDZEventModel) PushEvent(event desk.DeskEvent) {
	model.mu.Lock()
	if model.closed {
		model.mu.Unlock()
		return
	}
	model.event <- event
	model.mu.Unlock()
}

// PushRequest 压入玩家请求
func (model *DDZEventModel) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	entry := logrus.WithFields(logrus.Fields{
		"player_id": playerID,
		"msg_id":    head.GetMsgId(),
	})

	trans := GetTranslator()
	eventID, eventData, err := trans.Translate(playerID, head, bodyData)
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
	eventParams := desk.CreateEventParams(eventContext, playerID)
	event := desk.NewDeskEvent(eventID, fixed.NormalEvent, model.GetDesk(), eventParams)
	model.PushEvent(event)
}

func (model *DDZEventModel) processEvents(ctx context.Context) {
	logEntry := logrus.WithFields(logrus.Fields{
		"desk_uid": model.GetDesk().GetUid(),
		"game_id":  model.GetDesk().GetGameId(),
	})
	defer func() {
		if x := recover(); x != nil {
			logEntry.Errorln(x)
			debug.PrintStack()
		}
	}()
	playerModel := GetModelManager().GetPlayerModel(model.GetDesk().GetUid())
	playerEnterChannel := playerModel.getEnterChannel()
	playerLeaveChannel := playerModel.getLeaveChannel()
	tick := time.NewTicker(time.Millisecond * 200)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			{
				logEntry.Infoln("done")
				return
			}
		case enterInfo := <-playerEnterChannel:
			{
				model.handlePlayerEnter(enterInfo)
			}
		case leaveInfo := <-playerLeaveChannel:
			{
				model.handlePlayerLeave(leaveInfo)
			}
		case event := <-model.event:
			{
				eventContext := model.getEventContext(event)
				if model.processEvent(event.EventID, eventContext) {
					return
				}
			}
		case <-tick.C:
			{
				events := model.genTimerEvent()
				for _, event := range events {
					context := model.getEventContext(event)
					if model.processEvent(event.EventID, context) {
						return
					}
					model.recordTuoguanOverTimeCount(event)
				}
			}
		}
	}
}

// handlePlayerEnter 处理玩家进入牌桌
func (model *DDZEventModel) handlePlayerEnter(enterInfo playerIDWithChannel) {
	entry := logrus.WithField("player_id", enterInfo.playerID)

	modelMgr := GetModelManager()
	modelMgr.GetPlayerModel(model.GetDesk().GetUid()).handlePlayerEnter(enterInfo.playerID)
	// 生成恢复对局事件
	eventMessage := &ddz.ResumeRequestEvent{
		Head: &ddz.RequestEventHead{PlayerId: enterInfo.playerID},
	}
	for {
		eventID := int(ddz.EventID_event_resume_request)
		eventContext, err := proto.Marshal(eventMessage)
		if err != nil {
			entry.WithError(err).Errorln("事件消息序列化失败")
			break
		}
		model.processEvent(eventID, eventContext)
		entry.Debugln("玩家进入")
		break
	}
	close(enterInfo.finishChannel)
}

// handlePlayerLeave 处理玩家离开牌桌
func (model *DDZEventModel) handlePlayerLeave(leaveInfo playerIDWithChannel) {
	modelMgr := GetModelManager()
	playerID := leaveInfo.playerID

	modelMgr.GetPlayerModel(model.GetDesk().GetUid()).handlePlayerLeave(playerID, true)
	logrus.WithField("player_id", playerID).Debugln("玩家退出")
	close(leaveInfo.finishChannel)
}

// getEventPlayerID  获取事件玩家
func (model *DDZEventModel) getEventPlayerID(event desk.DeskEvent) uint64 {
	return event.Params.Params[1].(uint64)
}

// getEventContext 获取事件现场
func (model *DDZEventModel) getEventContext(event desk.DeskEvent) []byte {
	return event.Params.Params[0].([]byte)
}

// recordTuoguanOverTimeCount 记录托管超时计数
func (model *DDZEventModel) recordTuoguanOverTimeCount(event desk.DeskEvent) {
	if event.EventType != fixed.OverTimeEvent {
		return
	}
	playerID := model.getEventPlayerID(event)
	if playerID == 0 {
		return
	}
	deskPlayer := player.GetPlayerMgr().GetPlayer(playerID)
	if deskPlayer != nil {
		deskPlayer.OnPlayerOverTime()
	}
}

func (model *DDZEventModel) getMessageSender() ddzmachine.MessageSender {
	messageModel := GetModelManager().GetMessageModel(model.GetDesk().GetUid())
	return func(players []uint64, msgID msgid.MsgID, body proto.Message) error {
		return messageModel.BroadCastDeskMessage(players, msgID, body, true)
	}
}

// processEvent 处理单个事件
// step 1. 调用麻将逻辑的接口来处理事件(返回最新麻将现场, 自动事件， 发送给玩家的消息)， 并且更新 mjContext
// step 2. 将消息发送给玩家
// step 3. 调用 room 的结算逻辑来处理结算
// step 4. 如果有自动事件， 将自动事件写入自动事件通道
// step 5. 如果当前状态是游戏结束状态， 调用 cancel 终止游戏
// 返回： 游戏是否结束
func (model *DDZEventModel) processEvent(eventID int, eventContext []byte) bool {
	entry := logrus.WithFields(logrus.Fields{
		"event_id": eventID,
	})
	gameContext := model.GetGameContext().(*context2.DDZDeskContext)
	params := procedure.Params{
		Context:      gameContext.DDZContext,
		Sender:       model.getMessageSender(),
		EventID:      eventID,
		EventContext: eventContext,
	}

	result := procedure.HandleEvent(params)
	if !result.Succeed {
		entry.Errorln("处理事件失败")
		return false
	}
	gameContext.DDZContext = result.Context

	// 自动事件不为空，继续处理事件
	if result.HasAutoEvent {
		if result.AutoEventDuration == time.Duration(0) {
			return model.processEvent(result.AutoEventID, result.AutoEventContext)
		}
		go func() {
			timer := time.NewTimer(result.AutoEventDuration)
			<-timer.C
			eventParams := desk.CreateEventParams(result.AutoEventContext, 0)
			model.PushEvent(desk.NewDeskEvent(result.AutoEventID, fixed.NormalEvent, model.GetDesk(), eventParams))
		}()
	}
	return model.checkGameOver(entry)
}

// checkGameOver 检查游戏结束
func (model *DDZEventModel) checkGameOver(logEntry *logrus.Entry) bool {
	gameContext := model.GetGameContext().(*context2.DDZDeskContext)
	ddzContext := gameContext.DDZContext

	if ddzContext.GetCurState() == ddz.StateID_state_over {
		continueModel := GetContinueModel(model.GetDesk().GetUid())
		players := ddzContext.GetPlayers()
		statistics := make(map[uint64]int64, len(players))

		for _, player := range players {
			if player.GetWin() {
				statistics[player.GetPlayerId()] = 1
			} else {
				statistics[player.GetPlayerId()] = -1
			}
		}
		continueModel.ContinueDesk(false, 0, statistics)
		return true
	}
	return false
}

// genTimerEvent 生成计时事件
func (model *DDZEventModel) genTimerEvent() []desk.DeskEvent {
	playerModel := GetModelManager().GetPlayerModel(model.GetDesk().GetUid())
	dContext := model.GetDesk().GetConfig().Context.(*context2.DDZDeskContext)
	ddzContext := &dContext.DDZContext

	deskPlayers := playerModel.GetDeskPlayers()
	robotLvs := make(map[uint64]int, len(deskPlayers))
	for _, deskPlayer := range deskPlayers {
		robotLv := deskPlayer.GetRobotLv()
		if robotLv != 0 {
			robotLvs[deskPlayer.GetPlayerID()] = robotLv
		}
	}
	startTime := time.Time{}
	startTime.UnmarshalBinary(ddzContext.GetStartTime())
	// 产生AI事件
	result := ai.GetAtEvent().GenerateV2(&ai.AutoEventGenerateParams{
		Desk:      model.GetDesk(),
		StartTime: startTime,
		RobotLv:   robotLvs,
	})
	return result.Events
}
