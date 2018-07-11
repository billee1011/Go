package mgr

import (
	"sync"
	"sync/atomic"
)

type DeskMgr struct {
	deskMap sync.Map // deskID -> desk
	maxID uint64
}

var deskMgr DeskMgr

func init(){
	deskMgr = DeskMgr{maxID:0}
}

func GetDeskMgr() DeskMgr{
	return deskMgr
}

func (mgr DeskMgr) AllocDeskID() (uint64, error) {
	return atomic.AddUint64(&mgr.maxID, 1), nil
}
