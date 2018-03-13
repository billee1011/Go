package configuration

type Env string

const (
	Product = Env("product")
	Stage   = Env("stage")
	Dev     = Env("dev")
)

// 获取配置的参数
type ConfigGetParam struct {
	// 配置所处环境
	Env Env
	// 配置所处版本号，空代表使用最新版本号
	Version string
	// 配置的键值
	Key string
	// 是否以Key为前缀取所有配置
	Prefix bool
}

type Configuration interface {
	// 获取最新的配置版本
	GetConfigVer(env Env) (string, error)
	GetConfig(param *ConfigGetParam) (map[string]string, error)
}
