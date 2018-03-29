# serviceloader 配置指南

## 配置选项表

配置项 | 类型 | 作用 | 默认值 | 备注
----- | ---- |-------- | ----- | ------
log_level | string | log日志等级。 | info |  可取值为 debug, info, warning, error, fatal
log_dir   | string | log日志目录 | 空|  为空时不记录日志文件 
log_file | string | log日志文件名前缀 | 空 | 实际的log日志文件为 prefix_年_月_日_时.log， 为空时不记录日志 
log_stderr | bool | 是否输出日志到标准错误 | true | 和输出到文件互不影响，可以同时输出标准错误和文件

