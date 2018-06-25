package mjoption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardTypeOptionManager_loadOption(t *testing.T) {
	com := NewCardTypeOptionManager("testdata/cardtype")
	opt1 := com.GetCardTypeOption(1)
	assert.Equal(t, opt1.EnableQidui, true)

	opt2 := com.GetCardTypeOption(2)
	assert.Equal(t, opt2.EnableQidui, false)
}
