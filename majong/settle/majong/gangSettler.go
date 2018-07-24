package majong

import (
	"steve/common/mjoption"
	"steve/majong/interfaces"
	majongpb "steve/entity/majong"

	"github.com/Sirupsen/logrus"
)

// GangSettle 杠结算
type GangSettle struct {
}

// Settle  杠结算方法
func (gangSettle *GangSettle) Settle(params interfaces.GangSettleParams) *majongpb.SettleInfo {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":      "GangSettle",
		"settleOptionID": params.SettleOptionID,
		"gangType":       params.GangType,
		"gangPlayer":     params.GangPlayer,
		"srcPlayer":      params.SrcPlayer,
		"allPlayers":     params.AllPlayers,
		"hasHuPlayers":   params.HasHuPlayers,
		"quitPlayers":    params.QuitPlayers,
	})
	logEntry.Debugln("杠结算信息")
	// 游戏结算玩法
	settleOption := mjoption.GetSettleOption(params.SettleOptionID)
	if !settleOption.EnableGang {
		return nil
	}
	// 底数
	ante := GetDi()
	// 杠倍数
	gangValue := GetGangValue(settleOption, params.GangType)
	// 总分 (杠倍数*底分)
	total := int(gangValue) * int(ante)
	// 结算信息
	gangSettleInfo := gangSettle.newGangSettleInfo(&params, gangValue)
	if params.GangType == majongpb.GangType_gang_minggang {
		gangSettleInfo.Scores[params.GangPlayer] = int64(total)
		gangSettleInfo.Scores[params.SrcPlayer] = 0 - int64(total)
	} else if params.GangType == majongpb.GangType_gang_bugang || params.GangType == majongpb.GangType_gang_angang {
		win := 0
		for _, playerID := range params.AllPlayers {
			if playerID != params.GangPlayer && CanGangSettle(playerID, params.GiveupPlayers, params.HasHuPlayers, params.QuitPlayers, settleOption) {
				gangSettleInfo.Scores[playerID] = 0 - int64(total)
				win = win + total
			}
		}
		gangSettleInfo.Scores[params.GangPlayer] = int64(win)
	}
	return gangSettleInfo
}

// GetGangValue 获取杠对应倍数
func GetGangValue(settleOption *mjoption.SettleOption, gangType majongpb.GangType) uint32 {
	if gangType == majongpb.GangType_gang_bugang {
		return settleOption.GangValue.BuGangValue
	} else if gangType == majongpb.GangType_gang_angang {
		return settleOption.GangValue.AnGangValue
	} else if gangType == majongpb.GangType_gang_minggang {
		return settleOption.GangValue.MingGangValue
	}
	return 1
}

// newGangSettleInfo 杠的结算信息
func (gangSettle *GangSettle) newGangSettleInfo(params *interfaces.GangSettleParams, gangValue uint32) *majongpb.SettleInfo {
	params.SettleID = params.SettleID + 1
	return &majongpb.SettleInfo{
		Id:         params.SettleID,
		Scores:     make(map[uint64]int64),
		HuType:     -1,
		SettleType: gangSettle.gangType2SettleType(params.GangType),
		CardValue:  gangValue,
		HuPlayers:  params.HasHuPlayers,
	}
}

// gangType2SettleType 杠类型转SettleType
func (gangSettle *GangSettle) gangType2SettleType(gangType majongpb.GangType) majongpb.SettleType {
	return map[majongpb.GangType]majongpb.SettleType{
		majongpb.GangType_gang_angang:   majongpb.SettleType_settle_angang,
		majongpb.GangType_gang_minggang: majongpb.SettleType_settle_minggang,
		majongpb.GangType_gang_bugang:   majongpb.SettleType_settle_bugang,
	}[gangType]
}

// CanGangSettle 玩家能否参与杠结算
func CanGangSettle(playerID uint64, givePlayers, hasHuPlayers, quitPlayers []uint64, settleOption *mjoption.SettleOption) bool {
	for _, giveupPlayer := range givePlayers {
		if giveupPlayer != playerID {
			continue
		}
		return settleOption.GiveUpPlayerSettle.GiveUpPlayerGangSettle
	}
	for _, hasHupalyer := range hasHuPlayers {
		if hasHupalyer != playerID {
			continue
		}
		for _, quitPlayer := range quitPlayers {
			if quitPlayer != playerID {
				continue
			}
			return settleOption.HuQuitPlayerSettle.HuQuitPlayerGangSettle
		}
		return settleOption.HuPlayerSettle.HuPlayerGangSettle
	}
	return true
}
