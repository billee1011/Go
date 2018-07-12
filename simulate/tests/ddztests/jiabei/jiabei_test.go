package jiabei

import (
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//TestJiaodizhu 叫地主测试
//发牌完成后，服务器指定的玩家发起叫地主请求
//期望：
//     1. 所有玩家都收到，那个玩家的叫地主广播
func TestJiabei(t *testing.T) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "TestJiabei()",
	})

	// 正常开始游戏
	params := global.NewStartDDZGameParams()
	deskData, err := utils.StartDDZGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 当前状态的时间间隔
	logEntry.Infof("当前状态 = %v, 进入下一状态等待时间 = %v", deskData.DDZData.NextState.GetStage(), deskData.DDZData.NextState.GetTime())

}
