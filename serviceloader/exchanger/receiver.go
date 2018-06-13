package exchanger

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	iexchanger "steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

var byteSliceType = reflect.TypeOf([]byte{})

type handler struct {
	// handleFunc 消息处理函数，具体类型参考 iexchanger.Exchanger 函数
	handleFunc interface{}
	// msgType 通过反射获取到的消息实际类型
	msgType reflect.Type
}

type receiver struct {
	handlerMap sync.Map
}

func (r *receiver) register(msgID uint32, handlerFunc interface{}) error {
	// TODO 判断消息 ID 范围
	msgType, err := r.checkHandlerFunc(handlerFunc)
	if err != nil {
		return err
	}
	if _, loaded := r.handlerMap.LoadOrStore(msgID, handler{
		handleFunc: handlerFunc,
		msgType:    msgType,
	}); loaded {
		return fmt.Errorf("该消息 ID 已经被注册过了")
	}
	return nil
}

// HandleClientMessage 处理客户端消息
func (r *receiver) HandleClientMessage(ctx context.Context, msg *steve_proto_gaterpc.ClientMessage) (*steve_proto_gaterpc.HandleResult, error) {
	header := msg.GetHeader()
	msgID := header.GetMsgId()
	handler := r.handlerOf(msgID)
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "receiver.HandleClientMessage",
		"msg_id":    msgID,
		"client_id": msg.GetClientId(),
	})
	if handler == nil {
		logEntry.Warnln("没有对应的消息处理器")
		return &steve_proto_gaterpc.HandleResult{}, nil
	}

	responses := r.callHandler(logEntry, handler, msg.GetClientId(), header, msg.GetRequestData())
	return r.packResults(responses)
}

// callHandler 根据消息类型反序列化消息体和回调处理器
func (r *receiver) callHandler(logEntry *logrus.Entry, handler *handler, clientID uint64,
	header *steve_proto_gaterpc.Header, body []byte) []iexchanger.ResponseMsg {

	var callResults []reflect.Value
	f := reflect.ValueOf(handler.handleFunc)

	if handler.msgType == byteSliceType {
		callResults = f.Call([]reflect.Value{
			reflect.ValueOf(clientID),
			reflect.ValueOf(header),
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
			reflect.ValueOf(header),
			reflect.ValueOf(bodyMsg).Elem(),
		})
	}
	if callResults == nil || len(callResults) == 0 || callResults[0].IsNil() {
		return []iexchanger.ResponseMsg{}
	}
	return callResults[0].Interface().([]iexchanger.ResponseMsg)
}

// packResults 将应答消息打包返回
func (r *receiver) packResults(responses []iexchanger.ResponseMsg) (*steve_proto_gaterpc.HandleResult, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "receiver.packResults",
	})
	resultMessages := []*steve_proto_gaterpc.ResponseMessage{}
	for _, resp := range responses {
		bodyData, err := proto.Marshal(resp.Body)
		if err != nil {
			logEntry.WithField("msg_id", resp.MsgID).Errorln("消息序列化失败")
			return nil, errors.New("消息序列化失败")
		}
		resultMessages = append(resultMessages, &steve_proto_gaterpc.ResponseMessage{
			Header: &steve_proto_gaterpc.Header{
				MsgId: resp.MsgID,
			},
			Body: bodyData,
		})
	}
	return &steve_proto_gaterpc.HandleResult{
		Responses: resultMessages,
	}, nil
}

func (r *receiver) handlerOf(msgID uint32) *handler {
	v, ok := r.handlerMap.Load(msgID)
	if !ok || v == nil {
		return nil
	}
	h := v.(handler)
	return &h
}

// checkHandlerParams 检查消息回调函数参数是否符合规范
func (r *receiver) checkHandlerParams(handlerType reflect.Type) (reflect.Type, error) {
	if handlerType.NumIn() != 3 {
		return nil, fmt.Errorf("处理函数需要接受 3 个参数")
	}
	if handlerType.In(0).Kind() != reflect.Uint64 {
		return nil, fmt.Errorf("处理函数的第 1 个参数必须是 uint64，表示客户端 ID")
	}
	var header *steve_proto_gaterpc.Header
	if handlerType.In(1) != reflect.TypeOf(header) {
		return nil, fmt.Errorf("处理函数的第 2 个参数必须是 *exchanger.MessageHeader")
	}
	msgType := handlerType.In(2)
	if byteSliceType != msgType {
		msg := reflect.New(msgType)
		if _, ok := msg.Interface().(proto.Message); !ok {
			return nil, fmt.Errorf("处理函数的第 3 个参数必须是 proto.Message 类型或者 []byte")
		}
	}
	return msgType, nil
}

// checkHandlerReturns 检查消息回调函数的返回值是否符合规范
func (r *receiver) checkHandlerReturns(handlerType reflect.Type) error {
	if handlerType.NumOut() != 1 {
		return fmt.Errorf("处理函数需要返回 1 个值")
	}
	retType := handlerType.Out(0)
	if !retType.ConvertibleTo(reflect.TypeOf([]iexchanger.ResponseMsg{})) {
		return fmt.Errorf("处理函数的返回值需要可以转换成 []ResponseMsg 类型")
	}
	return nil
}

// checkHandlerFunc 检查消息回调函数是否符合要求
func (r *receiver) checkHandlerFunc(handlerFunc interface{}) (reflect.Type, error) {
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		return nil, fmt.Errorf("参数错误，第二个参数需要是函数")
	}
	if err := r.checkHandlerReturns(handlerType); err != nil {
		return nil, err
	}
	return r.checkHandlerParams(handlerType)
}
