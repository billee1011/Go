package config

import (
	"flag"
	"os"
)

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
