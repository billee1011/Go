package interfaces

import (
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// MessageSender 消息发送器
type MessageSender interface {
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
