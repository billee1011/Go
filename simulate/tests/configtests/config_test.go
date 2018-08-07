package configtests

import (
	"steve/external/configclient"
	"steve/simulate/cheater"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_Config 测试配置获取
func Test_Config(t *testing.T) {
	assert.Nil(t, cheater.MockConfigClient())
	chargeMax, err := configclient.GetConfig("charge", "day_max")
	assert.Nil(t, err)
	assert.NotEmpty(t, chargeMax)
}
