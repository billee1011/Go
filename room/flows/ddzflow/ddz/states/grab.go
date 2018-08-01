package states

import (
	"steve/room/flows/ddzflow/machine"
	"steve/server_pb/ddz"

	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/server_pb/majong"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type grabState struct{}

func (s *grabState) OnEnter(m machine.Machine) {
	context := getDDZContext(m)
	context.CurStage = ddz.DDZStage_DDZ_STAGE_CALL
	context.CurrentPlayerId = context.CallPlayerId
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
	if event.EventID == int(majong.EventID_event_cartoon_finish_request) {
		return int(ddz.StateID_state_grab), nil
	}
	if event.EventID != int(ddz.EventID_event_grab_request) {
		return int(ddz.StateID_state_grab), global.ErrInvalidEvent
	}

	message := &ddz.GrabRequestEvent{}
	err := proto.Unmarshal(event.EventData, message)
	if err != nil {
		return int(ddz.StateID_state_grab), global.ErrUnmarshalEvent
	}

	context := getDDZContext(m)
	playerId := message.GetHead().GetPlayerId()
	grab := message.GetGrab()

	logEntry := logrus.WithFields(logrus.Fields{"playerId": playerId, "grab": grab})
	if !isValidPlayer(context, playerId) {
		logEntry.WithField("players", getPlayerIds(m)).Errorln("玩家不在本牌桌上!")
		return int(ddz.StateID_state_grab), global.ErrInvalidRequestPlayer
	}
	if context.CurrentPlayerId != playerId {
		logEntry.WithField("expected player:", context.CurrentPlayerId).Errorln("未到本玩家抢地主")
		return int(ddz.StateID_state_grab), global.ErrInvalidRequestPlayer
	}
	logEntry.Infoln("玩家叫/抢地主")

	GetPlayerByID(context.GetPlayers(), playerId).Grab = grab //记录该玩家已叫/弃地主
	context.GrabbedCount++                                    //记录完毕

	if context.FirstGrabPlayerId == 0 && grab { //第一次叫地主
		context.FirstGrabPlayerId = playerId
		context.TotalGrab = 1
		context.CurStage = ddz.DDZStage_DDZ_STAGE_GRAB
	} else if context.FirstGrabPlayerId != 0 && grab { //抢地主
		context.TotalGrab = context.TotalGrab * 2
		context.LastGrabPlayerId = playerId
		context.GrabbedPlayers = append(context.GrabbedPlayers, playerId)
	}

	nextPlayerId := GetNextPlayerByID(context.GetPlayers(), playerId).PlayerId
	lordPlayerId := uint64(0)      //不为0时确定地主
	if context.GrabbedCount == 3 { //三个玩家操作完毕
		if context.FirstGrabPlayerId == 0 { //没人叫地主
			context.AllAbandonCount++
			if context.AllAbandonCount <= 3 {
				context.CurStage = ddz.DDZStage_DDZ_STAGE_DEAL
			}
			nextPlayerId = 0 //重新发牌，没有操作玩家
		}

		if context.TotalGrab == 1 { //只有一个人叫，其他两个人弃时，地主为叫地主的人
			lordPlayerId = context.FirstGrabPlayerId
			nextPlayerId = 0 //确定地主，进入加倍阶段，没有操作玩家
		} else { //有人叫，则由叫地主玩家最后决定
			nextPlayerId = context.FirstGrabPlayerId
		}
	}

	if context.GrabbedCount >= 4 { //叫地主玩家第二次操作
		if grab { //叫地主玩家抢庄
			lordPlayerId = playerId
		} else { //叫地主玩家弃庄
			lordPlayerId = context.LastGrabPlayerId
		}
	}

	if nextPlayerId == 0 {
		context.Duration = 0 //清除倒计时
	} else {
		//更新当前操作用户并产生超时事件
		context.CurrentPlayerId = nextPlayerId
		context.CountDownPlayers = []uint64{context.CurrentPlayerId}
		context.StartTime, _ = time.Now().MarshalBinary()
		context.Duration = StageTime[room.DDZStage_DDZ_STAGE_GRAB]
	}

	if lordPlayerId != 0 || context.AllAbandonCount > 3 {
		context.CurStage = ddz.DDZStage_DDZ_STAGE_NONE //最后一个人的抢地主广播NextStage为NONE
	}

	totalGrab := context.TotalGrab
	if totalGrab == 0 { //产品要求不能显示0倍
		totalGrab = 1
	}
	broadcast(m, msgid.MsgID_ROOM_DDZ_GRAB_LORD_NTF, &room.DDZGrabLordNtf{
		PlayerId:     &playerId,
		Grab:         &grab,
		TotalGrab:    &totalGrab,
		NextPlayerId: &nextPlayerId,
		NextStage:    GenNextStage(room.DDZStage(int32(context.CurStage))),
	})

	if context.CurStage == ddz.DDZStage_DDZ_STAGE_DEAL {
		context.GrabbedCount = 0
		context.CallPlayerId = getRandPlayerId(context.GetPlayers())
		return int(ddz.StateID_state_deal), nil //重新发牌
	}

	if context.AllAbandonCount > 3 { //三轮重新发牌没人叫地主，随机确定庄家
		context.AllAbandonCount = 0
		lordPlayerId = getRandPlayerId(context.GetPlayers())
		context.TotalGrab = 1
	}

	if lordPlayerId != 0 {
		lordPlayer := GetPlayerByID(context.GetPlayers(), lordPlayerId)
		lordPlayer.Lord = true
		for _, card := range context.Dipai {
			lordPlayer.HandCards = append(lordPlayer.HandCards, card)
		}
		lordPlayer.HandCards = DDZSortDescend(lordPlayer.HandCards)
		context.LordPlayerId = lordPlayerId
		context.Duration = 0 //清除倒计时
		context.CurStage = ddz.DDZStage_DDZ_STAGE_DOUBLE
		broadcast(m, msgid.MsgID_ROOM_DDZ_LORD_NTF, &room.DDZLordNtf{
			PlayerId:  &lordPlayerId,
			TotalGrab: &context.TotalGrab,
			Dipai:     context.Dipai,
			NextStage: GenNextStage(room.DDZStage_DDZ_STAGE_DOUBLE),
		})
		return int(ddz.StateID_state_double), nil
	} else {
		return int(ddz.StateID_state_grab), nil
	}
}
