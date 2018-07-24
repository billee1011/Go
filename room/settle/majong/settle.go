package majong

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/gutils"
	"steve/majong/utils"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	majongpb "steve/entity/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// majongCoin   key:playerID value:score
type majongCoin map[uint64]int64

// majongSettle 麻将结算
type majongSettle struct {
	settleMap map[uint64]majongCoin // setttleInfo实际扣分 key:结算id value:majongCoin

	roundScore map[uint64]int64 // 每个玩家单局实际总扣分 key:玩家id value:分数

	handleSettle map[uint64]bool // setttleInfo扣分 key:结算id value:true为已扣分，false为未扣分

	handleRevert map[uint64]bool // 退税处理

	revertScore map[uint64]majongCoin // revertScore  退稅分数 key:退税结算id value:majongCoin

	lastGangSettleID uint64 // 呼叫转移
}

// NewMajongSettle 初始化麻将结算
func NewMajongSettle() *majongSettle {
	return &majongSettle{
		settleMap:    make(map[uint64]majongCoin),
		handleSettle: make(map[uint64]bool),
		handleRevert: make(map[uint64]bool),
		roundScore:   make(map[uint64]int64),
		revertScore:  make(map[uint64]majongCoin),
	}
}

func (majongSettle *majongSettle) GetStatistics() map[uint64]int64 {
	statistics := make(map[uint64]int64, 4)
	for _, settleMap := range majongSettle.settleMap {
		for playerID, value := range settleMap {
			statistics[playerID] = statistics[playerID] + value
		}
	}
	return statistics
}

// Settle 单次结算
func (majongSettle *majongSettle) Settle(desk interfaces.Desk, mjContext majongpb.MajongContext) {

	settleOption := GetSettleOption(int(mjContext.GetGameId())) // 游戏结算玩法

	allSettleInfos := mjContext.SettleInfos // 结算信息

	deskPlayers := desk.GetDeskPlayers() // 牌局玩家

	giveUpPlayers := getGiveupPlayers(deskPlayers, mjContext) // 认输玩家

	revertIds := mjContext.RevertSettles   // 退税id
	for _, sInfo := range allSettleInfos { // 遍历
		if majongSettle.handleSettle[sInfo.Id] {
			continue
		}
		if IsGangSettle(sInfo.SettleType) {
			majongSettle.lastGangSettleID = sInfo.Id
		}
		if sInfo.SettleType == majongpb.SettleType_settle_calldiver {
			sInfo.Scores = majongSettle.handleCallDiver(majongSettle.lastGangSettleID, sInfo, allSettleInfos, mjContext)
		}
		score := make(map[uint64]int64, 0) // 玩家输赢分数

		brokerPlayers := make([]uint64, 0) // 破产的玩家id

		huQuitPlayers := majongSettle.getHuSettleQuitPlayers(deskPlayers, mjContext, sInfo.HuPlayers) // 胡牌且退出房间后的玩家

		groupID := len(sInfo.GroupId) // 关联的一组结算id
		if groupID <= 1 {
			score, brokerPlayers = calcCoin(deskPlayers, mjContext.GetPlayers(), huQuitPlayers, sInfo.Scores)
			majongSettle.settleMap[sInfo.Id] = score
			majongSettle.handleSettle[sInfo.Id] = true
		} else {
			groupSInfos, masterSInfo := mergeSettle(mjContext.SettleInfos, sInfo)
			score, brokerPlayers = calcCoin(deskPlayers, mjContext.GetPlayers(), huQuitPlayers, masterSInfo.Scores)
			majongSettle.apartScore2Settle(groupSInfos, score)
		}
		if CanInstantSettle(sInfo.SettleType, settleOption) { // 立即结算
			majongSettle.instantSettle(desk, sInfo, score, brokerPlayers, giveUpPlayers)
		}
		// 生成结算完成事件
		GenerateSettleEvent(desk, sInfo.SettleType, brokerPlayers)
	}
	if len(revertIds) != 0 {
		for _, revertID := range revertIds {
			if majongSettle.handleRevert[revertID] {
				continue
			}
			huQuitPlayers := majongSettle.getHuQuitPlayers(deskPlayers, mjContext) // 胡牌且退出房间后的玩家
			// 退稅结算信息
			gangSettle := GetSettleInfoByID(allSettleInfos, revertID)
			rSettleInfo := majongSettle.generateRevertSettle2(revertID, gangSettle, huQuitPlayers, giveUpPlayers, revertIds, mjContext)
			if rSettleInfo != nil {
				// 扣费并设置玩家金币数
				majongSettle.chargeCoin(deskPlayers, rSettleInfo.Scores)
				billInfo := majongSettle.getBillPlayerInfos(deskPlayers, rSettleInfo, rSettleInfo.Scores)
				facade.BroadCastDeskMessageExcept(desk, []uint64{}, true, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
					BillPlayersInfo: billInfo,
				})
			}
			majongSettle.handleRevert[revertID] = true
		}
	}

}

