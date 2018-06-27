package mjoption

import "sync"

var std struct {
	gameOptionMgr     *GameOptionManager
	gameOptionMgrInit sync.Once

	settleOptionMgr     *SettleOptionManager
	settleOptionMgrInit sync.Once

	xingPaiOptionMgr     *XingPaiOptionManager
	xingPaiOptionMgrInit sync.Once

	cardtypeOptionMgr     *CardTypeOptionManager
	cardtypeOptionMgrInit sync.Once
}

// GetGameOptions 获取游戏选项
func GetGameOptions(gameID int) *GameOptions {
	std.gameOptionMgrInit.Do(initStdGameOptionMgr)
	return std.gameOptionMgr.GetGameOptions(gameID)
}

func initStdGameOptionMgr() {
	std.gameOptionMgr = NewGameOptionManager("optionconfig/mjoption.yaml")
}

// GetSettleOption 获取结算选项
func GetSettleOption(settleOptID int) *SettleOption {
	std.settleOptionMgrInit.Do(initStdSettleOptionMgr)
	return std.settleOptionMgr.GetSettleOption(settleOptID)
}

func initStdSettleOptionMgr() {
	std.settleOptionMgr = NewSettleOptionManager("optionconfig/settle/")
}

// GetXingpaiOption 获取行牌选项
func GetXingpaiOption(xingpaiOptID int) *XingPaiOption {
	std.xingPaiOptionMgrInit.Do(initStdXingpaiOptionMgr)
	return std.xingPaiOptionMgr.GetXingPaiOption(xingpaiOptID)
}

func initStdXingpaiOptionMgr() {
	std.xingPaiOptionMgr = NewXingPaiOptionManager("optionconfig/xingpai")
}

// GetCardTypeOption 获取牌型选项
func GetCardTypeOption(cardtypeOptID int) *CardTypeOption {
	std.cardtypeOptionMgrInit.Do(initCardTypeOptionMgr)
	return std.cardtypeOptionMgr.GetCardTypeOption(cardtypeOptID)
}

func initCardTypeOptionMgr() {
	std.cardtypeOptionMgr = NewCardTypeOptionManager("optionconfig/cardtype")
}
