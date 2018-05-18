package settle

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// GangSettle 杠结算
type GangSettle struct {
}

// Settle  杠结算方法
func (gangSettle *GangSettle) Settle(params interfaces.GangSettleParams) *majongpb.SettleInfo {
	entry := logrus.WithFields(logrus.Fields{
		"name":       "GangSettle",
		"gangType":   params.GangType,
		"gangPlayer": params.GangPlayer,
		"srcPlayer":  params.SrcPlayer,
	})
	// 底数
	ante := GetDi()
	// 杠倍数
	gangScore := getGangScore(params.GangType)

	total := gangScore * int(ante)

	// 结算信息
	gangSettleInfo := NewSettleInfo(params.SettleID)
	gangSettleInfo.SettleType = majongpb.SettleType_settle_gang
	if params.GangType == majongpb.GangType_gang_minggang {
		gangSettleInfo.Scores[params.GangPlayer] = int64(total)
		gangSettleInfo.Scores[params.SrcPlayer] = 0 - int64(total)
	} else if params.GangType == majongpb.GangType_gang_bugang || params.GangType == majongpb.GangType_gang_angang {
		win := 0
		for _, playerID := range params.AllPlayers {
			if playerID != params.GangPlayer {
				gangSettleInfo.Scores[playerID] = 0 - int64(total)
				win = win + total
			}
		}
		gangSettleInfo.Scores[params.GangPlayer] = int64(win)
	}
	gangSettleInfo.CardValue = uint32(gangScore)
	entry.Info("杠结算")
	return gangSettleInfo
}

// getGangScore 获取杠对应分数
func getGangScore(gangType majongpb.GangType) int {
	if gangType == majongpb.GangType_gang_bugang {
		return 1
	} else if gangType == majongpb.GangType_gang_angang {
		return 2
	} else if gangType == majongpb.GangType_gang_minggang {
		return 2
	}
	return 1
}
