package utils

import (
	msgid "steve/client_pb/msgId"
	"steve/simulate/interfaces"
)

func createMsgHead(msgID msgid.MsgID) interfaces.SendHead {
	return interfaces.SendHead{
		Head: interfaces.Head{
			MsgID: uint32(msgID),
		},
	}
}

// CreateMsgHead 创建消息头
func CreateMsgHead(msgID msgid.MsgID) interfaces.SendHead {
	return createMsgHead(msgID)
}
