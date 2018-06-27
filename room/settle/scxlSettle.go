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
	deskPlayers := desk.GetDeskPlayers()
	// 若存在未处理的结算信息，进行处理
	if len(contextSInfos) != 0 {
		for _, SInfo := range contextSInfos {
			if !s.handleSettle[SInfo.Id] {
				// 玩家结算信息
				billplayerInfos := make([]*room.BillPlayerInfo, 0)
				// 玩家输赢分数
				pidScore := make(map[uint64]int64, 0)
				// 破产的玩家id
				brokerPlayers := make([]uint64, 0)
				// 若存在相关联的一组SettleInfo(一炮多响情况)
				if len(SInfo.GroupId) > 1 {
					// 合并该组settleInfo,计算实际输赢分
					groupSInfos, sumSInfo := s.sumSettleInfo(mjContext.SettleInfos, SInfo)
					pidScore, brokerPlayers = s.calcCoin(deskPlayers, mjContext.GetPlayers(), sumSInfo.Scores)
					s.resolveScore(groupSInfos, pidScore)
				} else {
					// 单条settleInfo直接计算输赢分
					pidScore, brokerPlayers = s.calcCoin(deskPlayers, mjContext.GetPlayers(), SInfo.Scores)
					s.settleMap[SInfo.Id] = pidScore
					s.handleSettle[SInfo.Id] = true
				}
				// 结算信息
				billplayerInfos = s.getBillPlayerInfos(deskPlayers, SInfo, pidScore)
				// 广播结算信息
				NotifyMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
					BillPlayersInfo: billplayerInfos,
				})
				if len(brokerPlayers) != 0 {
					// 广播认输信息
					NotifyMessage(desk, msgid.MsgID_ROOM_PLAYER_GIVEUP_NTF, &room.RoomGiveUpNtf{
						PlayerId: brokerPlayers,
					})
				}
				// 结算完生成事件
				s.generateSettleEvent(desk, SInfo.SettleType, brokerPlayers)
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
	deskPlayers := desk.GetDeskPlayers()
	for i := 0; i < len(deskPlayers); i++ {
		if !deskPlayers[i].IsQuit() {
			pid := deskPlayers[i].GetPlayerID()
			//记录该玩家单局结算信息
			balanceRsp := &room.RoomBalanceInfoRsp{
				Pid:             proto.Uint64(pid),
				BillDetail:      make([]*room.BillDetail, 0),
				BillPlayersInfo: make([]*room.BillPlayerInfo, 0),
			}
			// 记录该玩家单局结算总倍数
			cardValue := int32(0)
			// 遍历牌局所有结算信息，获取所有与该玩家有关的结算，获取结算详情列表
			for _, sInfo := range contextSInfos {
				if sInfo.Scores[pid] != 0 {
					bd := s.getBillDetail(pid, sInfo)
					cardValue = cardValue + bd.GetFanValue()*int32(len(bd.GetRelatedPid()))
					balanceRsp.BillDetail = append(balanceRsp.BillDetail, bd)
				}
				// 退税结算详情
				if s.isRevertID(mjContext.RevertSettles, sInfo.Id) {
					pScore := s.settleMap[sInfo.Id][pid]
					if pScore != 0 {
						revertbd := s.getRevertbillDetail(pid, pScore, sInfo)
						cardValue = cardValue + revertbd.GetFanValue()*int32(len(revertbd.GetRelatedPid()))
						balanceRsp.BillDetail = append(balanceRsp.BillDetail, revertbd)
					}
				}

			}
			// 获取玩家单局结算详情
			balanceRsp.BillPlayersInfo = s.getRoundBillPlayerInfo(pid, cardValue, mjContext)
			// 通知该玩家单局结算信息
			NotifyPlayersMessage(desk, []uint64{pid}, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp)
		}
	}
}

// isRevertID 是否是退税id
func (s *scxlSettle) isRevertID(revertIds []uint64, settleID uint64) bool {
	for _, rID := range revertIds {
		if rID == settleID {
			return true
		}
	}
	return false
}

