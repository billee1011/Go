package watchdog

import (
	"fmt"
	"runtime/debug"
	"steve/structs/proto/base"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type recvdata struct {
	header *steve_proto_base.Header
	data   []byte
}

type senddata struct {
	header *steve_proto_base.Header
	body   []byte
}

type clientV2 struct {

	// 交互接口
	exchanger

	// 回调接口
	callback clientCallback

	// 关闭通道，外部通过该通道关闭客户端
	finish chan struct{}

	// 接收数据状态读写锁
	recvStateMutex sync.RWMutex
	// 接收数据状态， 最近接收的消息序号
	recvSeq uint64
	// 接收数据状态， 最近接收到消息的时间
	recvTimestamp time.Time

	// 数据发送通道
	csend chan *senddata

	// 发送序号
	sendSeq uint64
}

func newClientV2(e exchanger, callback clientCallback) *clientV2 {
	return &clientV2{
		exchanger: e,
		callback:  callback,
	}
}

func (c *clientV2) pushMessage(head *steve_proto_base.Header, body []byte) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("client closed")
			debug.PrintStack()
		}
	}()

	if body == nil || head == nil {
		panic("发送数据为空")
	}

	// 为了保证消息有序，不能使用 goroutine， 而要使用带缓冲区的通道
	c.csend <- &senddata{
		header: proto.Clone(head).(*steve_proto_base.Header),
		body:   body,
	}
	return
}

// 关闭
func (c *clientV2) close() {
	close(c.finish)
}

// run 返回后表示 client 的生命周期结束
func (c *clientV2) run(onFinish func()) error {
	defer func() {
		if x := recover(); x != nil {
			c.callback.onError(fmt.Errorf("客户端执行错误: %v. 堆栈： %s", x, string(debug.Stack())))
		}
	}()

	// done 通道， 在 run 退出后关闭， 其他 goroutine 在 done 通道关闭后停止工作
	done := make(chan struct{})
	defer close(done)

	c.csend = make(chan *senddata, 5)
	defer close(c.csend)

	c.finish = make(chan struct{})
	defer close(c.finish)
	// 接收数据 goroutine
	rfinish := c.recvLoop(done)

	// 发送数据 goroutine
	sfinish := c.sendLoop(done)

	if onFinish != nil {
		defer onFinish()
	}

	// 在接收数据、发送数据 goroutine 或者是 c.finish (外部关闭)关闭后，返回
	select {
	case <-rfinish:
		return nil
	case <-sfinish:
		return nil
	case <-c.finish:
		return nil
	}
}

// recvLoop 启动接收数据 goroutine
func (c *clientV2) recvLoop(done <-chan struct{}) (finish chan struct{}) {
	finish = make(chan struct{})
	go func() {
		defer func() {
			if x := recover(); x != nil {
				c.callback.onError(fmt.Errorf("接收数据错误: %v. 堆栈： %s", x, string(debug.Stack())))
			}
		}()

		defer close(finish)
		for {
			select {
			case <-done:
				return
			default:
				data := c.recv()
				if data == nil {
					return
				}
				c.handleRecv(data)
			}
		}
	}()
	return finish
}

func (c *clientV2) sendLoop(done <-chan struct{}) (finish chan struct{}) {
	finish = make(chan struct{})

	go func() {
		defer func() {
			if x := recover(); x != nil {
				c.callback.onError(fmt.Errorf("发送据错误: %v. 堆栈： %s", x, string(debug.Stack())))
			}
		}()
		defer close(finish)
		for {
			select {
			case data := <-c.csend:
				if data == nil {
					return
				}
				if err := c.send(data); err != nil {
					c.callback.onError(fmt.Errorf("数据发送失败: %v", err))
					return
				}
			case <-done:
				return
			}
		}
	}()

	return finish
}

func (c *clientV2) send(data *senddata) error {
	bodySz := uint32(len(data.body))

	recvSeq, recvTimeStamp := c.getRecvState()

	timeStamp := uint64(recvTimeStamp.Unix())
	sendSeq := c.sendSeq + 1

	head := data.header

	head.BodyLength = &bodySz
	head.RecvSeq = &recvSeq
	head.StampTime = &timeStamp
	head.SendSeq = &sendSeq

	var headBytes []byte
	var err error
	if headBytes, err = proto.Marshal(head); err != nil {
		return fmt.Errorf("消息头序列化失败")
	}
	if len(headBytes) > 0xff {
		return fmt.Errorf("消息头过长")
	}

	c.sendSeq++
	wholedata := make([]byte, 1+len(headBytes)+int(bodySz))
	wholedata[0] = byte(len(headBytes))
	copy(wholedata[1:len(headBytes)+1], headBytes)
	copy(wholedata[len(headBytes)+1:], data.body)
	return c.exchanger.Send(wholedata)
}

// recv 从 exchanger 中接收数据。 接收过程中不会中断
func (c *clientV2) recv() *recvdata {
	d, err := c.exchanger.Recv()
	if err != nil {
		c.callback.onError(fmt.Errorf("接受数据失败： %v", err))
		return nil
	}
	if len(d) == 0 {
		c.callback.onError(fmt.Errorf("0字节消息包"))
		return nil
	}
	headsz := uint8(d[0])
	if uint32(headsz) > uint32(len(d))-uint32(1) {
		c.callback.onError(fmt.Errorf("消息头大小超过消息包数据大小"))
		return nil
	}
	header := new(steve_proto_base.Header)
	err = proto.Unmarshal(d[1:1+headsz], header)
	if err != nil {
		c.callback.onError(fmt.Errorf("消息头反序列化失败"))
		return nil
	}

	left := uint32(len(d)) - uint32(1) - uint32(headsz)

	if header.GetBodyLength() != left {
		c.callback.onError(fmt.Errorf("消息体大小错误， 剩余:%v 需要:%v", left, header.GetBodyLength()))
		return nil
	}
	return &recvdata{
		header: header,
		data:   d[1+headsz:],
	}
}

// handleRecv 处理返回消息
func (c *clientV2) handleRecv(data *recvdata) {
	c.recvStateMutex.Lock()
	c.recvSeq = data.header.GetSendSeq()
	c.recvTimestamp = time.Now()
	c.recvStateMutex.Unlock()

	c.callback.onRecvPkg(data.header, data.data)
}

func (c *clientV2) getRecvState() (uint64, time.Time) {
	c.recvStateMutex.RLock()
	defer c.recvStateMutex.RUnlock()
	return c.recvSeq, c.recvTimestamp
}
