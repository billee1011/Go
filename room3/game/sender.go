package game

import (
	"github.com/gogo/protobuf/proto"
	"steve/client_pb/msgid"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
)

type Sender struct {
	e exchanger.Exchanger
}

var DefaultSender = NewSender()

func NewSender() *Sender {
	return new(Sender)
}

func (s *Sender) SetSender(e exchanger.Exchanger) {
	s.e = e
}

func (s *Sender) SendMessageByPlayer(playerID uint64, msgID msgid.MsgID, body proto.Message) error {
	return s.e.SendPackageByPlayerID(playerID, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgID),
	}, body)
}

func (s *Sender) BroadCastMessageBare(playerIDs []uint64, msgID msgid.MsgID, body []byte) error {
	return s.e.BroadcastPackageBareByPlayerID(playerIDs, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgID),
	}, body)
}
