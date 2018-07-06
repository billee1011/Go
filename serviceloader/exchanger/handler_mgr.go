package exchanger

import (
	"fmt"
	"reflect"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"sync"

	"github.com/golang/protobuf/proto"
)

type handlerMgr struct {
	handlerMap sync.Map
}

func (hm *handlerMgr) RegisterHandle(msgID uint32, handlerFunc interface{}) error {
	msgType, err := hm.checkHandlerFunc(handlerFunc)
	if err != nil {
		return err
	}
	if _, loaded := hm.handlerMap.LoadOrStore(msgID, exchanger.Handler{
		HandlerFunc: handlerFunc,
		MsgType:     msgType,
	}); loaded {
		return fmt.Errorf("该消息 ID 已经被注册过了")
	}
	return nil
}

func (hm *handlerMgr) GetHandler(msgID uint32) *exchanger.Handler {
	iHandler, ok := hm.handlerMap.Load(msgID)
	if !ok {
		return nil
	}
	handler := iHandler.(exchanger.Handler)
	return &handler
}

// checkHandlerParams 检查消息回调函数参数是否符合规范
func (hm *handlerMgr) checkHandlerParams(handlerType reflect.Type) (reflect.Type, error) {
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
func (hm *handlerMgr) checkHandlerReturns(handlerType reflect.Type) error {
	if handlerType.NumOut() != 1 {
		return fmt.Errorf("处理函数需要返回 1 个值")
	}
	retType := handlerType.Out(0)
	if !retType.ConvertibleTo(reflect.TypeOf([]exchanger.ResponseMsg{})) {
		return fmt.Errorf("处理函数的返回值需要可以转换成 []ResponseMsg 类型")
	}
	return nil
}

// checkHandlerFunc 检查消息回调函数是否符合要求
func (hm *handlerMgr) checkHandlerFunc(handlerFunc interface{}) (reflect.Type, error) {
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		return nil, fmt.Errorf("参数错误，第二个参数需要是函数")
	}
	if err := hm.checkHandlerReturns(handlerType); err != nil {
		return nil, err
	}
	return hm.checkHandlerParams(handlerType)
}
