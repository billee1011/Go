package fantype

import (
	"steve/majong/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCombines(t *testing.T) {
	cards := []utils.Card{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19}
	//cards := []Card{11, 11, 11, 13, 14, 15, 16, 17, 18, 19, 19, 19, 19}
	cardCombines := utils.FastCheckTingV2(cards, nil)

	for card, combines := range cardCombines {
		assert.Zero(t, card)

		assert.Nil(t, newCombines(combines))
	}
}
