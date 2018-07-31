package models

import (
	"steve/room2/desk"
	"steve/room2/fixed"
)

// CreateModel 创建 model
func CreateModel(name string, desk *desk.Desk) DeskModel {
	var result DeskModel
	switch name {
	case fixed.Event:
		result = NewMjEventModel(desk)
	case fixed.Player:
		result = NewPlayertModel(desk)
	case fixed.Request:
		result = NewRequestModel(desk)
	case fixed.Message:
		result = NewMessageModel(desk)
	case fixed.Chat:
		result = NewChatModel(desk)
	case fixed.Continue:
		result = NewContinueModel(desk)
	}
	return result
}
