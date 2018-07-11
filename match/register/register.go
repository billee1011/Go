package register

import (
	"steve/structs/exchanger"
	"steve/client_pb/msgid"
	"steve/match/match"
)

// RegisterHandles 注册消息处理
func RegisterHandles(e exchanger.Exchanger) error {
	register := func(id msgid.MsgID, handler interface{}) {
		err := e.RegisterHandle(uint32(id), handler)
		if err != nil {
			panic(err)
		}
	}

	register(msgid.MsgID_MATCH_REQ, match.HandleMatchReq) // 匹配请求消息

	return nil
}