// handleCallDiver 处理呼叫转移
func (majongSettle *majongSettle) handleCallDiver(lastGangSettleID uint64, sinfo *majongpb.SettleInfo, allSinfo []*majongpb.SettleInfo, mjContext majongpb.MajongContext) map[uint64]int64 {
	gangSettle := GetSettleInfoByID(allSinfo, lastGangSettleID) // 杠的结算信息

	_, gangWinScore := getWinners(majongSettle.settleMap[lastGangSettleID]) // 杠实际赢的钱

	dianGangPlayer, _ := getLosers(gangSettle.Scores) // 点杠者

	dianPaoPlayer, _ := getLosers(sinfo.Scores) // 点炮者

	huPlayers, _ := getWinners(sinfo.Scores) // 赢家

	winSum := int64(len(huPlayers))

	callDiverScore := make(map[uint64]int64, 0)

	if winSum == 1 {
		callDiverScore[huPlayers[0]] = gangWinScore
		callDiverScore[dianPaoPlayer[0]] = -gangWinScore
	} else {
		// 一炮多响
		if gangSettle.SettleType == majongpb.SettleType_settle_minggang {
			contain := false
			for _, huPlayerID := range huPlayers {
				if dianGangPlayer[0] != huPlayerID {
					continue
				}
				contain = true
				break
			}
			if contain {
				callDiverScore[dianGangPlayer[0]] = gangWinScore
				callDiverScore[dianPaoPlayer[0]] = -gangWinScore
			} else {
				// 平分
				callDiverScore = majongSettle.divideScore(dianPaoPlayer[0], huPlayers, gangWinScore, winSum, callDiverScore, mjContext)
			}
		} else if gangSettle.SettleType == majongpb.SettleType_settle_angang || gangSettle.SettleType == majongpb.SettleType_settle_bugang {
			// （暗杠、补杠）先收杠钱,平分,杠钱后还有多余，多余的杠钱按位置给第一个胡牌玩家
			// 平分
			callDiverScore = majongSettle.divideScore(dianPaoPlayer[0], huPlayers, gangWinScore, winSum, callDiverScore, mjContext)
		}
	}
	return callDiverScore
}

func (majongSettle *majongSettle) divideScore(dianPaoPlayer uint64, huPlayers []uint64, gangScore, winSum int64, callDiverScore map[uint64]int64, mjContext majongpb.MajongContext) map[uint64]int64 {
	// 平分
	equallyTotal := gangScore / winSum
	// 剩余分数
	surplusTotal := gangScore - (equallyTotal * int64(winSum))
	// 所有玩家
	allPlayers := make([]uint64, 0)

	for _, player := range mjContext.GetPlayers() {
		allPlayers = append(allPlayers, player.GetPalyerId())
	}
	for _, huPlayerID := range huPlayers {
		callDiverScore[huPlayerID] = equallyTotal
		callDiverScore[dianPaoPlayer] = callDiverScore[dianPaoPlayer] - equallyTotal
	}
	if surplusTotal != 0 {
		startIndex, _ := utils.GetPlayerIDIndex(dianPaoPlayer, allPlayers)
		firstPlayerID := utils.GetPalyerCloseFromTarget(startIndex, allPlayers, huPlayers)
		if firstPlayerID != 0 {
			callDiverScore[firstPlayerID] = callDiverScore[firstPlayerID] + surplusTotal
			callDiverScore[dianPaoPlayer] = callDiverScore[dianPaoPlayer] - surplusTotal
		}
	}
	return callDiverScore
}

