package fantype

import (
	"steve/client_pb/room"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 平胡
func TestCalculateAndCardTypeValuePingHu(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 11, 12, 13, 14, 15, 16, 14, 15, 16, 17}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(17)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, genCount, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_PINGHU))
	assert.Equal(t, genCount, uint32(0))
}

// 清一色
func TestCalculateAndCardTypeValueQingYiSe(t *testing.T) {
	handUtilCards := []utils.Card{21, 22, 23, 21, 22, 23, 24, 25, 26, 27, 28, 29, 24}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(24)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, genCount, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_QINGYISE))
	assert.Equal(t, genCount, uint32(0))

}

// 七对
func TestCalculateAndCardTypeValueQiDui(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 11, 12, 13, 24, 25, 26, 24, 25, 26, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_QIDUI))

}

// 龙七对
func TestCalculateAndCardTypeValueLongQiDui(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 11, 12, 13, 24, 25, 25, 24, 24, 24, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_LONGQIDUI))

}

// 清七对
func TestCalculateAndCardTypeValueQingQiDui(t *testing.T) {
	handUtilCards := []utils.Card{21, 22, 23, 21, 22, 23, 24, 25, 26, 24, 25, 26, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_QINGQIDUI))
}

// 清龙七对
func TestCalculateAndCardTypeValueQingLongQiDui(t *testing.T) {
	handUtilCards := []utils.Card{21, 22, 23, 21, 22, 23, 24, 25, 25, 24, 24, 24, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Equal(t, cardTypes, int(room.FanType_FT_QINGLONGQIDUI))
}

// 碰碰胡
func TestCalculateAndCardTypeValuePengPengHu(t *testing.T) {
	handUtilCards := []utils.Card{22, 22, 22, 23, 23, 23, 15, 15, 15, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{21}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_PENGPENGHU))

}

// 清碰
func TestCalculateAndCardTypeValueQingPeng(t *testing.T) {
	handUtilCards := []utils.Card{23, 23, 23, 25, 25, 25, 27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{21, 22}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Equal(t, cardTypes, int(room.FanType_FT_QINGPENG))
}

// 金钩钓
func TestCalculateAndCardTypeValueJingGouDiao(t *testing.T) {
	handUtilCards := []utils.Card{27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{21, 22, 23, 15}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Equal(t, cardTypes, int(room.FanType_FT_JINGGOUDIAO))

}

// 清金钩钓
func TestCalculateAndCardTypeValueQingJingGouDiao(t *testing.T) {
	handUtilCards := []utils.Card{27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{21, 22, 23, 25}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Equal(t, cardTypes, int(room.FanType_FT_QINGJINGGOUDIAO))

}

// 十八罗汉
func TestCalculateAndCardTypeValueShiBaLuoHan(t *testing.T) {
	handUtilCards := []utils.Card{27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{21, 22, 23, 15}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
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
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Equal(t, cardTypes, int(room.FanType_FT_SHIBALUOHAN))
}

// 清十八罗汉
func TestCalculateAndCardTypeValueQingShiBaLuoHan(t *testing.T) {
	handUtilCards := []utils.Card{27}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{21, 22, 23, 25}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(27)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
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
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Equal(t, cardTypes, int(room.FanType_FT_QINGSHIBALUOHAN))
}

// 根
func TestCardGenSum(t *testing.T) {
	handUtilCards := []utils.Card{23, 23, 23, 23, 24}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{22}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{21, 24}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(21)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	_, genCount, _ := calculate(playerParams)
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
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(41)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_DASIXI))
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
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{
			&majongpb.PengCard{
				Card: pengCards[0],
			},
		},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_DASANYUAN))
}

// TestJiuLianBaoDeng 九莲宝灯
func TestJiuLianBaoDeng(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{46}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_JIULIANBAODENG))
}

