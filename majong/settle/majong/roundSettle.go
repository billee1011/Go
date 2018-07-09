package majong

import (
	"steve/common/mjoption"
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
		"func_name":       "roundSettle",
		"gameId":          params.GameID,
		"flowPigPlayers":  params.FlowerPigPlayers,
		"huPlayers":       params.HuPlayers,
		"quitPlayers":     params.QuitPlayers,
		"notTinPlayers":   params.NotTingPlayers,
		"tingPlayersInfo": params.TingPlayersInfo,
	})
	logEntry.Debugln("单局结算信息")

	// 游戏结算玩法
	settleOption := GetSettleOption(int(params.GameID))
	// 结算信息
	setletInfos := make([]*majongpb.SettleInfo, 0)
	// 查花猪
	if settleOption.EnableChaDaJiao {
		flowerPigSettleInfos := roundSettle.flowerPigSettle(&params, settleOption)
		for _, flowerPigSettleInfo := range flowerPigSettleInfos {
			setletInfos = append(setletInfos, flowerPigSettleInfo)
		}
	}
	// 查大叫
	if settleOption.EnableChaDaJiao {
		yellSettleInfos := roundSettle.yellSettle(&params)
		for _, yellSettleInfo := range yellSettleInfos {
			setletInfos = append(setletInfos, yellSettleInfo)
		}
	}
	// 退税
	tuisuiIds := make([]uint64, 0)
	if settleOption.EnableTuisui {
		tuisuiIds = roundSettle.GetTuisuiIds(params)
	}
	return setletInfos, tuisuiIds
}

// 查大叫结算
// 未上听者需赔上听者最大可能番数（自摸、杠后炮、杠上开花、抢杠胡、海底捞、海底炮不参与）的牌型钱。
// 注：查大叫时，若上听者牌型中有根，则根也要未上听者包给上听者。
func (roundSettle *RoundSettle) yellSettle(params *interfaces.RoundSettleParams) []*majongpb.SettleInfo {
	//底注
	ante := GetDi()
	// 查大叫结算信息
	yellSettleInfos := make([]*majongpb.SettleInfo, 0)

	for _, noTingPlayer := range params.NotTingPlayers {
		// 关联
		groupIds := make([]uint64, 0)
		groupyellSettles := make([]*majongpb.SettleInfo, 0)
		// 听玩家结算处理
		for playerID, value := range params.TingPlayersInfo {
			settleInfoMap := map[uint64]int64{
				playerID:     int64(value) * ante,
				noTingPlayer: -(int64(value) * ante)}
			yellSettleInfo := newRoundSettleInfo(params, settleInfoMap, majongpb.SettleType_settle_yell, value)
			groupyellSettles = append(groupyellSettles, yellSettleInfo)
			groupIds = append(groupIds, yellSettleInfo.Id)
		}
		// 结算信息记录
		for _, groupyellSettle := range groupyellSettles {
			groupyellSettle.GroupId = groupIds
			yellSettleInfos = append(yellSettleInfos, groupyellSettle)
		}
	}
	return yellSettleInfos
}

