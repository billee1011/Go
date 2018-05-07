package settle

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
)

type factory struct{}

func (f *factory) CreateDeskSettler(gameID int) interfaces.DeskSettler {
	switch gameID {
	default:
		return new(nullSettler)
	}
}

func init() {
	global.SetDeskSettleFactory(new(factory))
}
