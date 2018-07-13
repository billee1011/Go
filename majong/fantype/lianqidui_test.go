package fantype

import (
	"steve/majong/global"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_LianQiDui 连七对:由一种花色序数牌组成序数相连的 7 个对子组成的胡牌;
func Test_LianQiDui(t *testing.T) {
	handCard := []*majongpb.Card{
		&global.Card1W, &global.Card1W, &global.Card2W, &global.Card2W, &global.Card3W, &global.Card3W, &global.Card4W, &global.Card4W,
		&global.Card5W, &global.Card5W, &global.Card6W, &global.Card6W, &global.Card7W,
	}
	tc := &typeCalculator{
		handCards: handCard,
		huCard: &majongpb.HuCard{
			Card: &global.Card7W,
		},
		cache: make(map[int]bool, 0),
	}

	assert.Equal(t, true, checkLianQiDui(tc))
}
