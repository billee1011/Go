package config

import "github.com/spf13/viper"

const (
	// ListenClientAddr 代表监听客户端的IP地址，默认值为 127.0.0.1
	ListenClientAddr = "lis_client_addr"

	// ListenClientPort 代表监听客户端的端口， 默认值为 36001
	ListenClientPort = "lis_client_port"

	// ListenClientAddrInquire 客户端连接地址
	ListenClientAddrInquire = "lis_client_addr_inquire"

	// AuthKey 认证秘钥
	AuthKey = "auth_key"
)

func init() {
	viper.SetDefault(ListenClientAddr, "127.0.0.1")
	viper.SetDefault(ListenClientPort, 36001)
	viper.SetDefault(ListenClientAddrInquire, "127.0.0.1")
	viper.SetDefault(AuthKey, "stevegame.cn")
}

// GetRPCAddr 获取 RPC 服务地址
func GetRPCAddr() string {
	return viper.GetString("rpc_addr")
}

// GetRPCPort 获取 RPC 服务端口
func GetRPCPort() int {
	return viper.GetInt("rpc_port")
}

// GetListenClientAddr 获取客户端监听地址
func GetListenClientAddr() string {
	return viper.GetString(ListenClientAddr)
}

// GetListenClientPort 获取监听客户端端口
func GetListenClientPort() int {
	return viper.GetInt(ListenClientPort)
}
