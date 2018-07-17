package register

import (
	"steve/client_pb/msgid"
	"steve/match/matchv2"
	"steve/structs/exchanger"
)

// RegisterHandles 注册消息处理
func RegisterHandles(e exchanger.Exchanger) error {
	register := func(id msgid.MsgID, handler interface{}) {
		err := e.RegisterHandle(uint32(id), handler)
		if err != nil {
			panic(err)
		}
	}

	register(msgid.MsgID_MATCH_REQ, matchv2.HandleMatchReq)             // 匹配请求消息
	register(msgid.MsgID_MATCH_CONTINUE_REQ, matchv2.HandleContinueReq) // 续局请求

	return nil
}
