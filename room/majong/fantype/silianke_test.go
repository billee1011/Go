package fantype

import (
	majongpb "steve/entity/majong"
	"steve/room/majong/global"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SiLianKe 四连刻:胡牌时,含有一种花色 4 副依次递增一位数的刻子;
func Test_SiLianKe(t *testing.T) {
	handCard := []*majongpb.Card{
		&global.Card5W, &global.Card5W, &global.Card6W, &global.Card6W, &global.Card5W, &global.Card6W, &global.Card1Z, &global.Card1Z,
	}
	pengCard := []*majongpb.PengCard{
		&majongpb.PengCard{
			Card: &global.Card3W,
		},
		&majongpb.PengCard{
			Card: &global.Card4W,
		},
	}
	tc := &typeCalculator{
		handCards: handCard,
		huCard: &majongpb.HuCard{
			Card: &global.Card7W,
		},
		combines: []Combine{
			Combine{
				jiang: 41,
				kes:   []int{15, 16},
			},
		},
		playerID: 1,
		player: &majongpb.Player{
			PlayerId:  1,
			PengCards: pengCard,
		},
		cache: make(map[int]bool, 0),
	}

	assert.Equal(t, true, checkSiLianKe(tc))
}
