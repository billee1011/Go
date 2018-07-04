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
	gangCards := make([]*majongpb.GangCard, 0)
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
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
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
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(24)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
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
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
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
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
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
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
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
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
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
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{21}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
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

// 清碰
func TestCalculateAndCardTypeValueQingPeng(t *testing.T) {
	handUtilCards := []utils.Card{23, 23, 23, 25, 25, 25, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{21, 22}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.SCXLGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	testFanTypes := []majongpb.CardType{majongpb.CardType_QingPeng}

	assert.Equal(t, cardTypes, testFanTypes)
	assert.Equal(t, genCount, uint32(0))
	valuer, gen := global.GetCardTypeCalculator().CardTypeValue(gutils.SCXLGameID, cardTypes, genCount)
	assert.Equal(t, valuer, uint32(8))
	assert.Equal(t, gen, uint32(0))
}

// 金钩钓
func TestCalculateAndCardTypeValueJingGouDiao(t *testing.T) {
	handUtilCards := []utils.Card{27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{21, 22, 23, 15}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
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
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{21, 22, 23, 25}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
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
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[2],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[3],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.SCXLGameID,
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
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[2],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[3],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.SCXLGameID,
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
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.SCXLGameID,
	}
	_, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Equal(t, genCount, uint32(4))
}

// TestDasixi 大四喜
func TestDasixi(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 41, 41, 42, 42, 42, 43, 43, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{44}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(41)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_DaSiXi)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(88))
}

// TestDasanyuan 大三元
func TestDasanyuan(t *testing.T) {
	handUtilCards := []utils.Card{11, 45, 45, 45, 15, 16, 17}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{47}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{46}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_DaSanYuan)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(88))
}

// TestJiuLianBaoDeng 九莲宝灯
func TestJiuLianBaoDeng(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{46}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_JiuLianBaoDeng)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(88))
}

// TestDayuwu 大于五
func TestDayuwu(t *testing.T) {
	handUtilCards := []utils.Card{16, 17, 18, 16, 17, 18, 16, 17, 18, 16, 17, 18, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_DaYuWu)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(88))
}

// TestXiaoyuwu 小于五
func TestXiaoyuwu(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 11, 12, 13, 11, 12, 13, 11, 12, 13, 14}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(14)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_XiaoYuWu)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(88))
}

// TestDaqixing 大七星
func TestDaqixing(t *testing.T) {
	handUtilCards := []utils.Card{41, 41, 46, 45, 45, 45, 45, 43, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{47}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(46)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_DaQiXing)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(88))
}

// TestLianqidui 连七对
func TestLianqidui(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(17)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_LianQiDui)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(88))
}

// TestSiGang 四杠
func TestSiGang(t *testing.T) {
	handUtilCards := []utils.Card{11}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{41, 42, 43, 44}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[2],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[3],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_SiGang)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(88))
}

// TestXiaosixi 小四喜
func TestXiaosixi(t *testing.T) {
	handUtilCards := []utils.Card{44, 12, 12, 12, 43, 43, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{41, 42}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(44)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_XiaoSiXi)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(64))
}

// TestXiaosanyuan 小三元
func TestXiaosanyuan(t *testing.T) {
	handUtilCards := []utils.Card{45, 46, 46, 46, 47, 47, 47}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{11, 12}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_XiaoSanYuan)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(64))
}

// TestShuanglonghui 双龙会
func TestShuanglonghui(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 11, 12, 13, 15, 17, 18, 19, 17, 18, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{11, 12}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(15)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_ShuangLongHui)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(64))
}

// TestZiyise 字一色
func TestZiyise(t *testing.T) {
	handUtilCards := []utils.Card{41, 41, 42, 42, 42, 43, 43, 43, 44, 44, 44, 45, 45}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{11, 12}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_ZiYiSe)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(64))
}

// TestSianke 四暗刻
func TestSianke(t *testing.T) {
	handUtilCards := []utils.Card{45}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{11, 12, 41, 42}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[2],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[3],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_SiAnKe)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(64))
}

// TestSitongshun 四同顺
func TestSitongshun(t *testing.T) {
	handUtilCards := []utils.Card{45, 11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{11, 12, 41, 42}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_SiTongShun)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(48))
}

// TestSanyuanqidui 三元七对
func TestSanyuanqidui(t *testing.T) {
	handUtilCards := []utils.Card{45, 41, 41, 42, 42}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{46, 47}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_SanYuanQiDui)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(48))
}

// TestSixiqidui 四喜七对
func TestSixiqidui(t *testing.T) {
	handUtilCards := []utils.Card{45, 41, 41, 42, 42}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{43, 44}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_SiXiQiDui)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(48))
}

// TestSilianke 四连刻
func TestSilianke(t *testing.T) {
	handUtilCards := []utils.Card{12, 12, 12, 13, 13, 13, 14, 14, 14, 15, 15, 15, 16}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{43, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(16)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_SiLianKe)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(48))
}

// TestSibugao 四步高
func TestSibugao(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 12, 13, 14, 13, 14, 15, 14, 15, 16, 41}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{43, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(41)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_SiBuGao)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(32))
}

// TestHunyaojiu 混幺九
func TestHunyaojiu(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 19, 19, 19, 41, 41, 41, 42, 42, 42, 46}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{43, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(46)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{},
		HuCard:   HuCard,
		GameID:   gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_HunYaoJiu)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(32))
}

// TestSangang 三杠
func TestSangang(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{43, 44, 42}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := interfaces.CardCalcParams{
		HandCard: handCards,
		PengCard: pengCards,
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_minggang,
			},
			&majongpb.GangCard{
				Card: gangCards[2],
				Type: majongpb.GangType_gang_bugang,
			},
		},
		HuCard: HuCard,
		GameID: gutils.ERGameID,
	}
	cardTypes, genCount := global.GetCardTypeCalculator().Calculate(playerParams)
	assert.Contains(t, cardTypes, majongpb.CardType_SanGang)
	value, _ := global.GetCardTypeCalculator().CardTypeValue(gutils.ERGameID, cardTypes, genCount)
	assert.Equal(t, value, uint32(32))
}

func calculate(params *interfaces.CardCalcParams) {
	// mjContext := &majongpb.MajongContext{}

}
