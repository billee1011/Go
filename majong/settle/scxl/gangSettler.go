package scxl

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
	params.SettleID = params.SettleID + 1
	gangSettleInfo := newGangSettleInfo(params)
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

// newGangSettleInfo 初始化生成一条新的杠结算信息
func newGangSettleInfo(params interfaces.GangSettleParams) *majongpb.SettleInfo {
	settleType := majongpb.SettleType(-1)
	if params.GangType == majongpb.GangType_gang_angang {
		settleType = majongpb.SettleType_settle_angang
	} else if params.GangType == majongpb.GangType_gang_minggang {
		settleType = majongpb.SettleType_settle_minggang
	} else if params.GangType == majongpb.GangType_gang_bugang {
		settleType = majongpb.SettleType_settle_bugang
	}
	return &majongpb.SettleInfo{
		Id:         params.SettleID,
		Scores:     make(map[uint64]int64),
		HuType:     -1,
		SettleType: settleType,
	}
}