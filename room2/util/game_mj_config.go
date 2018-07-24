package util

import (
	"github.com/spf13/viper"
	"steve/room2/fixed"
)

func init() {
	viper.SetDefault(fixed.ListenClientAddr, "127.0.0.1")
	viper.SetDefault(fixed.ListenClientPort, 36001)
	viper.SetDefault(fixed.ListenPeipaiAddr, "")
	viper.SetDefault(fixed.XingPaiTimeOut, 10)
	viper.SetDefault(fixed.MaxFapaiCartoonTime, 10*1000)
	viper.SetDefault(fixed.MaxHuansanzhangCartoonTime, 10*1000)
}