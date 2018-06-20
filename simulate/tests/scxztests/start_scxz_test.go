package tests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_StartGame_NoHsz 测试游戏开始
// 游戏开始流程包括： 登录，加入房间，配牌，洗牌，发牌,定缺
// 期望不出现换三张
func Test_SCXZ_StartGame_NoHsz(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId(2)
	params.PeiPaiGame = "scxz"
	params.IsHsz = false // 不换三张
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
}
