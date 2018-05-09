package exchanger

import (
	"context"
	"fmt"
	"reflect"
	"steve/structs"
	"steve/structs/common"
	iexchanger "steve/structs/exchanger"
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
		return fmt.Errorf("处理函数的返回值需要可以转换成 []ResponseMsg 类型")
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

	g := structs.GetGlobalExposer()
	if g == nil {
		entry.Error("获取全局对象失败")
		return fmt.Errorf("获取全局对象失败")
	}
	// TODO 网关服务绑定， 不同的网关分开发送
	cc, err := g.RPCClient.GetClientConnByServerName(common.GateServiceName)
	if err != nil {
		entry.WithError(err).Warn("获取客户端连接失败")
		return fmt.Errorf("获取客户端连接失败： %v", err)
	}

	data, err := proto.Marshal(body)
	if err != nil {
		entry.WithError(err).Warn("消息序列化失败")
		return fmt.Errorf("消息序列化失败： %v", err)
	}

	mc := steve_proto_gaterpc.NewMessageSenderClient(cc)
	r, err := mc.SendMessage(context.Background(), &steve_proto_gaterpc.SendMessageRequest{
		ClientId: clientIDs,
		Header:   head,
		Data:     data,
	})

	if err != nil || r == nil {
		entry.WithError(err).Error("调用 RPC 接口失败")
		return fmt.Errorf("调用 RPC 接口失败: %v", err)
	}
	if !r.GetOk() {
		entry.Info("网关发消息返回失败")
		return fmt.Errorf("网关发送消息返回失败")
	}
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
