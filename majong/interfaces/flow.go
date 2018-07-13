package interfaces

import (
	"steve/client_pb/msgId"
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
	PushMessages(playerIDs []uint64, msgs ...ToClientMessage)
	GetMessages() []majongpb.ReplyClientMessage
}

// BroadcaseMessage 将消息广播给牌桌所有玩家
func BroadcaseMessage(flow MajongFlow, msgID msgid.MsgID, msg proto.Message) {
	mjContext := flow.GetMajongContext()
	players := []uint64{}

	for _, player := range mjContext.GetPlayers() {
		players = append(players, player.GetPalyerId())
	}
	flow.PushMessages(players, ToClientMessage{
		MsgID: int(msgID),
		Msg:   msg,
	})
}
