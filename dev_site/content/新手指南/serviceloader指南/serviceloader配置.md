---
title: "serviceloader配置指南"
date: 2018-03-13T11:57:59+08:00
author: 安佳玮
draft: false
---

## 概要介绍

* 配置可以使用 yml 或者 json 格式
* 使用 serviceloader 加载 plugin 时，使用 config 参数指定配置文件路径。详情参考 `serviceloader --help`


## log 配置

* serviceloader 使用 logrus 作为日志框架
* plugin 可以使用 logrus 中的所有接口
* 配置表：

配置项 | 类型 | 作用 | 默认值 | 示例 |备注 
----- | ---- |-------- | ----- | ----- | ------
log_level | string | log日志等级。 | info | debug|  可取值为 debug, info, warning, error, fatal
log_dir   | string | log日志目录 | 空| ./log |为空时不记录日志文件 
log_file | string | log日志文件名前缀 | 空 | mylog | 实际的log日志文件为 prefix_年_月_日_时.log， 为空时不记录日志 
log_stderr | bool | 是否输出日志到标准错误 | true | true | 和输出到文件互不影响，可以同时输出标准错误和文件


* yml 格式配置示例

```
log_level: debug 
log_dir: ./log
log_file: gateway
log_stderr: true 
```

## RPC 配置

* RPC 为其它进程提供接口
* 配置好 RPC 后， serviceloader 会默认启动几个特殊的服务
    - 提供给 consul 的健康检查服务（暂未实现）
    - 提供给 gateway 使用的消息派发服务（服务原型参考： steve/proto/gate_rpc/gate_rpc.proto）。 注意： 这个服务是 serviceloader 自动启动的， 不需要 plugin 再去实现。 具体使用方法参考 [服务编码指南](/serviceloader指南/服务编码指南.md)
* 正确配置 RPC 后， serviceloader 将会自动使用服务名称注册到 consul。

* 配置表： 

配置项 | 类型 | 作用 | 默认值 | 示例 |备注 
----- | ---- |-------- | --- | ----- | ------
rpc_addr | string | 服务 RPC 接口监听地址 | 空 | 127.0.0.1 | 为空时不启动 RPC 服务
rpc_port | int | 服务 RPC 监听端口 | 0 | 36001 | 为 0 时不启动 RPC 服务
rpc_server_name | string | RPC 服务名称 | 空 | hall | RPC 名称将会注册到 consul， 其他服务使用该名称来查找相应的服务

* yml 配置示例：

```
rpc_addr: 127.0.0.1
rpc_port: 37001
rpc_server_name: hall
```

