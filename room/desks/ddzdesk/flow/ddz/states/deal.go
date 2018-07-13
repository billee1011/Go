package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"math/rand"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"time"
)

type dealState struct{}

var wallCards = []uint32{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D,
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D,
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D,
	0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D,
	0x0E, 0x0F,
}

func (s *dealState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入发牌状态")
	rand.Shuffle(len(wallCards), func(i, j int) {
		wallCards[i], wallCards[j] = wallCards[j], wallCards[i]
	})
	PeiPai(wallCards, getDDZContext(m).Peipai)
	context := getDDZContext(m)
	context.CurStage = ddz.DDZStage_DDZ_STAGE_DEAL
	players := context.GetPlayers()
	for i := range players {
		players[i].HandCards = DDZSortDescend(wallCards[i*17 : (i+1)*17])
		players[i].OutCards = make([]uint32, 0)
		sendToPlayer(m, players[i].PalyerId, msgid.MsgID_ROOM_DDZ_DEAL_NTF, &room.DDZDealNtf{
			Cards:     players[i].HandCards,
			NextStage: GenNextStage(room.DDZStage_DDZ_STAGE_CALL),
		})
	}
	context.WallCards = wallCards[51:]
	delay := StageTime[room.DDZStage_DDZ_STAGE_DEAL]
	duration := time.Second * time.Duration(delay)
	setMachineAutoEvent(m, machine.Event{EventID: int(ddz.EventID_event_deal_finish), EventData: nil}, duration)
}

func (s *dealState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开发牌状态")
}

func (s *dealState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	logrus.WithFields(logrus.Fields{
		"context": getDDZContext(m),
		"event":   event,
	}).Debugln("发牌完成")
	if event.EventID == int(ddz.EventID_event_deal_finish) {
		return int(ddz.StateID_state_grab), nil
	}
	return int(ddz.StateID_state_deal), global.ErrInvalidEvent
}
