# pprof

*加上了根路径查看所有信息的入口页面，直接ip+对应端口就可以。*

`http://localhost:9909/`

用于输出pprof调试信息，第一版网关的socket连接接入了自定义的profile，能查看连接数，其他是一些内置的profile

## 配置

config.yml新增配置：

### pprofExposeType 输出类型

#### file 
保存自定义profile到运行目录，每次修改时写入到文件，不建议使用

#### http 
pprof默认的开一个http服务端口输出

通过 http://localhost:9909/debug/pprof/ 查看

#### svg 

在http基础上将输出页面的链接重新封装一次输出为svg格式的图片，必须先安装graphviz，建议开发环境使用

默认输出通过 http://localhost:9909/debug/pprof/ 查看

svg图片输出通过 http://localhost:9909/debug/pprofsvg/ 查看

### pprofHttpPort
输出类型为http或svg时，指定http服务的端口

### 配置示例
```
pprofExposeType: svg
pprofHttpPort: 9909
```

## pprof使用
如果没有使用svg的类型输出，可以手动调用pprof命令查看并生成图片

`go tool pprof -png http://localhost:9909/debug/pprof/my_experiment_thing?debug=1 > prof.png`

参考：https://medium.com/@cep21/creating-custom-go-profiles-with-pprof-b737dfc58e11