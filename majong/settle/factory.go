package settle

import (
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/settle/null"
	"steve/majong/settle/scxl"
)

type gameSettlerFactory struct {
	factories map[int]interfaces.SettlerFactory
}

func (f *gameSettlerFactory) CreateSettlerFactory(gameID int) interfaces.SettlerFactory {
	factory, exist := f.factories[gameID]
	if !exist {
		return &null.SettlerFactory{}
	}
	return factory
}

func init() {
	factories := map[int]interfaces.SettlerFactory{
		gutils.SCXLGameID: &scxl.SettlerFactory{},
		gutils.SCXZGameID: &scxl.SettlerFactory{},
	}
	global.SetGameSettlerFactory(&gameSettlerFactory{
		factories: factories,
	})
}
