package settle

import (
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/settle/majong"
)

type gameSettlerFactory struct {
	factories map[int]interfaces.SettlerFactory
}

func (f *gameSettlerFactory) CreateSettlerFactory(gameID int) interfaces.SettlerFactory {
	factory, exist := f.factories[gameID]
	if !exist {
		return &majong.SettlerFactory{}
	}
	return factory
}

func init() {
	factories := map[int]interfaces.SettlerFactory{
		gutils.SCXLGameID: &majong.SettlerFactory{},
		gutils.SCXZGameID: &majong.SettlerFactory{},
	}
	global.SetGameSettlerFactory(&gameSettlerFactory{
		factories: factories,
	})
}