// TestDayuwu 大于五
func TestDayuwu(t *testing.T) {
	handUtilCards := []utils.Card{16, 17, 18, 16, 17, 18, 16, 17, 18, 16, 17, 18, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_DAYUWU))
}

// TestXiaoyuwu 小于五
func TestXiaoyuwu(t *testing.T) {
	handUtilCards := []utils.Card{12, 13, 14, 12, 13, 14, 12, 13, 14, 12, 13, 14, 11}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_XIAOYUWU))
}

// TestDaqixing 大七星
func TestDaqixing(t *testing.T) {
	handUtilCards := []utils.Card{41, 41, 46, 45, 45, 45, 45, 43, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{47}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(46)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
		},
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_DAQIXING))
}

// TestLianqidui 连七对
func TestLianqidui(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(17)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_LIANQIDUI))
}

// TestSiGang 四杠
func TestSiGang(t *testing.T) {
	handUtilCards := []utils.Card{11}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{41, 42, 43, 44}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
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
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SIGANG))
}

// TestXiaosixi 小四喜
func TestXiaosixi(t *testing.T) {
	handUtilCards := []utils.Card{44, 12, 12, 12, 43, 43, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{41, 42}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(44)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
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
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_XIAOSIXI))
}

// TestXiaosanyuan 小三元
func TestXiaosanyuan(t *testing.T) {
	handUtilCards := []utils.Card{45, 46, 46, 46, 47, 47, 47}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{11, 12}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
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
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_XIAOSANYUAN))
}

// TestShuanglonghui 双龙会
func TestShuanglonghui(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 11, 12, 13, 15, 17, 18, 19, 17, 18, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{11, 12}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(15)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SHUANGLONGHUI))

}

// TestZiyise 字一色
func TestZiyise(t *testing.T) {
	handUtilCards := []utils.Card{41, 41, 42, 42, 42, 43, 43, 43, 44, 44, 44, 45, 45}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{11, 12}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_ZIYISE))
}

// TestSianke 四暗刻
func TestSianke(t *testing.T) {
	handUtilCards := []utils.Card{45}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{11, 12, 41, 42}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
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
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SIANKE))
}

// TestSitongshun 四同顺
func TestSitongshun(t *testing.T) {
	handUtilCards := []utils.Card{45, 11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{11, 12, 41, 42}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SITONGSHUN))
}

// TestSanyuanqidui 三元七对
func TestSanyuanqidui(t *testing.T) {
	handUtilCards := []utils.Card{45, 41, 41, 42, 42}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{46, 47}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
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
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SANYUANQIDUI))
}

// TestSixiqidui 四喜七对
func TestSixiqidui(t *testing.T) {
	handUtilCards := []utils.Card{45, 41, 41, 42, 42}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{43, 44}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(45)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
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
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SIXIQIDUI))
}

// TestSilianke 四连刻
func TestSilianke(t *testing.T) {
	handUtilCards := []utils.Card{12, 12, 12, 13, 13, 13, 14, 14, 14, 15, 15, 15, 16}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{43, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(16)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SILIANKE))
}

// TestSibugao 四步高
func TestSibugao(t *testing.T) {
	handUtilCards := []utils.Card{11, 12, 13, 12, 13, 14, 13, 14, 15, 14, 15, 16, 41}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{43, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(41)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SIBUGAO))
}

// TestHunyaojiu 混幺九
func TestHunyaojiu(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 19, 19, 19, 41, 41, 41, 42, 42, 42, 46}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{43, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(46)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_HUNYAOJIU))
}

// TestSangang 三杠
func TestSangang(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{43, 44, 42}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	//pengUtilCards := []utils.Card{}
	//pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
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
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SANGANG))

}

// TestSizike 四字刻
func TestSizike(t *testing.T) {
	handUtilCards := []utils.Card{19, 42, 42, 42}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{43, 44}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{47}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{
			&majongpb.PengCard{
				Card: pengCards[0],
			},
		},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_minggang,
			},
		},
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SIZIKE))

}

