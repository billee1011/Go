package transition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadTransition(t *testing.T) {
	tt, err := loadTransition("./transition.yaml")
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(tt))
}
