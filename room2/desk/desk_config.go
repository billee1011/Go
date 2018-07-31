package desk

import (
	"steve/room2/fixed"
	"steve/room2/util"
)

type DeskConfig struct {
	Models    []string
	Context   interface{} //预留gameContext
	Settle    DeskSettler
	PlayerIds []uint64
	Num       int
}

//默认自带的
var defaultModels = []string{fixed.Player, fixed.Message, fixed.Request, fixed.Chat, fixed.Event, fixed.Continue}

//麻将
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

//斗地主
func NewDDZMDeskCreateConfig(context interface{}, num int) DeskConfig {
	merage := [][]string{defaultModels}
	names := util.MergeStringArray(merage)
	return DeskConfig{
		Models:  names,
		Context: context,
		Num:     num,
	}
}

//包含基础model
func NewDeskCreateConfigDefault(context interface{}, models ...string) DeskConfig {
	merage := [][]string{defaultModels, models}
	names := util.MergeStringArray(merage)
	return DeskConfig{
		Models:  models,
		Context: names,
	}
}
