package arg
/*
  功能： App启动命令行参数和yml配置文件参数获取
 */

type Option interface {

	// yml文件选项参数
	StringOption(key string) string
	IntOption(key string) int64

	// 命令行启动参数
	StringArg(key string) string
	IntArg(key string) int64
}