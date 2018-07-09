package majong

import (
	"steve/common/mjoption"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// HuSettle 胡结算
type HuSettle struct {
}

// Settle  胡结算方法
// 胡牌分=番型*底分
func (huSettle *HuSettle) Settle(params interfaces.HuSettleParams) []*majongpb.SettleInfo {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":    "HuSettle",
		"GameID":       params.GameID,
		"winnersID":    params.HuPlayers,
		"settleType":   params.SettleType,
		"huType":       params.HuType,
		"allPlayers":   params.HuPlayers,
		"hasHuPlayers": params.HasHuPlayers,
		"quitPlayers":  params.QuitPlayers,
		"cardTypes":    params.CardTypes,
		"cardValues":   params.CardValues,
		"genCount":     params.GenCount,
	})
	logEntry.Debugln("胡结算信息")
	// 游戏结算玩法
	settleOption := GetSettleOption(int(params.GameID))
	// 结算信息
	settleInfos := make([]*majongpb.SettleInfo, 0)
	// 底数
	ante := GetDi()

	if params.SettleType == majongpb.SettleType_settle_zimo {
		huSettleInfo := new(majongpb.SettleInfo)
		// 赢家
		huPlayerID := params.HuPlayers[0]
		// 赢分
		win := int64(0)
		// 倍数
		toalValue := uint32(params.CardValues[huPlayerID])
		// 总分 (底分*倍数)
		total := int64(toalValue) * ante
		// 玩家输赢分
		scoreInfo := make(map[uint64]int64)
		// 自摸全赔
		for _, playerID := range params.AllPlayers {
			huPlayerID := params.HuPlayers[0]
			if playerID != huPlayerID && huSettle.canHuSettle(playerID, params.GiveupPlayers, params.HasHuPlayers, params.QuitPlayers, settleOption) {
				scoreInfo[playerID] = -total
				win = win + total
			}
		}
		scoreInfo[huPlayerID] = win
		huSettleInfo = newHuSettleInfo(&params, scoreInfo, huPlayerID, toalValue)
		settleInfos = append(settleInfos, huSettleInfo)
	} else if params.SettleType == majongpb.SettleType_settle_dianpao {
		groupIds := make([]uint64, 0)
		huSettleInfos := make([]*majongpb.SettleInfo, 0)
		for _, huPlayerID := range params.HuPlayers {
			scoreInfo := make(map[uint64]int64)
			huSettleInfo := new(majongpb.SettleInfo)
			// 倍数(番型倍数*胡牌倍数)
			toalValue := uint32(params.CardValues[huPlayerID])
			// 总分 (底分*倍数)
			total := int64(toalValue) * ante
			// 点炮一家给
			scoreInfo[huPlayerID] = total
			scoreInfo[params.SrcPlayer] = -total
			huSettleInfo = newHuSettleInfo(&params, scoreInfo, huPlayerID, toalValue)
			huSettleInfos = append(huSettleInfos, huSettleInfo)
			groupIds = append(groupIds, huSettleInfo.Id)
		}
		for _, huSettleInfo := range huSettleInfos { // 一炮多响结算信息为一组
			huSettleInfo.GroupId = groupIds
			settleInfos = append(settleInfos, huSettleInfo)
		}
	}
	if params.HuType == majongpb.HuType_hu_ganghoupao { // 杠后炮需呼叫转移
		cSettleInfo := huSettle.newCallTransferSettleInfo(&params, settleOption)
		settleInfos = append(settleInfos, cSettleInfo)
	}
	return settleInfos
}

