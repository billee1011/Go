package redisfactory

import "github.com/go-redis/redis"

// RedisFactory redis 工厂，用来创建 redis 客户端
// RedisFactory 会从配置中加载 redis 服务的地址信息，所以不用在参数中指定具体地址
// 考虑到后面可能使用集群，所以不指定 Database
// 多次调用 NewClient 会使用之前已经创建的 redis.Client 返回
// NewClient 会校验是否能连接成功， 如果连接不成功会返回错误
type RedisFactory interface {
	NewClient() (*redis.Client, error)
}
