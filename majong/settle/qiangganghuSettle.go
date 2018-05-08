package settle

import (
	"steve/majong/settle/fan"
	"steve/majong/utils"
	"steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// QiangGangHuSettle 抢杠胡的结算
type QiangGangHuSettle struct {
}

// SettleQiangGangHu 抢杠胡结算	operatorIds 抢杠胡玩家 ,beiQiangGangPlayerID 被抢杠的玩家, gangCard 杠的牌, settleType 结算类型, huType 胡牌类型
func (qiangGangHuSettle *QiangGangHuSettle) SettleQiangGangHu(context *majong.MajongContext, operatorIds []uint64, beiQiangGangPlayerID uint64, gangCard *majong.Card, settleType majong.SettleType, huType majong.HuType) ([]*majong.SettleInfo, error) {
	entry := logrus.WithFields(logrus.Fields{
		"name":                 "SettleQiangGangHu",
		"operatorIds":          operatorIds,
		"beiQiangGangPlayerID": beiQiangGangPlayerID,
		"settleType":           settleType,
		"gangCard":             gangCard,
		"huType":               huType,
	})

	settleInfos := make([]*majong.SettleInfo, 0)
	for i := 0; i < len(operatorIds); i++ {
		winner := utils.GetPlayerByID(context.Players, operatorIds[i])
		fansMap := make(map[string]uint32)
		gen := uint32(0)
		for i := 0; i < len(fan.ScxlFan); i++ {
			if fan.ScxlFan[i].Condition(*context, huType, winner) {
				fansMap[fan.ScxlFan[i].GetFanName()] = fan.ScxlFan[i].GetFanValue()
			}
		}
		fansMap, gen = scxlFanMutex(fansMap, fan.GetGenCount(winner))

		fanTotal := 1
		for _, value := range fansMap {
			if value != 0 {
				fanTotal = fanTotal * int(value)
			}
		}
		//底数
		ante := GetDi()
		total := int64(fanTotal) * ante * (1 << gen)
		// 结算信息
		settleInfo := NewSettleInfo(context, settleType)
		for _, player := range context.Players {
			if winner.PalyerId == player.PalyerId {
				settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] + total
			} else if beiQiangGangPlayerID == player.PalyerId {
				settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] - total
			} else {
				settleInfo.Scores[player.PalyerId] = 0
			}
		}
		settleInfos = append(settleInfos, settleInfo)
	}

	beiQiangGangPlayer := utils.GetPlayerByID(context.Players, beiQiangGangPlayerID)
	for i := len(beiQiangGangPlayer.GangCards) - 1; i > 0; i-- {
		if utils.CardEqual(beiQiangGangPlayer.GangCards[i].Card, gangCard) {
			if beiQiangGangPlayer.GangCards[i].Type == majong.GangType_gang_bugang { // 抢杠胡成功移除被抢杠玩家补杠的钱
				for _, settleInfo := range settleInfos {
					settleInfo.Scores[beiQiangGangPlayer.PalyerId] = settleInfo.Scores[beiQiangGangPlayer.PalyerId] - 3
				}
			}
		}
	}
	entry.Info("点炮结算")
	return settleInfos, nil
}
