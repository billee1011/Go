package room2

import "steve/room2/abs"

type deskImpl struct {
	uid            uint64
	gameID         int
	models []abs.DeskModel
}

func NewDesk(uid uint64, gameId int) abs.Desk {
	return deskImpl{uid: uid,
		gameID: gameId,
	}
}

func (desk deskImpl) GetUid() uint64 {
	return desk.uid
}

func (desk deskImpl) GetGameId() int {
	return desk.gameID
}

func (desk deskImpl) Start() {
	for _,v := range desk.models{
		v.Start()
	}
}

func (desk deskImpl) Stop() {
	for _,v := range desk.models{
		v.Stop()
	}
}