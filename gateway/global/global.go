package global

import "steve/gateway/interfaces"

var gConnectionManager interfaces.ConnectionManager
var gPlayerManager interfaces.PlayerManager

// GetConnectionManager 获取全局 ConnectionManager
func GetConnectionManager() interfaces.ConnectionManager {
	return gConnectionManager
}

// SetConnectionManager 设置全局 ConnectionManager
func SetConnectionManager(cm interfaces.ConnectionManager) {
	gConnectionManager = cm
}

// GetPlayerManager 获取全局 PlayerManager
func GetPlayerManager() interfaces.PlayerManager {
	return gPlayerManager
}

// SetPlayerManager 设置全局 PlayerManager
func SetPlayerManager(pm interfaces.PlayerManager) {
	gPlayerManager = pm
}
