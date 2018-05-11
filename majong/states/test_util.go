package states

import (
	"fmt"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

//FmtPlayerInfo 打印玩家信息
func FmtPlayerInfo(player *majongpb.Player) logrus.Fields {
	fields := logrus.Fields{
		"玩家ID":      player.GetPalyerId(),
		"手牌":        FmtMajongpbCards(player.HandCards),
		"问询下可以有的操作": player.PossibleActions,
		"杠过的牌":      FmtGangCards(player.GangCards),
		"胡过的牌":      FmtHuCards(player.HuCards),
		"碰过的牌":      FmtPengCards(player.PengCards),
		"出过的牌":      FmtMajongpbCards(player.OutCards),
	}
	// results := ""
	// results += fmt.Sprintln("playerInfo:")
	// results += fmt.Sprintf("playerID:%v\n", player.GetPalyerId())
	// results += fmt.Sprintf("handCards:%s\n", FmtMajongpbCards(player.HandCards))
	// results += fmt.Sprintf("possibleActions: %v\n", player.PossibleActions)
	// results += fmt.Sprintf("gangCard: %v\n", FmtGangCards(player.GangCards))
	// results += fmt.Sprintf("huCard: %v\n", FmtHuCards(player.HuCards))
	// results += fmt.Sprintf("pengCard: %v\n", FmtPengCards(player.PengCards))
	// results += fmt.Sprintf("outCard: %v\n", FmtMajongpbCards(player.OutCards))

	return fields
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
		result += fmt.Sprintf("杠的类型:%v ", gangCard.Type.String())
		result += fmt.Sprintf("杠的牌:%v%v ", gangCard.Card.Point, getColor(gangCard.Card.Color))
		result += fmt.Sprintf("来自玩家:%v ", gangCard.SrcPlayer)
	}
	return result
}

//FmtPengCards 打印pengCards
func FmtPengCards(pengCards []*majongpb.PengCard) string {
	result := ""
	for _, pengCard := range pengCards {
		result += fmt.Sprintf("碰的牌:%v%v ", pengCard.Card.Point, getColor(pengCard.Card.Color))
		result += fmt.Sprintf("来自玩家:%v ", pengCard.SrcPlayer)
	}
	return result
}

//FmtHuCards 打印hucards
func FmtHuCards(huCards []*majongpb.HuCard) string {
	result := ""
	for _, huCard := range huCards {
		result += fmt.Sprintf("胡的类型:%v ", huCard.Type.String())
		result += fmt.Sprintf("胡的牌:%v%v ", huCard.Card.Point, getColor(huCard.Card.Color))
		result += fmt.Sprintf("来自玩家:%v ", huCard.SrcPlayer)
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
