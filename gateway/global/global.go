package global

import "steve/gateway/interfaces"

var gConnectionManager interfaces.ConnectionManager

// GetConnectionManager 获取全局 ConnectionManager
func GetConnectionManager() interfaces.ConnectionManager {
	return gConnectionManager
}

// SetConnectionManager 设置全局 ConnectionManager
func SetConnectionManager(cm interfaces.ConnectionManager) {
	gConnectionManager = cm
}
