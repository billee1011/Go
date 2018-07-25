package models

import (
	"sync"
	"github.com/Sirupsen/logrus"
	"steve/room2/desk"
	"steve/room2/fixed"
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
	modelMap := make(map[string]interface{}, len(modelName))
	for _, name := range modelName {
		model := CreateModel(name, desk)
		if model == nil {
			logrus.Error("创建Model失败[" + name + "]")
			continue
		}
		model.Start()
		modelMap[name] = model
	}
	manager.modelMap.Store(deskId, modelMap)
}

func (manager *ModelManager) RemoveDeskModel(deskId uint64){
	manager.modelMap.Delete(deskId)
}

func (manager *ModelManager) GetChatModel(deskId uint64) *ChatModel {
	model := manager.getModel(deskId)[fixed.Chat]
	return model.(*ChatModel)
}
func (manager *ModelManager) GetRequestModel(deskId uint64) *RequestModel {
	model := manager.getModel(deskId)[fixed.Request]
	return model.(*RequestModel)
}
func (manager *ModelManager) GetMessageModel(deskId uint64) *MessageModel {
	model := manager.getModel(deskId)[fixed.Message]
	return model.(*MessageModel)
}
func (manager *ModelManager) GetMjEventModel(deskId uint64) *MjEventModel {
	model := manager.getModel(deskId)[fixed.Event]
	return model.(*MjEventModel)
}
func (manager *ModelManager) GetPlayerModel(deskId uint64) *PlayerModel {
	model := manager.getModel(deskId)[fixed.Player]
	return model.(*PlayerModel)
}

func (manager *ModelManager) getModel(deskId uint64) map[string]interface{} {
	model, _ := manager.modelMap.Load(deskId)
	return model.(map[string]interface{})
}
