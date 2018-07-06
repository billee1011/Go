package majong

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/gutils"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	majongpb "steve/server_pb/majong"

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
	// 游戏结算玩法
	// settleOption := GetSettleOption(int(mjContext.GetGameId()))
	// 牌局所有结算信息
	contextSInfos := mjContext.SettleInfos
	// 牌局玩家
	deskPlayers := desk.GetDeskPlayers()
	// 胡牌且退出房间后的玩家
	huQuitPlayers := majongSettle.getHuQuitPlayers(deskPlayers, mjContext)
	// 认输玩家
	giveUpPlayers := majongSettle.getGiveupPlayers(deskPlayers, mjContext)
	// 退税id
	revertIds := mjContext.RevertSettles

	// 若存在未处理的结算信息，进行处理
	if len(contextSInfos) != 0 {
		for _, sInfo := range contextSInfos {
			if !majongSettle.handleSettle[sInfo.Id] {
				// 玩家实际输赢分数
				score := make(map[uint64]int64, 0)
				// 破产的玩家id
				brokerPlayers := make([]uint64, 0)
				// 若存在相关联的一组SettleInfo(一炮多响情况)
				if len(sInfo.GroupId) > 1 {
					groupSInfos, sumSInfo := majongSettle.sumSettleInfo(mjContext.SettleInfos, sInfo)
					score, brokerPlayers = majongSettle.calcCoin(deskPlayers, mjContext.GetPlayers(), huQuitPlayers, sumSInfo.Scores)
					majongSettle.apartSettle(groupSInfos, score)
				} else {
					score, brokerPlayers = majongSettle.calcCoin(deskPlayers, mjContext.GetPlayers(), huQuitPlayers, sInfo.Scores)
					majongSettle.settleMap[sInfo.Id] = score
					majongSettle.handleSettle[sInfo.Id] = true
				}
				// 扣费并设置玩家金币数
				majongSettle.chargeCoin(deskPlayers, score)
				// 广播结算信息
				NotifyMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
					BillPlayersInfo: majongSettle.getBillPlayerInfos(deskPlayers, sInfo, score),
				})
				// 查花猪，大叫，退税过程不需要认输
				if !IsRoundSettle(sInfo.SettleType) {
					needSend := make([]uint64, 0)
					for _, brokerPlayer := range brokerPlayers {
						if !giveUpPlayers[brokerPlayer] {
							needSend = append(needSend, brokerPlayer)
						}
					}
					// 广播认输信息
					NotifyMessage(desk, msgid.MsgID_ROOM_PLAYER_GIVEUP_NTF, &room.RoomGiveUpNtf{
						PlayerId: needSend,
					})
				}
				// 生成结算完成事件
				GenerateSettleEvent(desk, sInfo.SettleType, brokerPlayers)
			}
		}
	}
	if len(revertIds) != 0 {
		for _, revertID := range revertIds {
			// 退稅结算信息
			gangSettle := GetSettleInfoByID(contextSInfos, revertID)
			rSettleInfo := majongSettle.generateRevertSettle2(revertID, gangSettle, huQuitPlayers, giveUpPlayers, revertIds, mjContext)
			if rSettleInfo != nil {
				NotifyMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
					BillPlayersInfo: majongSettle.getBillPlayerInfos(deskPlayers, rSettleInfo, rSettleInfo.Scores),
				})
			}
		}
	}
}

// RoundSettle 单局结算
func (majongSettle *majongSettle) RoundSettle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	majongSettle.Settle(desk, mjContext)
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
					bd := majongSettle.getBillDetail(pid, sInfo)
					cardValue = cardValue + bd.GetFanValue()
					balanceRsp.BillDetail = append(balanceRsp.BillDetail, bd)
				}

			}
			// 退税结算详情
			for _, sInfo := range contextSInfos {
				for rID, rScore := range majongSettle.revertScore {
					if rID == sInfo.Id && rScore[pid] != 0 {
						revertbd := majongSettle.getRevertbillDetail(pid, rScore)
						balanceRsp.BillDetail = append(balanceRsp.BillDetail, revertbd)
					}
				}
			}
			// 获取玩家单局结算详情
			balanceRsp.BillPlayersInfo = majongSettle.getRoundBillPlayerInfo(pid, cardValue, mjContext)
			// 通知该玩家单局结算信息
			NotifyPlayersMessage(desk, []uint64{pid}, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp)
		}
	}
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

