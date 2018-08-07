
# 重要信息
* 配置存储在数据库库中，配置服负责发布配置更新消息
* 对应的数据库表为 config.t_common_config
* 增加/删除配置时，请更新 configuration/doc/配置项.md 
* 推荐在 entity/constant/configs.go 中定义键值常量，方便统一管理
* 如果配置需要初始化，请在 entity/sql/config.sql 增加初始化语句
* 使用 external/configclient 包获取配置信息
* 订阅 entity/constant.UpdateConfig， 配置更新后将收到通知