package mjoption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXingPaiOptionManager_loadOption(t *testing.T) {
	xom := NewXingPaiOptionManager("testdata/xingpai")
	opt1 := xom.GetXingPaiOption(1)
	assert.Equal(t, opt1.HuGameOver, false)

	opt2 := xom.GetXingPaiOption(2)
	assert.Equal(t, opt2.HuGameOver, true)
}
