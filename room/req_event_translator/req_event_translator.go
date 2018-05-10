package rtoet

import (
	"errors"
	"reflect"
	"steve/client_pb/msgId"
	"steve/room/interfaces/global"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type msgTranslator struct {
	f        interface{}
	bodyType reflect.Type
}

type translator struct {
	msgTranslators map[msgid.MsgID]msgTranslator
}

var errTranslatorNotExists = errors.New("转换器不存在")
var errUnmarshalReqFailed = errors.New("反序列化请求消息体失败")

func (t *translator) Translate(playerID uint64, header *steve_proto_gaterpc.Header, bodyData []byte) (eventID server_pb.EventID, eventContext proto.Message, err error) {
	f, ok := t.msgTranslators[msgid.MsgID(header.GetMsgId())]
	if !ok {
		err = errTranslatorNotExists
		return
	}
	return t.callTranslator(f, playerID, header, bodyData)
}

func (t *translator) callTranslator(msgTranslator msgTranslator, playerID uint64,
	header *steve_proto_gaterpc.Header, bodyData []byte) (eventID server_pb.EventID, eventContext proto.Message, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "translator.callTranslator",
		"player_id": playerID,
		"msg_id":    header.GetMsgId(),
	})

	bodyMsg := reflect.New(msgTranslator.bodyType).Interface()
	if err = proto.Unmarshal(bodyData, bodyMsg.(proto.Message)); err != nil {
		logEntry.WithError(err).Errorln(errUnmarshalReqFailed)
		err = errUnmarshalReqFailed
		return
	}

	f := reflect.ValueOf(msgTranslator.f)
	callResults := f.Call([]reflect.Value{
		reflect.ValueOf(playerID),
		reflect.ValueOf(header),
		reflect.ValueOf(bodyMsg).Elem(),
	})

	eventID = callResults[0].Interface().(server_pb.EventID)
	eventContext = callResults[1].Interface().(proto.Message)
	err = callResults[2].Interface().(error)
	return
}

func (t *translator) addTranslator(msgID msgid.MsgID, f interface{}) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "translator.addTranslator",
		"msg_id":    msgID,
	})

	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		logEntry.Panic("需要函数类型")
	}
	if fType.NumIn() != 3 || fType.In(0).Kind() != reflect.Uint64 ||
		fType.In(1) != reflect.TypeOf(&steve_proto_gaterpc.Header{}) {
		logEntry.Panic("处理函数的参数错误")
	}
	bodyType := fType.In(2)
	msg := reflect.New(bodyType)
	if _, ok := msg.Interface().(proto.Message); !ok {
		logEntry.Panic("处理函数的第 3 个参数必须是 proto.Message 类型")
	}

	typeOfEventID := reflect.TypeOf(server_pb.EventID(0))
	typeOfErr := reflect.TypeOf(errors.New(""))

	if fType.NumOut() != 3 || fType.Out(0) != typeOfEventID || fType.Out(2).Name() != "error" {
		logEntry.WithFields(logrus.Fields{
			"num_out":          fType.NumOut(),
			"out_0_type":       fType.Out(0),
			"out_2_type":       fType.Out(2),
			"type_of_event_id": typeOfEventID,
			"type_of_err":      typeOfErr,
		}).Panic("处理函数的返回值类型错误")
	}
	eventContextType := fType.Out(1)
	if eventContextType.Name() != "Message" || !strings.HasSuffix(eventContextType.PkgPath(), "proto") {
		logEntry.WithFields(logrus.Fields{
			"event_context_type_name": eventContextType.Name(),
			"event_context_type_pkg":  eventContextType.PkgPath(),
		}).Panic("处理函数的第 2 个返回值类型错误")
	}

	t.msgTranslators[msgID] = msgTranslator{
		f:        f,
		bodyType: bodyType,
	}
}

func (t *translator) addTranslators() {
	// TODO 添加所有请求转事件表
	t.addTranslator(msgid.MsgID_room_huansanzhang_req, translateHuansanzhangReq)
	t.addTranslator(msgid.MsgID_room_zimo_req, translateZimoReq)
	t.addTranslator(msgid.MsgID_room_bugang_req, translateBugangReq)
	t.addTranslator(msgid.MsgID_room_angang_req, translateAngangReq)
}

func init() {
	t := &translator{
		msgTranslators: make(map[msgid.MsgID]msgTranslator, 1),
	}
	t.addTranslators()

	global.SetReqEventTranslator(t)
}
