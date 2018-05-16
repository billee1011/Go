package tests

import (
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_StartGame 测试游戏开始
// 游戏开始流程包括： 登录，加入房间，配牌，洗牌，发牌，
func Test_StartGame(t *testing.T) {
	deskData, err := utils.StartGame(commonStartGameParams)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
}
