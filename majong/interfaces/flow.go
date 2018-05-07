package interfaces

import (
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// ToClientMessage 要发送给客户端的消息
type ToClientMessage struct {
	MsgID int
	Msg   proto.Message
}

// MajongFlow 麻将逻辑
type MajongFlow interface {
	GetMajongContext() *majongpb.MajongContext
	SetAutoEvent(autoEvent majongpb.AutoEvent)
	GetAutoEvent() *majongpb.AutoEvent
	ProcessEvent(eventID majongpb.EventID, eventContext []byte) error
	GetSettler(settlerType SettlerType) Settler
	PushMessages(playerIDs []uint64, msgs ...ToClientMessage)
	GetMessages() []majongpb.ReplyClientMessage
}
