package config

import "github.com/spf13/viper"

const (
	// ListenClientAddr 代表监听客户端的IP地址，默认值为 127.0.0.1
	ListenClientAddr = "lis_client_addr"

	// ListenClientPort 代表监听客户端的端口， 默认值为 36001
	ListenClientPort = "lis_client_port"

	// ListenPeipaiAddr 配牌監聽地址
	ListenPeipaiAddr = "peipai_addr"

	// XingPaiTimeOut 行牌超时时间，单位为second，默认值为 10
	XingPaiTimeOut = "xp_timeout"

	// TingStateTimeOut 听牌状态下的超时时间，单位为second，默认值为 1
	TingStateTimeOut = "ts_timeout"

	// HuStateTimeOut 胡牌状态下的超时时间，单位为second，默认值为 1
	HuStateTimeOut = "hs_timeout"

	// MaxFapaiCartoonTime 发牌的动画时间
	MaxFapaiCartoonTime = "fp_cartoontime"

	// MaxHuansanzhangCartoonTime 换三张动画时间
	MaxHuansanzhangCartoonTime = "hsz_cartoontime"
)

func init() {
	viper.SetDefault(ListenClientAddr, "127.0.0.1")
	viper.SetDefault(ListenClientPort, 36001)
	viper.SetDefault(ListenPeipaiAddr, "")
	viper.SetDefault(XingPaiTimeOut, 10)
	viper.SetDefault(TingStateTimeOut, 1)
	viper.SetDefault(HuStateTimeOut, 1)
	viper.SetDefault(MaxFapaiCartoonTime, 6*1000)
	viper.SetDefault(MaxHuansanzhangCartoonTime, 4*1000)
}
