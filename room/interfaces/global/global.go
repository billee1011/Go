package global

import (
	"steve/room/interfaces"
)

var gPlayerMgr interfaces.PlayerMgr

// SetPlayerMgr 设置玩家管理器
func SetPlayerMgr(pm interfaces.PlayerMgr) {
	gPlayerMgr = pm
}

// GetPlayerMgr 获取玩家管理器
func GetPlayerMgr() interfaces.PlayerMgr {
	return gPlayerMgr
}
