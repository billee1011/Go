package util

import (
	"github.com/spf13/viper"
	"steve/room2/common"
)

func init() {
	viper.SetDefault(common.ListenClientAddr, "127.0.0.1")
	viper.SetDefault(common.ListenClientPort, 36001)
	viper.SetDefault(common.ListenPeipaiAddr, "")
	viper.SetDefault(common.XingPaiTimeOut, 10)
	viper.SetDefault(common.MaxFapaiCartoonTime, 10*1000)
	viper.SetDefault(common.MaxHuansanzhangCartoonTime, 10*1000)
}