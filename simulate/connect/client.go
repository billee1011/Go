package connect

import (
	"container/list"
	"reflect"
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
	SendSeq uint64 // 发送序号
}

// RecvHead 收包消息头
type RecvHead struct {
	Head
	RspSeq        uint64 // 服务器的回复序号
	ServerVersion string // 服务器版本号
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

	// RegisterHandle 注册消息处理函数
	// handler 为处理函数，声明类型必须是： func (head RecvHead, body YourProtoType)
	RegisterHandle(msgID uint32, handler interface{}) error

	// Request 发送一个请求,阻塞返回响应消息
	Request(header SendHead, body proto.Message, timeOut time.Duration) (*Response, error)

	// GetResponse 获取服务端发送消息
	GetResponse(msgID uint32, index int) (*Response, error)
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
