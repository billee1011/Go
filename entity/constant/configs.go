package constant

// 配置 key 和 subkey 的定义

// ConfigKey 配置键值
type ConfigKey struct {
	Key, SubKey string
}

var (
	// ChargeItemListKey 充值系统商品列表配置
	ChargeItemListKey = ConfigKey{Key: "charge", SubKey: "item_list"}
	// ChargeDayMaxKey 每日充值上限配置
	ChargeDayMaxKey = ConfigKey{Key: "charge", SubKey: "day_max"}
)
