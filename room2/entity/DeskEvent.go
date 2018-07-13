package entity

import "steve/room2/abs"

type DeskEvent struct {
	EventID int
	ParamLen int
	Params interface{}
	Desk abs.Desk
}