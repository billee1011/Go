package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"math/rand"
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
	players := getPlayers(m)
	i := rand.Intn(len(players))
	callPlayer := players[i+1]//叫地主玩家
	var stageTime uint32 = 4;
	broadcast(m, msgid.MsgID_ROOM_DDZ_START_GAME_NTF, &room.DDZStartGameNtf{
		PlayerId:&callPlayer,
		NextStage:&room.NextStage{Stage: room.DDZStage_DDZ_STAGE_DEAL.Enum(), Time: &stageTime},
	})
	return int(ddz.StateID_state_deal), nil
}
