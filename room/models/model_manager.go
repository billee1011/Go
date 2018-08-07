package models

import (
	"fmt"
	"steve/room/desk"
	"steve/room/fixed"
	"sync"

	"github.com/Sirupsen/logrus"
)

type ModelManager struct {
	modelMap sync.Map //deskid-[model] //models    map[string]models.DeskModel
}

var manager *ModelManager

func init() {
	manager = &ModelManager{}
}

func GetModelManager() *ModelManager {
	return manager
}

func (manager *ModelManager) InitDeskModel(deskId uint64, modelName []string, desk *desk.Desk) {
	modelMap := make(map[string]DeskModel, len(modelName))
	manager.modelMap.Store(deskId, modelMap)
	for _, name := range modelName {
		model := CreateModel(name, desk)
		if model == nil {
			logrus.Error("创建Model失败[" + name + "]")
			continue
		}
		modelMap[name] = model
	}
	for _, model := range modelMap {
		model.Active()
	}
}

// StartDeskModel 启动所有 model
func (manager *ModelManager) StartDeskModel(deskID uint64) error {
	_models, ok := manager.modelMap.Load(deskID)
	if !ok {
		return fmt.Errorf("牌桌(%d)不存在", deskID)
	}
	models := _models.(map[string]DeskModel)
	for _, model := range models {
		model.Start()
	}
	return nil
}

// StopDeskModel 停止 models
func (manager *ModelManager) StopDeskModel(deskID uint64) error {
	_models, ok := manager.modelMap.Load(deskID)
	manager.modelMap.Delete(deskID)
	if !ok {
		return fmt.Errorf("无对象")
	}
	models := _models.(map[string]DeskModel)
	for _, model := range models {
		model.Stop()
	}
	return nil
}

func (manager *ModelManager) GetChatModel(deskId uint64) *ChatModel {
	model := manager.getModel(deskId)[fixed.ChatModelName]
	return model.(*ChatModel)
}

// GetRequestModel 获取 request model
func (manager *ModelManager) GetRequestModel(deskID uint64) *RequestModel {
	_model := GetModelManager().GetModelByName(deskID, fixed.RequestModelName)
	if _model == nil {
		logrus.WithField("desk_id", deskID).Warningln("request model 不存在")
		return nil
	}
	if model, ok := _model.(*RequestModel); ok {
		return model
	}
	logrus.WithField("desk_id", deskID).Warningln("request model 不存在")
	return nil
}

// GetMessageModel 获取 message model
func (manager *ModelManager) GetMessageModel(deskID uint64) *MessageModel {
	_model := GetModelManager().GetModelByName(deskID, fixed.MessageModelName)
	if _model == nil {
		logrus.WithField("desk_id", deskID).Warningln("message model 不存在")
		return nil
	}
	if model, ok := _model.(*MessageModel); ok {
		return model
	}
	logrus.WithField("desk_id", deskID).Warningln("message model 不存在")
	return nil
}

func (manager *ModelManager) GetPlayerModel(deskId uint64) *PlayerModel {
	model := manager.getModel(deskId)[fixed.PlayerModelName]
	return model.(*PlayerModel)
}

// GetModelByName 根据名字获取 model
func (manager *ModelManager) GetModelByName(deskID uint64, modelName string) DeskModel {
	models := manager.getModel(deskID)
	if models == nil {
		return nil
	}
	return models[modelName]
}

func (manager *ModelManager) getModel(deskId uint64) map[string]DeskModel {
	model, ok := manager.modelMap.Load(deskId)
	if !ok {
		return nil
	}
	return model.(map[string]DeskModel)
}

// GetContinueModel 获取续局 model
func GetContinueModel(deskID uint64) *ContinueModel {
	_model := GetModelManager().GetModelByName(deskID, fixed.ContinueModelName)
	if _model == nil {
		return nil
	}
	if model, ok := _model.(*ContinueModel); ok {
		return model
	}
	return nil
}

// GetEventModel 获取 event model
func GetEventModel(deskID uint64) DeskEventModel {
	_model := GetModelManager().GetModelByName(deskID, fixed.EventModelName)
	if _model == nil {
		return nil
	}
	if model, ok := _model.(DeskEventModel); ok {
		return model
	}
	return nil
}

// GetMjEventModel 获取麻将 Event model
func GetMjEventModel(deskID uint64) *MjEventModel {
	_model := GetEventModel(deskID)
	if _model == nil {
		return nil
	}
	if model, ok := _model.(*MjEventModel); ok {
		return model
	}
	return nil
}
