package fantype

import (
	"steve/majong/global"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SiGang 四杠:胡牌时,含有 4 个杠(明杠、暗杠);
func Test_SiGang(t *testing.T) {
	handCard := []*majongpb.Card{
		&global.Card2Z, &global.Card2Z,
	}
	gangCard := []*majongpb.GangCard{
		&majongpb.GangCard{
			Card: &global.Card1Z,
		},
		&majongpb.GangCard{
			Card: &global.Card3W,
		},
		&majongpb.GangCard{
			Card: &global.Card1W,
		},
		&majongpb.GangCard{
			Card: &global.Card2W,
		},
	}
	tc := &typeCalculator{
		handCards: handCard,
		huCard: &majongpb.HuCard{
			Card: &global.Card7W,
		},
		playerID: 1,
		player: &majongpb.Player{
			PalyerId:  1,
			GangCards: gangCard,
		},
		cache: make(map[int]bool, 0),
	}

	assert.Equal(t, true, checkSiGang(tc))
}
