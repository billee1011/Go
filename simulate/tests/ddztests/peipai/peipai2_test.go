package peipai

import (
	"steve/simulate/utils/doudizhu"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sirupsen/logrus"
)

//TestPeipai2 配牌测试
func TestPeipai2(t *testing.T) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "TestPeipai2()",
	})

	// 配牌测试1
	assert.Nil(t, doudizhu.PeipaiTest2(t))

	logEntry.Info("配牌2测试正常结束")
}
