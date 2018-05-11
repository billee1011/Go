package settle

import (
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// FlowerPigSettle 查花猪
type FlowerPigSettle struct {
}

// SettleFlowerPig 查花猪结算
// 1.—花猪赔给未听牌玩家和胡牌玩家16*底分
// 2.—花猪赔给听牌未胡玩家（查大叫倍数+16）*底分
func (flowerPigSettle *FlowerPigSettle) SettleFlowerPig(context *majongpb.MajongContext) []*majongpb.SettleInfo {
	// 花猪
	flowerPigPlayers := make([]*majongpb.Player, 0)
	// 胡过的玩家
	huPlayers := make([]*majongpb.Player, 0)
	// 未听玩家
	noTingPlayers := make([]*majongpb.Player, 0)
	// 听牌未胡玩家信息
	tingPlayersInfo := make(map[uint64]int64)
	for _, player := range context.Players {
		if isFlowerPig(player) {
			flowerPigPlayers = append(flowerPigPlayers, player)
		}
	}
	// 查花猪信息
	settleInfos := make([]*majongpb.SettleInfo, 0)
	if len(flowerPigPlayers) > 0 {
		for _, player := range context.Players {
			if len(player.HuCards) > 0 {
				huPlayers = append(huPlayers, player)
			}
			if isNoTingPlayers(player) {
				noTingPlayers = append(noTingPlayers, player)
			}
		}
		tingPlayersInfo, _ = getTingPlayerInfo(context)
	}
	// 进行查花猪结算
	for _, flowerPig := range flowerPigPlayers {
		ante := GetDi()
		lose := int64(0)
		settleInfoMap := make(map[uint64]int64)
		// 胡玩家结算处理
		for j := 0; j < len(huPlayers); j++ {
			settleInfoMap[huPlayers[j].PalyerId] = ante * 16
			lose = lose - (ante * 16)
		}
		// 不是花猪的未听玩家结算处理
		for n := 0; n < len(noTingPlayers); n++ {
			settleInfoMap[noTingPlayers[n].PalyerId] = ante * 16
			lose = lose - (ante * 16)
		}
		// 听玩家结算处理
		for playerID, multiple := range tingPlayersInfo {
			tingPlayer := utils.GetPlayerByID(context.Players, playerID)
			// 16*di + multiple*di = (16+multiple)*di
			win := (16 + multiple) * ante
			settleInfoMap[tingPlayer.PalyerId] = win
			lose = lose - win
		}
		// 查花猪玩家结算信息
		if len(settleInfoMap) > 0 {
			settleInfoMap[flowerPig.PalyerId] = lose
			flowerSettleInfo := NewSettleInfo(context, majongpb.SettleType_settle_flowerpig)
			flowerSettleInfo.Scores = settleInfoMap
			settleInfos = append(settleInfos, flowerSettleInfo)
			context.SettleInfos = append(context.SettleInfos, flowerSettleInfo)
		}
	}
	return settleInfos
}

//isFlowerPig 判断玩家是否是花猪,如果玩家从开局到牌局结束打出的牌全为定缺牌（碰、杠不影响），结束后该玩家手上还有定缺牌，此时该玩家不被查花猪
func isFlowerPig(player *majongpb.Player) bool {
	outCardDingQue := true
	for _, card := range player.OutCards {
		if card.Color != player.DingqueColor {
			outCardDingQue = false
		}
	}
	if !outCardDingQue {
		// 玩家手牌中是否存在定缺牌
		for _, card := range player.HandCards {
			if card.Color == player.DingqueColor {
				return true
			}
		}
	}
	return false
}

// isNoTingPlayers 判断玩家是否未听，不包括花猪，因为查花猪包括了查大叫，所以未听玩家，中是花猪的，都不用再进行查大叫
func isNoTingPlayers(player *majongpb.Player) bool {
	// 胡过的不算
	if len(player.HuCards) > 0 {
		return false
	}
	// 查听
	tingCards, _ := utils.GetTingCards(player.HandCards)
	// 不能听
	if len(tingCards) == 0 && !isFlowerPig(player) {
		return true
	}
	return false
}
