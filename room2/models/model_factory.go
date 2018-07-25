package models

import (
	"steve/room2/fixed"
	"steve/room2/desk"
)


func CreateModel(name string, desks *desk.Desk) DeskModel {
	var result DeskModel = nil
	switch name {
	case fixed.Event:
		result = NewMjEventModel(desks)
	case fixed.Player:
		result = NewPlayertModel(desks)
	case fixed.Request:
		result = NewRequestModel(desks)
	case fixed.Message:
		result = NewMessageModel(desks)
	case fixed.Chat:
		result = NewChatModel(desks)
	}
	print(result.GetDesk() == nil)
	return result
}
