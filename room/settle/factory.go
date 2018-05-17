package settle

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
)

type factory struct{}

func (f *factory) CreateDeskSettler(gameID int) interfaces.DeskSettler {
	switch gameID {
	case 1:
		return newScxlSettle()
	default:
		return new(nullSettler)
	}
}

func init() {
	global.SetDeskSettleFactory(new(factory))
}
