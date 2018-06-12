package global

import (
	"steve/room/interfaces"
)

var gPlayerMgr interfaces.PlayerMgr
var gDeskMgr interfaces.DeskMgr
var gDeskFactory interfaces.DeskFactory
var gMessageSender interfaces.MessageSender
var gReqEventTranslator interfaces.ReqEventTranslator
var gDeskIDAllocator interfaces.DeskIDAllocator
var gSettleFactory interfaces.DeskSettlerFactory
var gDeskAutoEventGenerator interfaces.DeskAutoEventGenerator

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

// SetReqEventTranslator 设置请求到事件的转换器
func SetReqEventTranslator(ret interfaces.ReqEventTranslator) {
	gReqEventTranslator = ret
}

// GetReqEventTranslator 获取请求到事件的转换器
func GetReqEventTranslator() interfaces.ReqEventTranslator {
	return gReqEventTranslator
}

// SetDeskIDAllocator 设置桌子 ID 分配器
func SetDeskIDAllocator(alloc interfaces.DeskIDAllocator) {
	gDeskIDAllocator = alloc
}

// GetDeskIDAllocator 获取牌桌 ID 分配器
func GetDeskIDAllocator() interfaces.DeskIDAllocator {
	return gDeskIDAllocator
}

// SetDeskSettleFactory 设置牌桌结算工厂
func SetDeskSettleFactory(f interfaces.DeskSettlerFactory) {
	gSettleFactory = f
}

// GetDeskSettleFactory 获取牌桌结算工厂
func GetDeskSettleFactory() interfaces.DeskSettlerFactory {
	return gSettleFactory
}

// SetDeskAutoEventGenerator 设置自动事件产生器
func SetDeskAutoEventGenerator(g interfaces.DeskAutoEventGenerator) {
	gDeskAutoEventGenerator = g
}

// GetDeskAutoEventGenerator 获取自动事件产生器
func GetDeskAutoEventGenerator() interfaces.DeskAutoEventGenerator {
	return gDeskAutoEventGenerator
}
