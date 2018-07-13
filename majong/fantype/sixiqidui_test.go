package fantype

import (
	"steve/majong/global"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SiXiQiDui 四喜七对:胡牌为七对,并且包含“东南西北”
func Test_SiXiQiDui(t *testing.T) {
	handCard := []*majongpb.Card{
		&global.Card1Z, &global.Card1Z, &global.Card2Z, &global.Card2Z, &global.Card3Z, &global.Card3Z, &global.Card4Z, &global.Card4Z,
		&global.Card5W, &global.Card5W, &global.Card6W, &global.Card6W, &global.Card7W,
	}
	tc := &typeCalculator{
		handCards: handCard,
		huCard: &majongpb.HuCard{
			Card: &global.Card7W,
		},
		cache: make(map[int]bool, 0),
	}

	assert.Equal(t, true, checkSiXiQiDui(tc))
}
