package connect

import (
	"errors"
	"steve/base/socket"
	"steve/client_pb/msgid"
	"steve/simulate/interfaces"
	"steve/structs/proto/base"
	"sync/atomic"

	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// sendingData 发送的数据
type sendingData struct {
	header  interfaces.SendHead
	body    proto.Message
	sendSeq uint64 // 发送序号
}

// requestInfo 请求信息
type requestInfo struct {
	sendSeq  uint64        // 发送序号
	rspMsgID uint32        // 回复的消息 ID
	rspBody  proto.Message // 回复的消息体
	rspChan  chan error
}

// messageExpector 消息期望
type messageExpector struct {
	ch     chan []byte
	closer func()
}

func (me *messageExpector) Recv(timeOut time.Duration, body proto.Message) error {
	timer := time.NewTimer(timeOut)
	select {
	case <-timer.C:
		{
			return errors.New("等待超时")
		}
	case bodyData, ok := <-me.ch:
		{
			if !ok {
				return errors.New("已经被关闭")
			}
			if body == nil {
				return nil
			}
			if err := proto.Unmarshal(bodyData, body); err != nil {
				return fmt.Errorf("消息反序列化失败： %v", err)
			}
			return nil
		}
	}
}

func (me *messageExpector) Clear() {
	for {
		cont := false

		select {
		case <-me.ch:
			cont = true
		default:
			cont = false
		}

		if !cont {
			break
		}
	}
}

func (me *messageExpector) Close() {
	me.closer()
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

	// requestInfos 请求表
	requestInfos sync.Map

	// expectInfos 期望表
	expectInfos sync.Map
}

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

func (c *client) Closed() bool {
	select {
	case <-c.finish:
		return true
	default:
		return false
	}
}

func (c *client) SendPackage(header interfaces.SendHead, body proto.Message) (result *interfaces.SendResult, err error) {
	sendSeq := c.allocSeq()
	sendTimestamp := time.Now().UnixNano()

	c.sendingChan <- sendingData{
		header:  header,
		body:    body,
		sendSeq: sendSeq,
	}
	return &interfaces.SendResult{
		SendSeq:       sendSeq,
		SendTimestamp: sendTimestamp,
	}, nil
}

func (c *client) Request(header interfaces.SendHead, body proto.Message, timeOut time.Duration, rspMsgID uint32, rspBody proto.Message) error {
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
}

func (c *client) ExpectMessage(msgID msgid.MsgID) (interfaces.MessageExpector, error) {
	me := messageExpector{
		closer: func() {
			c.expectInfos.Delete(msgID)
		},
		ch: make(chan []byte, 256),
	}
	old, loaded := c.expectInfos.LoadOrStore(msgID, me)
	if loaded {
		meOld := old.(messageExpector)
		return &meOld, errors.New("已经存在该消息的期望")
	}
	return &me, nil
}

// RemoveMsgExpect 移除消息期望
func (c *client) RemoveMsgExpect(msgID msgid.MsgID) {
	c.expectInfos.Delete(msgID)
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

func (c *client) checkRequests(header *base.Header, body []byte) {
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
		header := new(base.Header)
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
		logrus.WithField("msg_id", header.GetMsgId()).Debugln("收到消息")
		c.checkRequests(header, data[1+headsz:])
		c.checkExpects(header, data[1+headsz:])
	}
}

func (c *client) checkExpects(header *base.Header, bodyData []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"fun_name": "client.checkExpects",
	})
	msgID := header.GetMsgId()
	iExpector, ok := c.expectInfos.Load(msgid.MsgID(msgID))
	if iExpector == nil || !ok {
		logEntry.WithField("msgID", msgID).Infoln("没有对应的Expector，需要添加")
		return
	}
	me := iExpector.(messageExpector)
	select {
	case me.ch <- bodyData:
	}
}

func (c *client) send(data sendingData) error {
	bodyData, err := proto.Marshal(data.body)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %v", err)
	}
	wireHeader := base.Header{
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

// NewClient 创建客户端接口
func NewClient() (interfaces.Client, error) {
	c := &client{}
	return c, nil
}

// NewTestClient 创建测试客户端
func NewTestClient(target, version string) interfaces.Client {
	c, _ := NewClient()
	if err := c.Start(target, version); err != nil {
		panic(err)
	}
	return c
}
