package majong

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/gutils"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// majongCoin   key:playerID value:score
type majongCoin map[uint64]int64

// majongSettle 麻将结算
type majongSettle struct {
	settleMap map[uint64]majongCoin // setttleInfo实际扣分 key:结算id value:majongCoin

	roundScore map[uint64]int64 // 每个玩家单局实际总扣分 key:玩家id value:分数

	handleSettle map[uint64]bool // setttleInfo扣分 key:结算id value:true为已扣分，false为未扣分

	revertScore map[uint64]majongCoin // revertScore  退稅分数 key:退税结算id value:majongCoin
}

// NewMajongSettle 初始化麻将结算
func NewMajongSettle() *majongSettle {
	return &majongSettle{
		settleMap:    make(map[uint64]majongCoin),
		handleSettle: make(map[uint64]bool),
		roundScore:   make(map[uint64]int64),
		revertScore:  make(map[uint64]majongCoin),
	}
}

// Settle 单次结算
func (majongSettle *majongSettle) Settle(desk interfaces.Desk, mjContext majongpb.MajongContext) {

	settleOption := GetSettleOption(int(mjContext.GetGameId())) // 游戏结算玩法

	allSettleInfos := mjContext.SettleInfos // 结算信息

	deskPlayers := desk.GetDeskPlayers() // 牌局玩家

	huQuitPlayers := majongSettle.getHuQuitPlayers(deskPlayers, mjContext) // 胡牌且退出房间后的玩家

	for _, sInfo := range allSettleInfos { // 遍历
		if majongSettle.handleSettle[sInfo.Id] {
			continue
		}
		score := make(map[uint64]int64, 0) // 玩家输赢分数

		brokerPlayers := make([]uint64, 0) // 破产的玩家id

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
			majongSettle.instantSettle(desk, sInfo, score, brokerPlayers)
		}
		// 生成结算完成事件
		GenerateSettleEvent(desk, sInfo.SettleType, brokerPlayers)
	}
	// 退税id
	revertIds := mjContext.RevertSettles
	if len(revertIds) != 0 {
		// 退稅结算信息
		rSettleInfo := majongSettle.generateRevertSettle(deskPlayers, huQuitPlayers, revertIds, settleOption)
		// 扣费并设置玩家金币数
		majongSettle.chargeCoin(deskPlayers, rSettleInfo.Scores)
		// 广播退税信息
		facade.BroadCastDeskMessageExcept(desk, []uint64{}, true, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
			BillPlayersInfo: majongSettle.getBillPlayerInfos(deskPlayers, rSettleInfo, rSettleInfo.Scores),
		})
	}
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
		} else if len(contextSInfos) != 0 {
			sinfo := contextSInfos[0]
			cardOption := mjoption.GetCardTypeOption(int(mjContext.GetCardtypeOptionId()))
			fans := make([]*room.Fan, 0)
			fans = makeFanType(sinfo.CardType, sinfo.HuaCount, cardOption)
			wiiners, _ := getWinners(sinfo.Scores)
			balanceRsp.BillPlayersInfo = majongSettle.makeBillPlayerInfo(wiiners[0], int32(sinfo.CardValue), fans, mjContext)
		}
		// 通知该玩家单局结算信息
		facade.BroadCastDeskMessage(desk, []uint64{pid}, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp, true)
	}
}

// instantSettle 立即结算并扣费
func (majongSettle *majongSettle) instantSettle(desk interfaces.Desk, sInfo *majongpb.SettleInfo, score map[uint64]int64, brokerPlayers []uint64) {
	// 扣费并设置玩家金币数
	majongSettle.chargeCoin(desk.GetDeskPlayers(), score)
	// 广播结算
	facade.BroadCastDeskMessageExcept(desk, []uint64{}, true, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
		BillPlayersInfo: majongSettle.getBillPlayerInfos(desk.GetDeskPlayers(), sInfo, score),
	})
	// 广播认输
	facade.BroadCastDeskMessageExcept(desk, []uint64{}, true, msgid.MsgID_ROOM_PLAYER_GIVEUP_NTF, &room.RoomGiveUpNtf{
		PlayerId: brokerPlayers,
	})
}

func (majongSettle *majongSettle) makeBillDetails(pid uint64, contextSInfos []*majongpb.SettleInfo) (billDetails []*room.BillDetail, totalValue int32) {
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

// getHuQuitPlayers  获取牌局胡牌且退出房间后的玩家id
func (majongSettle *majongSettle) getHuQuitPlayers(dPlayers []interfaces.DeskPlayer, mjContext majongpb.MajongContext) map[uint64]bool {
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

// generateRevertSettle 获取退税的结算信息
func (majongSettle *majongSettle) generateRevertSettle(deskPlayers []interfaces.DeskPlayer, huQuitPlayers map[uint64]bool, revertIds []uint64, settleOption *mjoption.SettleOption) *majongpb.SettleInfo {
	revertScore := make(map[uint64]int64, 0)
	for _, revertID := range revertIds {
		// 需要退钱的玩家
		rlosePid := uint64(0)
		// 需要退的分
		rloseScore := int64(0)
		for pid, score := range majongSettle.settleMap[revertID] {
			if score < 0 { // 胡牌玩家已退出，不用退钱给它
				if !CanRoundSettle(pid, huQuitPlayers, settleOption) {
					continue
				} else {
					revertScore[pid] = revertScore[pid] - score
					rloseScore = rloseScore + score
					majongSettle.revertScore[revertID] = map[uint64]int64{
						pid: -score,
					}
				}
			}
			if score > 0 {
				rlosePid = pid
			}
		}
		revertScore[rlosePid] = revertScore[rlosePid] + rloseScore
		majongSettle.revertScore[revertID] = map[uint64]int64{
			rlosePid: rloseScore,
		}
	}
	return &majongpb.SettleInfo{
		Scores:     revertScore,
		SettleType: majongpb.SettleType_settle_taxrebeat,
	}
}

// makeBillDetail 获得玩家单次结算详情，包括番型，分数，倍数，以及输赢玩家
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
func (majongSettle *majongSettle) getRevertbillDetail(pid uint64, revertScore int64, revertSInfos []*majongpb.SettleInfo) *room.BillDetail {
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
			CurrentScore: proto.Int64(coin),
		}
		if playerID == currentPid {
			billPlayerInfo.CardsGroup = gutils.GetCardsGroup(player)
			billPlayerInfo.Fan = fans
		}
		billPlayerInfos = append(billPlayerInfos, billPlayerInfo)
	}
	return billPlayerInfos
}
