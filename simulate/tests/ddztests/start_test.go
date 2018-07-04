package ddztests

import (
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
)

func Test_StartGame(t *testing.T) {

	// 扑克游戏的参数
	params := global.NewStartPokeGameParams()

	// 开始扑克游戏
	//ddzGame, err := game.StartGame(params)

	//assert.NotNil(t, ddzGame)
	//assert.Nil(t, err)

	utils.StartPokeGame(params)
	//fmt.Println("%v", err)
	//deskData, err := utils.StartPokeGame(params)
	//assert.NotNil(t, deskData)
	//assert.Nil(t, err)
}
