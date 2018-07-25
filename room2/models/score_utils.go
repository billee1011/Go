package models

import (
	"steve/gutils"
	majongpb "steve/server_pb/majong"
	"steve/room2/player"
)

// mergeSettle 合并一组SettleInfo
// 返回参数:	[]*majongpb.SettleInfo(该组settleInfo) / *majongpb.SettleInfo(合并后的settleInfo)
func MergeSettle(contextSInfo []*majongpb.SettleInfo, settleInfo *majongpb.SettleInfo) ([]*majongpb.SettleInfo, *majongpb.SettleInfo) {
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

// GetSettleInfoBySid 根据settleID获取对应settleInfo的下标index
func GetSettleInfoBySid(settleInfos []*majongpb.SettleInfo, ID uint64) int {
	for index, s := range settleInfos {
		if s.Id == ID {
			return index
		}
	}
	return -1
}

// calcCoin 计算扣除的金币
// 如果出现一炮多响的情况：
// 1.玩家身上的钱够赔付胡牌玩家的话,直接赔付
// 2.玩家身上的钱不够赔付胡牌玩家的话,那么该玩家身上的钱平分给胡牌玩家，,按逆时针方向,从点炮者数起,余 1 情况赔付于第一胡牌玩家,
//	 余 2 情况赔付于第一、第二胡牌玩家;
func CalcCoin(deskPlayer []*player.Player, contextPlayer []*majongpb.Player, huQuitPlayers map[uint64]bool, score map[uint64]int64) (map[uint64]int64, []uint64) {
	// 赢豆上限
	maxScore := getMaxScore(deskPlayer, huQuitPlayers, score)
	// 赢家
	winPlayers := make([]uint64, 0)
	// 输家
	losePlayers := make([]uint64, 0)
	// 输的分数(总共)
	totalose := int64(0)

	winPlayers, _ = getWinners(maxScore)
	losePlayers, totalose = getLosers(maxScore)
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
		coinCost, brokePlayers = calcSocreWinner1(winPlayer, losePlayers, maxScore)

	} else if loseSum == 1 { // 1个输家
		// 输家
		losePlayer := losePlayers[0]
		loseScore := abs(totalose) // 输家输的分
		coinCost, brokePlayers = calcSocrelose1(winPlayers, losePlayer, loseScore, maxScore, contextPlayer)
	}
	return coinCost, brokePlayers
}

