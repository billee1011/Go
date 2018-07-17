package models

import (
	"steve/room2/desk"
	"steve/room2/desk/models/mj"
)

const (
	Event       = "EventModel"
	Player      = "PlayerModel"
	Request     = "RequestModel"
	Message     = "MessageModel"
	Trusteeship = "TrusteeshipModel"
)

func CreateModel(name string, desk *desk.Desk) DeskModel {
	var result DeskModel = nil
	switch name {
	case Event:
		result = mj.NewMjEventModel(desk)
	case Player:
	case Request:
	case Message:
	case Trusteeship:
	}
	return result
}
