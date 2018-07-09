package interfaces

import (
	"steve/structs/proto/gate_rpc"
)

// ReqEventTranslator 请求到事件的转换器
type ReqEventTranslator interface {
	Translate(playerID uint64, header *steve_proto_gaterpc.Header, bodyData []byte) (eventID int, eventContext interface{}, err error)
}