// TestDasanfeng 大三风
func TestDasanfeng(t *testing.T) {
	handUtilCards := []utils.Card{19, 41, 41, 41}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{47, 44}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{42}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{
			&majongpb.PengCard{
				Card: pengCards[0],
			},
		},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_angang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_minggang,
			},
		},
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_DASANFENG))
}

// TestSantongshun 三同顺
func TestSantongshun(t *testing.T) {
	handUtilCards := []utils.Card{12, 13, 11, 12, 13, 41, 41, 42, 42, 42}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{42}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	chiUtilCards := []utils.Card{11}
	chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		ChiCard: []*majongpb.ChiCard{
			&majongpb.ChiCard{
				Card:    chiCards[0],
				OprCard: chiCards[0],
			},
		},
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SANTONGSHUN))
}

// TestSanlianke 三连刻
func TestSanlianke(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 12, 13, 17, 17, 17, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	pengUtilCards := []utils.Card{18}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{11}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{
			&majongpb.PengCard{
				Card: pengCards[0],
			},
		},
		GangCard: []*majongpb.GangCard{},
		ChiCard:  []*majongpb.ChiCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SANLIANKE))
}

// TestQinglong 清龙
func TestQinglong(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 14, 15, 16, 17, 18, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{18}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	chiUtilCards := []utils.Card{11}
	chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		ChiCard: []*majongpb.ChiCard{
			&majongpb.ChiCard{Card: chiCards[0]},
		},
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_QINGLONG))
}

// TestSanbugao 三步高
func TestSanbugao(t *testing.T) {
	handUtilCards := []utils.Card{15, 16, 17, 17, 18, 19, 41, 41, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{47, 44}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{18}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	chiUtilCards := []utils.Card{13}
	chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		ChiCard: []*majongpb.ChiCard{
			&majongpb.ChiCard{Card: chiCards[0]},
		},
		HuCard: &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SANBUGAO))
}

// TestXiaosanfeng 小三风
func TestXiaosanfeng(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 12, 12, 12, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{41}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{42}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(43)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{
			&majongpb.PengCard{Card: pengCards[0]},
		},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0]},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_XIAOSANFENG))
}

// TestHunyise 混一色
func TestHunyise(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 12, 12, 12, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{41}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{42}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(43)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{
			&majongpb.PengCard{Card: pengCards[0]},
		},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0]},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_HUNYISE))
}

//天胡
//地胡
//人胡
//天听

//全花
// TestSananke 三暗刻
func TestSananke(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 12, 12, 12, 13, 13, 13, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{41}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	// pengUtilCards := []utils.Card{42}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(43)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0]},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SANANKE))
}

//妙手回春
//海底捞月
//杠上开花
//抢杠胡

//双箭刻
// TestShuangjianke 双箭刻
func TestShuangjianke(t *testing.T) {
	handUtilCards := []utils.Card{46, 46, 46, 12, 12, 12, 13, 13, 13, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{45}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	// pengUtilCards := []utils.Card{42}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(43)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0]},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SHUANGJIANKE))
}

// TestShuangangang 双暗杠
func TestShuangangang(t *testing.T) {
	handUtilCards := []utils.Card{12, 12, 12, 13, 13, 13, 43}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{45, 46}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	// pengUtilCards := []utils.Card{42}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(43)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0]},
			&majongpb.GangCard{Card: gangCards[1]},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SHUANGANGANG))
}

// TestQuanqiuren 全求人
func TestQuanqiuren(t *testing.T) {
	handUtilCards := []utils.Card{11}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{41, 43, 12}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	pengUtilCards := []utils.Card{42}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{
			&majongpb.PengCard{Card: pengCards[0]},
		},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0]},
			&majongpb.GangCard{Card: gangCards[1]},
			&majongpb.GangCard{Card: gangCards[2]},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_dianpao},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_QUANQIUREN))
}

