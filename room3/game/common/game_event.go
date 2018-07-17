package common

import (
	mahjong "steve/server_pb/majong"
)

type GameEventResult interface {
	ReplyMsgs() []interface{}
	IsSuccess() bool
	NewContext() interface{}
}

type MajongEventResult struct {
	newContext mahjong.MajongContext        // 处理后的牌局现场
	autoEvent  *mahjong.AutoEvent           // 自动事件
	replyMsgs  []mahjong.ReplyClientMessage // 回复给客户端的消息
	success    bool                         // 是否成功
}

func (mer *MajongEventResult) ReplyMsgs() []interface{} {
	msgs := make([]interface{}, len(mer.replyMsgs))
	for i, replyMsg := range mer.replyMsgs {
		msgs[i] = replyMsg
	}
	return msgs
}

func (mer *MajongEventResult) IsSuccess() bool {
	return mer.success
}

func (mer *MajongEventResult) NewContext() interface{} {
	return mer.newContext
}
