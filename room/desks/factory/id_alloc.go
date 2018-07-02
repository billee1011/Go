package factory

import (
	"steve/room/interfaces/global"
	"sync/atomic"
)

type idAlloc struct {
	maxID uint64
}

func (ia *idAlloc) AllocDeskID() (uint64, error) {
	return atomic.AddUint64(&ia.maxID, 1), nil
}

func init() {
	global.SetDeskIDAllocator(&idAlloc{
		maxID: 0,
	})
}
