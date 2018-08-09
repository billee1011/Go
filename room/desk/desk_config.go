package desk

import (
	"steve/room/fixed"
	"steve/room/util"
)

type DeskConfig struct {
	Models    []string
	Context   interface{} //预留gameContext
	Settle    DeskSettler
	PlayerIds []uint64
	Num       int
	MinScore  uint64 // 金豆准入下线
	MaxScore  uint64 // 金豆准入上限
	BaseScore uint64 // 底分
}

//默认自带的
var defaultModels = []string{fixed.PlayerModelName, fixed.MessageModelName, fixed.RequestModelName, fixed.ChatModelName, fixed.EventModelName, fixed.ContinueModelName}

// NewMjDeskCreateConfig 麻将
func NewMjDeskCreateConfig(context interface{}, settle DeskSettler, num int) DeskConfig {
	merage := [][]string{defaultModels}
	names := util.MergeStringArray(merage)
	return DeskConfig{
		Models:  names,
		Context: context,
		Num:     num,
		Settle:  settle,
	}
}

// NewDDZMDeskCreateConfig 斗地主
func NewDDZMDeskCreateConfig(context interface{}, num int) DeskConfig {
	merage := [][]string{defaultModels}
	names := util.MergeStringArray(merage)
	return DeskConfig{
		Models:  names,
		Context: context,
		Num:     num,
	}
}

// NewDeskCreateConfigDefault 包含基础model
func NewDeskCreateConfigDefault(context interface{}, models ...string) DeskConfig {
	merage := [][]string{defaultModels, models}
	names := util.MergeStringArray(merage)
	return DeskConfig{
		Models:  models,
		Context: names,
	}
}
