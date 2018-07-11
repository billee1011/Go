package facade

import (
	msgid "steve/client_pb/msgid"
	"steve/simulate/interfaces"
	"time"

	"github.com/golang/protobuf/proto"
)

// Request 发起请求，并等待响应
func Request(client interfaces.Client, msgID msgid.MsgID, body proto.Message, timeOut time.Duration, responseMsgID msgid.MsgID, responseBody proto.Message) error {
	return client.Request(interfaces.SendHead{
		Head: interfaces.Head{
			MsgID: uint32(msgID),
		},
	}, body, timeOut, uint32(responseMsgID), responseBody)
}
