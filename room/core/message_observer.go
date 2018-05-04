package core

import (
	"reflect"
	iexchanger "steve/structs/exchanger"
	"steve/structs/net"
	"steve/structs/proto/base"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type messageObserver struct {
	core *roomCore
}

var _ net.MessageObserver = new(messageObserver)

func (o *messageObserver) OnRecv(clientID uint64, header *steve_proto_base.Header, body []byte) {
	logEntry := logrus.WithField("name", "msgHandler.HandleClientMessage")

	msgID := header.GetMsgId()
	logEntry = logEntry.WithFields(logrus.Fields{
		"msg_id":    msgID,
		"client_id": clientID,
	})

	handler := o.core.exchanger.getHandler(msgID)
	if handler == nil {
		logEntry.Warnln("未处理的客户端消息")
		return
	}

	f := reflect.ValueOf(handler.handleFunc)
	bodyMsg := reflect.New(handler.msgType).Interface()
	if err := proto.Unmarshal(body, bodyMsg.(proto.Message)); err != nil {
		logEntry.WithError(err).Errorln("反序列化消息体失败")
		return
	}
	callHeader := steve_proto_gaterpc.Header{
		MsgId: header.GetMsgId(),
	}

	results := f.Call([]reflect.Value{
		reflect.ValueOf(clientID),
		reflect.ValueOf(&callHeader),
		reflect.ValueOf(bodyMsg).Elem(),
	})
	result := results[0]

	if result.IsNil() {
		return
	}
	retMessages, _ := result.Interface().([]iexchanger.ResponseMsg)
	for _, retMessage := range retMessages {
		responseHeader := steve_proto_base.Header{
			MsgId:  proto.Uint32(retMessage.MsgID),
			RspSeq: proto.Uint64(header.GetSendSeq()),
		}
		bodyData, err := proto.Marshal(retMessage.Body)
		if err != nil {
			logEntry.WithField("ret_msg_id", retMessage.MsgID).Errorln("消息反序列化失败")
			continue
		}
		o.core.dog.SendPackage(clientID, &responseHeader, bodyData)
	}

}
