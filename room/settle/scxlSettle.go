package settle

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/utils"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	majongpb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// scxlSettle 血流麻将结算
type scxlSettle struct {
	// 每条setttleInfo中每个玩家实际输赢分 key:settleId value:playerCoin
	settleMap map[uint64]playerCoin
	// 汇总setttleInfo中每个玩家输赢总分 key:playerID value:score
	roundScore map[uint64]int64
	// setttleInfo处理情况 		key:settleId value:true为已处理，false为未处理
	handleSettle map[uint64]bool
}

// newScxlSettle 创建四川血流结算
func newScxlSettle() *scxlSettle {
	return &scxlSettle{
		settleMap:    make(map[uint64]playerCoin),
		handleSettle: make(map[uint64]bool),
		roundScore:   make(map[uint64]int64),
	}
}

// playerCoin 玩家实际输赢分   key:playerID value:score
type playerCoin map[uint64]int64

// Settle 结算信息扣分并通知客户端
func (s *scxlSettle) Settle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	// 单局所有结算信息
	settleInfos := mjContext.SettleInfos
	// 牌局玩家
	deskPlayers := desk.GetPlayers()
	if len(settleInfos) != 0 {
		for _, settleInfo := range mjContext.SettleInfos {
			if !s.handleSettle[settleInfo.Id] {
				// 玩家结算信息
				billplayerInfos := make([]*room.BillPlayerInfo, 0)
				realScore := make(map[uint64]int64, 0)
				if len(settleInfo.GroupId) > 1 {
					// 合并一炮多响多条结算信息
					groupsInfos, combineSInfo := s.combineSettleInfo(mjContext.SettleInfos, settleInfo)
					realScore = s.calcScore(deskPlayers, combineSInfo)
					for _, sinfo := range groupsInfos {
						singleCost := make(map[uint64]int64, 0)
						cost := int64(0)
						for pid, score := range sinfo.Scores {
							if score > 0 {
								cost = realScore[pid]
								singleCost[pid] = realScore[pid]
							} else if score < 0 {
								singleCost[pid] = 0 - cost

							}
						}
						s.settleMap[sinfo.Id] = singleCost
					}
				} else {
					realScore = s.calcScore(deskPlayers, settleInfo)
					s.settleMap[settleInfo.Id] = realScore
					s.handleSettle[settleInfo.Id] = true
				}
				billplayerInfos = s.calcPlayerSettle(deskPlayers, settleInfo, realScore)
				// 广播即时结算消息
				notifyDeskMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
					BillPlayersInfo: billplayerInfos,
				})
			}
		}
	}
	// 退税
	revertIds := mjContext.RevertSettles
	if len(revertIds) != 0 {
		billplayerInfos := make([]*room.BillPlayerInfo, 0)
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
			// 设置玩家分数
			global.GetPlayerMgr().GetPlayer(pid).SetCoin(deskPlayers[i].GetCoin())
			billplayerInfo.CurrentScore = proto.Int64(int64(*deskPlayers[i].Coin))
			billplayerInfos = append(billplayerInfos, billplayerInfo)
		}
		// 即时结算消息
		notifyDeskMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
			BillPlayersInfo: billplayerInfos,
		})
	}
}

// combineSettleInfo 合并一炮多响的一组SettleInfo成一条
func (s *scxlSettle) combineSettleInfo(allSInfo []*majongpb.SettleInfo, settleInfo *majongpb.SettleInfo) ([]*majongpb.SettleInfo, *majongpb.SettleInfo) {
	combineSInfo := &majongpb.SettleInfo{
		Scores: make(map[uint64]int64, 0),
	}
	groupsInfos := make([]*majongpb.SettleInfo, 0)
	for _, id := range settleInfo.GroupId {
		index := settleInfoIndexByID(allSInfo, id)
		groupsInfos = append(groupsInfos, allSInfo[index])
		combineSInfo.SettleType = allSInfo[index].SettleType
		s.handleSettle[id] = true
	}
	for _, groupsInfo := range groupsInfos {
		for pid, score := range groupsInfo.Scores {
			combineSInfo.Scores[pid] = combineSInfo.Scores[pid] + score
		}
	}
	return groupsInfos, combineSInfo
}