// 查花猪结算
// 1.—花猪赔给未听牌玩家和胡牌玩家16*底分
// 2.—花猪赔给听牌未胡玩家（查大叫倍数+16）*底分
func (roundSettle *RoundSettle) flowerPigSettle(params *interfaces.RoundSettleParams, settleOption *mjoption.SettleOption) []*majongpb.SettleInfo {
	//底注
	ante := GetDi()
	// 查花猪信息
	flowwePigSettleInfos := make([]*majongpb.SettleInfo, 0)
	for _, flowerPig := range params.FlowerPigPlayers {
		if !roundSettle.canRoundSettle(flowerPig, params.GiveupPlayers, params.HuPlayers, params.QuitPlayers, settleOption) {
			continue
		}
		// 关联
		groupIds := make([]uint64, 0)
		groupflowerSettles := make([]*majongpb.SettleInfo, 0)
		// 胡玩家结算处理
		for j := 0; j < len(params.HuPlayers); j++ {
			if !roundSettle.canRoundSettle(params.HuPlayers[j], params.GiveupPlayers, params.HuPlayers, params.QuitPlayers, settleOption) {
				continue
			}
			settleInfoMap := map[uint64]int64{
				params.HuPlayers[j]: ante * 16,
				flowerPig:           -(ante * 16),
			}
			flowerSettleInfo := newRoundSettleInfo(params, settleInfoMap, majongpb.SettleType_settle_flowerpig, 16)
			groupflowerSettles = append(groupflowerSettles, flowerSettleInfo)
			groupIds = append(groupIds, flowerSettleInfo.Id)
		}
		// 不是花猪的未听玩家结算处理
		for n := 0; n < len(params.NotTingPlayers); n++ {
			settleInfoMap := map[uint64]int64{
				params.NotTingPlayers[n]: ante * 16,
				flowerPig:                -(ante * 16),
			}
			flowerSettleInfo := newRoundSettleInfo(params, settleInfoMap, majongpb.SettleType_settle_flowerpig, 16)
			groupflowerSettles = append(groupflowerSettles, flowerSettleInfo)
			groupIds = append(groupIds, flowerSettleInfo.Id)
		}
		// 听玩家结算处理
		for playerID, value := range params.TingPlayersInfo {
			settleInfoMap := map[uint64]int64{
				playerID:  (16 + value) * ante,
				flowerPig: 0 - ((16 + value) * ante),
			}
			flowerSettleInfo := newRoundSettleInfo(params, settleInfoMap, majongpb.SettleType_settle_flowerpig, 16+value)
			groupflowerSettles = append(groupflowerSettles, flowerSettleInfo)
			groupIds = append(groupIds, flowerSettleInfo.Id)
		}
		for _, groupyellSettle := range groupflowerSettles {
			groupyellSettle.GroupId = groupIds
			flowwePigSettleInfos = append(flowwePigSettleInfos, groupyellSettle)
		}
	}
	return flowwePigSettleInfos
}

// GetTuisuiIds 退稅ID
func (roundSettle *RoundSettle) GetTuisuiIds(params interfaces.RoundSettleParams) []uint64 {
	tuiSuiiIds := make([]uint64, 0)
	gangSettleType := map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_angang:   true,
		majongpb.SettleType_settle_minggang: true,
		majongpb.SettleType_settle_bugang:   true,
	}
	tuiSuiPlayers := make([]uint64, 0)
	for _, notingPlayer := range params.NotTingPlayers {
		tuiSuiPlayers = append(tuiSuiPlayers, notingPlayer)
	}
	for _, flowerPigPlayer := range params.FlowerPigPlayers {
		tuiSuiPlayers = append(tuiSuiPlayers, flowerPigPlayer)
	}
	for _, tuiSuiPlayer := range tuiSuiPlayers {
		for _, s := range params.SettleInfos {
			if gangSettleType[s.SettleType] {
				score := s.Scores[tuiSuiPlayer]
				if score > 0 && !s.CallTransfer {
					tuiSuiiIds = append(tuiSuiiIds, s.Id)
				}
			}
		}
	}
	return tuiSuiiIds
}

// newRoundSettleInfo 初始化生成一条新的结算信息
func newRoundSettleInfo(params *interfaces.RoundSettleParams, scoreMap map[uint64]int64, settleType majongpb.SettleType, cardValue int64) *majongpb.SettleInfo {
	params.SettleID = params.SettleID + 1

	return &majongpb.SettleInfo{
		Id:         params.SettleID,
		Scores:     scoreMap,
		HuType:     -1,
		SettleType: settleType,
		CardValue:  uint32(cardValue),
	}
}

// canRoundSettle 玩家能否参与总结算
func (roundSettle *RoundSettle) canRoundSettle(playerID uint64, givePlayers, hasHuPlayers, quitPlayers []uint64, settleOption *mjoption.SettleOption) bool {
	for _, giveupPlayer := range givePlayers {
		if giveupPlayer != playerID {
			break
		}
		return settleOption.GiveUpPlayerSettle.GiveUpPlayerRoundSettle
	}
	for _, hasHupalyer := range hasHuPlayers {
		if hasHupalyer != playerID {
			break
		}
		for _, quitPlayer := range quitPlayers {
			if quitPlayer != playerID {
				break
			}
			return settleOption.HuQuitPlayerSettle.HuQuitPlayerRoundSettle
		}
		return settleOption.HuPlayerSettle.HuPlayerRoundSettle
	}
	return true
}
