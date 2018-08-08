package util

import (
	"steve/room/fixed"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(fixed.ListenClientAddr, "127.0.0.1")
	viper.SetDefault(fixed.ListenClientPort, 36001)
	viper.SetDefault(fixed.ListenPeipaiAddr, "")
	viper.SetDefault(fixed.XingPaiTimeOut, 10)
	viper.SetDefault(fixed.MaxFapaiCartoonTime, 10*1000)
	viper.SetDefault(fixed.MaxHuansanzhangCartoonTime, 10*1000)
	viper.SetDefault(fixed.TingStateTimeOut, 1)
	viper.SetDefault(fixed.HuStateTimeOut, 1)
	viper.SetDefault(fixed.MaxFapaiCartoonTime, 6*1000)
	viper.SetDefault(fixed.MaxHuansanzhangCartoonTime, 4*1000)
}
