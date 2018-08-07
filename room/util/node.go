package util

import (
	"steve/gutils"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var idAllocObject *gutils.Node

func init() {
	node := viper.GetInt("node")
	var err error
	idAllocObject, err = gutils.NewNode(int64(node))
	if err != nil {
		logrus.Panicf("创建 id 生成器失败: %v", err)
	}
}

// GenUniqueID 生成全局唯一ID
func GenUniqueID() gutils.ID {
	return idAllocObject.Generate()
}
