package mjoption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSettleOptionManager_loadOption 测试结算选项加载
func TestSettleOptionManager_loadOption(t *testing.T) {
	som := NewSettleOptionManager("testdata/settle")
	opt1 := som.GetSettleOption(1)
	assert.Equal(t, opt1.EnableTuisui, true)

	opt2 := som.GetSettleOption(2)
	assert.Equal(t, opt2.EnableTuisui, false)
}