// TestQuandaiyao 全带幺
func TestQuandaiyao(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 11, 12, 13, 17, 18, 19, 17, 18, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{41}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{42}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		ChiCard:  []*majongpb.ChiCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_QUANDAIYAO))
}

// TestShuangminggang 双明杠
func TestShuangminggang(t *testing.T) {
	handUtilCards := []utils.Card{17, 18, 19, 17, 18, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{11, 12}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	// pengUtilCards := []utils.Card{42}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{
				Card: gangCards[0],
				Type: majongpb.GangType_gang_minggang,
			},
			&majongpb.GangCard{
				Card: gangCards[1],
				Type: majongpb.GangType_gang_minggang,
			},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SHUANGMINGGANG))
}

// TestBuqiuren 不求人
func TestBuqiuren(t *testing.T) {
	handUtilCards := []utils.Card{11, 11, 11, 11, 12, 13, 17, 18, 19, 17, 18, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{11, 12}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{42}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(19)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		ChiCard:  []*majongpb.ChiCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_zimo},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_BUQIUREN))
}

// TestJuezhang 绝张
func TestJuezhang(t *testing.T) {
	handUtilCards := []utils.Card{12, 13, 17, 18, 19, 17, 18, 19, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{11, 12}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	pengUtilCards := []utils.Card{11}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{
			&majongpb.PengCard{Card: pengCards[0]},
		},
		GangCard: []*majongpb.GangCard{},
		ChiCard:  []*majongpb.ChiCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_zimo},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_JUEZHANG))
}

//门风刻
//圈风刻

// TestJianke 箭刻
func TestJianke(t *testing.T) {
	handUtilCards := []utils.Card{12, 13, 45, 45, 45, 46, 46, 46, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{11, 12}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	pengUtilCards := []utils.Card{47}
	pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{13}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{
			&majongpb.PengCard{Card: pengCards[0]},
		},
		GangCard: []*majongpb.GangCard{},
		ChiCard:  []*majongpb.ChiCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_zimo},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_JIANKE))
}

// TestSiguiyi 四归一
func TestSiguiyi(t *testing.T) {
	handUtilCards := []utils.Card{12, 13, 45, 45, 45, 46, 46, 46, 19, 19}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{47}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0]},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_zimo},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SIGUIYI))
}

// TestDuanyao 断幺
func TestDuanyao(t *testing.T) {
	handUtilCards := []utils.Card{12, 12, 15, 15, 15, 16, 16, 16, 18, 18}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{13}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(12)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0]},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_zimo},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_DUANYAO))
}

// TestShuanganke 双暗刻
func TestShuanganke(t *testing.T) {
	handUtilCards := []utils.Card{12, 12, 15, 15, 15, 16, 16, 16, 18, 18}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{13}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(12)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0]},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_zimo},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_SHUANGANKE))
}

// TestAngang 暗杠
func TestAngang(t *testing.T) {
	handUtilCards := []utils.Card{12, 12, 15, 15, 15, 16, 16, 16, 18, 18}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{13}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	assert.Nil(t, err)
	HuCard, err := utils.IntToCard(12)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0], Type: majongpb.GangType_gang_angang},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_zimo},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_ANGANG))
}

// TestMenqianqing 门前清
func TestMenqianqing(t *testing.T) {
	handUtilCards := []utils.Card{12, 12, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{13}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	// assert.Nil(t, err)
	HuCard, err := utils.IntToCard(12)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		ChiCard:  []*majongpb.ChiCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_dianpao},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_MENQIANQING))
}

//报听

// TestYibangao 一般高
func TestYibangao(t *testing.T) {
	handUtilCards := []utils.Card{12, 12, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{13}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	// assert.Nil(t, err)
	HuCard, err := utils.IntToCard(12)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		ChiCard:  []*majongpb.ChiCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_dianpao},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_YIBANGAO))
}