// RoundSettle 单局结算
func (majongSettle *majongSettle) RoundSettle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	// 牌局所有结算信息
	contextSInfos := mjContext.SettleInfos
	// 牌局玩家
	deskPlayers := desk.GetDeskPlayers()
	// 游戏结算玩法
	settleOption := GetSettleOption(int(mjContext.GetGameId()))

	for _, sInfo := range contextSInfos {
		if !CanInstantSettle(sInfo.SettleType, settleOption) {
			// 扣费并设置玩家金币数
			majongSettle.chargeCoin(deskPlayers, majongSettle.settleMap[sInfo.Id])
		}
	}
	majongSettle.sendRounSettleMessage(contextSInfos, desk, mjContext)
}

func (majongSettle *majongSettle) sendRounSettleMessage(contextSInfos []*majongpb.SettleInfo, desk interfaces.Desk, mjContext majongpb.MajongContext) {
	// 牌局玩家
	deskPlayers := desk.GetDeskPlayers()
	// 结算配置
	settleOption := mjoption.GetSettleOption(int(mjContext.SettleOptionId))
	// 番型配置
	cardOption := mjoption.GetCardTypeOption(int(mjContext.GetCardtypeOptionId()))

	for i := 0; i < len(deskPlayers); i++ {
		if deskPlayers[i].IsQuit() {
			continue
		}
		pid := deskPlayers[i].GetPlayerID()
		//记录该玩家单局结算信息
		balanceRsp := &room.RoomBalanceInfoRsp{
			Pid:             proto.Uint64(pid),
			BillDetail:      make([]*room.BillDetail, 0),
			BillPlayersInfo: make([]*room.BillPlayerInfo, 0),
		}
		if settleOption.NeedBillDetails {
			billDetail, totalValue := majongSettle.makeBillDetails(pid, contextSInfos)
			billPlayersInfo := majongSettle.makeBillPlayerInfo(pid, totalValue, nil, mjContext)

			balanceRsp.BillDetail = append(balanceRsp.BillDetail, billDetail...)
			balanceRsp.BillPlayersInfo = append(balanceRsp.BillPlayersInfo, billPlayersInfo...)
		} else {
			// 一条结算记录
			if len(contextSInfos) != 1 {
				// 通知该玩家单局结算信息
				facade.BroadCastDeskMessageExcept(desk, []uint64{}, true, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp)
				return
			}
			sinfo := contextSInfos[0]
			winers, _ := getWinners(sinfo.Scores)
			if len(winers) == 0 {
				return
			}
			fans := getFans(sinfo.CardType, sinfo.HuaCount, cardOption)
			billPlayersInfo := majongSettle.makeBillPlayerInfo(winers[0], int32(sinfo.CardValue), fans, mjContext)
			balanceRsp.BillPlayersInfo = append(balanceRsp.BillPlayersInfo, billPlayersInfo...)
		}
		// 通知该玩家单局结算信息
		facade.BroadCastDeskMessage(desk, []uint64{pid}, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp, true)
		logrus.WithFields(logrus.Fields{
			"func_name":  "sendRounSettleMessage",
			"pid":        pid,
			"balanceRsp": balanceRsp,
		}).Debugln("通知玩家单局结算信息")
	}
}

