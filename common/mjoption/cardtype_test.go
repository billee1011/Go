package mjoption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkFanType(t *testing.T, fantype FanType, opt *CardTypeOption) {
	ft, ok := opt.Fantypes[fantype.ID]
	assert.True(t, ok)
	assert.Equal(t, ft.FuncID, fantype.FuncID)
	assert.Equal(t, ft.ID, fantype.ID)
	assert.Equal(t, ft.Mutex, fantype.Mutex)
	assert.Equal(t, ft.Method, fantype.Method)
	assert.Equal(t, ft.Score, fantype.Score)
}

func TestCardTypeOptionManager_loadOption(t *testing.T) {
	com := NewCardTypeOptionManager("testdata/cardtype")
	opt1 := com.GetCardTypeOption(4)
	checkFanType(t, FanType{
		ID:     0,
		FuncID: 0,
		Mutex:  []int{},
		Method: 0,
		Score:  1,
	}, opt1)

	checkFanType(t, FanType{
		ID:     1,
		FuncID: 1,
		Mutex:  []int{0},
		Method: 1,
		Score:  1,
	}, opt1)
}
