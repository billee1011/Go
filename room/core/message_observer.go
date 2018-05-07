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

var byteSliceType = reflect.TypeOf([]byte{})

// callHandler 根据消息类型反序列化消息体和回调处理器
func (o *messageObserver) callHandler(logEntry *logrus.Entry, handler *wrapHandler, clientID uint64,
	header *steve_proto_base.Header, body []byte) []iexchanger.ResponseMsg {

	callHeader := steve_proto_gaterpc.Header{
		MsgId: header.GetMsgId(),
	}
	var callResults []reflect.Value
	f := reflect.ValueOf(handler.handleFunc)

	if handler.msgType == byteSliceType {
		callResults = f.Call([]reflect.Value{
			reflect.ValueOf(clientID),
			reflect.ValueOf(&callHeader),
			reflect.ValueOf(body),
		})
	} else {
		bodyMsg := reflect.New(handler.msgType).Interface()
		if err := proto.Unmarshal(body, bodyMsg.(proto.Message)); err != nil {
			logEntry.WithError(err).Errorln("反序列化消息体失败")
			return []iexchanger.ResponseMsg{}
		}
		callResults = f.Call([]reflect.Value{
			reflect.ValueOf(clientID),
			reflect.ValueOf(&callHeader),
			reflect.ValueOf(bodyMsg).Elem(),
		})
	}
	if callResults == nil || len(callResults) == 0 || callResults[0].IsNil() {
		return []iexchanger.ResponseMsg{}
	}
	return callResults[0].Interface().([]iexchanger.ResponseMsg)
}

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
	retMessages := o.callHandler(logEntry, handler, clientID, header, body)
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
		if bodyData == nil {
			logEntry.Panic("bodyData nil")
		}

		o.core.dog.SendPackage(clientID, &responseHeader, bodyData)
	}
}
