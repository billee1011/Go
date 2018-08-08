package conn

import "github.com/spf13/viper"

type Config struct {
	Address string
	Port    int
}

//向大数据端发送日志通用
func GetReportClientConfig() Config{
	return Config{
		Address:viper.GetString("report_server_ip"),
		Port:viper.GetInt("report_server_port"),
	}
}