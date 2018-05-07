package core

import (
	"errors"
	"fmt"
	"reflect"
	iexchanger "steve/structs/exchanger"
	"steve/structs/net"
	"steve/structs/proto/base"
	"steve/structs/proto/gate_rpc"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// wrapHandler 包装了的消息处理器
type wrapHandler struct {
	// handleFunc 消息处理函数，具体类型参考 iexchanger.Exchanger 函数
	handleFunc interface{}
	// msgType 通过反射获取到的消息实际类型
	msgType reflect.Type
}

type exchangerImpl struct {
	// handleMap 存储注册的消息处理器
	// key 为消息 ID uint32
	// value 为消息处理器 wrapHandler
	handleMap sync.Map

	// watchDog
	watchDog net.WatchDog
}

var _ iexchanger.Exchanger = new(exchangerImpl)

func (e *exchangerImpl) RegisterHandle(msgID uint32, handler interface{}) error {
	// TODO 判断消息 ID 范围
	funcType := reflect.TypeOf(handler)
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("参数错误，第二个参数需要是函数")
	}
	if funcType.NumIn() != 3 {
		return fmt.Errorf("处理函数需要接受 3 个参数")
	}
	if funcType.NumOut() != 1 {
		return fmt.Errorf("处理函数需要返回 1 个值")
	}
	if funcType.In(0).Kind() != reflect.Uint64 {
		return fmt.Errorf("处理函数的第 1 个参数必须是 uint64，表示客户端 ID")
	}
	var header *steve_proto_gaterpc.Header
	if funcType.In(1) != reflect.TypeOf(header) {
		return fmt.Errorf("处理函数的第 2 个参数必须是 *exchanger.MessageHeader")
	}
	msgType := funcType.In(2)
	msg := reflect.New(msgType)
	if _, ok := msg.Interface().(proto.Message); !ok {
		return fmt.Errorf("处理函数的第 3 个参数必须是 proto.Message 类型")
	}

	retType := funcType.Out(0)
	if !retType.ConvertibleTo(reflect.TypeOf([]iexchanger.ResponseMsg{})) {
		return fmt.Errorf("处理函数的返回值需要可以转换成 []proto.Message 类型")
	}

	if _, loaded := e.handleMap.LoadOrStore(msgID, wrapHandler{
		handleFunc: handler,
		msgType:    msgType,
	}); loaded {
		return fmt.Errorf("该消息 ID 已经被注册过了")
	}
	return nil
}

func (e *exchangerImpl) SendPackage(clientID uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	return e.BroadcastPackage([]uint64{clientID}, head, body)
}

func (e *exchangerImpl) BroadcastPackage(clientIDs []uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	entry := logrus.WithFields(logrus.Fields{
		"name":      "exchangerImpl.SendPackage",
		"client_id": clientIDs,
		"msg_id":    head.MsgId,
	})
	bodyData := []byte{}
	if err := proto.Unmarshal(bodyData, body); err != nil {
		var errUnmarshal = errors.New("序列化消息体失败")
		entry.WithError(err).Errorln(errUnmarshal)
		return errUnmarshal
	}
	e.BroadcastPackageBare(clientIDs, head, bodyData)
	return nil
}

// SendPackage 发送消息给指定客户端 clientID
// head 为消息头
// body 为任意 序列化 消息
func (e *exchangerImpl) SendPackageBare(clientID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	return e.BroadcastPackageBare([]uint64{clientID}, head, bodyData)
}

// BraodcastPackage 和 SendPackage 类似， 但将消息发给多个用户。 clientIDs 为客户端连接 ID 数组
func (e *exchangerImpl) BroadcastPackageBare(clientIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	header := steve_proto_base.Header{
		MsgId: proto.Uint32(head.MsgId),
	}
	e.watchDog.BroadPackage(clientIDs, &header, bodyData)
	return nil
}

func (e *exchangerImpl) getHandler(msgID uint32) *wrapHandler {
	v, ok := e.handleMap.Load(msgID)
	if !ok || v == nil {
		return nil
	}
	h := v.(wrapHandler)
	return &h
}