// getMaxScore 计算玩家输赢上限
// 赢豆上限 = max(进房豆子数,当前豆子数)
// 胡牌且退出房间后不参与牌局的所有结算
func getMaxScore(deskPlayer []*player.Player, huQuitPlayers map[uint64]bool, score map[uint64]int64) (maxScore map[uint64]int64) {
	maxScore = make(map[uint64]int64, 0)
	losePids := make([]uint64, 0)
	winnPids := make([]uint64, 0)
	for pid, pscore := range score {
		if pscore > 0 {
			if huQuitPlayers[pid] {
				maxScore[pid] = 0
			} else {
				maxScore[pid] = getWinMax(GetDeskPlayer(deskPlayer, pid), pscore)
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
			winMax := getWinMax(GetDeskPlayer(deskPlayer, winnPid), score[winnPid])
			if score[winnPid] >= winMax {
				maxScore[winnPid] = winMax
			}
			maxScore[winnPid] = score[winnPid]
			maxScore[losePids[0]] = maxScore[losePids[0]] - maxScore[winnPid]
		}
	} else if len(losePids) > 1 {
		for _, losePid := range losePids {
			winMax := getWinMax(GetDeskPlayer(deskPlayer, winnPids[0]), score[losePid])
			if abs(score[losePid]) >= winMax {
				maxScore[losePid] = 0 - winMax
			}
			maxScore[losePid] = score[losePid]
			maxScore[winnPids[0]] = maxScore[winnPids[0]] - maxScore[losePid]
		}
	}
	return
}

// GetDeskPlayer 获取指定id的room Player
func GetDeskPlayer(deskPlayers []*player.Player, pid uint64) *player.Player {
	for _, p := range deskPlayers {
		if p.GetPlayerID() == pid {
			return p
		}
	}
	return nil
}


func getWinMax(winPlayer *player.Player, winScore int64) (winMax int64) {
	winMax = int64(0)
	currentCoin := int64(winPlayer.GetCoin()) // 当前豆子数
	enterCoin := int64(winPlayer.GetEcoin())                                // 进房豆子数
	if currentCoin >= enterCoin {
		winMax = currentCoin
	} else {
		winMax = enterCoin
	}
	return
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// getWinners 获取赢家及赢得总分
func getWinners(score map[uint64]int64) ([]uint64, int64) {
	winPlayers := make([]uint64, 0)
	totalWin := int64(0)
	for playerID, pScore := range score {
		if pScore > 0 {
			totalWin = totalWin + pScore
			winPlayers = append(winPlayers, playerID)
		}
	}
	return winPlayers, totalWin
}

// getLosers 获取输家及输得总分
func getLosers(score map[uint64]int64) ([]uint64, int64) {
	losePlayers := make([]uint64, 0)
	totalose := int64(0)
	for playerID, pScore := range score {
		if pScore < 0 {
			totalose = totalose + pScore
			losePlayers = append(losePlayers, playerID)
		}
	}
	return losePlayers, totalose
}

// calcSocreWinner1 赢家唯一时扣分
func calcSocreWinner1(winPlayer uint64, losePlayers []uint64, maxScore map[uint64]int64) (map[uint64]int64, []uint64) {
	// 每个玩家扣除的金币数
	coinCost := make(map[uint64]int64, 0)
	// 破产玩家
	brokePlayers := make([]uint64, 0)
	for _, losePid := range losePlayers {
		loseScore := abs(maxScore[losePid])                                   // 输家输的分
		roomPlayer := player.GetPlayerMgr().GetPlayer(losePid)
		loseCoin := int64(roomPlayer.GetCoin()) // 输家金币数
		if loseScore < loseCoin {
			coinCost[losePid] = -loseScore
		} else {
			coinCost[losePid] = -loseCoin
			brokePlayers = append(brokePlayers, losePid)
		}
		coinCost[winPlayer] = coinCost[winPlayer] - coinCost[losePid]
	}
	return coinCost, brokePlayers
}

// calcSocrelose1 输家唯一时扣分
func calcSocrelose1(winPlayers []uint64, losePlayer uint64, loseScore int64, maxScore map[uint64]int64, contextPlayer []*majongpb.Player) (map[uint64]int64, []uint64) {
	// 每个玩家扣除的金币数
	coinCost := make(map[uint64]int64, 0)
	// 破产玩家
	brokePlayers := make([]uint64, 0)
	// 输家金币数
	roomPlayer := player.GetPlayerMgr().GetPlayer(losePlayer)
	loseCoin := int64(roomPlayer.GetCoin())
	// 赢家人数
	winSum := len(winPlayers)
	if loseScore < loseCoin {
		// 金币数够扣
		for _, win := range winPlayers {
			coinCost[win] = maxScore[win]
		}
		coinCost[losePlayer] = maxScore[losePlayer]
	} else if winSum == 1 {
		coinCost[winPlayers[0]] = loseCoin // 金币数不够扣，赢家为1时直接输家的金币全部给赢家，否则平分
		coinCost[losePlayer] = -loseCoin
		brokePlayers = append(brokePlayers, losePlayer)
	} else {
		coinCost, brokePlayers = divideScore(losePlayer, winPlayers, maxScore, contextPlayer)
	}
	return coinCost, brokePlayers
}

// divideScore 输家豆子数不够时平分给多个赢家，剩余豆子再分
func divideScore(losePlayer uint64, winPlayers []uint64, maxScore map[uint64]int64, contextPlayer []*majongpb.Player) (map[uint64]int64, []uint64) {
	// 每个玩家扣除的金币数
	coinCost := make(map[uint64]int64, 0)
	// 破产玩家
	brokePlayers := make([]uint64, 0)
	// 赢家人数
	winSum := len(winPlayers)
	// 输家金币数
	roomPlayer := player.GetPlayerMgr().GetPlayer(losePlayer)
	loseCoin := int64(roomPlayer.GetCoin())
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
	surplusScore := loseCoin - coinCost[losePlayer]
	resortWinnerPlayers := resortWinnerPlayers(losePlayer, winPlayers, contextPlayer)
	firstWinner := resortWinnerPlayers[0]
	if surplusScore%2 == 0 {
		secondWinner := resortWinnerPlayers[1]
		coinCost[firstWinner] = coinCost[firstWinner] + surplusScore/2
		coinCost[secondWinner] = coinCost[secondWinner] + surplusScore/2
	} else {
		coinCost[firstWinner] = coinCost[firstWinner] + surplusScore
	}
	coinCost[losePlayer] = coinCost[losePlayer] - surplusScore
	brokePlayers = append(brokePlayers, losePlayer)
	return coinCost, brokePlayers
}

// resortWinnerPlayers 返回按输家位置重新排序的赢家位置列表
func resortWinnerPlayers(losePlayer uint64, winPlayers []uint64, contextPlayer []*majongpb.Player) []uint64 {
	loseIndex := gutils.GetPlayerIndex(losePlayer, contextPlayer)
	resortPlayers := make([]uint64, 0)
	for i := 0; i < len(contextPlayer); i++ {
		index := (loseIndex + i) % len(contextPlayer)
		resortPlayers = append(resortPlayers, contextPlayer[index].GetPalyerId())
	}
	resortWinnerPlayers := make([]uint64, 0)
	for _, resortPID := range resortPlayers {
		for _, winPlayer := range winPlayers {
			if resortPID == winPlayer {
				resortWinnerPlayers = append(resortWinnerPlayers, resortPID)
			}
		}
	}
	return resortWinnerPlayers
}
