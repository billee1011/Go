package settle

import (
	"steve/majong/interfaces"
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

	//底数
	ante := GetDi()

	win := int64(0)
	if params.SettleType == majongpb.SettleType_settle_zimo {
		huSettleInfo := NewSettleInfo(params.SettleID)
		huSettleInfo.HuType = params.HuType
		huSettleInfo.SettleType = params.SettleType
		// 倍数
		value := int64(params.CardValues[params.HuPlayers[0]]) * int64(getHuTypeValue(params.HuType))
		// 总分
		total := value * ante
		for _, playerID := range params.AllPlayers {
			if playerID != params.HuPlayers[0] {
				huSettleInfo.Scores[playerID] = 0 - total
				win = win + total
			}
		}
		huSettleInfo.Scores[params.HuPlayers[0]] = huSettleInfo.Scores[params.HuPlayers[0]] + win
		huSettleInfo.CardType = params.CardTypes[params.HuPlayers[0]]
		huSettleInfo.GenCount = params.GenCount[params.HuPlayers[0]]
		huSettleInfo.CardValue = uint32(value)
		settleInfos = append(settleInfos, huSettleInfo)
	} else if params.SettleType == majongpb.SettleType_settle_dianpao {
		for _, playerID := range params.HuPlayers {
			huSettleInfo := NewSettleInfo(params.SettleID)
			params.SettleID++
			huSettleInfo.HuType = params.HuType
			huSettleInfo.SettleType = params.SettleType

			value := int64(params.CardValues[playerID]) * int64(getHuTypeValue(params.HuType))
			// 总分
			total := value * ante

			huSettleInfo.Scores[playerID] = total
			huSettleInfo.Scores[params.SrcPlayer] = -total
			huSettleInfo.CardType = params.CardTypes[playerID]
			huSettleInfo.GenCount = params.GenCount[playerID]
			huSettleInfo.CardValue = uint32(value)
			settleInfos = append(settleInfos, huSettleInfo)
		}
	}
	if params.HuType == majongpb.SettleHuType_settle_hu_ganghoupao { // 需呼叫转移
		callTransferS := callTransferSettle(params)
		callTransferS.Id++
		settleInfos = append(settleInfos, callTransferS)
	}
	entry.Info("胡结算")
	return settleInfos
}

func callTransferSettle(params interfaces.HuSettleParams) *majongpb.SettleInfo {
	callTransferS := NewSettleInfo(params.SettleID)
	callTransferS.SettleType = majongpb.SettleType_settle_calldiver

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
				startIndex := getPlayerIndex(params.SrcPlayer, params.AllPlayers)
				firstPlayerID := getPalyerCloseIndex(startIndex, params.AllPlayers, params.HuPlayers)
				if firstPlayerID != 0 {
					callTransferS.Scores[firstPlayerID] = surplusTotal
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

// NewSettleInfo 初始化生成一条新的结算信息
func NewSettleInfo(settleID uint64) *majongpb.SettleInfo {
	return &majongpb.SettleInfo{
		Id:     settleID + 1,
		Scores: make(map[uint64]int64),
		HuType: -1,
	}
}

func getPlayerIndex(playerID uint64, allPlayer []uint64) int {
	for index, player := range allPlayer {
		if playerID == player {
			return index
		}
	}
	return -1
}

func getPalyerCloseIndex(index int, allPlayer, huPlayers []uint64) uint64 {
	for i := 0; i <= len(allPlayer); i++ {
		nextIndex := (index + i) % len(allPlayer)
		for _, huPlayer := range huPlayers {
			index := getPlayerIndex(huPlayer, allPlayer)
			if index == nextIndex {
				return huPlayer
			}
		}
	}
	return 0
}

func getHuTypeValue(huType majongpb.SettleHuType) uint32 {
	huTypeValues := map[majongpb.SettleHuType]uint32{
		majongpb.SettleHuType_settle_hu_noramaldianpao:    1,
		majongpb.SettleHuType_settle_hu_zimo:              2,
		majongpb.SettleHuType_settle_hu_gangkai:           2 * 2,
		majongpb.SettleHuType_settle_hu_ganghoupao:        2,
		majongpb.SettleHuType_settle_hu_qiangganghu:       2,
		majongpb.SettleHuType_settle_hu_haidilao:          2,
		majongpb.SettleHuType_settle_hu_gangshanghaidilao: 4,
		majongpb.SettleHuType_settle_hu_tianhu:            32 * 2,
		majongpb.SettleHuType_settle_hu_dihu:              32 * 2,
	}
	return huTypeValues[huType]
}
