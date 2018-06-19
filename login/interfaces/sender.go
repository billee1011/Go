package interfaces

import (
	"steve/structs/proto/base"
)

// MessageSender 消息发送器
type MessageSender interface {
	SendMessage(clientID uint64, header *steve_proto_base.Header, body []byte) error
}
