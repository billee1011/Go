package ddztests

import (
	"steve/simulate/utils"
	"steve/simulate/utils/doudizhu"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StartGame(t *testing.T) {

	// 配牌1
	params := doudizhu.NewStartDDZGameParamsTest1()

	// 开始游戏
	deskData, err := utils.StartDDZGame(params)

	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	utils.ClearPeiPai(params.PeiPaiGame)
}
