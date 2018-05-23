package autoevent

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
)

type factory struct{}

func (f *factory) CreateGenerator() interfaces.DeskAutoEventGenerator {
	return newEventGenerator()
}

func init() {
	global.SetDeskAutoEventGeneratorFacotry(&factory{})
}
