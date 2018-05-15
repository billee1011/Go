package cardtype

import (
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

var playerParams = interfaces.CardCalcParams{
	GameID: 0,
}

func init() {
	handCards := make([]*majongpb.Card, 0)
	gangCards := make([]*majongpb.Card, 0)
	pengCards := make([]*majongpb.Card, 0)
	HuCard := new(majongpb.Card)
	playerParams = interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   playerParams.GameID,
	}
}

// 平胡
func TestCalculateAndCardTypeValuePingHu(t *testing.T) {
	handUtilCards := []utils.Card{21, 22, 23, 11, 12, 13, 24, 25, 26, 14, 15, 16, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_PingHu}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(1))
}

// 清一色
func TestCalculateAndCardTypeValueQingYiSe(t *testing.T) {
	handUtilCards := []utils.Card{21, 22, 23, 21, 22, 23, 24, 25, 26, 27, 28, 29, 24}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(24)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingYiSe}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(4))
}

// 七对
func TestCalculateAndCardTypeValueQiDui(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 11, 12, 13, 24, 25, 26, 24, 25, 26, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QiDui}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(4))
}

// 龙七对
func TestCalculateAndCardTypeValueLongQiDui(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 11, 12, 13, 24, 25, 25, 24, 24, 24, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_LongQiDui}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(8))
}

// 清七对
func TestCalculateAndCardTypeValueQingQiDui(t *testing.T) {
	handUtilCards := []utils.Card{21, 22, 23, 21, 22, 23, 24, 25, 26, 24, 25, 26, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingQiDui}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(16))
}

// 清龙七对
func TestCalculateAndCardTypeValueQingLongQiDui(t *testing.T) {
	handUtilCards := []utils.Card{21, 22, 23, 21, 22, 23, 24, 25, 25, 24, 24, 24, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingLongQiDui}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(32))
}

// 碰碰胡
func TestCalculateAndCardTypeValuePengPengHu(t *testing.T) {
	handUtilCards := []utils.Card{22, 22, 22, 23, 23, 23, 15, 15, 15, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{21}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_PengPengHu}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(2))
}

// 清碰
func TestCalculateAndCardTypeValueQingPeng(t *testing.T) {
	handUtilCards := []utils.Card{23, 23, 23, 25, 25, 25, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{21, 22}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingPeng}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(8))
}

// 金钩钓
func TestCalculateAndCardTypeValueJingGouDiao(t *testing.T) {
	handUtilCards := []utils.Card{27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{21, 22, 23, 15}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_JingGouDiao}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(4))
}

// 清金钩钓
func TestCalculateAndCardTypeValueQingJingGouDiao(t *testing.T) {
	handUtilCards := []utils.Card{27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{21, 22, 23, 25}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingJingGouDiao}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(16))
}

// 十八罗汉
func TestCalculateAndCardTypeValueShiBaLuoHan(t *testing.T) {
	handUtilCards := []utils.Card{27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{21, 22, 23, 15}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_ShiBaLuoHan}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(64))
}

// 清十八罗汉
func TestCalculateAndCardTypeValueQingShiBaLuoHan(t *testing.T) {
	handUtilCards := []utils.Card{27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{21, 22, 23, 25}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	cardTypes, genCount := new(ScxlCardTypeCalculator).Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingShiBaLuoHan}
	testCardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range testFanTypes {
		testCardTypes = append(testCardTypes, interfaces.CardType(cardType))
	}
	assert.Equal(t, cardTypes, testCardTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer := new(ScxlCardTypeCalculator).CardTypeValue(cardTypes, genCount)
	assert.Equal(t, valuer, uint32(256))
}

// 根
func TestCardGenSum(t *testing.T) {
	handUtilCards := []utils.Card{23, 23, 23, 23, 24}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{22}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{21, 24}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(21)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   0,
	}
	genSum := new(ScxlCardTypeCalculator).CardGenSum(playerParams)
	assert.Equal(t, genSum, uint32(4))
}
