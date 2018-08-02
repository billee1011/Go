package models

import (
	"encoding/json"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/entity/gamelog"
	majongpb "steve/entity/majong"
	"steve/gutils"
	"steve/gutils/topics"
	"steve/room/contexts"
	"steve/room/desk"
	"steve/room/fixed"
	"steve/room/majong/utils"
	playerpkg "steve/room/player"
	"steve/room/util"
	"steve/structs"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// MajongCoin   key:playerID value:score
type MajongCoin map[uint64]int64

// MajongSettle 麻将结算
type MajongSettle struct {
	settleMap map[uint64]MajongCoin // setttleInfo实际扣分 key:结算id value:MajongCoin

	roundScore map[uint64]int64 // 每个玩家单局实际总扣分 key:玩家id value:分数

	handleSettle map[uint64]bool // setttleInfo扣分 key:结算id value:true为已扣分，false为未扣分

	handleRevert map[uint64]bool // 退税处理

	revertScore map[uint64]MajongCoin // revertScore  退稅分数 key:退税结算id value:MajongCoin

	lastGangSettleID uint64 // 呼叫转移
}

// NewMajongSettle 初始化麻将结算
func NewMajongSettle() *MajongSettle {
	return &MajongSettle{
		settleMap:    make(map[uint64]MajongCoin),
		handleSettle: make(map[uint64]bool),
		handleRevert: make(map[uint64]bool),
		roundScore:   make(map[uint64]int64),
		revertScore:  make(map[uint64]MajongCoin),
	}
}

// GetStatistics 获取统计信息
func (majongSettle *MajongSettle) GetStatistics() map[uint64]int64 {
	return majongSettle.roundScore
}

// Settle 单次结算
func (majongSettle *MajongSettle) Settle(desk *desk.Desk, config *desk.DeskConfig) {
	mjContext := config.Context.(*contexts.MajongDeskContext).MjContext

	settleOption := GetSettleOption(int(mjContext.GetGameId())) // 游戏结算玩法

	allSettleInfos := mjContext.SettleInfos // 结算信息

	modelMgr := GetModelManager()
	deskID := desk.GetUid()
	deskPlayers := modelMgr.GetPlayerModel(deskID).GetDeskPlayers()

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
			score, brokerPlayers = CalcCoin(deskPlayers, mjContext.GetPlayers(), huQuitPlayers, sInfo.Scores)
			majongSettle.settleMap[sInfo.Id] = score
			majongSettle.handleSettle[sInfo.Id] = true
		} else {
			groupSInfos, masterSInfo := MergeSettle(mjContext.SettleInfos, sInfo)
			score, brokerPlayers = CalcCoin(deskPlayers, mjContext.GetPlayers(), huQuitPlayers, masterSInfo.Scores)
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
				modelMgr.GetMessageModel(deskID).BroadCastDeskMessageExcept([]uint64{}, true, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
					BillPlayersInfo: billInfo,
				})
			}
			majongSettle.handleRevert[revertID] = true
		}
	}

}

