package majong

import (
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/room/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

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

// mergeSettle 合并一组SettleInfo
// 返回参数:	[]*majongpb.SettleInfo(该组settleInfo) / *majongpb.SettleInfo(合并后的settleInfo)
func mergeSettle(contextSInfo []*majongpb.SettleInfo, settleInfo *majongpb.SettleInfo) ([]*majongpb.SettleInfo, *majongpb.SettleInfo) {
	sumSInfo := &majongpb.SettleInfo{
		Scores: make(map[uint64]int64, 0),
	}
	groupSInfos := make([]*majongpb.SettleInfo, 0)
	for _, id := range settleInfo.GroupId {
		sIndex := GetSettleInfoBySid(contextSInfo, id)
		groupSInfos = append(groupSInfos, contextSInfo[sIndex])
		sumSInfo.SettleType = contextSInfo[sIndex].SettleType
	}
	for _, singleSInfo := range groupSInfos {
		for pid, score := range singleSInfo.Scores {
			sumSInfo.Scores[pid] = sumSInfo.Scores[pid] + score
		}
	}
	return groupSInfos, sumSInfo
}

func settleType2BillType(settleType majongpb.SettleType) room.BillType {
	return map[majongpb.SettleType]room.BillType{
		majongpb.SettleType_settle_angang:    room.BillType_BILL_GANG,
		majongpb.SettleType_settle_bugang:    room.BillType_BILL_GANG,
		majongpb.SettleType_settle_minggang:  room.BillType_BILL_GANG,
		majongpb.SettleType_settle_dianpao:   room.BillType_BILL_DIANPAO,
		majongpb.SettleType_settle_zimo:      room.BillType_BILL_ZIMO,
		majongpb.SettleType_settle_yell:      room.BillType_BILL_CHECKSHOUT,
		majongpb.SettleType_settle_flowerpig: room.BillType_BILL_CHECKPIG,
		majongpb.SettleType_settle_calldiver: room.BillType_BILL_TRANSFER,
		majongpb.SettleType_settle_taxrebeat: room.BillType_BILL_REFUND,
	}[settleType]
}

func getFans(fanTypes []int64, huaCount uint32, cardOption *mjoption.CardTypeOption) (fan []*room.Fan) {
	fan = make([]*room.Fan, 0)
	for _, fanType := range fanTypes {
		rfan := &room.Fan{
			Name:  room.FanType(int32(fanType)).Enum(),
			Value: proto.Int32(int32(cardOption.Fantypes[int(fanType)].Score)),
			Type:  proto.Uint32(uint32(cardOption.Fantypes[int(fanType)].Type)),
		}
		fan = append(fan, rfan)
	}
	if huaCount != 0 {
		fan = append(fan, &room.Fan{
			Name:  room.FanType_FT_HUAPAI.Enum(),
			Value: proto.Int32(int32(huaCount)),
			Type:  proto.Uint32(0),
		})
	}
	return
}

// getGiveupPlayers  获取认输的玩家id
func getGiveupPlayers(dPlayers []interfaces.DeskPlayer, mjContext majongpb.MajongContext) map[uint64]bool {
	giveupPlayers := make(map[uint64]bool, 0)
	for _, cPlayer := range mjContext.Players {
		if cPlayer.GetXpState() == 2 {
			giveupPlayers[cPlayer.GetPalyerId()] = true
		}
	}
	return giveupPlayers
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
