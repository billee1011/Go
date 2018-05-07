package interfaces

import (
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	proto "github.com/golang/protobuf/proto"
)

// ReqEventTranslator 请求到事件的转换器
type ReqEventTranslator interface {
	Translate(playerID uint64, header *steve_proto_gaterpc.Header, body proto.Message) (eventID server_pb.EventID, eventContext []byte, err error)
}
