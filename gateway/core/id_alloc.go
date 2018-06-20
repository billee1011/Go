package core

import "steve/common/data/connect"

type idAllocator struct{}

func (ida *idAllocator) NewClientID() uint64 {
	id, _ := connect.AllocConnectID()
	return id
}
