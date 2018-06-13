package core

import (
	"context"
	"steve/client_pb/msgId"
	"steve/structs/proto/base"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type sender struct {
	core *gatewayCore
}

var _ steve_proto_gaterpc.MessageSenderServer = new(sender)

func (mss *sender) SendMessage(ctx context.Context, req *steve_proto_gaterpc.SendMessageRequest) (*steve_proto_gaterpc.SendMessageResult, error) {
	msgID := req.GetHeader().GetMsgId()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "sender.SendMessage",
		"msg_id":    msgid.MsgID(msgID),
		"clients":   req.GetClientId(),
	})
	header := steve_proto_base.Header{
		MsgId:   proto.Uint32(msgID),
		Version: proto.String("1.0"), // TODO
	}
	result := &steve_proto_gaterpc.SendMessageResult{}

	err := mss.core.dog.BroadPackage(req.GetClientId(), &header, req.GetData())
	if err != nil {
		result.Ok = false
	} else {
		result.Ok = true
	}
	return result, nil
}
