package core

import (
	"sync/atomic"
)

type idAllocator struct {
	maxID uint64
}

func (ida *idAllocator) NewClientID() uint64 {
	return atomic.AddUint64(&ida.maxID, 1)
}
