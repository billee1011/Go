package connect

import (
	"container/list"
	"fmt"
	"io"
	"reflect"
	"steve/base/socket"
	"steve/structs/proto/base"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// sendingData 发送的数据
type sendingData struct {
	header  SendHead
	body    proto.Message
	sendSeq uint64 // 发送序号
}

type interceptor struct {
	msgID        uint32
	rspSeq       uint64
	responseChan chan *Response
}

// handler 消息处理器
type handler struct {
	bodyType reflect.Type
	f        interface{}
}

type client struct {

	// 客户端版本号
	version string

	// 套接字
	sock socket.Socket

	// 待发送队列
	sendingChan chan sendingData

	// 关闭通道
	finish chan struct{}

	// 是否正在关闭中
	closing bool

	// closing 的读写锁
	closingMutex sync.RWMutex

	// 最近收到的消息序号
	lastRecvSeq uint64

	// 最大发送的消息序号
	maxSeq uint64

	// maxSeq 的读写锁
	seqMutex sync.RWMutex

	// 处理函数
	handlerMap sync.Map

	// 服务端返回消息缓存, 按接收顺序存储
	reponseList *list.List

	// 响应消息拦截器
	interceptors sync.Map
}

var _ Client = new(client)

func (c *client) Start(addr string, version string) error {
	c.version = version
	var err error
	c.sock, err = connectServer(addr)
	if err != nil {
		return fmt.Errorf("连接服务器失败：%v", err)
	}
	c.sendingChan = make(chan sendingData, 5)
	c.finish = make(chan struct{}, 0)

	go func() {
		c.sendLoop()
		c.Stop()
	}()

	go func() {
		c.recvLoop()
		c.Stop()
	}()

	go func() {
		<-c.finish
		close(c.sendingChan)
	}()
	return nil
}

func (c *client) Stop() error {
	c.closingMutex.Lock()
	defer c.closingMutex.Unlock()
	if c.closing {
		return nil
	}
	c.closing = true
	close(c.finish)
	return c.sock.Close()
}

func (c *client) SendPackage(header SendHead, body proto.Message) (result *SendResult, err error) {
	sendSeq := c.allocSeq()

	c.sendingChan <- sendingData{
		header:  header,
		body:    body,
		sendSeq: sendSeq,
	}
	return &SendResult{
		SendSeq: sendSeq,
	}, nil
}

func (c *client) Request(header SendHead, body proto.Message, timeOut time.Duration) (*Response, error) {
	entry := logrus.WithField("name", "client.Request")
	sendSeq := c.allocSeq()
	ch := make(chan *Response)
	it := &interceptor{
		msgID:        header.MsgID,
		rspSeq:       sendSeq,
		responseChan: ch,
	}

	// 增加一个消息拦截器
	c.interceptors.Store(header.MsgID, it)

	c.sendingChan <- sendingData{
		header:  header,
		body:    body,
		sendSeq: sendSeq,
	}

	timer := time.NewTimer(timeOut)
	select {
	case response := <-ch:
		return response, nil
	case <-timer.C:
		entry.Error("请求超时")
		return nil, fmt.Errorf("请求超时")
	}
}
func (c *client) GetResponse(msgID uint32, index int) (*Response, error) {
	return nil, nil
}

// allocSeq 分配发送序号
func (c *client) allocSeq() uint64 {
	defer c.seqMutex.Unlock()
	c.seqMutex.Lock()

	c.maxSeq++
	return c.maxSeq
}

func (c *client) RegisterHandle(msgID uint32, handlerFunc interface{}) error {
	ft := reflect.TypeOf(handlerFunc)
	if ft.Kind() != reflect.Func {
		return fmt.Errorf("handler 需要是一个函数")
	}
	if ft.NumIn() != 2 {
		return fmt.Errorf("handler 需要接收 2 个参数")
	}
	if ft.In(0) != reflect.TypeOf(RecvHead{}) {
		return fmt.Errorf("handler 的第一个参数需要是 RecvHead 类型")
	}
	bodyType := ft.In(1)
	if _, ok := reflect.New(bodyType).Interface().(proto.Message); !ok {
		return fmt.Errorf("handler 的第 2 个参数需要可以转换为 proto.Message")
	}
	c.handlerMap.Store(msgID, handler{
		f:        handlerFunc,
		bodyType: bodyType,
	})
	return nil
}

func (c *client) sendLoop() {
	entry := logrus.WithField("name", "client.sendLoop")

	defer func() {
		if x := recover(); x != nil {
			entry.WithField("error", x).Error("发包循环检测到异常")
		}
	}()
forstart:
	for {
		select {
		case data, ok := <-c.sendingChan:
			if !ok {
				break forstart
			}
			if err := c.send(data); err != nil {
				entry.WithField("msg_id", data.header.MsgID).WithError(err).Errorln("发送数据失败")
				break forstart
			}
		case <-c.finish:
			break forstart
		}
	}
	entry.Infoln("发送循环完成")
}

func (c *client) recvLoop() {
	entry := logrus.WithField("name", "client.recvLoop")
	defer func() {
		if x := recover(); x != nil {
			entry.WithField("error", x).Error("收包循环检测到异常")
		}
	}()

	for {
		data, err := c.sock.RecvPackage()
		if err != nil {
			if err == io.EOF {
				return
			}
			entry.WithError(err).Errorln("收包失败")
			return
		}
		if len(data) == 0 {
			entry.Errorln("收到的包长度为0")
			break
		}
		headsz := uint8(data[0])
		if uint32(headsz) > uint32(len(data))-uint32(1) {
			entry.Errorln("消息头大小超过消息包数据大小")
			break
		}
		header := new(steve_proto_base.Header)
		err = proto.Unmarshal(data[1:1+headsz], header)
		if err != nil {
			entry.WithError(err).Errorln("消息头反序列化失败")
			break
		}
		c.lastRecvSeq = header.GetSendSeq()

		left := uint32(len(data)) - uint32(1) - uint32(headsz)

		if header.GetBodyLength() != left {
			entry.WithFields(logrus.Fields{
				"left_size":   left,
				"body_length": header.GetBodyLength(),
			}).Errorln("消息体大小错误")
			break
		}

		c.intercept(header, data[1+headsz:])
		c.dispatch(header, data[1+headsz:])
	}
}

func (c *client) dispatch(header *steve_proto_base.Header, body []byte) {
	entry := logrus.WithFields(logrus.Fields{
		"name":   "client.dispatch",
		"msg_id": header.GetMsgId(),
	})
	h, ok := c.handlerMap.Load(header.GetMsgId())
	if !ok {
		entry.Warn("未处理的消息")
		return
	}
	hh := h.(handler)
	paramHeader := RecvHead{
		Head: Head{
			MsgID: header.GetMsgId(),
		},
		RspSeq:        header.GetRspSeq(),
		ServerVersion: header.GetVersion(),
	}

	msg := reflect.New(hh.bodyType).Interface()
	if err := proto.Unmarshal(body, msg.(proto.Message)); err != nil {
		entry.Error("消息体反序列化失败")
		return
	}

	f := reflect.ValueOf(hh.f)
	f.Call([]reflect.Value{
		reflect.ValueOf(paramHeader),
		reflect.ValueOf(msg).Elem(),
	})
}

func (c *client) send(data sendingData) error {
	bodyData, err := proto.Marshal(data.body)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %v", err)
	}
	wireHeader := steve_proto_base.Header{
		MsgId:      proto.Uint32(data.header.MsgID),
		SendSeq:    proto.Uint64(data.sendSeq),
		RecvSeq:    proto.Uint64(c.lastRecvSeq),
		StampTime:  proto.Uint64(uint64(time.Now().Unix())),
		BodyLength: proto.Uint32(uint32(len(bodyData))),
		RspSeq:     proto.Uint64(0),
		Version:    proto.String(c.version),
	}

	headData, err := proto.Marshal(&wireHeader)
	if err != nil {
		return fmt.Errorf("消息头反序列化失败: %v", err)
	}
	if len(headData) > 0xff {
		return fmt.Errorf("包头长度超过 1 字节")
	}
	totalSize := 2 + 1 + len(headData) + len(bodyData)
	if totalSize > 0xffff {
		return fmt.Errorf("总包长超过 2 字节")
	}

	fmt.Println(totalSize)

	wholeData := make([]byte, totalSize)

	wholeData[0] = byte((totalSize & 0xff00) >> 8)
	wholeData[1] = byte(totalSize & 0xff)
	wholeData[2] = byte(len(headData))
	copy(wholeData[3:3+len(headData)], headData)
	copy(wholeData[3+len(headData):], bodyData)

	err = c.sock.SendPackage(wholeData)
	if err != nil {
		return fmt.Errorf("发送消息失败： %v", err)
	}
	return nil
}

func (c *client) intercept(header *steve_proto_base.Header, body []byte) {
	entry := logrus.WithField("name", "client.intercept")

	// 检测是否存在该消息的拦截设置
	if meta, ok := metaByID[*header.MsgId]; ok {
		msg := reflect.New(meta.Type).Interface()
		if err := proto.Unmarshal(body, msg.(proto.Message)); err != nil {
			entry.Error(err)
		}
		recvHead := RecvHead{
			Head: Head{
				MsgID: header.GetMsgId(),
			},
			RspSeq:        header.GetRspSeq(),
			ServerVersion: header.GetVersion(),
		}

		response := &Response{
			Head: recvHead,
			Body: msg.(proto.Message),
		}

		if v, ok := c.interceptors.Load(header.GetMsgId()); ok {
			if v.(*interceptor).rspSeq == header.GetSendSeq() {
				v.(*interceptor).responseChan <- response
				c.interceptors.Delete(header.GetMsgId())
			}
		}
		c.reponseList.PushBack(response)
	} else {
		entry.Warningln("响应消息未注册")
	}

}
