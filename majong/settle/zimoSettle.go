package settle

import (
	"steve/majong/settle/fan"
	"steve/majong/utils"
	"steve/server_pb/majong"
	"strconv"
	"strings"

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
	settleInfo := NewSettleInfo(context, settleType, operatorID)

	winner := utils.GetPlayerByID(context.Players, operatorID)

	fansMap := make(map[string]uint32)
	gen := uint32(0)
	for i := 0; i < len(fan.ScxlFan); i++ {
		if fan.ScxlFan[i].Condition(*context, huType, winner) {
			fansMap[fan.ScxlFan[i].GetFanName()] = fan.ScxlFan[i].GetFanValue()
		}
	}
	fansMap, gen = scxlFanMutex(fansMap, fan.GetGenCount(winner))

	fanValues := 1
	fanNames := make([]string, 0)
	if gen != 0 {
		fanNames = append(fanNames, strconv.Itoa(int(gen))+"根")
	}
	for name, value := range fansMap {
		if value != 0 {
			fanValues = fanValues * int(value)
			fanNames = append(fanNames, name)
		}
	}
	// 自摸2倍
	fanValues = fanValues * 2
	//底数
	ante := GetDi()
	total := int64(fanValues) * (1 << gen) * ante

	for _, player := range context.Players {
		if winner.PalyerId == player.PalyerId {
			settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] + total
		} else {
			settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] - total
		}
	}
	settleInfo.Type = strings.Join(fanNames, ",")
	settleInfo.Times = int32(fanValues)
	entry.Info("自摸结算")
	return settleInfo, nil
}
