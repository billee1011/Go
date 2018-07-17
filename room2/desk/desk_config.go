package desk

import (
	"steve/room2/desk/models"
	"steve/room2/util"
)

type DeskConfig struct {
	Models []string
	Context interface{} //预留gameContext
	Num int
}

//默认自带的
var defaultModels = []string{models.Event,models.Message,models.Request,models.Player,models.Trusteeship}

//麻将
func NewMjDeskCreateConfig(context interface{},num int) DeskConfig {
	merage := [][]string{defaultModels}
	names := util.MergeStringArray(merage)
	return DeskConfig{
		Models:[]string{models.Event,models.Message,models.Request,models.Player,models.Trusteeship},
		Context:context,
		Num:num,
	}
}

//斗地主
func NewDDZMDeskCreateConfig(context interface{},num int) DeskConfig {
	merage := [][]string{defaultModels}
	names := util.MergeStringArray(merage)
	return DeskConfig{
		Models:names,
		Context:context,
		Num:num,
	}
}

//包含基础model
func NewDeskCreateConfigDefault(context interface{},models...string) DeskConfig {
	merage := [][]string{defaultModels,models}
	names := util.MergeStringArray(merage)
	return DeskConfig{
		Models:models,
		Context:names,
	}
}