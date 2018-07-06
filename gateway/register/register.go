package register

import (
	"steve/client_pb/msgId"
	"steve/gateway/auth"
	"steve/structs/exchanger"
)

// RegisteHandlers 注册消息处理器
func RegisteHandlers(e exchanger.Exchanger) {
	registe := func(id msgid.MsgID, handler interface{}) {
		err := e.RegisterHandle(uint32(id), handler)
		if err != nil {
			panic(err)
		}
	}

	registe(msgid.MsgID_GATE_AUTH_REQ, auth.HandleAuthReq)
}
