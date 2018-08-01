package settle

import (
	"steve/gutils"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/room/settle/majong"
)

type factory struct{}

func (f *factory) CreateDeskSettler(gameID int) interfaces.DeskSettler {
	switch gameID {
	case gutils.SCXLGameID:
		return majong.NewMajongSettle()
	case gutils.SCXZGameID:
		return majong.NewMajongSettle()
	default:
		return majong.NewMajongSettle()
	}
}

func init() {
	global.SetDeskSettleFactory(new(factory))
}
