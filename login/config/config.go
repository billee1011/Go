package config

import "github.com/spf13/viper"

const (
	// ListenClientAddr 代表监听客户端的IP地址，默认值为 127.0.0.1
	ListenClientAddr = "lis_client_addr"

	// ListenClientPort 代表监听客户端的端口， 默认值为 36201
	ListenClientPort = "lis_client_port"

	// AuthKey 认证秘钥
	AuthKey = "auth_key"
)

func init() {
	viper.SetDefault(ListenClientAddr, "127.0.0.1")
	viper.SetDefault(ListenClientPort, 36201)
	viper.SetDefault(AuthKey, "stevegame.cn")
}
