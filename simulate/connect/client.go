package connect

import (
	"container/list"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
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

// RecvHead 收包消息头
type RecvHead struct {
	Head
	RspSeq        uint64 // 服务器的回复序号
	ServerVersion string // 服务器版本号
	RecvTimestamp int64  // 接收时间戳
}

// Response 服务端推送消息
type Response struct {
	Head RecvHead
	Body proto.Message
}

// MessageMeta 消息元信息
type MessageMeta struct {
	Type reflect.Type
	ID   uint32
}

// Client 客户端接口
type Client interface {

	// 启动
	Start(addr string, version string) error

	// 停止
	Stop() error

	// SendPackage 发送数据包
	SendPackage(header SendHead, body proto.Message) (*SendResult, error)

	// Request 发送一个请求,阻塞返回响应消息
	// header: 发送的序号
	// body: 发送的消息体
	// timeOut: 超时事件
	// rspMsgID: 期望收到的响应消息 ID
	// rspBody: 收到的消息将反序列化到此接口中
	Request(header SendHead, body proto.Message, timeOut time.Duration, rspMsgID uint32, rspBody proto.Message) error

	// WaitMessage 等服务端的某条消息
	// timestamp 客户端发送某条消息的时间戳,纳秒
	WaitMessage(ctx context.Context, msgID uint32, timestamp int64) (*Response, error)
}

// NewClient 创建客户端接口
func NewClient() (Client, error) {
	c := &client{
		reponseList: list.New(),
	}
	return c, nil
}

// NewTestClient 创建测试客户端
func NewTestClient(target, version string) Client {
	c, _ := NewClient()
	if err := c.Start(target, version); err != nil {
		panic(err)
	}
	return c
}
