package desk

import (
	"sync"
	"sync/atomic"
)

type DeskMgr struct {
	deskMap sync.Map // deskID -> desk
	maxID uint64
}

const (
	GameId_GAMEID_XUELIU   = 1
	GameId_GAMEID_XUEZHAN  = 2
	GameId_GAMEID_DOUDIZHU = 3
	GameId_GAMEID_ERRENMJ  = 4
)

var deskMgr DeskMgr

func init(){
	deskMgr = DeskMgr{maxID:0}
}

func GetDeskMgr() DeskMgr{
	return deskMgr
}


//创建桌子并初始化所有model
func (mgr DeskMgr) CreateDesk(players []uint64, gameID int) Desk{
	var config DeskConfig
	var context interface{}
	switch gameID{
	case GameId_GAMEID_DOUDIZHU:
		config = NewMjDeskCreateConfig(context)
	default:
		config = NewDDZMDeskCreateConfig(context)
	}

	id,_ := mgr.allocDeskID()
	desk := NewDesk(id,gameID,config)
	desk.InitModel()
	return desk
}

func (mgr DeskMgr) allocDeskID() (uint64, error) {
	return atomic.AddUint64(&mgr.maxID, 1), nil
}
