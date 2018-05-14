package cardtype

import (
	"fmt"
	"steve/majong/interfaces"
	"testing"
	// "github.com/stretchr/testify/assert"
	majongpb "steve/server_pb/majong"
	"steve/majong/utils"
)

func TestCalculate(t *testing.T) {
	handUtilCards := []utils.Card{11,11,12,12,13,13,14,14,15,15,16,16,17,17}
	handCards,_ := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	pengCards := make([]*majongpb.Card, 0)
	gangCards := make([]*majongpb.Card, 0)
	huCard := new(majongpb.Card)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   huCard,
		GameID:   0,
	}
	cardTypes, genCount := new(cardTypeCalculator).Calculate(playerParams)
	fmt.Println(cardTypes, genCount)
}

func TestCardTypeValue(t *testing.T) {
	cardTypes := []interfaces.CardType{}
	genCount := 2
	valuer := new(cardTypeCalculator).CardTypeValue(cardTypes, genCount)
	fmt.Println(valuer)
}