// instantSettle 立即结算并扣费
func (majongSettle *majongSettle) instantSettle(desk interfaces.Desk, sInfo *majongpb.SettleInfo, score map[uint64]int64, brokerPlayers []uint64, giveUpPlayers map[uint64]bool) {
	// 扣费并设置玩家金币数
	majongSettle.chargeCoin(desk.GetDeskPlayers(), score)
	// 广播结算
	billInfo := majongSettle.getBillPlayerInfos(desk.GetDeskPlayers(), sInfo, score)
	facade.BroadCastDeskMessageExcept(desk, []uint64{}, true, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
		BillPlayersInfo: billInfo,
	})
	needSend := make([]uint64, 0)
	for _, brokerPlayer := range brokerPlayers {
		if !giveUpPlayers[brokerPlayer] {
			needSend = append(needSend, brokerPlayer)
		}
	}
	// 查花猪、查大叫、退税阶段不需要发送认输
	notNeedSend := map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_yell:      true,
		majongpb.SettleType_settle_flowerpig: true,
		majongpb.SettleType_settle_taxrebeat: true,
	}
	if !notNeedSend[sInfo.SettleType] {
		// 广播认输
		facade.BroadCastDeskMessageExcept(desk, []uint64{}, true, msgid.MsgID_ROOM_PLAYER_GIVEUP_NTF, &room.RoomGiveUpNtf{
			PlayerId: needSend,
		})
	}
}

func (majongSettle *majongSettle) makeBillDetails(pid uint64, contextSInfos []*majongpb.SettleInfo) (billDetails []*room.BillDetail, totalValue int32) {
	// 记录该玩家单局结算总倍数
	totalValue = int32(0)

	billDetails = make([]*room.BillDetail, 0)
	// 遍历牌局所有结算信息，获取所有与该玩家有关的结算，获取结算详情列表
	for _, sInfo := range contextSInfos {
		if sInfo.Scores[pid] != 0 {
			billdetail := majongSettle.makeBillDetail(pid, sInfo)
			billDetails = append(billDetails, billdetail)
			if billdetail.GetScore() > 0 {
				billValue := billdetail.GetFanValue() * int32(len(billdetail.GetRelatedPid()))
				totalValue = totalValue + billValue
			} else if billdetail.GetScore() < 0 {
				totalValue = totalValue + billdetail.GetFanValue()
			}
		}
	}
	// 获取退税结算详情
	for _, sInfo := range contextSInfos {
		for rID, rScore := range majongSettle.revertScore {
			if rID == sInfo.Id && rScore[pid] != 0 {
				revertbd := majongSettle.getRevertbillDetail(pid, rScore)
				billDetails = append(billDetails, revertbd)
			}
		}
	}
	return
}

// getHuSettleQuitPlayers  获取牌局已结算胡且退出的玩家
func (majongSettle *majongSettle) getHuSettleQuitPlayers(dPlayers []interfaces.DeskPlayer, mjContext majongpb.MajongContext, huPlayers []uint64) map[uint64]bool {
	huQuitPids := make(map[uint64]bool, 0)
	huPids := make(map[uint64]bool, 0)
	for _, hplayer := range huPlayers {
		huPids[hplayer] = true
	}
	for _, dPlayer := range dPlayers {
		if dPlayer.IsQuit() && huPids[dPlayer.GetPlayerID()] {
			huQuitPids[dPlayer.GetPlayerID()] = true
		}
	}

	return huQuitPids
}

// getHuQuitPlayers  获取牌局已胡牌且退出的玩家
func (majongSettle *majongSettle) getHuQuitPlayers(dPlayers []interfaces.DeskPlayer, mjContext majongpb.MajongContext) map[uint64]bool {
	huPids := make(map[uint64]bool, 0)
	for _, contextPlayer := range mjContext.GetPlayers() {
		huCard := contextPlayer.GetHuCards()
		if len(huCard) != 0 {
			huPids[contextPlayer.GetPalyerId()] = true
		}
	}

	huQuitPids := make(map[uint64]bool, 0)
	for _, dPlayer := range dPlayers {
		pid := dPlayer.GetPlayerID()
		if dPlayer.IsQuit() && huPids[pid] {
			huQuitPids[pid] = true
		}
	}

	return huQuitPids
}

