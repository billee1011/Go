package connect

import (
	"errors"
	"steve/base/socket"
	"sync/atomic"

	"container/list"
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// sendingData 发送的数据
type sendingData struct {
	header  SendHead
	body    proto.Message
	sendSeq uint64 // 发送序号
}

type interceptor struct {
	msgID         uint32
	rspSeq        uint64
	sendTimestamp int64
	responseChan  chan *Response
}

// 检测是否为拦截的消息
func (it *interceptor) check(resp *Response) bool {
	if it.msgID == resp.Head.MsgID {
		if it.rspSeq != 0 {
			// 需要验证响应序号
			return it.rspSeq == resp.Head.RspSeq
		}
		// 不需要验证响应序号,验证发送时时间
		return it.sendTimestamp < resp.Head.RecvTimestamp
	}
	return false
}

// requestInfo 请求信息
type requestInfo struct {
	sendSeq  uint64        // 发送序号
	rspMsgID uint32        // 回复的消息 ID
	rspBody  proto.Message // 回复的消息体
	rspChan  chan error
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

	// 处理函数
	handlerMap sync.Map

	// 服务端返回消息缓存, 按接收顺序存储
	reponseList *list.List

	// 响应消息拦截器
	interceptors sync.Map

	// 消息拦截锁
	interceptorMutex sync.Mutex

	// requestInfos 请求表
	requestInfos sync.Map
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
	sendTimestamp := time.Now().UnixNano()

	c.sendingChan <- sendingData{
		header:  header,
		body:    body,
		sendSeq: sendSeq,
	}
	return &SendResult{
		SendSeq:       sendSeq,
		SendTimestamp: sendTimestamp,
	}, nil
}

func (c *client) Request(header SendHead, body proto.Message, timeOut time.Duration, rspMsgID uint32, rspBody proto.Message) error {
	sendResult, err := c.SendPackage(header, body)
	if err != nil {
		return err
	}
	reqInfo := requestInfo{
		sendSeq:  sendResult.SendSeq,
		rspMsgID: rspMsgID,
		rspBody:  rspBody,
		rspChan:  make(chan error),
	}
	c.requestInfos.Store(sendResult.SendSeq, &reqInfo)
	defer c.requestInfos.Delete(sendResult.SendSeq)

	timer := time.NewTimer(timeOut)

	select {
	case err := <-reqInfo.rspChan:
		{
			return err
		}
	case <-timer.C:
		{
			return errors.New("请求超时")
		}
	}
	return nil
}

func (c *client) WaitMessage(ctx context.Context, msgID uint32, timestamp int64) (*Response, error) {
	ch := make(chan *Response)
	it := &interceptor{
		msgID:         msgID,
		sendTimestamp: timestamp,
		responseChan:  ch,
	}

	go func() {
		c.interceptorMutex.Lock()
		defer c.interceptorMutex.Unlock()

		// 先消息缓存里找
		var find bool
		e := c.reponseList.Back()
		for e != nil && e.Value != nil {
			// 最新接收的消息时间戳小鱼期望时间,提前退出
			if e.Value.(*Response).Head.RecvTimestamp < it.sendTimestamp {
				break
			}
			if it.check(e.Value.(*Response)) {
				find = true
				ch <- e.Value.(*Response)
				break
			}
			e = e.Prev()
		}

		// 没找到注册消息拦截器
		if !find {
			c.interceptors.Store(it.msgID, it)
		}

	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case response := <-ch:
		return response, nil
	}
}

// allocSeq 分配发送序号
func (c *client) allocSeq() uint64 {
	return atomic.AddUint64(&c.maxSeq, 1)
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

func (c *client) checkRequests(header *steve_proto_base.Header, body []byte) {
	rspSeq := header.GetRspSeq()
	d, ok := c.requestInfos.Load(rspSeq)
	if !ok {
		return
	}
	reqInfo := d.(*requestInfo)
	if reqInfo.rspMsgID != header.GetMsgId() {
		return
	}
	if err := proto.Unmarshal(body, reqInfo.rspBody); err != nil {
		reqInfo.rspChan <- fmt.Errorf("消息反序列化失败： %v", err)
	} else {
		reqInfo.rspChan <- nil
	}
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
		c.checkRequests(header, data[1+headsz:])
		c.intercept(header, data[1+headsz:])
	}
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
	c.interceptorMutex.Lock()
	defer c.interceptorMutex.Unlock()
	entry := logrus.WithField("name", "client.intercept")
	// 是否注册了该消息
	if meta, ok := metaByID[*header.MsgId]; ok {
		msg := reflect.New(meta.Type).Interface()
		if err := proto.Unmarshal(body, msg.(proto.Message)); err != nil {
			entry.Error(err)
		}
		recvHead := RecvHead{
			Head: Head{
				MsgID: header.GetMsgId(),
			},
			RspSeq:        header.GetSendSeq(),
			ServerVersion: header.GetVersion(),
			RecvTimestamp: time.Now().UnixNano(),
		}

		response := &Response{
			Head: recvHead,
			Body: msg.(proto.Message),
		}

		// 检测是否存在该消息的拦截设置
		if v, ok := c.interceptors.Load(header.GetMsgId()); ok {
			if v.(*interceptor).check(response) { // 匹配拦截
				v.(*interceptor).responseChan <- response
				c.interceptors.Delete(header.GetMsgId())
			}
		}
		c.reponseList.PushBack(response)
	} else {
		entry.Warningln("响应消息未注册")
	}

}
