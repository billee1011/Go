package settle

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/gutils"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// scxlSettle 血流麻将结算
type scxlSettle struct {
	currentIds []uint64
	// 每条setttleInfo中每个玩家实际输赢分 key:settleId value:playerCoin
	settleMap map[uint64]scxlplayerCoin
	// 汇总setttleInfo中每个玩家输赢总分 key:playerID value:score
	roundScore map[uint64]int64
	// setttleInfo处理情况 		key:settleId value:true为已处理，false为未处理
	handleSettle map[uint64]bool
}

// newScxlSettle 创建四川血流结算
func newScxlSettle() *scxlSettle {
	return &scxlSettle{
		settleMap:    make(map[uint64]scxlplayerCoin),
		handleSettle: make(map[uint64]bool),
		roundScore:   make(map[uint64]int64),
	}
}

// scxlplayerCoin 玩家实际输赢分   key:playerID value:score
type scxlplayerCoin map[uint64]int64

// Settle 单次结算
// 将玩家输赢分数及实际金币数进行计算，生成实际输赢的分数并记录，广播结算信息给牌局
func (s *scxlSettle) Settle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	// 牌局所有结算信息
	contextSInfos := mjContext.SettleInfos
	// 牌局玩家
	deskPlayers := desk.GetPlayers()
	// 若存在未处理的结算信息，进行处理
	if len(contextSInfos) != 0 {
		for _, SInfo := range contextSInfos {
			if !s.handleSettle[SInfo.Id] {
				// 记录玩家结算信息
				billplayerInfos := make([]*room.BillPlayerInfo, 0)
				// 记录玩家输赢分数
				pidScore := make(map[uint64]int64, 0)
				// 记录金币为0的玩家id
				giveupPlayers := make([]uint64, 0)
				// 若存在相关联的一组SettleInfo(一炮多响情况)
				if len(SInfo.GroupId) > 1 {
					// 合并该组settleInfo,计算实际输赢分
					groupSInfos, sumSInfo := s.sumSettleInfo(mjContext.SettleInfos, SInfo)
					pidScore, giveupPlayers = s.calcScore(deskPlayers, sumSInfo)
					s.resolveScore(groupSInfos, pidScore)
				} else {
					// 单条settleInfo直接计算输赢分
					pidScore, giveupPlayers = s.calcScore(deskPlayers, SInfo)
					s.settleMap[SInfo.Id] = pidScore
					s.handleSettle[SInfo.Id] = true
				}
				// 结算信息
				billplayerInfos = s.getBillPlayerInfos(deskPlayers, SInfo, pidScore)
				// 广播结算信息
				NotifyMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
					BillPlayersInfo: billplayerInfos,
				})
				if len(giveupPlayers) != 0 {
					// 广播认输信息
					NotifyMessage(desk, msgid.MsgID_ROOM_PLAYER_GIVEUP_NTF, &room.RoomGiveUpNtf{
						PlayerId: giveupPlayers,
					})
				}
				// 结算完生成事件
				s.generateSettleEvent(desk, giveupPlayers)
			}
		}
	}
	// 退税ids
	revertIds := mjContext.RevertSettles
	if len(revertIds) != 0 {
		// 退税的结算信息
		revertBillInfos := s.getRevertBillPlayerInfos(deskPlayers, revertIds)
		// 广播退税结算信息
		NotifyMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
			BillPlayersInfo: revertBillInfos,
		})
	}
}

