package settle

import (
	"steve/majong/settle/fan"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
	"strconv"
)

// YellSettle 查大叫
type YellSettle struct {
}

// SettleYell 查大叫结算
func (yellSettle *YellSettle) SettleYell(context *majongpb.MajongContext) []*majongpb.SettleInfo {
	// 未听玩家
	noTingPlayers := getNoTingPlayers(context)
	// 听牌玩家信息
	tingPlayersInfo, err := getTingPlayerInfo(context)
	if err != nil {
		return nil
	}
	// 查大叫结算信息
	settleInfos := make([]*majongpb.SettleInfo, 0)
	if len(noTingPlayers) > 0 {
		for _, noTingPlayer := range noTingPlayers {
			settleInfoMap := make(map[uint64]int64)
			win := int64(0)
			lose := int64(0)
			// 听玩家结算处理
			for playerID, fanMultiple := range tingPlayersInfo {
				tingPlayer := utils.GetPlayerByID(context.GetPlayers(), playerID)
				win = fanMultiple * GetDi()
				lose = lose - win
				settleInfoMap[tingPlayer.PalyerId] = win
			}
			// 结算信息记录
			if len(settleInfoMap) > 0 {
				settleInfoMap[noTingPlayer.PalyerId] = lose
				for _, player := range context.Players {
					if _, ok := settleInfoMap[player.PalyerId]; !ok {
						settleInfoMap[player.PalyerId] = 0
					}
				}
				yellSettleInfo := NewSettleInfo(context, majongpb.SettleType_settle_yell)
				yellSettleInfo.Scores = settleInfoMap
				settleInfos = append(settleInfos, yellSettleInfo)
				context.SettleInfos = append(context.SettleInfos, yellSettleInfo)
			}
		}
	}
	return settleInfos
}

// getNoTingPlayers 所有未听玩家，不包括花猪，因为查花猪包括了查大叫，所以未听玩家，中是花猪的，都不用再进行查大叫
func getNoTingPlayers(context *majongpb.MajongContext) []*majongpb.Player {
	noTingPlayers := make([]*majongpb.Player, 0)
	for i := 0; i < len(context.Players); i++ {
		// 胡过的不算
		if len(context.Players[i].HuCards) > 0 {
			continue
		}
		// 查听
		tingCards, err := utils.GetTingCards(context.Players[i].HandCards)
		if err != nil {
			return nil
		}
		// 不能听
		if len(tingCards) == 0 && !isFlowerPig(context.Players[i]) {
			noTingPlayers = append(noTingPlayers, context.Players[i])
		}
	}
	return noTingPlayers
}

// getTingPlayerInfo 判断玩家是否能听,和返回能听玩家的最大倍数
// 未上听者需赔上听者最大可能番数（自摸、杠后炮、杠上开花、抢杠胡、海底捞、海底炮不参与）的牌型钱。注：查大叫时，若上听者牌型中有根，则根也要未上听者包给上听者。
func getTingPlayerInfo(context *majongpb.MajongContext) (map[uint64]int64, error) {
	players := context.Players
	tingPlayers := make(map[uint64]int64, 0)
	for i := 0; i < len(players); i++ {
		// 胡过的不算
		if len(players[i].HuCards) > 0 {
			continue
		}
		handCardSum := len(players[i].HandCards)
		var maxMulti int64
		//只差1张牌就能胡，并且玩家手牌不存在花牌
		if handCardSum%3 == 1 && !hasDingQueCard(players[i].HandCards, players[i].DingqueColor) {
			tingCards, err := utils.GetTingCards(players[i].HandCards)
			if err != nil {
				return nil, err
			}
			for j := 0; j < len(tingCards); j++ {
				// 获取最大番型*根数
				multiple := getMulti(context, players[i], tingCards[0])
				if maxMulti < int64(multiple) {
					maxMulti = int64(multiple)
				}
			}
			if len(tingCards) != 0 {
				tingPlayers[players[i].GetPalyerId()] = maxMulti
			}
		}
	}
	return tingPlayers, nil
}

// getMulti 返回番的倍数*gen的倍数
func getMulti(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) uint32 {
	huCard := &majongpb.HuCard{
		Card: card,
	}
	fansMap := make(map[string]uint32)
	gen := uint32(0)
	for i := 0; i < len(fan.ScxlCardFan); i++ {
		if fan.ScxlCardFan[i].Condition(*context, majongpb.HuType_hu_dianpao, *huCard, player) {
			fansMap[fan.ScxlCardFan[i].GetFanName()] = fan.ScxlCardFan[i].GetFanValue()
		}
	}
	fansMap, gen = scxlFanMutex(fansMap, fan.GetGenCount(player, *huCard))

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
	return uint32(fanValues) * (1 << gen)
}

//hasDingQueCard 检查牌里面是否含有定缺的牌
func hasDingQueCard(cards []*majongpb.Card, color majongpb.CardColor) bool {
	for _, card := range cards {
		if card.Color == color {
			return true
		}
	}
	return false
}
