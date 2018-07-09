package ddztests

import (
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StartGame(t *testing.T) {

	// 扑克游戏的参数
	params := global.NewStartDDZGameParams()

	// 开始扑克游戏
	//ddzGame, err := game.StartGame(params)

	//assert.NotNil(t, ddzGame)
	//assert.Nil(t, err)

	deskData, err := utils.StartDDZGame(params)
	//fmt.Println("%v", err)
	//deskData, err := utils.StartDDZGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
}
