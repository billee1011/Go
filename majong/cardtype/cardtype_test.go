package cardtype

import (
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

var playerParams = interfaces.CardCalcParams{
	GameID: gutils.SCXLGameID,
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
		GameID:   gutils.SCXLGameID,
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_PingHu}
	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(1))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingYiSe}
	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(4))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QiDui}

	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(4))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_LongQiDui}
	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(8))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingQiDui}

	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(16))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingLongQiDui}
	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(32))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_PengPengHu}
	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(2))
	assert.Equal(t, gen, uint32(0))
}

// 清碰 加 一根
func TestCalculateAndCardTypeValueQingPeng(t *testing.T) {
	handUtilCards := []utils.Card{21, 21, 21, 21}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{25, 24, 23}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(21)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   HuCard,
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingPeng}

	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(1))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(16))
	assert.Equal(t, gen, uint32(1))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_JingGouDiao}
	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(4))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingJingGouDiao}
	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(16))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_ShiBaLuoHan}
	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(64))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingShiBaLuoHan}
	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(256))
	assert.Equal(t, gen, uint32(0))
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
		GameID:   gutils.SCXLGameID,
	}
	_, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Equal(t, genCount, uint32(4))
}