// TestLianliu 连六
func TestLianliu(t *testing.T) {
	handUtilCards := []utils.Card{13, 14, 15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{13}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	// assert.Nil(t, err)
	HuCard, err := utils.IntToCard(12)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		ChiCard:  []*majongpb.ChiCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_dianpao},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_LIANLIU))
}

// TestLaoshaofu 老少副
func TestLaoshaofu(t *testing.T) {
	handUtilCards := []utils.Card{12, 13, 15, 15, 15, 16, 16, 16, 17, 18, 19, 18, 18}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	// gangUtilCards := []utils.Card{13}
	// gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	// assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	// assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{},
		ChiCard:  []*majongpb.ChiCard{},
		HuCard:   &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_dianpao},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_LAOSHAOFU))
}

//花牌

// TestMinggang 明杠
func TestMinggang(t *testing.T) {
	handUtilCards := []utils.Card{12, 13, 16, 16, 16, 17, 18, 19, 18, 18}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{15}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	// assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0], Type: majongpb.GangType_gang_minggang},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_dianpao},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_MINGGANG))
}

//边张
//坎张
//单吊将
// TestZimo 自摸
func TestZimo(t *testing.T) {
	handUtilCards := []utils.Card{12, 13, 16, 16, 16, 17, 18, 19, 18, 18}
	handCards, err := utils.CheckHuUtilCardsToHandCards(handUtilCards)
	assert.Nil(t, err)
	gangUtilCards := []utils.Card{15}
	gangCards, err := utils.CheckHuUtilCardsToHandCards(gangUtilCards)
	assert.Nil(t, err)
	// pengUtilCards := []utils.Card{47}
	// pengCards, err := utils.CheckHuUtilCardsToHandCards(pengUtilCards)
	// chiUtilCards := []utils.Card{47}
	// chiCards, err := utils.CheckHuUtilCardsToHandCards(chiUtilCards)
	// assert.Nil(t, err)
	HuCard, err := utils.IntToCard(11)
	assert.Nil(t, err)
	playerParams := CardCalcParams{
		HandCard: handCards,
		PengCard: []*majongpb.PengCard{},
		GangCard: []*majongpb.GangCard{
			&majongpb.GangCard{Card: gangCards[0], Type: majongpb.GangType_gang_minggang},
		},
		ChiCard: []*majongpb.ChiCard{},
		HuCard:  &majongpb.HuCard{Card: HuCard, Type: majongpb.HuType_hu_zimo},
	}
	cardTypes, _, _ := calculate(playerParams)
	assert.Contains(t, cardTypes, int(room.FanType_FT_ZIMO))
}

func calculate(params CardCalcParams) ([]int, int, int) {
	player := &majongpb.Player{
		PalyerId:  1,
		HandCards: params.HandCard,
		PengCards: params.PengCard,
		GangCards: params.GangCard,
		ChiCards:  params.ChiCard,
		HuCards:   []*majongpb.HuCard{params.HuCard},
	}
	mjContext := &majongpb.MajongContext{
		XingpaiOptionId:  4,
		CardtypeOptionId: 4,
		SettleOptionId:   4,
		Players:          []*majongpb.Player{player},
	}
	return CalculateFanTypes(mjContext, 1, params.HandCard, params.HuCard)
}

// CardCalcParams 计算牌型的参数
type CardCalcParams struct {
	HandCard []*majongpb.Card
	PengCard []*majongpb.PengCard
	GangCard []*majongpb.GangCard
	ChiCard  []*majongpb.ChiCard
	HuCard   *majongpb.HuCard
}

func TestNewCombines(t *testing.T) {
	// cards := []utils.Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19}
	// //cards := []Card{11, 11, 11, 13, 14, 15, 16, 17, 18, 19, 19, 19, 19}
	// cardCombines := utils.FastCheckTingV2(cards, nil)

	// for card, combines := range cardCombines {
	// 	assert.Zero(t, card)

	// 	assert.Nil(t, newCombines(combines))
	// }
}