// RoundSettle 单局结算
func (s *scxlSettle) RoundSettle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	// 牌局所有结算信息
	contextSInfos := mjContext.SettleInfos
	// 牌局玩家
	deskPlayers := desk.GetPlayers()
	for i := 0; i < len(deskPlayers); i++ {
		pid := deskPlayers[i].GetPlayerId()
		//记录该玩家单局结算信息
		balanceRsp := &room.RoomBalanceInfoRsp{
			Pid:             deskPlayers[i].PlayerId,
			BillDetail:      make([]*room.BillDetail, 0),
			BillPlayersInfo: make([]*room.BillPlayerInfo, 0),
		}
		// 记录该玩家单局结算总倍数
		cardValue := int32(0)
		// 记录该玩家退税信息
		revertScore := int64(0)
		revertSInfos := make([]*majongpb.SettleInfo, 0)

		// 遍历牌局所有结算信息，获取所有与该玩家有关的结算，获取结算详情列表
		for _, sInfo := range contextSInfos {
			if sInfo.Scores[pid] != 0 {
				bd := s.getBillDetail(pid, sInfo)
				cardValue = cardValue + bd.GetFanValue()
				balanceRsp.BillDetail = append(balanceRsp.BillDetail, bd)
			}
			// 退税结算详情
			revertIds := mjContext.RevertSettles
			if len(revertIds) != 0 {
				for _, rID := range revertIds {
					pScore := s.settleMap[rID][pid]
					if rID == sInfo.Id && pScore != 0 {
						revertSInfos = append(revertSInfos, sInfo)
						revertScore = revertScore + pScore
					}
				}
			}
		}
		// 获取退税结算详情
		if revertScore != 0 {
			revertbd := s.getRevertbillDetail(pid, revertScore, revertSInfos)
			balanceRsp.BillDetail = append(balanceRsp.BillDetail, revertbd)
		}
		// 获取玩家单局结算详情
		balanceRsp.BillPlayersInfo = s.getRoundBillPlayerInfo(pid, cardValue, mjContext)
		// 通知该玩家单局结算信息
		NotifyPlayersMessage(desk, []uint64{pid}, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp)
	}
}

