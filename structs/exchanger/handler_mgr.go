package exchanger

import (
	"fmt"
	"reflect"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// Handler 消息处理器
type Handler struct {
	// HandlerFunc 消息处理函数，具体类型参考 exchanger.RegisterHandle 函数
	HandlerFunc interface{}
	// MsgType 通过反射获取到的消息实际类型
	MsgType reflect.Type
}

// HandlerMgr 消息处理器管理器
type HandlerMgr interface {
	// RegisterHandle 注册指定消息 ID 的回调函数， 当收到消息时， 会回调 handler 处理
	// handler 的声明可以是 func(clientID uint64, head *steve_proto_gaterpc.Header, body YourProtoType) []ResponseMsg
	// 		handler 的参数中 clientID 在网关服表示客户端连接 ID， 在其他应用服中为 玩家 ID
	//		head 为消息头， YourProtoType 可以为任意 proto 类型,
	// 		handler 的返回值 []proto.Message 表示需要回复的数据， 为 nil 或者空切片时则表示不需要回复， 此时服务仍可以通过 SendPackage 或者 BroadcastPackage 来回复消息
	// handler 的声明也可以是 func(clientID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) []ResponseMsg
	// 		和前一个类似，但是 bodyData 是字节数组，由应用层自己反序列化
	// 通过返回值和通过 SendPackage 方法回复消息的区别：
	// 		通过返回值回复，网关会填充回复序号为客户端的请求序号， 客户端可以根据序号作相应处理
	// 		通过 SendPackage 或者 BroadcastPackage 回复消息时， 网关不再拥有之前的请求现场， 因此不能作特定处理。但服务仍可以使用
	//		消息体自定义的一些直通数据来作相应的处理。
	// 		总而言之， 如果没有特殊需求， 使用返回值来回复可以得到更简洁的效果， 但使用 SendPackage 或者 BroadcastPackage 来回复消息可以作更多的自定义处理
	RegisterHandle(msgID uint32, handlerFunc interface{}) error

	// GetHandler 获取消息处理器
	GetHandler(msgID uint32) *Handler
}

var byteSliceType = reflect.TypeOf([]byte{})

// CallHandler 调用消息处理器
func CallHandler(handler *Handler, clientID uint64, header *steve_proto_gaterpc.Header, body []byte) ([]ResponseMsg, error) {

	var callResults []reflect.Value
	f := reflect.ValueOf(handler.HandlerFunc)

	if handler.MsgType == byteSliceType {
		callResults = f.Call([]reflect.Value{
			reflect.ValueOf(clientID),
			reflect.ValueOf(header),
			reflect.ValueOf(body),
		})
	} else {
		bodyMsg := reflect.New(handler.MsgType).Interface()
		if err := proto.Unmarshal(body, bodyMsg.(proto.Message)); err != nil {
			return []ResponseMsg{}, fmt.Errorf("反序列化消息体失败: %v", err)
		}
		callResults = f.Call([]reflect.Value{
			reflect.ValueOf(clientID),
			reflect.ValueOf(header),
			reflect.ValueOf(bodyMsg).Elem(),
		})
	}
	if callResults == nil || len(callResults) == 0 || callResults[0].IsNil() {
		return []ResponseMsg{}, nil
	}
	return callResults[0].Interface().([]ResponseMsg), nil
}
