package loader

import "strconv"

/*
  功能： 命令行启动参数解析和管理, 暂时只支持Int和string两种命令行参数
  作者： SkyWang
  日期: 2018-7-26
*/

// 定义可能会解析的命令行字段
/*
	// 添加通用的命令行启动参数
	mapArgs["port"] = rootCmd.Flags().String("port", "", "server rpc port")
	mapArgs["hport"] = rootCmd.Flags().String("hport", "", "server rpc health port")
	mapArgs["gid"] = rootCmd.Flags().String("gid", "", "group id")
	mapArgs["sid"] = rootCmd.Flags().String("sid", "", "server id")
	mapArgs["type"] = rootCmd.Flags().String("type", "", "server type")
	mapArgs["level"] = rootCmd.Flags().String("level", "", "server level")
	mapArgs["data"] = rootCmd.Flags().String("data", "", "server data")
*/
// serviceloader gold -sid=1 -gid=1000 -type=0 -data=333
var mapStringFlag = map[string]string{}

func IntArg(key string) (int64, bool) {
	s, ok := StringArg(key)
	if !ok {
		return 0, false
	}
	if len(s) == 0 {
		return 0, false
	}
	ret, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return ret, true
}

func StringArg(key string) (string, bool) {
	v , ok := mapStringFlag[key]
	return v, ok
}

func SetArg(key string, value string) {
	mapStringFlag[key] = value
}
