package public

import "steve/room2/desk"

type BaseModel struct {
	desk  *desk.Desk
}

func (model BaseModel) GetDesk() *desk.Desk {
	return model.desk
}

func (model BaseModel) SetDesk(desk *desk.Desk) {
	model.desk = desk
}