package utils

import (
	"sync/atomic"
)

var DefaultIDAlloc = NewIDAlloc()

func NewIDAlloc() *IDAlloc {
	return &IDAlloc{
		maxID: 0,
	}
}

type IDAlloc struct {
	maxID uint64
}

func (ia *IDAlloc) AllocDeskID() uint64 {
	return atomic.AddUint64(&ia.maxID, 1)
}
