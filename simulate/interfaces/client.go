package interfaces

import (
	msgid "steve/client_pb/msgId"
	"time"

	"github.com/golang/protobuf/proto"
)

// Head 公用消息头
type Head struct {
	MsgID uint32
}

// SendHead 发包消息头
type SendHead struct {
	Head
}

// SendResult 发包结果
type SendResult struct {
	SendSeq       uint64 // 发送序号
	SendTimestamp int64  // 发送时间戳
}

// MessageExpector 消息期望
type MessageExpector interface {
	Recv(timeOut time.Duration, body proto.Message) error
	Close()
	Clear() // 清空之前收到的消息
}

// Client 客户端接口
type Client interface {

	// 启动
	Start(addr string, version string) error

	// 停止
	Stop() error

	// Closed 是否已经关闭
	Closed() bool

	// SendPackage 发送数据包
	SendPackage(header SendHead, body proto.Message) (*SendResult, error)

	// Request 发送一个请求,阻塞返回响应消息
	// header: 发送的序号
	// body: 发送的消息体
	// timeOut: 超时事件
	// rspMsgID: 期望收到的响应消息 ID
	// rspBody: 收到的消息将反序列化到此接口中
	Request(header SendHead, body proto.Message, timeOut time.Duration, rspMsgID uint32, rspBody proto.Message) error

	// ExpectMessage 创建消息期望
	ExpectMessage(msgID msgid.MsgID) (MessageExpector, error)

	// RemoveMsgExpect 移除消息期望
	RemoveMsgExpect(msgID msgid.MsgID)
}
