package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"time"
)

type doubleState struct{}

func (s *doubleState) OnEnter(m machine.Machine) {
	context := getDDZContext(m)
	context.CurStage = ddz.DDZStage_DDZ_STAGE_DOUBLE
	context.CountDownPlayers = getPlayerIds(m)
	context.StartTime, _ = time.Now().MarshalBinary()
	context.Duration = StageTime[room.DDZStage_DDZ_STAGE_DOUBLE]
	logrus.WithField("context", context).Debugln("进入加倍状态")
}

func (s *doubleState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开加倍状态")
}

func (s *doubleState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID != int(ddz.EventID_event_double_request) {
		return int(ddz.StateID_state_double), global.ErrInvalidEvent
	}

	message := &ddz.DoubleRequestEvent{}
	err := proto.Unmarshal(event.EventData, message)
	if err != nil {
		return int(ddz.StateID_state_double), global.ErrUnmarshalEvent
	}

	context := getDDZContext(m)
	playerId := message.GetHead().GetPlayerId()
	if !isValidPlayer(context, playerId) {
		logrus.Error("玩家不在本牌桌上!")
		return int(ddz.StateID_state_double), global.ErrInvalidRequestPlayer
	}

	isDouble := message.IsDouble
	GetPlayerByID(context.GetPlayers(), playerId).IsDouble = isDouble //记录该玩家加倍
	context.DoubledPlayers = append(context.DoubledPlayers, playerId)
	if isDouble {
		context.TotalDouble = context.TotalDouble * 2
	}

	//删除该玩家倒计时
	context.CountDownPlayers = remove(context.CountDownPlayers, playerId)

	var nextStage *room.NextStage
	context.DoubledCount++
	if context.DoubledCount >= 3 {
		nextStage = GenNextStage(room.DDZStage_DDZ_STAGE_PLAYING)
	}
	broadcast(m, msgid.MsgID_ROOM_DDZ_DOUBLE_NTF, &room.DDZDoubleNtf{
		PlayerId:    &playerId,
		IsDouble:    &isDouble,
		TotalDouble: &context.TotalDouble,
		NextStage:   nextStage,
	})

	if context.DoubledCount >= 3 {
		context.CurrentPlayerId = context.LordPlayerId
		context.Duration = 0 //清除倒计时
		return int(ddz.StateID_state_playing), nil
	} else {
		return int(ddz.StateID_state_double), nil
	}
}