// callTransfer 呼叫转移结算信息
// 呼叫转移:
//  	非一炮多响:玩家杠牌后形成杠上炮，则需将自己这次的杠牌所得杠钱转给胡牌者
//      一炮多响:1. （暗杠、补杠）先收杠钱，然后把最后一个的杠钱平分给胡牌玩家，如果收到的杠钱平分后还有多余，则多余的杠钱按位置给第一个胡牌玩家
//	   		    2.  （直杠）如果胡家中包含点杠者，则转移给点杠者，否则平分
func (huSettle *HuSettle) newCallTransferSettleInfo(params *interfaces.HuSettleParams, settleOption *mjoption.SettleOption) *majongpb.SettleInfo {
	params.SettleID = params.SettleID + 1
	callTransferS := &majongpb.SettleInfo{
		Id:         params.SettleID,
		Scores:     make(map[uint64]int64),
		HuType:     -1,
		SettleType: majongpb.SettleType_settle_calldiver,
	}
	// 杠牌
	gangCard := params.GangCard
	// 杠倍数
	gangValue := GetGangValue(settleOption, gangCard.GetType())
	// 赢家人数
	winSum := int64(len(params.HuPlayers))
	// 底数
	ante := GetDi()
	// 杠的分数
	gangScore := ante * int64(gangValue)
	// 点炮者
	dianPaoPlayer := params.SrcPlayer
	if winSum == 1 {
		if gangCard.GetType() == majongpb.GangType_gang_angang || gangCard.GetType() == majongpb.GangType_gang_bugang {
			if len(params.HasHuPlayers) <= 1 {
				gangScore = gangScore * int64(len(params.AllPlayers)-1)
			} else {
				if settleOption.HuPlayerSettle.HuPlayerGangSettle {
					gangScore = gangScore * int64(len(params.AllPlayers)-1)
				} else {
					gangScore = gangScore * int64(len(params.AllPlayers)-len(params.HasHuPlayers))
				}
			}
		}
		callTransferS.Scores[params.HuPlayers[0]] = gangScore
		callTransferS.Scores[dianPaoPlayer] = -gangScore
	} else {
		// 一炮多响
		if gangCard.GetType() == majongpb.GangType_gang_minggang {
			dianGangPlayer := gangCard.GetSrcPlayer()
			contain := false
			for _, huPlayerID := range params.HuPlayers {
				if dianGangPlayer != huPlayerID {
					continue
				}
				contain = true
				break
			}
			if contain {
				callTransferS.Scores[dianGangPlayer] = gangScore
				callTransferS.Scores[dianPaoPlayer] = -gangScore
			} else {
				// 平分
				huSettle.divideScore(gangScore, winSum, params, callTransferS)
			}
		} else if gangCard.GetType() == majongpb.GangType_gang_angang || gangCard.GetType() == majongpb.GangType_gang_bugang {
			// （暗杠、补杠）先收杠钱,平分,杠钱后还有多余，多余的杠钱按位置给第一个胡牌玩家
			gangScore = gangScore * int64(len(params.AllPlayers)-1)
			// 平分
			huSettle.divideScore(gangScore, winSum, params, callTransferS)
		}
	}
	return callTransferS
}

// GetDi 获取底注
func GetDi() int64 {
	//return r.Option.(*pb.Option_SiChuangXueLiu).Di
	return 1
}

// newHuSettleInfo 生成胡结算信息
func newHuSettleInfo(params *interfaces.HuSettleParams, scoreInfo map[uint64]int64, huPlayerID uint64, cardValue uint32) *majongpb.SettleInfo {
	params.SettleID = params.SettleID + 1
	return &majongpb.SettleInfo{
		Id:         params.SettleID,
		Scores:     scoreInfo,
		SettleType: params.SettleType,
		HuType:     params.HuType,
		CardType:   params.CardTypes[huPlayerID],
		GenCount:   uint32(params.GenCount[huPlayerID]),
		HuaCount:   uint32(params.HuaCount[huPlayerID]),
		CardValue:  cardValue,
	}
}

func (huSettle *HuSettle) calcTotalValue(cardValue, huValue uint32) uint32 {
	return cardValue * huValue
}

func (huSettle *HuSettle) divideScore(gangScore, winSum int64, params *interfaces.HuSettleParams, callTransferS *majongpb.SettleInfo) {
	// 平分
	equallyTotal := gangScore / winSum
	// 剩余分数
	surplusTotal := gangScore - (equallyTotal * int64(winSum))
	// 点炮者
	dianPaoPlayer := params.SrcPlayer

	for _, huPlayerID := range params.HuPlayers {
		callTransferS.Scores[huPlayerID] = equallyTotal
		callTransferS.Scores[dianPaoPlayer] = callTransferS.Scores[dianPaoPlayer] - equallyTotal
	}
	if surplusTotal != 0 {
		startIndex, _ := utils.GetPlayerIDIndex(dianPaoPlayer, params.AllPlayers)
		firstPlayerID := utils.GetPalyerCloseFromTarget(startIndex, params.AllPlayers, params.HuPlayers)
		if firstPlayerID != 0 {
			callTransferS.Scores[firstPlayerID] = callTransferS.Scores[firstPlayerID] + surplusTotal
			callTransferS.Scores[dianPaoPlayer] = callTransferS.Scores[dianPaoPlayer] - surplusTotal
		}
	}
}

// canHuSettle 玩家能否参与胡结算
func (huSettle *HuSettle) canHuSettle(playerID uint64, givePlayers, hasHuPlayers, quitPlayers []uint64, settleOption *mjoption.SettleOption) bool {
	for _, giveupPlayer := range givePlayers {
		if giveupPlayer != playerID {
			break
		}
		return settleOption.GiveUpPlayerSettle.GiveUpPlayerHuSettle
	}
	for _, hasHupalyer := range hasHuPlayers {
		if hasHupalyer != playerID {
			break
		}
		for _, quitPlayer := range quitPlayers {
			if quitPlayer != playerID {
				break
			}
			return settleOption.HuQuitPlayerSettle.HuQuitPlayeHuSettle
		}
		return settleOption.HuPlayerSettle.HuPlayeHuSettle
	}
	return true
}
