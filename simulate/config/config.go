package config

import (
	"flag"
	"os"
)

// const (
// 	// ServerAddr 服务器地址
// 	ServerAddr = "127.0.0.1:36001"
// 	// ClientVersion 客户端版本号
// 	ClientVersion = "1.0"

// 	// MaJongConfigURL 配牌服务(选项，配牌，玩家金币数)地址
// 	MaJongConfigURL = "http://127.0.0.1:36102"
// )

var (
	clientVersion  = flag.String("client_version", "1.0", "客户端版本号")
	loinServerAddr = flag.String("login_server_addr", "127.0.0.1:36201", "登录服地址")
	peiPaiURL      = flag.String("peipai_url", "http://127.0.0.1:36102", "配牌服务地址")
	dbPath         *string
)

// GetLoginServerAddr 获取登录服地址
func GetLoginServerAddr() string {
	return *loinServerAddr
}

// GetPeipaiURL 获取配牌 URL
func GetPeipaiURL() string {
	return *peiPaiURL
}

// GetClientVersion 获取客户端版本号
func GetClientVersion() string {
	return *clientVersion
}

// GetDBPath 获取 DB 路径
func GetDBPath() string {
	return *dbPath
}

func init() {
	// db 目录，默认为 $GOPATH, 没有设置 GOPATH 时， 默认值为 ./
	defaultDBPath, ok := os.LookupEnv("GOPATH")
	if !ok {
		defaultDBPath = "./"
	}
	dbPath = flag.String("dbpath", defaultDBPath, "db path")

	flag.Parse()
}
