package play

import (
	"steve/simulate/utils"
	"steve/simulate/utils/doudizhu"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//TestPlaycard1 出牌测试
//游戏过程中0号玩家发起叫地主
//期望：
//     1. 所有玩家都收到，0号玩家的叫地主广播
func TestPlaycard1(t *testing.T) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "TestPlaycard1()",
	})

	// 配牌1
	params := doudizhu.NewStartDDZGameParamsTest1()

	deskData, err := utils.StartDDZGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 当前状态的时间间隔
	logEntry.Infof("当前状态 = %v, 进入下一状态等待时间 = %v", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())

	// 叫地主用例1
	assert.Nil(t, doudizhu.JiaodizhuTest1(deskData))

	// 加倍用例1
	assert.Nil(t, doudizhu.JiabeiTest1(deskData))

	// 出牌用例1
	assert.Nil(t, doudizhu.PlaycardTest1(deskData))
}