// generateSettleEvent 生成结算finish事件
func (s *scxlSettle) generateSettleEvent(desk interfaces.Desk, giveupPlayers []uint64) {
	// 序列化
	settlefinish := &majongpb.SettleFinishEvent{
		PlayerId: giveupPlayers,
	}
	eventContext, err := proto.Marshal(settlefinish)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"msg": settlefinish,
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

// sumSettleInfo 合并相关联的一组SettleInfo的Score分数为一条settleInfo
// 返回参数:	[]*majongpb.SettleInfo(该组settleInfo) / *majongpb.SettleInfo(合并后的settleInfo)
func (s *scxlSettle) sumSettleInfo(contextSInfo []*majongpb.SettleInfo, settleInfo *majongpb.SettleInfo) ([]*majongpb.SettleInfo, *majongpb.SettleInfo) {
	sumSInfo := &majongpb.SettleInfo{
		Scores: make(map[uint64]int64, 0),
	}
	groupSInfos := make([]*majongpb.SettleInfo, 0)
	for _, id := range settleInfo.GroupId {
		sIndex := GetSettleInfoBySid(contextSInfo, id)
		groupSInfos = append(groupSInfos, contextSInfo[sIndex])
		sumSInfo.SettleType = contextSInfo[sIndex].SettleType
		s.handleSettle[id] = true
	}
	for _, singleSInfo := range groupSInfos {
		for pid, score := range singleSInfo.Scores {
			sumSInfo.Scores[pid] = sumSInfo.Scores[pid] + score
		}
	}
	return groupSInfos, sumSInfo
}

// calcScore 计算实际扣除的分数(根据玩家实际的金币数)
// 如果出现一炮多响的情况：
// 1.玩家身上的钱够赔付胡牌玩家的话,直接赔付
// 2.玩家身上的钱不够赔付胡牌玩家的话,那么该玩家身上的钱平分给胡牌玩家，,按逆时针方向,从点炮者数起,余 1 情况赔付于赢钱最多的玩家,
//	 余 2 情况赔付于第一、第二胡牌玩家;
func (s *scxlSettle) calcScore(deskPlayer []*room.RoomPlayerInfo, settleInfo *majongpb.SettleInfo) (map[uint64]int64, []uint64) {
	winScore := int64(0)
	loseScore := int64(0)
	losePids := make([]uint64, 0)
	winPid := make([]uint64, 0)
	realCost := make(map[uint64]int64, 0)
	lessCoinPid := make([]uint64, 0) // 记录金币不足够扣费的玩家id
	for pid, score := range settleInfo.Scores {
		if score > 0 {
			winScore = winScore + score
			winPid = append(winPid, pid)
		} else if score < 0 {
			loseScore = loseScore + score
			losePids = append(losePids, pid)
		}
	}
	if len(losePids) > 1 {
		for _, losePid := range losePids {
			losePlayer := GetDeskPlayer(deskPlayer, losePid)
			cost := int64(0)
			if s.abs(settleInfo.Scores[losePid]) < int64(losePlayer.GetCoin()) {
				cost = settleInfo.Scores[losePid]
			} else {
				lessCoinPid = append(lessCoinPid, losePid)
				cost = int64(0 - losePlayer.GetCoin())
			}
			realCost[losePid] = cost
			realCost[winPid[0]] = realCost[winPid[0]] - realCost[losePid]
		}
	} else if len(losePids) == 1 {
		losePid := losePids[0]
		losePlayer := GetDeskPlayer(deskPlayer, losePid)
		if s.abs(loseScore) < int64(losePlayer.GetCoin()) {
			for _, win := range winPid {
				realCost[win] = settleInfo.Scores[win]
			}
			realCost[losePid] = settleInfo.Scores[losePid]
		} else {
			lessCoinPid = append(lessCoinPid, losePid)
			loseCoin := int64(losePlayer.GetCoin())
			if len(winPid) == 1 {
				realCost[winPid[0]] = loseCoin
				realCost[losePid] = -loseCoin
			} else {
				maxWinPid := winPid[0]
				// 多个赢家，按照赢钱的比例平分
				for _, win := range winPid {
					rank := float64(settleInfo.Scores[win]) / float64(winScore)
					realCost[win] = int64(rank * float64(loseCoin))
					realCost[losePid] = realCost[losePid] - int64(rank*float64(loseCoin))
					if settleInfo.Scores[win] > settleInfo.Scores[maxWinPid] {
						maxWinPid = win
					}
				}
				//剩余分数，给赢钱最多的玩家
				surplusTotal := loseCoin - realCost[losePid]
				if surplusTotal > 0 {
					realCost[maxWinPid] = realCost[maxWinPid] + surplusTotal
					realCost[losePid] = realCost[losePid] - surplusTotal
				}
			}
		}
	}
	return realCost, lessCoinPid
}

// resolveScore 将合并settleInfo计算出的totalScore分配到单独settleIn中
func (s *scxlSettle) resolveScore(groupsInfos []*majongpb.SettleInfo, totalScore map[uint64]int64) {
	for _, sinfo := range groupsInfos {
		singleCost := make(map[uint64]int64, 0)
		cost := int64(0)
		for pid, score := range sinfo.Scores {
			if score > 0 {
				cost = totalScore[pid]
				singleCost[pid] = cost
			} else if score < 0 {
				if cost != 0 {
					singleCost[pid] = 0 - cost
				} else {
					for _, tscore := range totalScore {
						if tscore > 0 {
							singleCost[pid] = 0 - tscore
						}
					}
				}
			}
		}
		s.settleMap[sinfo.Id] = singleCost
	}
}

// getBillPlayerInfos 获得玩家结算账单
func (s *scxlSettle) getBillPlayerInfos(deskPlayers []*room.RoomPlayerInfo, settleInfo *majongpb.SettleInfo, realScore map[uint64]int64) (billplayerInfos []*room.BillPlayerInfo) {
	billplayerInfos = make([]*room.BillPlayerInfo, 0)
	for i := 0; i < len(deskPlayers); i++ {
		pid := deskPlayers[i].GetPlayerId()
		score := realScore[pid]
		if score != 0 {
			billplayerInfo := s.newBillplayerInfo(pid, s.settleType2BillType(settleInfo.SettleType))
			// 玩家当前分数
			coin := int64(deskPlayers[i].GetCoin())
			// 玩家结算后的分数
			deskPlayers[i].Coin = proto.Uint64(uint64(coin + score))
			// 生成玩家结算账单
			billplayerInfo.Score = proto.Int64(score)
			billplayerInfo.CurrentScore = proto.Int64(int64(deskPlayers[i].GetCoin()))
			billplayerInfos = append(billplayerInfos, billplayerInfo)
		}
		s.roundScore[pid] = s.roundScore[pid] + realScore[pid]
		// 设置玩家分数
		global.GetPlayerMgr().GetPlayer(pid).SetCoin(deskPlayers[i].GetCoin())
	}
	return
}

// getRevertBillPlayerInfos 获得玩家退税结算账单
func (s *scxlSettle) getRevertBillPlayerInfos(deskPlayers []*room.RoomPlayerInfo, revertIds []uint64) (billplayerInfos []*room.BillPlayerInfo) {
	billplayerInfos = make([]*room.BillPlayerInfo, 0)
	for i := 0; i < len(deskPlayers); i++ {
		pid := deskPlayers[i].GetPlayerId()
		coin := int64(deskPlayers[i].GetCoin())
		billplayerInfo := &room.BillPlayerInfo{
			Pid:      deskPlayers[i].PlayerId,
			BillType: room.BillType_BILL_REFUND.Enum(),
			Score:    proto.Int64(0),
		}
		for _, revertID := range revertIds {
			if score, ok := s.settleMap[revertID][pid]; ok && score != 0 {
				billplayerInfo.Score = proto.Int64(billplayerInfo.GetScore() - score)
				deskPlayers[i].Coin = proto.Uint64(uint64(int64(coin) - score))
			}
		}
		billplayerInfo.CurrentScore = proto.Int64(int64(deskPlayers[i].GetCoin()))
		billplayerInfos = append(billplayerInfos, billplayerInfo)
		// 设置玩家分数
		global.GetPlayerMgr().GetPlayer(pid).SetCoin(deskPlayers[i].GetCoin())
	}
	return
}

// getBillDetail 获得玩家单次结算详情，包括番型，分数，倍数，以及输赢玩家
func (s *scxlSettle) getBillDetail(pid uint64, sInfo *majongpb.SettleInfo) *room.BillDetail {
	billDetail := &room.BillDetail{
		SetleType: room.SettleType(sInfo.SettleType).Enum(),
		HuType:    room.HuType(sInfo.HuType).Enum(),
		FanValue:  proto.Int32(int32(sInfo.CardValue)),
		GenCount:  proto.Uint32(sInfo.GenCount),
		Score:     proto.Int64(s.settleMap[sInfo.Id][pid]),
	}
	// 实际扣除分数
	realScore := s.settleMap[sInfo.Id]
	fanTypes := make([]room.FanType, 0)
	for _, cardType := range sInfo.CardType {
		fanTypes = append(fanTypes, room.FanType(cardType))
	}
	billDetail.FanType = fanTypes
	if realScore[pid] < 0 { // 输家结算倍数为负数
		billDetail.FanValue = proto.Int32(int32(0 - sInfo.GetCardValue()))
	}
	if realScore[pid] > 0 { // 赢家结算所关联玩家为所有输家
		for pid, score := range realScore {
			if score < 0 {
				billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
			}
		}
	} else if realScore[pid] < 0 { // 输家结算所关联玩家为赢家
		for pid, score := range realScore {
			if score > 0 {
				billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
			}
		}
	}
	return billDetail
}

// getRevertbd 获得玩家退税结算详情，包括分数以及输赢玩家
func (s *scxlSettle) getRevertbillDetail(pid uint64, revertScore int64, revertSInfos []*majongpb.SettleInfo) *room.BillDetail {
	billDetail := &room.BillDetail{
		SetleType: room.SettleType_ST_TAXREBEAT.Enum(),
		Score:     proto.Int64(-revertScore),
	}
	// 相关联玩家
	for _, revertSInfo := range revertSInfos {
		// 实际扣除分数
		realScore := s.settleMap[revertSInfo.Id]
		if realScore[pid] > 0 { // 赢家结算所关联玩家为所有输家
			for pid, score := range realScore {
				if score < 0 {
					billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
				}
			}
		} else if realScore[pid] < 0 { // 输家结算所关联玩家为赢家
			for pid, score := range realScore {
				if score > 0 {
					billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
				}
			}
		}
	}
	return billDetail
}

// getRoundBillPlayerInfo 获得单局结算玩家详情,包括玩家自己牌型,输赢分数，以及其余每个玩家的输赢分数
func (s *scxlSettle) getRoundBillPlayerInfo(currentPid uint64, cardValue int32, context majongpb.MajongContext) []*room.BillPlayerInfo {
	billPlayerInfos := make([]*room.BillPlayerInfo, 0)
	for _, player := range context.Players {
		playerID := player.GetPalyerId()
		billPlayerInfo := &room.BillPlayerInfo{
			Pid:       proto.Uint64(playerID),
			Score:     proto.Int64(s.roundScore[playerID]),
			CardValue: proto.Int32(cardValue),
		}
		if playerID == currentPid {
			billPlayerInfo.CardsGroup = gutils.GetCardsGroup(player)
		}
		billPlayerInfos = append(billPlayerInfos, billPlayerInfo)
	}
	return billPlayerInfos
}

func (s *scxlSettle) newBillplayerInfo(playID uint64, billType room.BillType) *room.BillPlayerInfo {
	return &room.BillPlayerInfo{
		Pid:      proto.Uint64(playID),
		BillType: billType.Enum(),
	}
}

func (s *scxlSettle) abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func (s *scxlSettle) settleType2BillType(settleType majongpb.SettleType) room.BillType {
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
