package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"math/rand"
	"steve/client_pb/msgid"
	"steve/client_pb/room"

	"github.com/Sirupsen/logrus"
	"errors"
)

type initState struct{}

func (s *initState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入初始状态")
}

func (s *initState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开初始状态")
}

func (s *initState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID == int(ddz.EventID_event_start_game) {
		return s.onStartGame(m)
	}
	return int(ddz.StateID_state_init), nil
}

// 开局随即确定一个叫地主玩家,然后广播通知游戏开始
func (s *initState) onStartGame(m machine.Machine) (int, error) {
	logrus.WithField("context", getDDZContext(m)).Debugln("开始游戏")
	players := getPlayerIds(m)
	if len(players) != 3 {
		return int(ddz.StateID_state_init), errors.New("玩家人数错误")
	}
	i := rand.Intn(len(players))
	callPlayer := players[i] //叫地主玩家
	context := getDDZContext(m)
	context.CurrentPlayerId = callPlayer
	context.GrabbedCount = 0
	context.FirstGrabPlayerId = 0
	broadcast(m, msgid.MsgID_ROOM_DDZ_START_GAME_NTF, &room.DDZStartGameNtf{
		PlayerId:  &callPlayer,
		NextStage: genNextStage(room.DDZStage_DDZ_STAGE_DEAL),
	})
	return int(ddz.StateID_state_deal), nil
}
