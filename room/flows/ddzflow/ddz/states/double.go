package states

import (
	"steve/entity/poker/ddz"
	"steve/room/fixed"
	"steve/room/flows/ddzflow/machine"

	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"time"

	"github.com/Sirupsen/logrus"
)

type doubleState struct{}

func (s *doubleState) OnEnter(m machine.Machine) {
	context := getDDZContext(m)
	context.CurStage = ddz.DDZStage_DDZ_STAGE_DOUBLE
	context.CountDownPlayers = getPlayerIds(m)
	context.StartTime = time.Now()
	context.Duration = StageTime[room.DDZStage_DDZ_STAGE_DOUBLE]
	logrus.WithField("context", context).Debugln("进入加倍状态")
}

func (s *doubleState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开加倍状态")
}

func (s *doubleState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID != int(ddz.EventID_event_double_request) {
		return int(ddz.StateID_state_double), fixed.ErrInvalidEvent
	}

	message := (event.EventData).(*ddz.DoubleRequestEvent)

	context := getDDZContext(m)
	playerID := message.GetHead().GetPlayerId()
	isDouble := message.IsDouble

	logEntry := logrus.WithFields(logrus.Fields{"playerId": playerID, "double": isDouble})
	if !IsValidPlayer(context, playerID) {
		logEntry.WithField("players", getPlayerIds(m)).Errorln("玩家不在本牌桌上!")
		return int(ddz.StateID_state_double), fixed.ErrInvalidRequestPlayer
	}
	for _, doubledPlayer := range context.DoubledPlayers {
		if doubledPlayer == playerID {
			logEntry.WithField("DoubledPlayers", context.DoubledPlayers).Warnln("玩家重复加倍")
			return int(ddz.StateID_state_double), nil
		}
	}
	logEntry.Infoln("斗地主玩家加倍")

	GetPlayerByID(context.GetPlayers(), playerID).IsDouble = isDouble //记录该玩家加倍
	context.DoubledPlayers = append(context.DoubledPlayers, playerID)
	if isDouble {
		context.TotalDouble = context.TotalDouble * 2
	}

	//删除该玩家倒计时
	context.CountDownPlayers = remove(context.CountDownPlayers, playerID)

	var nextStage *room.NextStage
	if len(context.DoubledPlayers) >= 3 {
		nextStage = GenNextStage(room.DDZStage_DDZ_STAGE_PLAYING)
	}
	broadcast(m, msgid.MsgID_ROOM_DDZ_DOUBLE_NTF, &room.DDZDoubleNtf{
		PlayerId:    &playerID,
		IsDouble:    &isDouble,
		TotalDouble: &context.TotalDouble,
		NextStage:   nextStage,
	})

	if len(context.DoubledPlayers) >= 3 {
		context.Duration = 0 //清除倒计时
		return int(ddz.StateID_state_playing), nil
	} else {
		return int(ddz.StateID_state_double), nil
	}
}
