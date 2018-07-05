package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"steve/majong/global"
	"github.com/gogo/protobuf/proto"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
)

type doubleState struct{}

func (s *doubleState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入加倍状态")
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

	context := getDDZContext(m);
	playerId := message.GetHead().GetPlayerId()
	isDouble := message.IsDouble
	GetPlayerByID(context.GetPlayers(), playerId).IsDouble = isDouble//记录该玩家加倍

	var nextStage room.NextStage
	context.DoubledCount++
	if context.DoubledCount >= 3 {
		nextStage = room.NextStage{
			Stage: room.DDZStage_DDZ_STAGE_PLAYING.Enum(),
			Time: proto.Uint32(15),
		}
	}
	totalDouble := GetTotalDouble(context.GetPlayers())
	broadcast(m, msgid.MsgID_ROOM_DDZ_DOUBLE_NTF, &room.DDZDoubleNtf{
		PlayerId: &playerId,
		IsDouble: &isDouble,
		TotalDouble: &totalDouble,
		NextStage: &nextStage,
	})

	if context.DoubledCount >= 3 {
		context.CurrentPlayerId = context.LordPlayerId
		context.CurCardType = ddz.CardType_CT_NONE
		context.TotalBomb = 1
		return int(ddz.StateID_state_playing), nil
	} else {
		return int(ddz.StateID_state_double), nil
	}
}
