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
	HandlerMgr
	// // SendPackage 发送消息给指定客户端 clientID
	// // head 为消息头
	// // body 为任意 proto 消息
	// SendPackage(clientID uint64, head *steve_proto_gaterpc.Header, body proto.Message) error

	// // BraodcastPackage 和 SendPackage 类似， 但将消息发给多个用户。 clientIDs 为客户端连接 ID 数组
	// BroadcastPackage(clientIDs []uint64, head *steve_proto_gaterpc.Header, body proto.Message) error

	// // SendPackage 发送消息给指定客户端 clientID
	// // head 为消息头
	// // body 为任意 序列化 消息
	// SendPackageBare(clientID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error

	// // BraodcastPackage 和 SendPackage 类似， 但将消息发给多个用户。 clientIDs 为客户端连接 ID 数组
	// BroadcastPackageBare(clientIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error

	// SendPackageByPlayerID 发送消息给玩家
	SendPackageByPlayerID(playerID uint64, head *steve_proto_gaterpc.Header, body proto.Message) error

	// BroadcastPackageByPlayerID 和 SendPackageByPlayerID 类似，但发送给多个玩家
	BroadcastPackageByPlayerID(playerIDs []uint64, head *steve_proto_gaterpc.Header, body proto.Message) error

	// SendPackageBareByPlayerID like SendPackageByPlayerID， 发送裸包给玩家
	SendPackageBareByPlayerID(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error

	// BroadcastPackageBareByPlayerID like SendPackageBareByPlayerID， 发送给多个玩家
	BroadcastPackageBareByPlayerID(playerIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error
}
