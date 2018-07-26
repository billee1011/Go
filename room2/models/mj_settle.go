package models

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	majongpb "steve/entity/majong"
	"steve/gutils"
	"steve/room2/contexts"
	"steve/room2/desk"
	"steve/room2/fixed"
	playerpkg "steve/room2/player"

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

// Settle 单次结算
func (majongSettle *MajongSettle) Settle(desk *desk.Desk, config *desk.DeskConfig) {
	mjContext := config.Context.(*contexts.MjContext).MjContext

	settleOption := GetSettleOption(int(desk.GetGameId())) // 游戏结算玩法

	allSettleInfos := mjContext.SettleInfos // 结算信息

	deskPlayers := GetModelManager().GetPlayerModel(desk.GetUid()).GetDeskPlayers() // 牌局玩家

	huQuitPlayers := majongSettle.getHuQuitPlayers(deskPlayers, mjContext) // 胡牌且退出房间后的玩家

	giveUpPlayers := getGiveupPlayers(deskPlayers, mjContext) // 认输玩家

	revertIds := mjContext.RevertSettles   // 退税id
	for _, sInfo := range allSettleInfos { // 遍历
		if majongSettle.handleSettle[sInfo.Id] {
			continue
		}
		score := make(map[uint64]int64, 0) // 玩家输赢分数

		brokerPlayers := make([]uint64, 0) // 破产的玩家id

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
			// 退稅结算信息
			gangSettle := GetSettleInfoByID(allSettleInfos, revertID)
			rSettleInfo := majongSettle.generateRevertSettle2(revertID, gangSettle, huQuitPlayers, giveUpPlayers, revertIds, mjContext)
			if rSettleInfo != nil {
				// 扣费并设置玩家金币数
				majongSettle.chargeCoin(deskPlayers, rSettleInfo.Scores)
				GetModelManager().GetMessageModel(desk.GetUid()).BroadCastDeskMessageExcept([]uint64{}, true, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
					BillPlayersInfo: majongSettle.getBillPlayerInfos(deskPlayers, rSettleInfo, rSettleInfo.Scores),
				})
			}
			majongSettle.handleRevert[revertID] = true
		}
	}

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
	coinCost = majongSettle.calcTaxbetCoin(rlosePid, rWinnerPids, scoreCost, mjContext.GetPlayers())
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
			desks.GetConfig().Context.(*contexts.MjContext).StateNumber,
			eventContext,
			0,
		))
		GetModelManager().GetMjEventModel(desks.GetUid()).PushEvent(event)
	}
}

// instantSettle 立即结算并扣费
func (majongSettle *MajongSettle) instantSettle(desk *desk.Desk, sInfo *majongpb.SettleInfo, score map[uint64]int64, brokerPlayers []uint64, giveUpPlayers map[uint64]bool) {
	// 扣费并设置玩家金币数
	players := GetModelManager().GetPlayerModel(desk.GetUid()).GetDeskPlayers()
	majongSettle.chargeCoin(players, score)
	// 广播结算
	GetModelManager().GetMessageModel(desk.GetUid()).BroadCastDeskMessageExcept([]uint64{}, true, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
		BillPlayersInfo: majongSettle.getBillPlayerInfos(players, sInfo, score),
	})
	needSend := make([]uint64, 0)
	for _, brokerPlayer := range brokerPlayers {
		if !giveUpPlayers[brokerPlayer] {
			needSend = append(needSend, brokerPlayer)
		}
	}
	// 广播认输
	GetModelManager().GetMessageModel(desk.GetUid()).BroadCastDeskMessageExcept([]uint64{}, true, msgid.MsgID_ROOM_PLAYER_GIVEUP_NTF, &room.RoomGiveUpNtf{
		PlayerId: needSend,
	})
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
	huQuitPids := make(map[uint64]bool, 0)
	for _, dPlayer := range dPlayers {
		if dPlayer.IsQuit() {
			pid := dPlayer.GetPlayerID()
			mjPlayers := mjContext.GetPlayers()
			mjPlayer := mjPlayers[gutils.GetPlayerIndex(pid, mjContext.GetPlayers())]
			if len(mjPlayer.HuCards) != 0 {
				huQuitPids[pid] = true
			}
		}
	}
	return huQuitPids
}

// RoundSettle 单局结算
func (majongSettle *MajongSettle) RoundSettle(desk *desk.Desk, config *desk.DeskConfig) {
	mjContext := config.Context.(*contexts.MjContext).MjContext
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
	// 牌局玩家
	deskPlayers := GetModelManager().GetPlayerModel(desk.GetUid()).GetDeskPlayers()

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
			balanceRsp.BillPlayersInfo = majongSettle.makeBillPlayerInfo(desk, pid, totalValue, nil, mjContext)
		} else if len(contextSInfos) != 0 {
			sinfo := contextSInfos[0]
			cardOption := mjoption.GetCardTypeOption(int(mjContext.GetCardtypeOptionId()))
			fans := make([]*room.Fan, 0)
			fans, totalValue = makeFanType(sinfo.CardType, cardOption)
			balanceRsp.BillPlayersInfo = majongSettle.makeBillPlayerInfo(desk, pid, totalValue, fans, mjContext)
		}
		// 通知该玩家单局结算信息
		GetModelManager().GetMessageModel(desk.GetUid()).BroadCastDeskMessage([]uint64{pid}, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp, true)
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
func (majongSettle *MajongSettle) makeBillPlayerInfo(desk *desk.Desk, currentPid uint64, cardValue int32, fans []*room.Fan, context majongpb.MajongContext) []*room.BillPlayerInfo {
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
	// 记录该玩家退税信息
	revertScore := int64(0)
	revertSInfos := make([]*majongpb.SettleInfo, 0)

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
		// 退税结算详情
		for rID, rScore := range majongSettle.revertScore {
			if rID == sInfo.Id && rScore[pid] != 0 {
				revertScore = revertScore + rScore[pid]
				revertSInfos = append(revertSInfos, sInfo)
			}
		}
	}
	// 获取退税结算详情
	if revertScore != 0 {
		revertbd := majongSettle.getRevertbillDetail(pid, revertScore, revertSInfos)
		billDetails = append(billDetails, revertbd)
	}
	return
}

// getRevertbd 获得玩家退税结算详情，包括分数以及输赢玩家
func (majongSettle *MajongSettle) getRevertbillDetail(pid uint64, revertScore int64, revertSInfos []*majongpb.SettleInfo) *room.BillDetail {
	billDetail := &room.BillDetail{
		SetleType: room.SettleType_ST_TAXREBEAT.Enum(),
		Score:     proto.Int64(-revertScore),
	}
	// 相关联玩家
	for _, revertSInfo := range revertSInfos {
		// 实际扣除分数
		realScore := majongSettle.settleMap[revertSInfo.Id]
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

// getBillDetail 获得玩家单次结算详情，包括番型，分数，倍数，以及输赢玩家
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
