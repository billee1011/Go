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
		"name":          "HuSettle",
		"winnersID":     params.HuPlayers,
		"settleType":    params.SettleType,
		"cardTypeVAlue": params.CardTypeValue,
	})

	settleInfos := make([]*majongpb.SettleInfo, 0)
	for i := 0; i < len(params.HuPlayers); i++ {
		huSettleInfo := NewSettleInfo(params.SettleID)
		//底数
		ante := GetDi()
		// 总分
		total := int64(params.CardTypeValue) * ante
		// 胡结算信息
		settleInfos := make([]*majongpb.SettleInfo, 0)
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
		huSettleInfo.CardValue = int32(params.CardTypeValue)
		huSettleInfo.CardType = params.CardTypes
		settleInfos = append(settleInfos, huSettleInfo)

	}
	entry.Info("胡结算")
	return settleInfos
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
