package scxl

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// RoundSettle 单局结算
type RoundSettle struct {
}

// Settle  单局结算方法 操作优先级：查花猪＞查大叫＞退税
func (roundSettle *RoundSettle) Settle(params interfaces.RoundSettleParams) ([]*majongpb.SettleInfo, []uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"name":            "roundSettle",
		"flowPigPlayers":  params.FlowerPigPlayers,
		"huPlayers":       params.HuPlayers,
		"notTinPlayers":   params.NotTingPlayers,
		"tingPlayersInfo": params.TingPlayersInfo,
	})
	setletInfos := make([]*majongpb.SettleInfo, 0)

	// 查花猪
	flowerPigSettleInfos, fSettleID := flowerPigSettle(params)
	params.SettleID = fSettleID
	if flowerPigSettleInfos != nil && len(flowerPigSettleInfos) > 0 {
		for _, s := range flowerPigSettleInfos {
			setletInfos = append(setletInfos, s)
		}
	}
	// 查大叫
	yellSettleInfos, ySettleID := yellSettle(params)
	params.SettleID = ySettleID
	if yellSettleInfos != nil && len(yellSettleInfos) > 0 {
		for _, s := range yellSettleInfos {
			setletInfos = append(setletInfos, s)
		}
	}
	// 退税
	taxRebeatIds := taxRebeat(params)
	logEntry.Infoln("单局结算")
	return setletInfos, taxRebeatIds
}

// 查大叫结算
// 未上听者需赔上听者最大可能番数（自摸、杠后炮、杠上开花、抢杠胡、海底捞、海底炮不参与）的牌型钱。
// 注：查大叫时，若上听者牌型中有根，则根也要未上听者包给上听者。
func yellSettle(params interfaces.RoundSettleParams) ([]*majongpb.SettleInfo, uint64) {
	//底注
	ante := GetDi()
	// 查大叫结算信息
	yellSettleInfos := make([]*majongpb.SettleInfo, 0)
	if len(params.NotTingPlayers) > 0 {
		for _, noTingPlayer := range params.NotTingPlayers {
			settleInfoMap := make(map[uint64]int64)
			win := int64(0)
			lose := int64(0)
			// 听玩家结算处理
			for playerID, value := range params.TingPlayersInfo {
				win = int64(value) * ante
				lose = lose - win
				settleInfoMap[playerID] = win
			}
			// 结算信息记录
			if len(settleInfoMap) > 0 {
				settleInfoMap[noTingPlayer] = lose
				yellSettleInfo := new(majongpb.SettleInfo)
				params.SettleID = params.SettleID + 1
				yellSettleInfo, params = newRoundSettleInfo(params, settleInfoMap, -1, majongpb.SettleType_settle_yell)
				yellSettleInfos = append(yellSettleInfos, yellSettleInfo)
			}
		}
	}
	return yellSettleInfos, params.SettleID
}

// 查花猪结算
// 1.—花猪赔给未听牌玩家和胡牌玩家16*底分
// 2.—花猪赔给听牌未胡玩家（查大叫倍数+16）*底分
func flowerPigSettle(params interfaces.RoundSettleParams) ([]*majongpb.SettleInfo, uint64) {
	//底注
	ante := GetDi()
	// 查花猪信息
	flowwePigSettleInfos := make([]*majongpb.SettleInfo, 0)
	for _, flowerPig := range params.FlowerPigPlayers {
		lose := int64(0)
		settleInfoMap := make(map[uint64]int64)
		// 胡玩家结算处理
		for j := 0; j < len(params.HuPlayers); j++ {
			settleInfoMap[params.HuPlayers[j]] = ante * 16
			lose = lose - (ante * 16)
		}
		// 不是花猪的未听玩家结算处理
		for n := 0; n < len(params.NotTingPlayers); n++ {
			settleInfoMap[params.NotTingPlayers[n]] = ante * 16
			lose = lose - (ante * 16)
		}
		// 听玩家结算处理
		for playerID, value := range params.TingPlayersInfo {
			win := int64(16+value) * ante
			settleInfoMap[playerID] = win
			lose = lose - win
		}
		// 查花猪玩家结算信息
		if len(settleInfoMap) > 0 {
			settleInfoMap[flowerPig] = lose
			flowerSettleInfo := new(majongpb.SettleInfo)
			params.SettleID = params.SettleID + 1
			flowerSettleInfo, params = newRoundSettleInfo(params, settleInfoMap, -1, majongpb.SettleType_settle_flowerpig)
			flowwePigSettleInfos = append(flowwePigSettleInfos, flowerSettleInfo)
		}
	}
	return flowwePigSettleInfos, params.SettleID
}

// 退稅结算
func taxRebeat(params interfaces.RoundSettleParams) []uint64 {
	taxRebeatIds := make([]uint64, 0)
	for _, notTingPlayer := range params.NotTingPlayers {
		for _, settleInfo := range params.SettleInfos {
			if settleInfo.SettleType == majongpb.SettleType_settle_angang || settleInfo.SettleType == majongpb.SettleType_settle_minggang || settleInfo.SettleType == majongpb.SettleType_settle_bugang {
				if (settleInfo.Scores[notTingPlayer] > 0) && !settleInfo.CallTransfer {
					taxRebeatIds = append(taxRebeatIds, settleInfo.Id)
				}
			}

		}
	}
	for _, flowerPigPlayer := range params.FlowerPigPlayers {
		for _, settleInfo := range params.SettleInfos {
			if settleInfo.SettleType == majongpb.SettleType_settle_angang || settleInfo.SettleType == majongpb.SettleType_settle_minggang || settleInfo.SettleType == majongpb.SettleType_settle_bugang {
				if (settleInfo.Scores[flowerPigPlayer] > 0) && !settleInfo.CallTransfer {
					taxRebeatIds = append(taxRebeatIds, settleInfo.Id)
				}
			}
		}
	}
	return taxRebeatIds
}

// newRoundSettleInfo 初始化生成一条新的结算信息
func newRoundSettleInfo(params interfaces.RoundSettleParams, scoreMap map[uint64]int64,
	huType majongpb.HuType, settleType majongpb.SettleType) (*majongpb.SettleInfo, interfaces.RoundSettleParams) {
	settleInfo := &majongpb.SettleInfo{
		Id:         params.SettleID,
		Scores:     scoreMap,
		HuType:     huType,
		SettleType: settleType,
	}
	return settleInfo, params
}
