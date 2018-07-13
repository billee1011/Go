package models

import (
	"steve/room2/desk"
)

type DeskModel interface{
	GetName() string
	Start()
	Stop()
	GetDesk() *desk.Desk
	SetDesk(desk *desk.Desk)
}