// getGiveupPlayers  获取认输的玩家id
func (majongSettle *majongSettle) getGiveupPlayers(dPlayers []interfaces.DeskPlayer, mjContext majongpb.MajongContext) map[uint64]bool {
	giveupPlayers := make(map[uint64]bool, 0)
	for _, cPlayer := range mjContext.Players {
		if cPlayer.GetXpState() == 2 {
			giveupPlayers[cPlayer.GetPalyerId()] = true
		}
	}
	return giveupPlayers
}

// generateSettleEvent 生成结算finish事件
func (majongSettle *majongSettle) generateSettleEvent(desk interfaces.Desk, settleType majongpb.SettleType, brokerPlayers []uint64) {
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
func (majongSettle *majongSettle) sumSettleInfo(contextSInfo []*majongpb.SettleInfo, settleInfo *majongpb.SettleInfo) ([]*majongpb.SettleInfo, *majongpb.SettleInfo) {
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

// apartSettle  将score分配到各自settleInfo中
func (majongSettle *majongSettle) apartSettle(groupSettleInfos []*majongpb.SettleInfo, allScores map[uint64]int64) {
	for _, sInfo := range groupSettleInfos {
		sID := sInfo.Id
		cost := int64(0)
		majongSettle.settleMap[sID] = make(map[uint64]int64)
		for pid, score := range sInfo.Scores {
			if score == 0 {
				continue
			} else if score > 0 {
				cost = allScores[pid]
				majongSettle.settleMap[sID][pid] = allScores[pid]
			} else {
				if cost != 0 {
					majongSettle.settleMap[sID][pid] = 0 - cost
				} else {
					for _, allScore := range allScores {
						if allScore > 0 {
							majongSettle.settleMap[sID][pid] = 0 - allScore
						}
					}
				}
			}
		}
		majongSettle.handleSettle[sID] = true
	}
}

// calcMaxScore 计算玩家输赢上限
// 赢豆上限 = max(进房豆子数,当前豆子数)
// 胡牌且退出房间后不参与牌局的所有结算
func (majongSettle *majongSettle) calcMaxScore(deskPlayer []interfaces.DeskPlayer, huQuitPlayers map[uint64]bool, score map[uint64]int64) (maxScore map[uint64]int64) {
	maxScore = make(map[uint64]int64, 0)
	losePids := make([]uint64, 0)
	winnPids := make([]uint64, 0)
	for pid, pscore := range score {
		if pscore > 0 {
			if huQuitPlayers[pid] {
				maxScore[pid] = 0
			} else {
				maxScore[pid] = majongSettle.getWinMax(GetDeskPlayer(deskPlayer, pid), pscore)
			}
		} else if pscore < 0 {
			losePids = append(losePids, pid)
		}
		if pscore > 0 {
			winnPids = append(winnPids, pid)
		}
		if huQuitPlayers[pid] {
			score[pid] = 0
		}
	}
	if len(losePids) == 1 {
		for _, winnPid := range winnPids {
			winMax := majongSettle.getWinMax(GetDeskPlayer(deskPlayer, winnPid), score[winnPid])
			maxScore[winnPid] = score[winnPid]
			if score[winnPid] >= winMax {
				maxScore[winnPid] = winMax
			}
			maxScore[losePids[0]] = maxScore[losePids[0]] - maxScore[winnPid]
		}
	} else if len(losePids) > 1 {
		for _, losePid := range losePids {
			winMax := majongSettle.getWinMax(GetDeskPlayer(deskPlayer, winnPids[0]), score[losePid])
			maxScore[losePid] = score[losePid]
			if majongSettle.abs(score[losePid]) >= winMax {
				maxScore[losePid] = 0 - winMax
			}
			maxScore[winnPids[0]] = maxScore[winnPids[0]] - maxScore[losePid]
		}
	}
	return
}

func (majongSettle *majongSettle) getWinMax(winPlayer interfaces.DeskPlayer, winScore int64) (winMax int64) {
	winMax = int64(0)
	winPid := winPlayer.GetPlayerID()
	currentCoin := int64(global.GetPlayerMgr().GetPlayer(winPid).GetCoin()) // 当前豆子数
	enterCoin := int64(winPlayer.GetEcoin())                                // 进房豆子数
	if currentCoin >= enterCoin {
		winMax = currentCoin
	} else {
		winMax = enterCoin
	}
	return
}

// calcCoin 计算扣除的金币
// 如果出现一炮多响的情况：
// 1.玩家身上的钱够赔付胡牌玩家的话,直接赔付
// 2.玩家身上的钱不够赔付胡牌玩家的话,那么该玩家身上的钱平分给胡牌玩家，,按逆时针方向,从点炮者数起,余 1 情况赔付于第一胡牌玩家,
//	 余 2 情况赔付于第一、第二胡牌玩家;
func (majongSettle *majongSettle) calcCoin(deskPlayer []interfaces.DeskPlayer, contextPlayer []*majongpb.Player, huQuitPlayers map[uint64]bool, score map[uint64]int64) (map[uint64]int64, []uint64) {
	// 赢豆上限
	maxScore := majongSettle.calcMaxScore(deskPlayer, huQuitPlayers, score)
	// 赢家
	winPlayers := make([]uint64, 0)
	// 输家
	losePlayers := make([]uint64, 0)
	// 赢的分数
	totalWin := int64(0)
	// 输的分数(总共)
	totalose := int64(0)
	for playerID, pScore := range maxScore {
		if pScore > 0 {
			totalWin = totalWin + pScore
			winPlayers = append(winPlayers, playerID)
		} else if pScore < 0 {
			totalose = totalose + pScore
			losePlayers = append(losePlayers, playerID)
		}
	}
	// 每个玩家扣除的金币数
	coinCost := make(map[uint64]int64, 0)
	// 破产玩家
	brokePlayers := make([]uint64, 0)
	// 输家人数
	loseSum := len(losePlayers)
	// 赢家人数
	winSum := len(winPlayers)
	if winSum == 1 && loseSum > 1 { // 有多个输家，最多不能赢超过输家的豆子数
		// 赢家
		winPlayer := winPlayers[0]
		for _, losePid := range losePlayers {
			loseScore := majongSettle.abs(maxScore[losePid])                      // 输家输的分
			loseCoin := int64(global.GetPlayerMgr().GetPlayer(losePid).GetCoin()) // 输家金币数
			if loseScore < loseCoin {
				coinCost[losePid] = -loseScore
			} else {
				coinCost[losePid] = -loseCoin
				brokePlayers = append(brokePlayers, losePid)
			}
			coinCost[winPlayer] = coinCost[winPlayer] - coinCost[losePid]
		}
	} else if loseSum == 1 { // 1个输家，多个赢家
		// 输家
		losePlayer := losePlayers[0]
		loseScore := majongSettle.abs(totalose)                                  // 输家输的分
		loseCoin := int64(global.GetPlayerMgr().GetPlayer(losePlayer).GetCoin()) // 输家金币数
		if loseScore < loseCoin {
			// 金币数够扣
			for _, win := range winPlayers {
				coinCost[win] = maxScore[win]
			}
			coinCost[losePlayer] = maxScore[losePlayer]
		} else {
			// 金币数不够扣，赢家为1时直接输家的金币全部给赢家，否则平分
			if winSum == 1 {
				coinCost[winPlayers[0]] = loseCoin
				coinCost[losePlayer] = -loseCoin
			} else {
				// 多个赢家，按照赢家人数平分
				for _, winPid := range winPlayers {
					winScore := int64(loseCoin / int64(winSum))
					if winScore >= maxScore[winPid] {
						winScore = maxScore[winPid]
					}
					coinCost[winPid] = winScore
					coinCost[losePlayer] = coinCost[losePlayer] - coinCost[winPid]
				}
				// 剩余分数，余 1 情况赔付于赢钱最多的玩家, 余 2 情况赔付于第一、第二胡牌玩家
				surplusScore := loseCoin - majongSettle.abs(coinCost[losePlayer])
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
			brokePlayers = append(brokePlayers, losePlayer)
		}
	}
	return coinCost, brokePlayers
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
		if costScore[pid] != 0 {
			billplayerInfo := majongSettle.newBillplayerInfo(pid, majongSettle.settleType2BillType(settleInfo.SettleType))
			billplayerInfo.Score = proto.Int64(costScore[pid])
			holdCoin := global.GetPlayerMgr().GetPlayer(pid).GetCoin()
			billplayerInfo.CurrentScore = proto.Int64(int64(holdCoin))
			billplayerInfos = append(billplayerInfos, billplayerInfo)
		}
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
	coinCost = majongSettle.calcTaxbetCoin(rlosePid, rWinnerPids, scoreCost, mjContext.GetPlayers())
	majongSettle.revertScore[revertID] = coinCost
	return &majongpb.SettleInfo{
		Scores:     coinCost,
		SettleType: majongpb.SettleType_settle_taxrebeat,
	}
}

func (majongSettle *majongSettle) calcTaxbetCoin(losePlayer uint64, winPlayers []uint64, score map[uint64]int64, contextPlayer []*majongpb.Player) (coinCost map[uint64]int64) {
	coinCost = make(map[uint64]int64, 0)
	loseCoin := int64(global.GetPlayerMgr().GetPlayer(losePlayer).GetCoin()) // 输家金币数
	loseScore := score[losePlayer]
	if majongSettle.abs(loseScore) < loseCoin {
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
			surplusScore := loseCoin - majongSettle.abs(coinCost[losePlayer])
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

// getBillDetail 获得玩家单次结算详情，包括番型，分数，倍数，以及输赢玩家
func (majongSettle *majongSettle) getBillDetail(pid uint64, sInfo *majongpb.SettleInfo) *room.BillDetail {
	billDetail := &room.BillDetail{
		SetleType: room.SettleType(sInfo.SettleType).Enum(),
		HuType:    room.HuType(sInfo.HuType).Enum(),
		FanValue:  proto.Int32(int32(sInfo.CardValue)),
		GenCount:  proto.Uint32(sInfo.GenCount),
		Score:     proto.Int64(majongSettle.settleMap[sInfo.Id][pid]),
	}
	// 实际扣除分数
	realScore := majongSettle.settleMap[sInfo.Id]
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

// getRoundBillPlayerInfo 获得单局结算玩家详情,包括玩家自己牌型,输赢分数，以及其余每个玩家的输赢分数
func (majongSettle *majongSettle) getRoundBillPlayerInfo(currentPid uint64, cardValue int32, context majongpb.MajongContext) []*room.BillPlayerInfo {
	billPlayerInfos := make([]*room.BillPlayerInfo, 0)
	for _, player := range context.Players {
		playerID := player.GetPalyerId()
		billPlayerInfo := &room.BillPlayerInfo{
			Pid:       proto.Uint64(playerID),
			Score:     proto.Int64(majongSettle.roundScore[playerID]),
			CardValue: proto.Int32(cardValue),
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

func (majongSettle *majongSettle) newBillplayerInfo(playID uint64, billType room.BillType) *room.BillPlayerInfo {
	return &room.BillPlayerInfo{
		Pid:      proto.Uint64(playID),
		BillType: billType.Enum(),
	}
}

func (majongSettle *majongSettle) abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func (majongSettle *majongSettle) settleType2BillType(settleType majongpb.SettleType) room.BillType {
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