// calcScore 计算分数
func (s *scxlSettle) calcScore(deskPlayer []*room.RoomPlayerInfo, settleInfo *majongpb.SettleInfo) map[uint64]int64 {
	winScore := int64(0)
	loseScore := int64(0)
	losePids := make([]uint64, 0)
	winPid := make([]uint64, 0)
	realCost := make(map[uint64]int64, 0)
	for pid, score := range settleInfo.Scores {
		p := pid
		if score > 0 {
			winScore = winScore + score
			winPid = append(winPid, p)
		} else if score <= 0 {
			loseScore = loseScore + score
			losePids = append(losePids, p)
		}
	}
	if len(losePids) > 1 {
		for _, losePid := range losePids {
			losePlayer := getDeskPlayer(deskPlayer, losePid)
			cost := int64(0)
			if abs(settleInfo.Scores[losePid]) <= int64(losePlayer.GetCoin()) {
				cost = settleInfo.Scores[losePid]
			} else {
				cost = int64(0 - losePlayer.GetCoin())
			}
			realCost[losePid] = realCost[losePid] + cost
			realCost[winPid[0]] = realCost[winPid[0]] - cost
		}
	} else {
		losePid := losePids[0]
		losePlayer := getDeskPlayer(deskPlayer, losePid)
		if abs(loseScore) < int64(losePlayer.GetCoin()) {
			for _, win := range winPid {
				realCost[win] = settleInfo.Scores[win]
			}
			realCost[losePid] = settleInfo.Scores[losePid]
		} else {
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
	return realCost
}
func (s *scxlSettle) calcPlayerSettle(deskPlayers []*room.RoomPlayerInfo, settleInfo *majongpb.SettleInfo, realScore map[uint64]int64) (billplayerInfos []*room.BillPlayerInfo) {
	billplayerInfos = make([]*room.BillPlayerInfo, 0)
	for i := 0; i < len(deskPlayers); i++ {
		pid := deskPlayers[i].GetPlayerId()
		score := realScore[pid]
		if score != 0 {
			billplayerInfo := newBillplayerInfo(pid, room.BillType(settleInfo.SettleType))
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

// RoundSettle 单局结算信息
func (s *scxlSettle) RoundSettle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	players := desk.GetPlayers()
	// 牌局所有settleInfo信息
	totalSInfos := mjContext.SettleInfos
	for i := 0; i < len(players); i++ {
		pid := players[i].GetPlayerId()
		// 玩家单局结算信息
		balanceRsp := &room.RoomBalanceInfoRsp{
			Pid:             players[i].PlayerId,
			BillDetail:      make([]*room.BillDetail, 0),
			BillPlayersInfo: make([]*room.BillPlayerInfo, 0),
		}
		// 玩家单局结算总倍数
		cardValue := int32(0)
		// 玩家退税SettleInfos
		revertIds := mjContext.RevertSettles
		revertSInfos := make([]*majongpb.SettleInfo, 0)
		// 玩家退税分数
		revertScore := int64(0)
		for _, sInfo := range totalSInfos {
			if sInfo.Scores[pid] != 0 {
				bd := s.createBillDetail(pid, sInfo)
				cardValue = cardValue + bd.GetFanValue()
				balanceRsp.BillDetail = append(balanceRsp.BillDetail, bd)
			}
			if len(revertIds) != 0 {
				for _, revertID := range revertIds {
					if revertID == sInfo.Id && s.settleMap[revertID][pid] != 0 {
						revertSInfos = append(revertSInfos, sInfo)
						revertScore = revertScore + s.settleMap[revertID][pid]
					}
				}
			}
		}
		if revertScore != 0 {
			revertbd := s.createRevertbd(pid, revertScore, revertSInfos)
			balanceRsp.BillDetail = append(balanceRsp.BillDetail, revertbd)
		}
		balanceRsp.BillPlayersInfo = s.createBillPInfo(pid, cardValue, mjContext)
		// 通知总结算
		notifyPlayerMessage(desk, pid, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp)
	}
}

// createBillDetail 生成玩家单次结算详情，包括番型，分数，倍数，以及输赢玩家
func (s *scxlSettle) createBillDetail(pid uint64, sInfo *majongpb.SettleInfo) *room.BillDetail {
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

// createRevertbd 生成玩家退税结算详情，包括分数以及输赢玩家
func (s *scxlSettle) createRevertbd(pid uint64, revertScore int64, revertSInfos []*majongpb.SettleInfo) *room.BillDetail {
	billDetail := &room.BillDetail{
		SetleType: room.SettleType_ST_TAXREBEAT.Enum(),
		Score:     proto.Int64(revertScore),
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

// createBillPInfo 生成单局结算玩家详情,包括玩家自己牌型,输赢分数，以及其余每个玩家的输赢分数
func (s *scxlSettle) createBillPInfo(currentPid uint64, cardValue int32, context majongpb.MajongContext) []*room.BillPlayerInfo {
	billPlayerInfos := make([]*room.BillPlayerInfo, 0)
	for _, player := range context.Players {
		playerID := player.GetPalyerId()
		billPlayerInfo := &room.BillPlayerInfo{
			Pid:       proto.Uint64(playerID),
			Score:     proto.Int64(s.roundScore[playerID]),
			CardValue: proto.Int32(cardValue),
		}
		if playerID == currentPid {
			billPlayerInfo.CardsGroup = utils.GetCardsGroup(player)
		}
		billPlayerInfos = append(billPlayerInfos, billPlayerInfo)
	}
	return billPlayerInfos
}

// settleInfoIndexByID 根据ettleID获取对应settleInfo的下标index
func settleInfoIndexByID(settleInfos []*majongpb.SettleInfo, ID uint64) int {
	for index, s := range settleInfos {
		if s.Id == ID {
			return index
		}
	}
	return -1
}

// calcCost 计算扣除的分数
func (s *scxlSettle) calcCost(deskPlayer *room.RoomPlayerInfo, settleInfo *majongpb.SettleInfo) int64 {
	pid := deskPlayer.GetPlayerId()
	score := settleInfo.Scores[pid]     // 输赢分数
	coin := int64(deskPlayer.GetCoin()) // 玩家剩余分数
	cost := int64(0)                    // 实际扣除分数
	if score != 0 {
		if abs(score) <= coin { // 剩余分数足够
			cost = score
		} else if score < 0 {
			cost = -coin
		}
	}
	return cost
}

func getDeskPlayer(deskPlayers []*room.RoomPlayerInfo, pid uint64) *room.RoomPlayerInfo {
	for _, p := range deskPlayers {
		if p.GetPlayerId() == pid {
			return p
		}
	}
	return nil
}

func newBillplayerInfo(playID uint64, billType room.BillType) *room.BillPlayerInfo {
	return &room.BillPlayerInfo{
		Pid:      proto.Uint64(playID),
		BillType: billType.Enum(),
	}
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func notifyDeskMessage(desk interfaces.Desk, msgid msgid.MsgID, message proto.Message) {
	players := desk.GetPlayers()
	clientIDs := []uint64{}

	playerMgr := global.GetPlayerMgr()
	for _, player := range players {
		playerID := player.GetPlayerId()
		p := playerMgr.GetPlayer(playerID)
		if p != nil {
			clientIDs = append(clientIDs, p.GetClientID())
		}
	}
	head := &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid)}
	ms := global.GetMessageSender()

	logrus.WithFields(logrus.Fields{
		"msg": message.String(),
	}).Debugln("通知立即结算")

	ms.BroadcastPackage(clientIDs, head, message)
}

func notifyPlayerMessage(desk interfaces.Desk, playerID uint64, msgid msgid.MsgID, message proto.Message) {
	clientID := global.GetPlayerMgr().GetPlayer(playerID).GetClientID()

	head := &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid)}
	ms := global.GetMessageSender()

	logrus.WithFields(logrus.Fields{
		"msg": message.String(),
	}).Debugln("通知总结算")
	ms.SendPackage(clientID, head, message)
}
