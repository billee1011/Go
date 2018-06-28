package mjoption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXingPaiOptionManager_loadOption(t *testing.T) {
	xom := NewXingPaiOptionManager("testdata/xingpai")
	opt1 := xom.GetXingPaiOption(1)
	assert.Equal(t, opt1.Hnz.Need, true)
	assert.Equal(t, opt1.Hnz.Num, 3)
	assert.Equal(t, opt1.ID, 1)
	assert.Equal(t, opt1.NeedAddflower, false)
	assert.Equal(t, opt1.NeedChi, false)
	assert.Equal(t, opt1.NeedDingque, true)
	assert.Equal(t, opt1.PlayerStates, []XingpaiState{})
	assert.Equal(t, len(opt1.WallCards), 136)
	assert.Equal(t, opt1.HuGameOver, false)

	opt2 := xom.GetXingPaiOption(2)
	assert.Equal(t, opt2.Hnz.Need, true)
	assert.Equal(t, opt2.Hnz.Num, 3)
	assert.Equal(t, opt2.ID, 2)
	assert.Equal(t, opt2.NeedAddflower, false)
	assert.Equal(t, opt2.NeedChi, false)
	assert.Equal(t, opt2.NeedDingque, true)
	assert.Equal(t, opt2.PlayerStates, []XingpaiState{Hu, Giveup})
	assert.Equal(t, len(opt2.WallCards), 136)
	assert.Equal(t, opt2.HuGameOver, false)
}
