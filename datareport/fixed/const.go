package fixed

import "github.com/spf13/viper"

var LOG_TABLE_NAME_ONLINE = "log_online_num" //在线日志
var LOG_TABLE_NAME_REALTIME = "log_realtime_report" //实时上报
var HEAD_PARAM = viper.GetInt("head_param")
var MAX_PARAM = viper.GetInt("max_param")
var MAX_CONN_NUM = viper.GetInt("max_connection")