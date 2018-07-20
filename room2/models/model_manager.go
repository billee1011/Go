package models

import (
	"sync"
	"github.com/Sirupsen/logrus"
	"steve/room2/desk"
)

type ModelManager struct {
	modelMap sync.Map //deskid-[model] //models    map[string]*models.DeskModel
}

var manager *ModelManager

func init() {
	manager = &ModelManager{}
}

func GetModelManager() *ModelManager {
	return manager
}

func (manager *ModelManager) InitDeskModel(deskId uint64, modelName []string, desk *desk.Desk) {
	modelMap := make(map[string]*DeskModel, len(modelName))
	for _, name := range modelName {
		model := CreateModel(name, desk)
		if model == nil {
			logrus.Error("创建Model失败[" + name + "]")
			continue
		}
		modelMap[name] = &model
	}
	manager.modelMap.Store(deskId, modelMap)
}

func (manager *ModelManager) RemoveDeskModel(deskId uint64){
	manager.modelMap.Delete(deskId)
}

func (manager *ModelManager) GetChatModel(deskId uint64) *ChatModel {
	return getModel(deskId).(*ChatModel)
}
func (manager *ModelManager) GetRequestModel(deskId uint64) *RequestModel {
	return getModel(deskId).(*RequestModel)
}
func (manager *ModelManager) GetMessageModel(deskId uint64) *MessageModel {
	return getModel(deskId).(*MessageModel)
}
func (manager *ModelManager) GetMjEventModel(deskId uint64) *MjEventModel {
	return getModel(deskId).(*MjEventModel)
}
func (manager *ModelManager) GetPlayerModel(deskId uint64) *PlayerModel {
	return getModel(deskId).(*PlayerModel)
}

func getModel(deskId uint64) interface{} {
	model, _ := manager.modelMap.Load(deskId)
	return model
}
