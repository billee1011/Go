package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"steve/client_pb/room"

	"errors"
	"github.com/Sirupsen/logrus"
	"steve/client_pb/msgid"
)

type initState struct{}

func (s *initState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入初始状态")
}

func (s *initState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开初始状态")
}

func (s *initState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID != int(ddz.EventID_event_start_game) {
		return int(ddz.StateID_state_init), nil
	}
	logrus.WithField("context", getDDZContext(m)).Debugln("开始游戏")
	context := getDDZContext(m)
	if len(context.GetPlayers()) != 3 {
		return int(ddz.StateID_state_init), errors.New("玩家人数错误")
	}

	// 开局随即确定一个叫地主玩家,然后广播通知游戏开始
	context.CallPlayerId = getRandPlayerId(context.GetPlayers())
	broadcast(m, msgid.MsgID_ROOM_DDZ_START_GAME_NTF, &room.DDZStartGameNtf{
		PlayerId:  &context.CallPlayerId,
		NextStage: GenNextStage(room.DDZStage_DDZ_STAGE_DEAL),
	})
	return int(ddz.StateID_state_deal), nil
}
