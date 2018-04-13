package redisfactory

// RedisFactory redis 工厂，用来创建 redis 客户端
// RedisFactory 会从配置中加载 redis 服务的地址信息，所以不用在参数中指定具体地址
type RedisFactory interface {
	NewClient()
}
