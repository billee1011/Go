package global

import "steve/gateway/interfaces"

var gConnectPlayerMap interfaces.ConnectPlayerMap

// GetConnectPlayerMap 获取全局 ConnectPlayerMap
func GetConnectPlayerMap() interfaces.ConnectPlayerMap {
	return gConnectPlayerMap
}

// SetConnectPlayerMap 设置全局 ConnectPlayerMap
func SetConnectPlayerMap(cpm interfaces.ConnectPlayerMap) {
	gConnectPlayerMap = cpm
}
