package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/golang/protobuf/proto"
	"github.com/Sirupsen/logrus"
	"steve/majong/global"
	"steve/client_pb/room"
	"steve/client_pb/msgId"
	"math/rand"
	"time"
)

type grabState struct{}

func (s *grabState) OnEnter(m machine.Machine) {
	context := getDDZContext(m)
	//产生超时事件
	context.CountDownPlayers = []uint64{context.CurrentPlayerId}
	context.StartTime, _ = time.Now().MarshalBinary()
	context.Duration = StageTime[room.DDZStage_DDZ_STAGE_GRAB]

	logrus.WithField("context", context).Debugln("进入叫/抢地主状态")
}

func (s *grabState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开叫/抢地主状态")
}

func (s *grabState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID != int(ddz.EventID_event_grab_request) {
		logrus.Error("grabState can only handle ddz.EventID_event_grab_request, invalid event")
		return int(ddz.StateID_state_grab), global.ErrInvalidEvent
	}

	message := &ddz.GrabRequestEvent{}
	err := proto.Unmarshal(event.EventData, message)
	if err != nil {
		logrus.Error("grabState unmarshal event error!")
		return int(ddz.StateID_state_grab), global.ErrUnmarshalEvent
	}

	context := getDDZContext(m);
	playerId := message.GetHead().GetPlayerId()
	if context.CurrentPlayerId != playerId {
		logrus.WithField("expected player:", context.CurrentPlayerId).WithField("fact player", playerId).Error("未到本玩家抢地主")
		return int(ddz.StateID_state_grab), global.ErrInvalidRequestPlayer
	}

	grab := message.GetGrab()
	GetPlayerByID(context.GetPlayers(), playerId).Grab = grab//记录该玩家已叫/弃地主

	nextStage := room.DDZStage_DDZ_STAGE_CALL //还没人叫地主，后面还是叫地主阶段
	if context.FirstGrabPlayerId != 0 {
		nextStage = room.DDZStage_DDZ_STAGE_GRAB //有人叫过地主，后面是抢地主阶段
	} else if grab {
		context.FirstGrabPlayerId = playerId;//记录第一次叫地主玩家
		nextStage = room.DDZStage_DDZ_STAGE_GRAB //当前玩家叫地主，后面是抢地主阶段
	}

	nextPlayerId := GetNextPlayerByID(context.GetPlayers(), playerId).PalyerId
	context.GrabbedCount++
	context.LastPlayerId = playerId
	context.CurrentPlayerId = nextPlayerId
	//产生超时事件
	context.CountDownPlayers = []uint64{context.CurrentPlayerId}
	context.StartTime, _ = time.Now().MarshalBinary()
	context.Duration = StageTime[room.DDZStage_DDZ_STAGE_GRAB]

	allAbandon := false;
	if context.GrabbedCount == 3 {//第三个人叫/弃地主时
		if IsAllAbandon(context.GetPlayers()){
			allAbandon = true
			context.AllAbandonCount++
			if context.AllAbandonCount < 3 {
				nextStage = room.DDZStage_DDZ_STAGE_DEAL
			}
			nextPlayerId = 0  //由DDZLordNtf通知下一个玩家
		} else {
			nextPlayerId = context.FirstGrabPlayerId //三个人都叫地主时，由第一个叫地主玩家最后决定
		}
	}

	if grab && context.FirstGrabPlayerId != 0 {
		context.TotalGrab = context.TotalGrab * 2
	}
	if context.GrabbedCount == 4 {
		nextPlayerId = 0 //由DDZLordNtf通知下一个玩家
	}
	broadcast(m, msgid.MsgID_ROOM_DDZ_GRAB_LORD_NTF, &room.DDZGrabLordNtf{
		PlayerId: &playerId,
		Grab: &grab,
		TotalGrab: &context.TotalGrab,
		NextPlayerId: &nextPlayerId,
		NextStage: genNextStage(nextStage),
	})

	lordPlayerId := uint64(0)
	if context.GrabbedCount == 4 {
		if grab {//叫地主玩家抢庄
			lordPlayerId = playerId
		} else {//叫地主玩家弃庄
			lordPlayerId = context.LastPlayerId
		}
		broadcast(m, msgid.MsgID_ROOM_DDZ_LORD_NTF, &room.DDZLordNtf{
			PlayerId: &lordPlayerId,
			TotalGrab: &context.TotalGrab,
			Dipai: context.WallCards,
			NextStage: genNextStage(room.DDZStage_DDZ_STAGE_DOUBLE),
		})
	}

	if allAbandon && context.AllAbandonCount < 3 {
		return int(ddz.StateID_state_deal), nil //重新发牌
	}
	context.AllAbandonCount = 0

	if allAbandon {//三轮没人叫地主，随机确定庄家
		i := rand.Intn(3)
		lordPlayerId = context.GetPlayers()[i+1].PalyerId

		broadcast(m, msgid.MsgID_ROOM_DDZ_LORD_NTF, &room.DDZLordNtf{
			PlayerId: &lordPlayerId,
			TotalGrab: &context.TotalGrab,
			Dipai: context.WallCards,
			NextStage: genNextStage(room.DDZStage_DDZ_STAGE_DOUBLE),
		})
	}

	if lordPlayerId != 0 {
		lordPlayer := GetPlayerByID(context.GetPlayers(), lordPlayerId)
		lordPlayer.Lord = true
		for _, card := range context.WallCards {
			lordPlayer.HandCards = append(lordPlayer.HandCards, card)
		}
		lordPlayer.HandCards = DDZSortDescend(lordPlayer.HandCards)
		context.WallCards = []uint32{}
		context.LordPlayerId = lordPlayerId
		return int(ddz.StateID_state_double), nil
	} else {
		return int(ddz.StateID_state_grab), nil
	}
}
