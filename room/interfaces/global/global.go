package global

import (
	"steve/room/interfaces"
)

var gPlayerMgr interfaces.PlayerMgr
var gDeskMgr interfaces.DeskMgr
var gDeskFactory interfaces.DeskFactory
var gMessageSender interfaces.MessageSender

// SetPlayerMgr 设置玩家管理器
func SetPlayerMgr(pm interfaces.PlayerMgr) {
	gPlayerMgr = pm
}

// GetPlayerMgr 获取玩家管理器
func GetPlayerMgr() interfaces.PlayerMgr {
	return gPlayerMgr
}

// SetDeskMgr 设置牌桌管理器
func SetDeskMgr(dm interfaces.DeskMgr) {
	gDeskMgr = dm
}

// GetDeskMgr 获取牌桌管理器
func GetDeskMgr() interfaces.DeskMgr {
	return gDeskMgr
}

// SetMessageSender 设置消息发送器
func SetMessageSender(ms interfaces.MessageSender) {
	gMessageSender = ms
}

// GetMessageSender 获取消息发送器
func GetMessageSender() interfaces.MessageSender {
	return gMessageSender
}

// SetDeskFactory 设置牌桌工厂
func SetDeskFactory(f interfaces.DeskFactory) {
	gDeskFactory = f
}

// GetDeskFactory 获取牌桌工厂
func GetDeskFactory() interfaces.DeskFactory {
	return gDeskFactory
}
