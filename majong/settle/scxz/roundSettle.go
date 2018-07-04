package scxz

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// RoundSettle 单局结算
type RoundSettle struct {
}

// Settle  单局结算方法 操作优先级：查花猪＞查大叫＞退税
// 查花猪:  胡牌玩家若已退出（不参与查花猪逻辑）
// 查大叫:  胡牌玩家（不参与查大叫逻辑）
// 退税:   胡牌玩家若已退出（不参与退税逻辑）
func (roundSettle *RoundSettle) Settle(params interfaces.RoundSettleParams) ([]*majongpb.SettleInfo, []uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"name":            "roundSettle",
		"flowPigPlayers":  params.FlowerPigPlayers,
		"huPlayers":       params.HuPlayers,
		"notTinPlayers":   params.NotTingPlayers,
		"giveupPlayers":   params.GiveupPlayers,
		"tingPlayersInfo": params.TingPlayersInfo,
	})
	setletInfos := make([]*majongpb.SettleInfo, 0)

	// 查花猪
	flowerPigSettleInfos := roundSettle.flowerPigSettle(&params)
	for _, flowerPigSettleInfo := range flowerPigSettleInfos {
		setletInfos = append(setletInfos, flowerPigSettleInfo)
	}

	// 查大叫
	yellSettleInfos := roundSettle.yellSettle(&params)
	for _, yellSettleInfo := range yellSettleInfos {
		setletInfos = append(setletInfos, yellSettleInfo)
	}
	// 退税
	taxRebeatIds := roundSettle.GetTaxRebeatIds(params)
	logEntry.Infoln("单局结算")
	return setletInfos, taxRebeatIds
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
		if isGiveUpPlayer(noTingPlayer, params.GiveupPlayers) {
			continue
		}
		// 关联
		groupIds := make([]uint64, 0)
		groupyellSettles := make([]*majongpb.SettleInfo, 0)
		// 听玩家结算处理
		for playerID, value := range params.TingPlayersInfo {
			if isGiveUpPlayer(playerID, params.GiveupPlayers) {
				continue
			}
			settleInfoMap := map[uint64]int64{
				playerID:     int64(value) * ante,
				noTingPlayer: -(int64(value) * ante)}
			params.SettleID = params.SettleID + 1
			yellSettleInfo := newRoundSettleInfo(params.SettleID, settleInfoMap, -1, majongpb.SettleType_settle_yell, value)
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
func (roundSettle *RoundSettle) flowerPigSettle(params *interfaces.RoundSettleParams) []*majongpb.SettleInfo {
	//底注
	ante := GetDi()
	// 查花猪信息
	flowwePigSettleInfos := make([]*majongpb.SettleInfo, 0)
	for _, flowerPig := range params.FlowerPigPlayers {
		if isGiveUpPlayer(flowerPig, params.GiveupPlayers) {
			continue
		}
		// 关联
		groupIds := make([]uint64, 0)
		groupflowerSettles := make([]*majongpb.SettleInfo, 0)
		// 胡玩家结算处理
		for j := 0; j < len(params.HuPlayers); j++ {
			if isGiveUpPlayer(params.HuPlayers[j], params.GiveupPlayers) {
				continue
			}
			settleInfoMap := map[uint64]int64{
				params.HuPlayers[j]: ante * 16,
				flowerPig:           -(ante * 16),
			}
			params.SettleID = params.SettleID + 1
			flowerSettleInfo := newRoundSettleInfo(params.SettleID, settleInfoMap, -1, majongpb.SettleType_settle_flowerpig, 16)
			groupflowerSettles = append(groupflowerSettles, flowerSettleInfo)
			groupIds = append(groupIds, flowerSettleInfo.Id)
		}
		// 不是花猪的未听玩家结算处理
		for n := 0; n < len(params.NotTingPlayers); n++ {
			if isGiveUpPlayer(params.NotTingPlayers[n], params.GiveupPlayers) {
				continue
			}
			settleInfoMap := map[uint64]int64{
				params.NotTingPlayers[n]: ante * 16,
				flowerPig:                -(ante * 16),
			}
			params.SettleID = params.SettleID + 1
			flowerSettleInfo := newRoundSettleInfo(params.SettleID, settleInfoMap, -1, majongpb.SettleType_settle_flowerpig, 16)
			groupflowerSettles = append(groupflowerSettles, flowerSettleInfo)
			groupIds = append(groupIds, flowerSettleInfo.Id)
		}
		// 听玩家结算处理
		for playerID, value := range params.TingPlayersInfo {
			if isGiveUpPlayer(playerID, params.GiveupPlayers) {
				continue
			}
			settleInfoMap := map[uint64]int64{
				playerID:  (int64(16+value) * ante),
				flowerPig: -(int64(16+value) * ante),
			}
			params.SettleID = params.SettleID + 1
			flowerSettleInfo := newRoundSettleInfo(params.SettleID, settleInfoMap, -1, majongpb.SettleType_settle_flowerpig, 16+value)
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

// GetTaxRebeatIds 退稅ID
func (roundSettle *RoundSettle) GetTaxRebeatIds(params interfaces.RoundSettleParams) []uint64 {
	taxRebeatIds := make([]uint64, 0)
	gangSettleType := map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_angang:   true,
		majongpb.SettleType_settle_minggang: true,
		majongpb.SettleType_settle_bugang:   true,
	}
	taxRebeatPlayers := make([]uint64, 0)
	for _, notingPlayer := range params.NotTingPlayers {
		taxRebeatPlayers = append(taxRebeatPlayers, notingPlayer)
	}
	for _, flowerPigPlayer := range params.FlowerPigPlayers {
		taxRebeatPlayers = append(taxRebeatPlayers, flowerPigPlayer)
	}
	for _, taxRebeatPlayer := range taxRebeatPlayers {
		if isGiveUpPlayer(taxRebeatPlayer, params.GiveupPlayers) {
			continue
		}
		for _, sInfo := range params.SettleInfos {
			if gangSettleType[sInfo.SettleType] == true {
				score := sInfo.Scores[taxRebeatPlayer]
				if score > 0 && !sInfo.CallTransfer {
					taxRebeatIds = append(taxRebeatIds, sInfo.Id)
				}
			}

		}
	}
	return taxRebeatIds
}

// newRoundSettleInfo 初始化生成一条新的结算信息
func newRoundSettleInfo(id uint64, scoreMap map[uint64]int64,
	huType majongpb.HuType, settleType majongpb.SettleType, cardValue int64) *majongpb.SettleInfo {
	return &majongpb.SettleInfo{
		Id:         id,
		Scores:     scoreMap,
		HuType:     huType,
		SettleType: settleType,
		CardValue:  uint32(cardValue),
	}
}

// GetDi 获取底注
func GetDi() int64 {
	//return r.Option.(*pb.Option_SiChuangXueLiu).Di
	return 1
}

func isGiveUpPlayer(playerID uint64, giveupPlayers []uint64) bool {
	for _, giveupPlayer := range giveupPlayers {
		if giveupPlayer == playerID {
			return true
		}
	}
	return false
}