// generateSettleEvent 生成结算finish事件
func (s *scxlSettle) generateSettleEvent(desk interfaces.Desk, settleType majongpb.SettleType, brokerPlayers []uint64) {
	needEvent := map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_angang:   true,
		majongpb.SettleType_settle_bugang:   true,
		majongpb.SettleType_settle_minggang: true,
		majongpb.SettleType_settle_dianpao:  true,
		majongpb.SettleType_settle_zimo:     true,
	}
	if needEvent[settleType] {
		// 序列化
		settlefinish := &majongpb.SettleFinishEvent{
			PlayerId: brokerPlayers,
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

// calcMaxScore 计算玩家输赢上限
// 赢豆上限 = max(进房豆子数,当前豆子数)
func (s *scxlSettle) calcMaxScore(deskPlayer []interfaces.DeskPlayer, score map[uint64]int64) (maxScore map[uint64]int64) {
	maxScore = make(map[uint64]int64, 0)
	losePids := make([]uint64, 0)
	loseScore := int64(0)
	for pid, pscore := range score {
		if pscore > 0 {
			maxScore[pid] = s.getWinMax(GetDeskPlayer(deskPlayer, pid), pscore)
		} else if pscore < 0 {
			losePids = append(losePids, pid)
		}
	}
	if len(losePids) == 1 {
		for _, mscore := range maxScore {
			loseScore = loseScore - mscore
		}
		maxScore[losePids[0]] = loseScore
	} else {
		for _, mscore := range maxScore {
			loseScore = loseScore - mscore
		}
		for _, losePid := range losePids {
			maxScore[losePid] = loseScore / int64(len(losePids))
		}
	}
	return
}

func (s *scxlSettle) getWinMax(winPlayer interfaces.DeskPlayer, winScore int64) (winMax int64) {
	winMax = int64(0)
	winPid := winPlayer.GetPlayerID()
	currentCoin := int64(global.GetPlayerMgr().GetPlayer(winPid).GetCoin()) // 当前豆子数
	enterCoin := int64(winPlayer.GetEcoin())                                // 进房豆子数
	if currentCoin >= enterCoin {
		winMax = currentCoin
	} else {
		winMax = enterCoin
	}
	if winScore <= winMax {
		winMax = winScore
	}
	return
}

// calcCoin 计算扣除的金币
// 如果出现一炮多响的情况：
// 1.玩家身上的钱够赔付胡牌玩家的话,直接赔付
// 2.玩家身上的钱不够赔付胡牌玩家的话,那么该玩家身上的钱平分给胡牌玩家，,按逆时针方向,从点炮者数起,余 1 情况赔付于第一胡牌玩家,
//	 余 2 情况赔付于第一、第二胡牌玩家;
func (s *scxlSettle) calcCoin(deskPlayer []interfaces.DeskPlayer, contextPlayer []*majongpb.Player, score map[uint64]int64) (map[uint64]int64, []uint64) {
	maxScore := s.calcMaxScore(deskPlayer, score)

	winPlayers := make([]uint64, 0)  // 所有赢家
	losePlayers := make([]uint64, 0) // 所有输家
	tWinScore := int64(0)            // 赢的分数
	tLoseScore := int64(0)           // 输的分数
	for playerID, playerScore := range maxScore {
		if playerScore > 0 {
			tWinScore = tWinScore + playerScore
			winPlayers = append(winPlayers, playerID)
		} else if playerScore < 0 {
			tLoseScore = tLoseScore + playerScore
			losePlayers = append(losePlayers, playerID)
		}
	}
	coinCost := make(map[uint64]int64, 0) // 每个玩家实际扣除的金币数
	brokePlayers := make([]uint64, 0)     // 已破产玩家id

	if len(losePlayers) > 1 {
		winPlayer := winPlayers[0] // 赢家
		for _, losePid := range losePlayers {
			loseScore := s.abs(maxScore[losePid])                                 // 输家输分
			loseCoin := int64(global.GetPlayerMgr().GetPlayer(losePid).GetCoin()) // 输家金币数
			if loseScore < loseCoin {
				coinCost[losePid] = -loseScore
			} else {
				coinCost[losePid] = -loseCoin
				brokePlayers = append(brokePlayers, losePid)
			}
			coinCost[winPlayer] = coinCost[winPlayer] - coinCost[losePid]
		}
	} else if len(losePlayers) == 1 {
		losePid := losePlayers[0]
		loseScore := s.abs(tLoseScore)                                        // 输家输分
		loseCoin := int64(global.GetPlayerMgr().GetPlayer(losePid).GetCoin()) // 输家金币数
		if loseScore < loseCoin {
			for _, win := range winPlayers {
				coinCost[win] = maxScore[win]
			}
			coinCost[losePid] = maxScore[losePid]
		} else {
			if len(winPlayers) == 1 {
				coinCost[winPlayers[0]] = loseCoin
				coinCost[losePid] = -loseCoin
			} else {
				// 多个赢家，按照赢钱的比例平分
				for _, winPid := range winPlayers {
					winScore := float64(maxScore[winPid])
					rank := winScore / float64(tWinScore)
					coinCost[winPid] = int64(rank * float64(loseCoin))
					coinCost[losePid] = coinCost[losePid] - coinCost[winPid]
				}
				// 剩余分数，余 1 情况赔付于赢钱最多的玩家, 余 2 情况赔付于第一、第二胡牌玩家
				surplusScore := loseCoin - coinCost[losePid]
				loseIndex := gutils.GetPlayerIndex(losePid, contextPlayer)
				resortPlayers := make([]uint64, 0)
				for i := 0; i < len(contextPlayer); i++ {
					index := (loseIndex + i) % len(contextPlayer)
					resortPlayers = append(resortPlayers, contextPlayer[index].GetPalyerId())
				}
				resortHuPlayers := make([]uint64, 0)
				for _, resortPID := range resortPlayers {
					for _, winPlayer := range winPlayers {
						if resortPID == winPlayer {
							resortHuPlayers = append(resortHuPlayers, resortPID)
						}
					}
				}
				if surplusScore%2 == 0 {
					coinCost[resortHuPlayers[0]] = coinCost[resortHuPlayers[0]] + surplusScore/2
					coinCost[resortHuPlayers[1]] = coinCost[resortHuPlayers[1]] + surplusScore/2
					coinCost[losePid] = coinCost[losePid] - surplusScore
				} else {
					coinCost[resortHuPlayers[0]] = coinCost[resortHuPlayers[0]] + surplusScore
					coinCost[losePid] = coinCost[losePid] - surplusScore
				}
			}
			brokePlayers = append(brokePlayers, losePid)
		}
	}
	return coinCost, brokePlayers
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
func (s *scxlSettle) getBillPlayerInfos(deskPlayers []interfaces.DeskPlayer, settleInfo *majongpb.SettleInfo, costScore map[uint64]int64) (billplayerInfos []*room.BillPlayerInfo) {
	billplayerInfos = make([]*room.BillPlayerInfo, 0)
	for i := 0; i < len(deskPlayers); i++ {
		pid := deskPlayers[i].GetPlayerID()
		score := costScore[pid]
		if score != 0 {
			billplayerInfo := s.newBillplayerInfo(pid, s.settleType2BillType(settleInfo.SettleType))
			// 玩家当前豆子数
			currentCoin := int64(global.GetPlayerMgr().GetPlayer(pid).GetCoin())
			// 玩家结算后的分数
			calcCoin := uint64(currentCoin + score)
			// 生成玩家结算账单
			billplayerInfo.Score = proto.Int64(score)
			billplayerInfo.CurrentScore = proto.Int64(int64(calcCoin))
			billplayerInfos = append(billplayerInfos, billplayerInfo)
			s.roundScore[pid] = s.roundScore[pid] + score
			// 设置玩家分数
			global.GetPlayerMgr().GetPlayer(pid).SetCoin(calcCoin)
		}
	}
	return
}

// getRevertBillPlayerInfos 获得玩家退税结算账单
func (s *scxlSettle) getRevertBillPlayerInfos(deskPlayers []interfaces.DeskPlayer, revertIds []uint64) (billplayerInfos []*room.BillPlayerInfo) {
	billplayerInfos = make([]*room.BillPlayerInfo, 0)
	for i := 0; i < len(deskPlayers); i++ {
		pid := deskPlayers[i].GetPlayerID()
		coin := int64(global.GetPlayerMgr().GetPlayer(pid).GetCoin())
		billplayerInfo := &room.BillPlayerInfo{
			Pid:      proto.Uint64(pid),
			BillType: room.BillType_BILL_REFUND.Enum(),
			Score:    proto.Int64(0),
		}
		for _, revertID := range revertIds {
			if score, ok := s.settleMap[revertID][pid]; ok && score != 0 {
				billplayerInfo.Score = proto.Int64(billplayerInfo.GetScore() - score)
				coin = coin - score
			}
		}
		billplayerInfo.CurrentScore = proto.Int64(coin)
		billplayerInfos = append(billplayerInfos, billplayerInfo)
		// 设置玩家分数
		global.GetPlayerMgr().GetPlayer(pid).SetCoin(uint64(coin))
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
func (s *scxlSettle) getRevertbillDetail(pid uint64, revertScore int64, revertSettle *majongpb.SettleInfo) *room.BillDetail {
	billDetail := &room.BillDetail{
		SetleType: room.SettleType_ST_TAXREBEAT.Enum(),
		Score:     proto.Int64(-revertScore),
		FanValue:  proto.Int32(int32(revertSettle.CardValue)),
	}
	// 实际扣除分数
	realScore := s.settleMap[revertSettle.Id]
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