// getHuSettleQuitPlayers  获取牌局已结算胡且退出的玩家
func (majongSettle *MajongSettle) getHuSettleQuitPlayers(dPlayers []*playerpkg.Player, mjContext majongpb.MajongContext, huPlayers []uint64) map[uint64]bool {
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

// handleCallDiver 处理呼叫转移
func (majongSettle *MajongSettle) handleCallDiver(lastGangSettleID uint64, sinfo *majongpb.SettleInfo, allSinfo []*majongpb.SettleInfo, mjContext majongpb.MajongContext) map[uint64]int64 {
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

func (majongSettle *MajongSettle) divideScore(dianPaoPlayer uint64, huPlayers []uint64, gangScore, winSum int64, callDiverScore map[uint64]int64, mjContext majongpb.MajongContext) map[uint64]int64 {
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

func (majongSettle *MajongSettle) generateRevertSettle2(revertID uint64, gangSettle *majongpb.SettleInfo, huQuitPlayers, giveUpPlayers map[uint64]bool, revertIds []uint64, mjContext majongpb.MajongContext) *majongpb.SettleInfo {
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

func (majongSettle *MajongSettle) calcTaxbetCoin(losePlayer uint64, winPlayers []uint64, score map[uint64]int64, contextPlayer []*majongpb.Player) (coinCost map[uint64]int64) {
	coinCost = make(map[uint64]int64, 0)
	loseCoin := int64(playerpkg.GetPlayerMgr().GetPlayer(losePlayer).GetCoin()) // 输家金币数
	loseScore := score[losePlayer]
	if abs(loseScore) < loseCoin {
		// 金币数够扣
		for _, win := range winPlayers {
			coinCost[win] = score[win]
		}
		coinCost[losePlayer] = score[losePlayer]
	} else {
		winSum := len(winPlayers)
		// 金币数不够扣，赢家为1时直接输家的金币全部给赢家，否则平分
		if winSum == 1 {
			coinCost[winPlayers[0]] = loseCoin
			coinCost[losePlayer] = -loseCoin
		} else if winSum > 1 {
			// 多个赢家，按照赢家人数平分
			for _, winPid := range winPlayers {
				winScore := int64(loseCoin / int64(winSum))
				coinCost[winPid] = winScore
				coinCost[losePlayer] = coinCost[losePlayer] - coinCost[winPid]
			}
			// 剩余分数，余 1 情况赔付于靠近的第一的玩家, 余 2 情况赔付于靠近第一、第二玩家
			surplusScore := loseCoin - abs(coinCost[losePlayer])
			if surplusScore > 0 {
				loseIndex := gutils.GetPlayerIndex(losePlayer, contextPlayer)
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
					coinCost[losePlayer] = coinCost[losePlayer] - surplusScore
				} else {
					coinCost[resortHuPlayers[0]] = coinCost[resortHuPlayers[0]] + surplusScore
					coinCost[losePlayer] = coinCost[losePlayer] - surplusScore
				}
			}
		}
	}
	return
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

// GenerateSettleEvent 结算finish事件
func GenerateSettleEvent(desks *desk.Desk, settleType majongpb.SettleType, brokerPlayers []uint64) {
	needEvent := map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_angang:   true,
		majongpb.SettleType_settle_bugang:   true,
		majongpb.SettleType_settle_minggang: true,
		majongpb.SettleType_settle_dianpao:  true,
		majongpb.SettleType_settle_zimo:     true,
	}
	if needEvent[settleType] {
		eventContext := &majongpb.SettleFinishEvent{
			PlayerId: brokerPlayers,
		}
		/*event := majongpb.AutoEvent{
			EventId:      majongpb.EventID_event_settle_finish,
			EventContext: eventContext,
		}*/

		/*interfaces.Event{
			ID:        event.GetEventId(),
			Context:   event.GetEventContext(),
			EventType: interfaces.NormalEvent,
			PlayerID:  0,
		}*/

		event := desk.NewDeskEvent(int(majongpb.EventID_event_settle_finish), fixed.NormalEvent, desks, desk.CreateEventParams(
			desks.GetConfig().Context.(*contexts.MajongDeskContext).StateNumber,
			eventContext,
			0,
		))
		GetMjEventModel(desks.GetUid()).PushEvent(event)
	}
}

// instantSettle 立即结算并扣费
func (majongSettle *MajongSettle) instantSettle(desk *desk.Desk, sInfo *majongpb.SettleInfo, score map[uint64]int64, brokerPlayers []uint64, giveUpPlayers map[uint64]bool) {
	modelMgr := GetModelManager()
	deskID := desk.GetUid()
	// 扣费并设置玩家金币数
	players := modelMgr.GetPlayerModel(deskID).GetDeskPlayers()
	majongSettle.chargeCoin(players, score)
	messageModel := modelMgr.GetMessageModel(deskID)
	// 广播结算
	messageModel.BroadCastDeskMessageExcept([]uint64{}, true, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
		BillPlayersInfo: majongSettle.getBillPlayerInfos(players, sInfo, score),
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
		messageModel.BroadCastDeskMessageExcept([]uint64{}, true, msgid.MsgID_ROOM_PLAYER_GIVEUP_NTF, &room.RoomGiveUpNtf{
			PlayerId: needSend,
		})
	}
}

// getBillPlayerInfos 获得玩家结算账单
func (majongSettle *MajongSettle) getBillPlayerInfos(deskPlayers []*playerpkg.Player, settleInfo *majongpb.SettleInfo, costScore map[uint64]int64) (billplayerInfos []*room.BillPlayerInfo) {
	billplayerInfos = make([]*room.BillPlayerInfo, 0)
	for i := 0; i < len(deskPlayers); i++ {
		pid := deskPlayers[i].GetPlayerID()
		holdCoin := playerpkg.GetPlayerMgr().GetPlayer(pid).GetCoin()
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

// apartScore2Settle  将score分配到各自settleInfo中
func (majongSettle *MajongSettle) apartScore2Settle(groupSettleInfos []*majongpb.SettleInfo, allScores map[uint64]int64) {
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

// getGiveupPlayers  获取认输的玩家id
func getGiveupPlayers(dPlayers []*playerpkg.Player, mjContext majongpb.MajongContext) map[uint64]bool {
	giveupPlayers := make(map[uint64]bool, 0)
	for _, cPlayer := range mjContext.Players {
		if cPlayer.GetXpState() == 2 {
			giveupPlayers[cPlayer.GetPalyerId()] = true
		}
	}
	return giveupPlayers
}

// getHuQuitPlayers  获取牌局胡牌且退出房间后的玩家id
func (majongSettle *MajongSettle) getHuQuitPlayers(dPlayers []*playerpkg.Player, mjContext majongpb.MajongContext) map[uint64]bool {
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

// RoundSettle 单局结算
func (majongSettle *MajongSettle) RoundSettle(desk *desk.Desk, config *desk.DeskConfig) {
	majongSettle.roundSettle(desk, config)
	majongSettle.gameLog(desk)
}

func (majongSettle *MajongSettle) roundSettle(desk *desk.Desk, config *desk.DeskConfig) {
	mjContext := config.Context.(*contexts.MajongDeskContext).MjContext
	// 牌局所有结算信息
	contextSInfos := mjContext.SettleInfos
	// 牌局玩家
	deskPlayers := GetModelManager().GetPlayerModel(desk.GetUid()).GetDeskPlayers()
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

func (majongSettle *MajongSettle) sendRounSettleMessage(contextSInfos []*majongpb.SettleInfo, desk *desk.Desk, mjContext majongpb.MajongContext) {
	deskID := desk.GetUid()
	modelMgr := GetModelManager()

	// 牌局玩家
	deskPlayers := modelMgr.GetPlayerModel(deskID).GetDeskPlayers()
	msgModel := modelMgr.GetMessageModel(deskID)

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
		totalValue := int32(0)
		needBillDetails := mjoption.GetSettleOption(int(mjContext.SettleOptionId)).NeedBillDetails
		if needBillDetails {
			balanceRsp.BillDetail, totalValue = majongSettle.makeBillDetails(pid, contextSInfos)
			balanceRsp.BillPlayersInfo = majongSettle.makeBillPlayerInfo(pid, totalValue, nil, mjContext)
		} else {
			// 一条结算记录
			if len(contextSInfos) != 1 {
				// 通知该玩家单局结算信息
				msgModel.BroadCastDeskMessageExcept([]uint64{}, true, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp)
				return
			}
			sinfo := contextSInfos[0]
			winers, _ := getWinners(sinfo.Scores)
			if len(winers) == 0 {
				return
			}
			cardOption := mjoption.GetCardTypeOption(int(mjContext.GetCardtypeOptionId()))
			fans := getFans(sinfo.CardType, sinfo.HuaCount, cardOption)
			billPlayersInfo := majongSettle.makeBillPlayerInfo(winers[0], int32(sinfo.CardValue), fans, mjContext)
			balanceRsp.BillPlayersInfo = append(balanceRsp.BillPlayersInfo, billPlayersInfo...)
		}
		// 通知该玩家单局结算信息
		msgModel.BroadCastDeskMessage([]uint64{pid}, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp, true)
	}
}

func makeFanType(fanTypes []int64, cardOption *mjoption.CardTypeOption) (fan []*room.Fan, totalValue int32) {
	fan = make([]*room.Fan, 0)
	totalValue = int32(0)
	for _, fanType := range fanTypes {
		rfan := &room.Fan{
			Name:  room.FanType(int32(fanType)).Enum(),
			Value: proto.Int32(int32(cardOption.Fantypes[int(fanType)].Score)),
			Type:  proto.Uint32(uint32(cardOption.Fantypes[int(fanType)].Type)),
		}
		totalValue = totalValue + int32(cardOption.Fantypes[int(fanType)].Score)
		fan = append(fan, rfan)
	}
	return
}

// makeBillPlayerInfo 获得单局结算玩家详情,包括玩家自己牌型,输赢分数，以及其余每个玩家的输赢分数
func (majongSettle *MajongSettle) makeBillPlayerInfo(currentPid uint64, cardValue int32, fans []*room.Fan, context majongpb.MajongContext) []*room.BillPlayerInfo {
	billPlayerInfos := make([]*room.BillPlayerInfo, 0)
	for _, player := range context.Players {
		playerID := player.GetPalyerId()
		roomPlayer := playerpkg.GetPlayerMgr().GetPlayer(playerID)

		coin := int64(roomPlayer.GetCoin())
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

func (majongSettle *MajongSettle) makeBillDetails(pid uint64, contextSInfos []*majongpb.SettleInfo) (billDetails []*room.BillDetail, totalValue int32) {
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

// getRevertbd 获得玩家退税结算详情，包括分数以及输赢玩家
func (majongSettle *MajongSettle) getRevertbillDetail(pid uint64, revertScore map[uint64]int64) *room.BillDetail {
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

// makeBillDetail 获得玩家单次结算详情，包括番型，分数，倍数，以及输赢玩家
func (majongSettle *MajongSettle) makeBillDetail(pid uint64, sInfo *majongpb.SettleInfo) *room.BillDetail {
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

func (majongSettle *MajongSettle) chargeCoin(players []*playerpkg.Player, payScore map[uint64]int64) {
	for _, player := range players {
		pid := player.GetPlayerID()
		// 玩家当前豆子数
		currentCoin := int64(player.GetCoin())
		// 扣费后豆子数
		realCoin := uint64(currentCoin + payScore[pid])
		// 设置玩家豆子数
		player.SetCoin(realCoin)
		// 记录玩家单局总输赢
		majongSettle.roundScore[pid] = majongSettle.roundScore[pid] + payScore[pid]
	}
}

// GetSettleOption 获取游戏的结算配置
func GetSettleOption(gameID int) *mjoption.SettleOption {
	return mjoption.GetSettleOption(mjoption.GetGameOptions(gameID).SettleOptionID)
}

// CanInstantSettle 能否立即结算
func CanInstantSettle(settleType majongpb.SettleType, settleOption *mjoption.SettleOption) bool {
	if IsGangSettle(settleType) {
		return settleOption.GangInstantSettle
	} else if IsHuSettle(settleType) {
		return settleOption.HuInstantSettle
	}
	return true
}

// IsHuSettle 是否是胡结算方式
func IsHuSettle(settleType majongpb.SettleType) bool {
	return map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_dianpao: true,
		majongpb.SettleType_settle_zimo:    true,
	}[settleType]
}

// IsGangSettle 是否是杠结算方式
func IsGangSettle(settleType majongpb.SettleType) bool {
	return map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_angang:   true,
		majongpb.SettleType_settle_bugang:   true,
		majongpb.SettleType_settle_minggang: true,
	}[settleType]
}

// CanRoundSettle 玩家是否可以单局结算
func CanRoundSettle(playerID uint64, huQuitPlayers map[uint64]bool, settleOption *mjoption.SettleOption) bool {
	if huQuitPlayers[playerID] {
		return settleOption.HuQuitPlayerSettle.HuQuitPlayerRoundSettle
	}
	return true
}

func (majongSettle *MajongSettle) gameLog(desk *desk.Desk) {
	summaryID := int64(util.GenUniqueID())
	majongSettle.genGameSummary(desk, summaryID)
	majongSettle.genGameDetail(desk, summaryID)
}

func (majongSettle *MajongSettle) genGameSummary(desk *desk.Desk, summaryID int64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "MajongSettle.genGameSummary",
		"player_id": desk.GetUid(),
	})

	gameSummary := gamelog.TGameSummary{
		Sumaryid: summaryID,
		Deskid:   int64(desk.GetUid()),
		Gameid:   desk.GetGameId(),
		// Levelid: todo,
		Playerids: desk.GetPlayerIds(),
		// Createtime: todo,
	}
	// scoreinfo and winners
	gameSummary.Scoreinfo, gameSummary.Winnerids = majongSettle.getScoreinfoWinners(desk)
	// RoundCurrency
	gameSummary.Roundcurrency = majongSettle.getRoundCurrency(desk.GetConfig())

	// 序列化
	data, err := json.Marshal(gameSummary)
	if err != nil {
		logEntry.WithError(err).Errorln("序列化失败")
	}
	publisher := structs.GetGlobalExposer().Publisher
	publisher.Publish(topics.GameSummaryRecord, data)
}

func (majongSettle *MajongSettle) genGameDetail(desk *desk.Desk, summaryID int64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "MajongSettle.genGameDetail",
		"player_id": desk.GetUid(),
	})
	roundScore := majongSettle.roundScore
	bigWinner := getBigWinner(roundScore)
	for _, playerID := range desk.GetPlayerIds() {
		gameDetail := gamelog.TGameDetail{
			Sumaryid: summaryID,
			Playerid: playerID,
			Deskid:   int64(desk.GetUid()),
			Gameid:   desk.GetGameId(),
			Amount:   roundScore[playerID],
			//Createtime: todo 含义
		}
		if playerID == bigWinner {
			gameDetail.Iswinner = 1
		}
		data, err := json.Marshal(gameDetail)
		if err != nil {
			logEntry.WithError(err).Errorln("序列化失败")
		}
		publisher := structs.GetGlobalExposer().Publisher
		publisher.Publish(topics.GameDetailRecord, data)
	}

}

func (majongSettle *MajongSettle) getRoundCurrency(config *desk.DeskConfig) (currencys []gamelog.RoundCurrency) {
	currencys = make([]gamelog.RoundCurrency, len(majongSettle.settleMap)+len(majongSettle.revertScore))

	id2typeMap := majongSettle.getSettleid2TypeMap(config)
	var index int
	for id, scoreMap := range majongSettle.settleMap {
		currencys[index].Settletype = int32(id2typeMap[id])
		for playerID, score := range scoreMap {
			detail := gamelog.SettleDetail{
				Playerid:  playerID,
				ChangeVal: score,
			}
			currencys[index].Settledetails = append(currencys[index].Settledetails, detail)
		}
		index++
	}

	for _, scoreMap := range majongSettle.revertScore {
		currencys[index].Settletype = int32(majongpb.SettleType_settle_taxrebeat)
		for playerID, score := range scoreMap {
			detail := gamelog.SettleDetail{
				Playerid:  playerID,
				ChangeVal: score,
			}
			currencys[index].Settledetails = append(currencys[index].Settledetails, detail)
		}
		index++
	}
	return
}

func (majongSettle *MajongSettle) getSettleid2TypeMap(config *desk.DeskConfig) (id2tyepMap map[uint64]majongpb.SettleType) {
	mjContext := config.Context.(*contexts.MajongDeskContext).MjContext
	contextSInfos := mjContext.SettleInfos
	id2tyepMap = make(map[uint64]majongpb.SettleType, len(contextSInfos))
	for _, info := range contextSInfos {
		id2tyepMap[info.GetId()] = info.GetSettleType()
	}
	return
}

func (majongSettle *MajongSettle) getScoreinfoWinners(desk *desk.Desk) (scoreInfo []int64, winners []uint64) {
	scoreInfo = make([]int64, len(desk.GetPlayerIds()))
	scoreStatistics := desk.GetConfig().Settle.GetStatistics()
	for index, playerID := range desk.GetPlayerIds() {
		if score, ok := scoreStatistics[playerID]; ok {
			scoreInfo[index] = score
			if score > 0 {
				winners = append(winners, playerID)
			}
		} else {
			scoreInfo[index] = 0
		}
	}
	return
}

func getBigWinner(roundScore map[uint64]int64) uint64 {
	var bigWinner uint64
	var maxScore int64
	for playerID, score := range roundScore {
		if score > maxScore {
			maxScore = score
			bigWinner = playerID
		}
	}
	return bigWinner
}
