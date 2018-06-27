package scxl

import (
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// HuSettle 胡结算
type HuSettle struct {
}

// Settle  胡结算方法
func (huSettle *HuSettle) Settle(params interfaces.HuSettleParams) []*majongpb.SettleInfo {
	entry := logrus.WithFields(logrus.Fields{
		"name":       "HuSettle",
		"winnersID":  params.HuPlayers,
		"settleType": params.SettleType,
		"huType":     params.HuType,
		"cardTypes":  params.CardTypes,
		"cardValues": params.CardValues,
		"genCount":   params.GenCount,
	})
	settleInfos := make([]*majongpb.SettleInfo, 0)
	scoreInfo := make(map[uint64]int64)
	// 底数
	ante := GetDi()

	if params.SettleType == majongpb.SettleType_settle_zimo {
		huSettleInfo := new(majongpb.SettleInfo)
		// 赢家
		huPlayerID := params.HuPlayers[0]
		// 赢分
		win := int64(0)
		// 倍数
		value := int64(params.CardValues[huPlayerID]) * int64(getHuTypeValue(params.HuType))
		// 总分
		total := value * ante
		// 自摸全赔
		for _, playerID := range params.AllPlayers {
			if playerID != huPlayerID {
				scoreInfo[playerID] = -total
				win = win + total
			}
		}
		scoreInfo[huPlayerID] = win
		huSettleInfo = newHuSettleInfo(&params, scoreInfo, huPlayerID)
		huSettleInfo.CardValue = uint32(value)
		settleInfos = append(settleInfos, huSettleInfo)
	} else if params.SettleType == majongpb.SettleType_settle_dianpao {
		groupIds := make([]uint64, 0)
		huSettleInfos := make([]*majongpb.SettleInfo, 0)
		for _, huPlayerID := range params.HuPlayers {
			scoreInfo := make(map[uint64]int64)
			huSettleInfo := new(majongpb.SettleInfo)
			// 倍数
			value := int64(params.CardValues[huPlayerID]) * int64(getHuTypeValue(params.HuType))
			// 总分
			total := value * ante
			// 输赢分
			scoreInfo[huPlayerID] = total
			scoreInfo[params.SrcPlayer] = -total
			huSettleInfo = newHuSettleInfo(&params, scoreInfo, huPlayerID)
			huSettleInfo.CardValue = uint32(value)
			huSettleInfos = append(huSettleInfos, huSettleInfo)
			groupIds = append(groupIds, huSettleInfo.Id)
		}
		for _, huSettleInfo := range huSettleInfos {
			huSettleInfo.GroupId = groupIds
			settleInfos = append(settleInfos, huSettleInfo)
		}
	}
	if params.HuType == majongpb.HuType_hu_ganghoupao { // 杠后炮需呼叫转移
		ctSettleInfo := huSettle.callTransferSettle(&params)
		settleInfos = append(settleInfos, ctSettleInfo)
	}
	entry.Info("胡结算")
	return settleInfos
}

// callTransferSettle 呼叫转移结算信息
func (huSettle *HuSettle) callTransferSettle(params *interfaces.HuSettleParams) *majongpb.SettleInfo {
	callTransferS := newCallTransferSettleInfo(params)

	gangCard := params.GangCard
	gangScore := getGangScore(gangCard.GetType())
	// 赢家人数
	winSum := len(params.HuPlayers)

	score := GetDi() * int64(gangScore)

	if winSum == 1 {
		if gangCard.GetType() == majongpb.GangType_gang_angang || gangCard.GetType() == majongpb.GangType_gang_bugang {
			score = score * int64(len(params.AllPlayers)-1)
		}
		callTransferS.Scores[params.HuPlayers[0]] = score
		callTransferS.Scores[params.SrcPlayer] = -score
	} else {
		// 一炮多响
		if gangCard.GetType() == majongpb.GangType_gang_minggang { // （直杠）如果胡家中包含点杠者，则转移给点杠者，否则平分
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
				callTransferS.Scores[dianGangPlayer] = score
				callTransferS.Scores[params.SrcPlayer] = -score
			} else {
				// 平分
				equallyTotal := score / int64(winSum)
				for _, huPlayerID := range params.HuPlayers {
					callTransferS.Scores[huPlayerID] = equallyTotal
					callTransferS.Scores[params.SrcPlayer] = callTransferS.Scores[params.SrcPlayer] - equallyTotal
				}
			}
		} else if gangCard.GetType() == majongpb.GangType_gang_angang || gangCard.GetType() == majongpb.GangType_gang_bugang {
			// （暗杠、补杠）先收杠钱,平分,杠钱后还有多余，多余的杠钱按位置给第一个胡牌玩家
			score = score * int64(len(params.AllPlayers)-1)
			// 平分
			equallyTotal := score / int64(winSum)
			for _, huPlayerID := range params.HuPlayers {
				callTransferS.Scores[huPlayerID] = equallyTotal
				callTransferS.Scores[params.SrcPlayer] = callTransferS.Scores[params.SrcPlayer] - equallyTotal
			}

			// 剩余分数
			surplusTotal := score % int64(winSum)

			if surplusTotal != 0 {
				startIndex, _ := utils.GetPlayerIDIndex(params.SrcPlayer, params.AllPlayers)
				firstPlayerID := utils.GetPalyerCloseFromTarget(startIndex, params.AllPlayers, params.HuPlayers)
				if firstPlayerID != 0 {
					callTransferS.Scores[firstPlayerID] = callTransferS.Scores[firstPlayerID] + surplusTotal
					callTransferS.Scores[params.SrcPlayer] = callTransferS.Scores[params.SrcPlayer] - surplusTotal
				}
			}
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
func newHuSettleInfo(params *interfaces.HuSettleParams, scoreMap map[uint64]int64, huPlayerID uint64) *majongpb.SettleInfo {
	params.SettleID = params.SettleID + 1
	return &majongpb.SettleInfo{
		Id:         params.SettleID,
		Scores:     scoreMap,
		SettleType: params.SettleType,
		HuType:     params.HuType,
		CardType:   params.CardTypes[huPlayerID],
		GenCount:   params.GenCount[huPlayerID],
	}
}

// newCallTransferSettleInfo 生成呼叫转移结算信息
func newCallTransferSettleInfo(params *interfaces.HuSettleParams) *majongpb.SettleInfo {
	params.SettleID = params.SettleID + 1
	return &majongpb.SettleInfo{
		Id:         params.SettleID,
		Scores:     make(map[uint64]int64),
		HuType:     -1,
		SettleType: majongpb.SettleType_settle_calldiver,
	}
}

func getHuTypeValue(huType majongpb.HuType) uint32 {
	huTypeValues := map[majongpb.HuType]uint32{
		majongpb.HuType_hu_dianpao:           1,
		majongpb.HuType_hu_zimo:              2,
		majongpb.HuType_hu_gangkai:           2 * 2,
		majongpb.HuType_hu_ganghoupao:        2,
		majongpb.HuType_hu_qiangganghu:       2,
		majongpb.HuType_hu_haidilao:          2,
		majongpb.HuType_hu_gangshanghaidilao: 4,
		majongpb.HuType_hu_tianhu:            32 * 2,
		majongpb.HuType_hu_dihu:              32 * 2,
	}
	return huTypeValues[huType]
}
