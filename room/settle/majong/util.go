package majong

import (
	msgid "steve/client_pb/msgid"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// NotifyMessage 将消息广播给牌桌所有玩家
func NotifyMessage(desk interfaces.Desk, msgid msgid.MsgID, message proto.Message) {
	facade.BroadCastDeskMessage(desk, nil, msgid, message, true)
}

// NotifyPlayersMessage 将消息广播给牌桌指定playerIds[]中的玩家
func NotifyPlayersMessage(desk interfaces.Desk, playerIds []uint64, msgid msgid.MsgID, message proto.Message) {
	facade.BroadCastDeskMessage(desk, playerIds, msgid, message, true)
}

// GetDeskPlayer 获取指定id的room Player
func GetDeskPlayer(deskPlayers []interfaces.DeskPlayer, pid uint64) interfaces.DeskPlayer {
	for _, p := range deskPlayers {
		if p.GetPlayerID() == pid {
			return p
		}
	}
	return nil
}

// GetSettleInfoBySid 根据settleID获取对应settleInfo的下标index
func GetSettleInfoBySid(settleInfos []*majongpb.SettleInfo, ID uint64) int {
	for index, s := range settleInfos {
		if s.Id == ID {
			return index
		}
	}
	return -1
}

// GenerateSettleEvent 结算finish事件
func GenerateSettleEvent(desk interfaces.Desk, settleType majongpb.SettleType, brokerPlayers []uint64) {
	needEvent := map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_angang:   true,
		majongpb.SettleType_settle_bugang:   true,
		majongpb.SettleType_settle_minggang: true,
		majongpb.SettleType_settle_dianpao:  true,
		majongpb.SettleType_settle_zimo:     true,
	}
	if needEvent[settleType] {
		eventContext, err := proto.Marshal(&majongpb.SettleFinishEvent{
			PlayerId: brokerPlayers,
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"func_name":     "GenerateSettleEvent",
				"settleType":    settleType,
				"brokerPlayers": brokerPlayers,
			}).WithError(err).Errorln("消息序列化失败")
			return
		}
		event := majongpb.AutoEvent{
			EventId:      majongpb.EventID_event_settle_finish,
			EventContext: eventContext,
		}
		desk.PushEvent(interfaces.Event{
			ID:        event.GetEventId(),
			Context:   event.GetEventContext(),
			EventType: interfaces.NormalEvent,
			PlayerID:  0,
		})
	}
}

// GetSettleInfoByID 根据settleID获取对应settleInfo
func GetSettleInfoByID(settleInfos []*majongpb.SettleInfo, ID uint64) *majongpb.SettleInfo {
	for _, s := range settleInfos {
		if s.Id == ID {
			return s
		}
	}
	return nil
}
