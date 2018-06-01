package exchanger

import (
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// ResponseMsg 消息处理器的回复消息
type ResponseMsg struct {
	MsgID uint32
	Body  proto.Message
}

// Exchanger 与客户端交互接口
type Exchanger interface {

	// RegisterHandle 注册指定消息 ID 的回调函数， 当收到消息时， 会回调 handler 处理
	// handler 的声明可以是 func(clientID uint64, head *steve_proto_gaterpc.Header, body YourProtoType) []ResponseMsg
	// 		handler 的参数中 clientID 为客户端连接 ID， head 为消息头， YourProtoType 可以为任意 proto 类型,
	// 		handler 的返回值 []proto.Message 表示需要回复的数据， 为 nil 或者空切片时则表示不需要回复， 此时服务仍可以通过 SendPackage 或者 BroadcastPackage 来回复消息
	// handler 的声明也可以是 func(clientID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) []ResponseMsg
	// 		和前一个类似，但是 bodyData 是字节数组，由应用层自己反序列化
	// 通过返回值和通过 SendPackage 方法回复消息的区别：
	// 		通过返回值回复，网关会填充回复序号为客户端的请求序号， 客户端可以根据序号作相应处理
	// 		通过 SendPackage 或者 BroadcastPackage 回复消息时， 网关不再拥有之前的请求现场， 因此不能作特定处理。但服务仍可以使用
	//		消息体自定义的一些直通数据来作相应的处理。
	// 		总而言之， 如果没有特殊需求， 使用返回值来回复可以得到更简洁的效果， 但使用 SendPackage 或者 BroadcastPackage 来回复消息可以作更多的自定义处理
	RegisterHandle(msgID uint32, handler interface{}) error

	// SendPackage 发送消息给指定客户端 clientID
	// head 为消息头
	// body 为任意 proto 消息
	SendPackage(clientID uint64, head *steve_proto_gaterpc.Header, body proto.Message) error

	// BraodcastPackage 和 SendPackage 类似， 但将消息发给多个用户。 clientIDs 为客户端连接 ID 数组
	BroadcastPackage(clientIDs []uint64, head *steve_proto_gaterpc.Header, body proto.Message) error

	// SendPackage 发送消息给指定客户端 clientID
	// head 为消息头
	// body 为任意 序列化 消息
	SendPackageBare(clientID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error

	// BraodcastPackage 和 SendPackage 类似， 但将消息发给多个用户。 clientIDs 为客户端连接 ID 数组
	BroadcastPackageBare(clientIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error
}
