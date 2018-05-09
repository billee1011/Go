package states

import (
	"fmt"
	majongpb "steve/server_pb/majong"
)

//FmtPlayerInfo 打印玩家信息
func FmtPlayerInfo(player *majongpb.Player) string {
	results := ""
	results += fmt.Sprintln("playerInfo:")
	results += fmt.Sprintf("playerID:%v\n", player.GetPalyerId())
	results += fmt.Sprintf("handCards:%s\n", FmtMajongpbCards(player.HandCards))
	results += fmt.Sprintf("possibleActions: %v\n", player.PossibleActions)
	results += fmt.Sprintf("gangCard: %v\n", FmtGangCards(player.GangCards))
	results += fmt.Sprintf("huCard: %v\n", FmtHuCards(player.HuCards))
	results += fmt.Sprintf("pengCard: %v\n", FmtPengCards(player.PengCards))
	results += fmt.Sprintf("outCard: %v\n", FmtMajongpbCards(player.OutCards))

	return results
}

//FmtMajongpbCards 打印牌组
func FmtMajongpbCards(cards []*majongpb.Card) string {
	results := ""
	for _, card := range cards {
		results += fmt.Sprintf("%v%v ", card.Point, getColor(card.Color))
	}
	return results
}

//FmtGangCards 打印gangCards
func FmtGangCards(gangCards []*majongpb.GangCard) string {
	result := ""
	for _, gangCard := range gangCards {
		result += fmt.Sprintf("Type:%v ", gangCard.Type.String())
		result += fmt.Sprintf("card:%v%v ", gangCard.Card.Point, getColor(gangCard.Card.Color))
		result += fmt.Sprintf("player:%v\n", gangCard.SrcPlayer)
	}
	return result
}

//FmtPengCards 打印pengCards
func FmtPengCards(pengCards []*majongpb.PengCard) string {
	result := ""
	for _, pengCard := range pengCards {
		result += fmt.Sprintf("card:%v%v ", pengCard.Card.Point, getColor(pengCard.Card.Color))
		result += fmt.Sprintf("player:%v\n", pengCard.SrcPlayer)
	}
	return result
}

//FmtHuCards 打印hucards
func FmtHuCards(huCards []*majongpb.HuCard) string {
	result := ""
	for _, huCard := range huCards {
		result += fmt.Sprintf("Type:%v", huCard.Type.String())
		result += fmt.Sprintf("card:%v%v ", huCard.Card.Point, getColor(huCard.Card.Color))
		result += fmt.Sprintf("player:%v\n", huCard.SrcPlayer)
	}
	return result
}

func getColor(srcColor majongpb.CardColor) string {
	if srcColor == majongpb.CardColor_ColorWan {
		return "w"
	}
	if srcColor == majongpb.CardColor_ColorTiao {
		return "t"
	}
	if srcColor == majongpb.CardColor_ColorTong {
		return "b"
	}
	return "none"
}
