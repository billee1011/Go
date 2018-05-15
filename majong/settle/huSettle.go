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
func (huSettle *HuSettle) Settle(params interfaces.HuSettleParams) *majongpb.SettleInfo {
	entry := logrus.WithFields(logrus.Fields{
		"name":       "HuSettle",
		"winnersID":  params.HuPlayers,
		"settleType": params.SettleType,
		"huType":     params.HuType,
		"cardTypes":  params.CardTypes,
		"cardValues": params.CardValues,
	})

	huSettleInfo := NewSettleInfo(params.SettleID)
	huSettleInfo.HuType = params.HuType
	for i := 0; i < len(params.HuPlayers); i++ {
		//底数
		ante := GetDi()
		// 总分
		total := int64(params.CardValues[params.HuPlayers[i]]) * ante
		win := int64(0)
		lose := int64(0)
		if params.SettleType == majongpb.SettleType_settle_zimo {
			for _, playerID := range params.AllPlayers {
				if playerID != params.HuPlayers[i] {
					huSettleInfo.Scores[playerID] = 0 - total
				} else {
					win = win + total
				}
			}
			huSettleInfo.Scores[params.HuPlayers[i]] = huSettleInfo.Scores[params.HuPlayers[i]] + win
		} else if params.SettleType == majongpb.SettleType_settle_dianpao {
			for _, playerID := range params.HuPlayers {
				huSettleInfo.Scores[playerID] = total
				lose = lose - total
			}
			huSettleInfo.Scores[params.SrcPlayer] = total
		}
		huSettleInfo.CardValue = params.CardValues[params.HuPlayers[i]]
		huSettleInfo.CardType = params.CardTypes[params.HuPlayers[i]]
		huSettleInfo.GenCount = params.GenCount[params.HuPlayers[i]]

	}
	entry.Info("胡结算")
	return huSettleInfo
}

func callTransferSettle(params interfaces.HuSettleParams) {
	// gangCard := params.GangCard
	// gangScore := getGangScore(gangCard.GetType())
	// // 赢家人数
	// winSum := len(params.HuPlayers)

	// if winSum == 1 {
	// 	win := GetDi() * int64(gangScore)
	// 	if gangCard.GetType() ==
	// } else { // 多个赢家平分杠钱

	// }
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
	}
}