// apartScore2Settle  将score分配到各自settleInfo中
func (majongSettle *majongSettle) apartScore2Settle(groupSettleInfos []*majongpb.SettleInfo, allScores map[uint64]int64) {
	for _, sInfo := range groupSettleInfos {
		sID := sInfo.Id
		cost := int64(0)
		majongSettle.settleMap[sID] = make(map[uint64]int64)
		losePid := uint64(0)
		for pid, score := range sInfo.Scores {
			if score == 0 {
				continue
			} else if score > 0 {
				cost = allScores[pid]
				majongSettle.settleMap[sID][pid] = allScores[pid]
			} else if score < 0 {
				losePid = pid
			}
		}
		if cost != 0 {
			majongSettle.settleMap[sID][losePid] = 0 - cost
		}
		majongSettle.handleSettle[sID] = true
	}
}

func (majongSettle *majongSettle) chargeCoin(deskPlayers []interfaces.DeskPlayer, payScore map[uint64]int64) {
	for _, deskPlayer := range deskPlayers {
		pid := deskPlayer.GetPlayerID()
		// 玩家当前豆子数
		currentCoin := int64(global.GetPlayerMgr().GetPlayer(pid).GetCoin())
		// 扣费后豆子数
		realCoin := uint64(currentCoin + payScore[pid])
		// 设置玩家豆子数
		global.GetPlayerMgr().GetPlayer(pid).SetCoin(realCoin)
		// 记录玩家单局总输赢
		majongSettle.roundScore[pid] = majongSettle.roundScore[pid] + payScore[pid]
	}
}

// getBillPlayerInfos 获得玩家结算账单
func (majongSettle *majongSettle) getBillPlayerInfos(deskPlayers []interfaces.DeskPlayer, settleInfo *majongpb.SettleInfo, costScore map[uint64]int64) (billplayerInfos []*room.BillPlayerInfo) {
	billplayerInfos = make([]*room.BillPlayerInfo, 0)
	for i := 0; i < len(deskPlayers); i++ {
		pid := deskPlayers[i].GetPlayerID()
		holdCoin := global.GetPlayerMgr().GetPlayer(pid).GetCoin()
		if costScore[pid] == 0 {
			continue
		}
		billplayerInfos = append(billplayerInfos, &room.BillPlayerInfo{
			Pid:          proto.Uint64(pid),
			BillType:     settleType2BillType(settleInfo.SettleType).Enum(),
			Score:        proto.Int64(costScore[pid]),
			CurrentScore: proto.Int64(int64(holdCoin)),
		})
	}
	return
}

func (majongSettle *majongSettle) generateRevertSettle2(revertID uint64, gangSettle *majongpb.SettleInfo, huQuitPlayers, giveUpPlayers map[uint64]bool, revertIds []uint64, mjContext majongpb.MajongContext) *majongpb.SettleInfo {
	// 扣除的豆子数
	coinCost := make(map[uint64]int64, 0)
	// 扣除的分数
	scoreCost := make(map[uint64]int64, 0)
	// 退钱的玩家
	rlosePid := uint64(0)
	// 赢钱的玩家
	rWinnerPids := make([]uint64, 0)
	// 退的钱
	rloseScore := int64(0)
	for pid, score := range gangSettle.Scores {
		if score > 0 {
			if huQuitPlayers[pid] || giveUpPlayers[pid] { // 胡牌玩家已退出/认输玩家，不用退税
				return nil
			}
			rlosePid = pid
		}
	}
	for pid, score := range majongSettle.settleMap[revertID] {
		if huQuitPlayers[pid] || giveUpPlayers[pid] { // 胡牌玩家已退出/认输玩家，不用退税
			continue
		}
		if score < 0 {
			scoreCost[pid] = scoreCost[pid] - score
			rloseScore = rloseScore + score
			rWinnerPids = append(rWinnerPids, pid)
		}
	}
	scoreCost[rlosePid] = scoreCost[rlosePid] + rloseScore
	coinCost = calcTaxbetCoin(rlosePid, rWinnerPids, scoreCost, mjContext.GetPlayers())
	majongSettle.revertScore[revertID] = coinCost

	return &majongpb.SettleInfo{
		Scores:     coinCost,
		SettleType: majongpb.SettleType_settle_taxrebeat,
	}
}

