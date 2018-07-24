package fantype

import (
	"steve/majong/global"
	majongpb "steve/entity/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_ZiYiSe 字一色:由字牌组成的胡牌;
func Test_ZiYiSe(t *testing.T) {
	handCard := []*majongpb.Card{
		&global.Card1Z, &global.Card1Z, &global.Card2Z, &global.Card2Z, &global.Card3Z, &global.Card3Z, &global.Card4Z, &global.Card4Z,
		&global.Card5Z, &global.Card5Z, &global.Card6Z, &global.Card6Z, &global.Card7Z,
	}
	tc := &typeCalculator{
		handCards: handCard,
		huCard: &majongpb.HuCard{
			Card: &global.Card7Z,
		},
		playerID: 1,
		player: &majongpb.Player{
			PalyerId:  1,
			GangCards: make([]*majongpb.GangCard, 0),
			PengCards: make([]*majongpb.PengCard, 0),
			ChiCards:  make([]*majongpb.ChiCard, 0),
		},
		cache: make(map[int]bool, 0),
	}

	assert.Equal(t, true, checkZiYiSe(tc))
}
