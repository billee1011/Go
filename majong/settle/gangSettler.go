package settle

import (
	"steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// GangSettle 杠的结算
type GangSettle struct {
}

// SettleGang  杠立即结算,相关玩家账单 operator 操作者 ,lastActionPlayer 上次操作的人, settleType 结算类型， huType 胡牌类型
func (gangSettle *GangSettle) SettleGang(context *majong.MajongContext, operator *majong.Player, lastPlayer *majong.Player, settleType majong.SettleType) *majong.SettleInfo {
	entry := logrus.WithFields(logrus.Fields{
		"name":       "SettleGang",
		"operator":   operator.PalyerId,
		"lastPlayer": lastPlayer.PalyerId,
		"settleType": settleType,
	})
	//TODO 底数
	ante := GetDi()
	// 杠倍数
	gangScore := getGangScore(settleType)

	baseScore := gangScore * int(ante)

	// 结算信息
	settleInfo := NewSettleInfo(context, settleType)

	for _, player := range context.Players {
		if settleType == majong.SettleType_settle_bugang || settleType == majong.SettleType_settle_angang { // 补杠||暗杠
			if operator.PalyerId == player.PalyerId {
				settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] + int64(baseScore)
			} else {
				settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] - int64(baseScore)
			}
		} else if settleType == majong.SettleType_settle_mingang { // 明杠
			if operator.PalyerId == player.PalyerId {
				settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] + int64(baseScore)
			} else if lastPlayer.PalyerId == player.PalyerId {
				settleInfo.Scores[player.PalyerId] = settleInfo.Scores[player.PalyerId] - int64(baseScore)
			} else {
				settleInfo.Scores[player.PalyerId] = 0
			}
		}
	}
	context.SettleInfos = append(context.SettleInfos, settleInfo)
	entry.Info("杠结算")
	return settleInfo
}

// getGangScore 获取杠对应分数
func getGangScore(settleType majong.SettleType) int {
	if settleType == majong.SettleType_settle_bugang {
		return 1
	} else if settleType == majong.SettleType_settle_angang {
		return 2
	} else if settleType == majong.SettleType_settle_mingang {
		return 2
	}
	return 0
}
