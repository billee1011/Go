package models

import (
	"steve/room/desk"
	"steve/room/fixed"
)

// CreateModel 创建 model
func CreateModel(name string, desk *desk.Desk) DeskModel {
	var result DeskModel
	switch name {
	case fixed.EventModelName:
		result = createEventModel(desk)
	case fixed.PlayerModelName:
		result = NewPlayertModel(desk)
	case fixed.RequestModelName:
		result = NewRequestModel(desk)
	case fixed.MessageModelName:
		result = NewMessageModel(desk)
	case fixed.ChatModelName:
		result = NewChatModel(desk)
	case fixed.ContinueModelName:
		result = NewContinueModel(desk)
	}
	return result
}

// createEventModel 创建 event model
func createEventModel(desk *desk.Desk) DeskModel {
	switch desk.GetGameId() {
	case GameId_GAMEID_DOUDIZHU:
		return NewDDZEventModel(desk)
	default:
		return NewMjEventModel(desk)
	}
}
