package mjoption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_loadGameOptions 测试游戏选项的加载
func Test_loadGameOptions(t *testing.T) {
	gom := NewGameOptionManager("testdata/mjoption.yaml")

	opt1 := gom.GetGameOptions(1)
	assert.Equal(t, 1, opt1.GameID)
	assert.Equal(t, 1, opt1.SettleOptionID)
	assert.Equal(t, 1, opt1.CardTypeOptionID)
	assert.Equal(t, 1, opt1.XingPaiOptionID)

	opt2 := gom.GetGameOptions(2)
	assert.Equal(t, 2, opt2.GameID)
	assert.Equal(t, 2, opt2.SettleOptionID)
	assert.Equal(t, 2, opt2.CardTypeOptionID)
	assert.Equal(t, 2, opt2.XingPaiOptionID)
}
