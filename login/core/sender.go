package core

import (
	"steve/structs/net"
	"steve/structs/proto/base"
)

type sender struct {
	watchDog net.WatchDog
}

func (s *sender) SendMessage(clientID uint64, header *steve_proto_base.Header, body []byte) error {
	return s.watchDog.SendPackage(clientID, header, body)
}