// getBillDetail 获得玩家单次结算详情，包括番型，分数，倍数，以及输赢玩家
func (majongSettle *majongSettle) makeBillDetail(pid uint64, sInfo *majongpb.SettleInfo) *room.BillDetail {
	billDetail := &room.BillDetail{
		SetleType: room.SettleType(sInfo.SettleType).Enum(),
		HuType:    room.HuType(sInfo.HuType).Enum(),
		FanType:   make([]room.FanType, 0),
		FanValue:  proto.Int32(int32(sInfo.CardValue)),
		GenCount:  proto.Uint32(sInfo.GenCount),
		Score:     proto.Int64(majongSettle.settleMap[sInfo.Id][pid]),
	}
	realScore := majongSettle.settleMap[sInfo.Id] // 实际扣除分数
	for _, cardType := range sInfo.CardType {
		billDetail.FanType = append(billDetail.FanType, room.FanType(cardType))
	}
	if realScore[pid] < 0 { // 输家结算倍数为负数
		billDetail.FanValue = proto.Int32(int32(0 - sInfo.GetCardValue()))
	}
	winnerIds := make([]uint64, 0)
	loseIds := make([]uint64, 0)
	for pid, score := range realScore {
		if score < 0 {
			loseIds = append(loseIds, pid)
		}
		if score > 0 {
			winnerIds = append(winnerIds, pid)
		}
	}
	if realScore[pid] > 0 { // 赢家结算所关联玩家为所有输家
		billDetail.RelatedPid = loseIds
	} else if realScore[pid] < 0 { // 输家结算所关联玩家为赢家
		billDetail.RelatedPid = winnerIds
	}
	return billDetail
}

// getRevertbd 获得玩家退税结算详情，包括分数以及输赢玩家
func (majongSettle *majongSettle) getRevertbillDetail(pid uint64, revertScore map[uint64]int64) *room.BillDetail {
	billDetail := &room.BillDetail{
		SetleType: room.SettleType_ST_TAXREBEAT.Enum(),
		Score:     proto.Int64(revertScore[pid]),
	}

	if revertScore[pid] > 0 { // 赢家结算所关联玩家为所有输家
		for pid, score := range revertScore {
			if score < 0 {
				billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
			}
		}
	} else if revertScore[pid] < 0 { // 输家结算所关联玩家为赢家
		for pid, score := range revertScore {
			if score > 0 {
				billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
			}
		}
	}
	return billDetail
}

// makeBillPlayerInfo 获得单局结算玩家详情,包括玩家自己牌型,输赢分数，以及其余每个玩家的输赢分数
func (majongSettle *majongSettle) makeBillPlayerInfo(currentPid uint64, cardValue int32, fans []*room.Fan, context majongpb.MajongContext) []*room.BillPlayerInfo {
	billPlayerInfos := make([]*room.BillPlayerInfo, 0)
	for _, player := range context.Players {
		playerID := player.GetPalyerId()
		coin := int64(global.GetPlayerMgr().GetPlayer(playerID).GetCoin())
		billPlayerInfo := &room.BillPlayerInfo{
			Pid:          proto.Uint64(playerID),
			Score:        proto.Int64(majongSettle.roundScore[playerID]),
			CardValue:    proto.Int32(cardValue),
			BillType:     room.BillType(-1).Enum(),
			Fan:          fans,
			CurrentScore: proto.Int64(coin),
		}
		if len(player.CardsGroup) != 0 {
			billPlayerInfo.CardsGroup = gutils.CardsGroupSvr2Client(player.CardsGroup)
		} else if playerID == currentPid {
			billPlayerInfo.CardsGroup = gutils.GetCardsGroup(player)
		}

		billPlayerInfos = append(billPlayerInfos, billPlayerInfo)
	}
	return billPlayerInfos
}
