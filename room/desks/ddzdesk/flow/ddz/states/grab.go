package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"steve/majong/global"
	"steve/client_pb/room"
	"steve/client_pb/msgId"
	"math/rand"
	"github.com/gogo/protobuf/proto"
)

type grabState struct{}

func (s *grabState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入叫/抢地主状态")
}

func (s *grabState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开叫/抢地主状态")
}

func (s *grabState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID != int(ddz.EventID_event_grab_request) {
		return int(ddz.StateID_state_grab), global.ErrInvalidEvent
	}

	message := &ddz.GrabRequestEvent{}
	err := proto.Unmarshal(event.EventData, message)
	if err != nil {
		return int(ddz.StateID_state_grab), global.ErrUnmarshalEvent
	}

	context := getDDZContext(m);
	playerId := message.GetHead().GetPlayerId()
	if context.CurrentPlayerId != playerId {
		return int(ddz.StateID_state_grab), global.ErrInvalidRequestPlayer
	}

	grab := message.GetGrab()
	GetPlayerByID(context.GetPlayers(), playerId).Grab = grab//记录该玩家已叫/弃地主

	nextStage := room.DDZStage_DDZ_STAGE_CALL.Enum() //还没人叫地主，后面还是叫地主阶段
	if context.FirstGrabPlayerId != 0 {
		nextStage = room.DDZStage_DDZ_STAGE_GRAB.Enum() //有人叫地主，后面是抢地主阶段
	} else if grab {
		context.FirstGrabPlayerId = playerId;//记录第一次叫地主玩家
	}

	nextPlayerId := GetNextPlayerByID(context.GetPlayers(), playerId).PalyerId
	if context.GrabbedCount == 2 {//第三个人叫/弃地主时，此时grabbedCount还没有++
		if IsAllAbandon(context.GetPlayers()){
			nextPlayerId = 0  //由DDZLordNtf通知下一个玩家
		} else {
			nextPlayerId = context.FirstGrabPlayerId //三个人都叫地主时，由第一个叫地主玩家最后决定
		}
	}

	broadcast(m, msgid.MsgID_ROOM_DDZ_GRAB_LORD_NTF, &room.DDZGrabLordNtf{
		PlayerId: &playerId,
		Grab: &grab,
		NextPlayerId: &nextPlayerId,
		NextStage: &room.NextStage{
			Stage: nextStage,
			Time: proto.Uint32(15),
		},
	})
	context.GrabbedCount++

	totalGrab := GetTotalGrab(context.GetPlayers())
	lordPlayerId := uint64(0)
	if context.GrabbedCount == 3 && IsAllAbandon(context.GetPlayers()) {//没人叫地主，随机确定庄家
		i := rand.Intn(3)
		lordPlayerId = context.GetPlayers()[i+1].PalyerId

		broadcast(m, msgid.MsgID_ROOM_DDZ_LORD_NTF, &room.DDZLordNtf{
			PlayerId: &lordPlayerId,
			TotalGrab: &totalGrab,
			Dipai: context.WallCards,
			NextStage: &room.NextStage{
				Stage: room.DDZStage_DDZ_STAGE_DOUBLE.Enum(),
				Time: proto.Uint32(15),
			},
		})
	}

	if context.GrabbedCount == 4 {
		if grab {//叫地主玩家抢庄
			lordPlayerId = playerId
			totalGrab = totalGrab * 2

			broadcast(m, msgid.MsgID_ROOM_DDZ_LORD_NTF, &room.DDZLordNtf{
				PlayerId: &lordPlayerId,
				TotalGrab: &totalGrab,
				Dipai: context.WallCards,
				NextStage: &room.NextStage{
					Stage: room.DDZStage_DDZ_STAGE_DOUBLE.Enum(),
					Time: proto.Uint32(15),
				},
			})
		} else {//叫地主玩家弃庄
			lordPlayerId = context.LastPlayerId

			broadcast(m, msgid.MsgID_ROOM_DDZ_LORD_NTF, &room.DDZLordNtf{
				PlayerId: &lordPlayerId,
				TotalGrab: &totalGrab,
				Dipai: context.WallCards,
				NextStage: &room.NextStage{
					Stage: room.DDZStage_DDZ_STAGE_DOUBLE.Enum(),
					Time: proto.Uint32(15),
				},
			})
		}
	}

	context.LastPlayerId = playerId
	context.CurrentPlayerId = nextPlayerId

	if lordPlayerId != 0 {
		lordPlayer := GetPlayerByID(context.GetPlayers(), lordPlayerId)
		lordPlayer.Lord = true
		for _, card := range context.WallCards {
			lordPlayer.HandCards = append(lordPlayer.HandCards, card)
		}
		lordPlayer.HandCards = ddzSort(lordPlayer.HandCards)
		context.WallCards = []uint32{}
		return int(ddz.StateID_state_double), nil
	} else {
		return int(ddz.StateID_state_grab), nil
	}
}
