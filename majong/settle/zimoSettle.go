package settle

import (
	"steve/majong/settle/fan"
	"steve/majong/utils"
	"steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// ZiMoSettle 自摸的结算
type ZiMoSettle struct {
}

// SettleZiMo 自摸立即结算
func (ziMoSettle *ZiMoSettle) SettleZiMo(context *majong.MajongContext, operatorID uint64, settleType majong.SettleType, huType majong.HuType) (*majong.SettleInfo, error) {
	entry := logrus.WithFields(logrus.Fields{
		"name":       "SettleZiMo",
		"operatorId": operatorID,
		"settleType": settleType,
		"huType":     huType,
	})

	// 结算信息
	settleInfo := NewSettleInfo(context, settleType)

	winner := utils.GetPlayerByID(context.Players, operatorID)

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
	total := int64(fanTotal) * (1 << gen) * ante

	for _, player := range context.Players {
		if winner.PalyerId == player.PalyerId {
			settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] + total
		} else {
			settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] - total
		}
	}
	entry.Info("自摸结算")
	return settleInfo, nil
}
