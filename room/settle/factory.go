package settle

import (
	"steve/gutils"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
)

type factory struct{}

func (f *factory) CreateDeskSettler(gameID int) interfaces.DeskSettler {
	switch gameID {
	case gutils.SCXLGameID:
		return newScxlSettle()
	case gutils.SCXZGameID:
		return newScxzSettle()
	default:
		return new(nullSettler)
	}
}

func init() {
	global.SetDeskSettleFactory(new(factory))
}
