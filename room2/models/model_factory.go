package models

import (
	"steve/room2/fixed"
	"steve/room2/desk"
)

const (

)

func CreateModel(name string, desk *desk.Desk) DeskModel {
	var result DeskModel = nil
	switch name {
	case fixed.Event:
		result = NewMjEventModel(desk)
	case fixed.Player:
	case fixed.Request:
	case fixed.Message:
	case fixed.Trusteeship:
	case fixed.Chat:
	}
	return result
}
