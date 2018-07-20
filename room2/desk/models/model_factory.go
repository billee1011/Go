package models

import (
	"steve/room2/desk"
)

const (
	Event       = "EventModel"
	Player      = "PlayerModel"
	Request     = "RequestModel"
	Message     = "MessageModel"
	Trusteeship = "TrusteeshipModel"
	Chat = "ChatModel"
)

func CreateModel(name string, desk *desk.Desk) DeskModel {
	var result DeskModel = nil
	switch name {
	case Event:
		result = NewMjEventModel(desk)
	case Player:
	case Request:
	case Message:
	case Trusteeship:
	case Chat:
	}
	return result
}
