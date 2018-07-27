package watchdog

import (
	"steve/structs/proto/base"
)

type clientCallback interface {
	onRecvPkg(header *base.Header, body []byte)
	afterSendPkg(header *base.Header, body []byte, err error)
	onClientClose()
	onError(err error)
}

// type sendingPkg struct {
// 	header *base.Header
// 	body   []byte
// }

// type client struct {
// 	exchanger
// 	callback clientCallback

// 	sendChan  chan sendingPkg
// 	sendSeq   uint64
// 	recvSeq   uint64
// 	stampTime time.Time

// 	finish chan struct{}
// 	closed bool
// 	mutex  sync.Mutex
// }

// func newClient(e exchanger, callback clientCallback) *client {
// 	return &client{
// 		exchanger: e,
// 		callback:  callback,
// 	}
// }

// func (c *client) run() error {
// 	var err error

// 	defer func() {
// 		if x := recover(); x != nil {
// 			err = fmt.Errorf("client occur an error:%v", x)
// 		}
// 	}()

// 	if c.sendChan == nil {
// 		c.sendChan = make(chan sendingPkg, 16)
// 	}
// 	if c.finish == nil {
// 		c.finish = make(chan struct{}, 0)
// 	}

// 	go func() {
// 		defer c.close()
// 		c.doSendLoop()
// 	}()

// 	go func() {
// 		defer c.close()
// 		c.doRecvLoop()
// 	}()

// 	<-c.finish
// 	close(c.sendChan)
// 	return err
// }

// func (c *client) pushMessage(head *base.Header, body []byte) (err error) {
// 	defer func() {
// 		if x := recover(); x != nil {
// 			err = fmt.Errorf("client closed")
// 		}
// 	}()

// 	c.sendChan <- sendingPkg{
// 		header: proto.Clone(head).(*base.Header),
// 		body:   body,
// 	}

// 	return
// }

// func (c *client) callError(err error) {
// 	if c.callback == nil {
// 		return
// 	}
// 	c.callback.onError(err)
// }

// func (c *client) doRecvLoop() {
// 	defer func() {
// 		if x := recover(); x != nil {
// 			err := fmt.Errorf("send loop occur an error:%v", x)
// 			c.callError(err)
// 		}
// 	}()
// 	for {
// 		data, err := c.exchanger.Recv()
// 		if err == io.EOF {
// 			if c.callback != nil {
// 				c.callback.onClientClose()
// 			}
// 			break
// 		}
// 		if err != nil {
// 			c.callError(err)
// 			break
// 		}
// 		if len(data) == 0 {
// 			c.callError(fmt.Errorf("0字节消息包"))
// 			break
// 		}
// 		headsz := uint8(data[0])
// 		if uint32(headsz) > uint32(len(data))-uint32(1) {
// 			c.callError(fmt.Errorf("消息头大小超过消息包数据大小"))
// 			break
// 		}
// 		header := new(base.Header)
// 		err = proto.Unmarshal(data[1:1+headsz], header)
// 		if err != nil {
// 			c.callError(fmt.Errorf("消息头反序列化失败"))
// 			break
// 		}
// 		// todo 序号校验
// 		c.recvSeq = header.GetSendSeq()
// 		c.stampTime = time.Now()

// 		left := uint32(len(data)) - uint32(1) - uint32(headsz)

// 		if header.GetBodyLength() != left {
// 			c.callError(fmt.Errorf("消息体大小错误， 剩余:%v 需要:%v", left, header.GetBodyLength()))
// 			break
// 		}
// 		if c.callback != nil {
// 			c.callback.onRecvPkg(header, data[1+headsz:])
// 		}
// 	}
// }

// func (c *client) send(data sendingPkg) error {
// 	bodySz := uint32(len(data.body))
// 	timeStamp := uint64(c.stampTime.Unix())
// 	sendSeq := c.sendSeq + 1

// 	head := data.header

// 	head.BodyLength = &bodySz
// 	head.RecvSeq = &(c.recvSeq)
// 	head.StampTime = &timeStamp
// 	head.SendSeq = &sendSeq

// 	var headBytes []byte
// 	var err error
// 	if headBytes, err = proto.Marshal(head); err != nil {
// 		return fmt.Errorf("消息头序列化失败")
// 	}
// 	if len(headBytes) > 0xff {
// 		return fmt.Errorf("消息头过长")
// 	}

// 	c.sendSeq++
// 	wholedata := make([]byte, 1+len(headBytes)+int(bodySz))
// 	wholedata[0] = byte(len(headBytes))
// 	copy(wholedata[1:len(headBytes)+1], headBytes)
// 	copy(wholedata[len(headBytes)+1:], data.body)
// 	return c.exchanger.Send(wholedata)
// }

// func (c *client) doSendLoop() {
// 	defer func() {
// 		if x := recover(); x != nil {
// 			err := fmt.Errorf("send loop occur an error:%v", x)
// 			c.callError(err)
// 		}
// 	}()

// forstart:
// 	for {
// 		select {
// 		case context, ok := <-c.sendChan:
// 			if !ok {
// 				break forstart
// 			}
// 			if err := c.send(context); err != nil {
// 				c.callError(err)
// 				break forstart
// 			}
// 		case <-c.finish:
// 			break forstart
// 		}
// 	}
// }

// func (c *client) close() {
// 	c.mutex.Lock()
// 	if c.closed {
// 		c.mutex.Unlock()
// 		return
// 	}
// 	c.closed = true
// 	close(c.finish)
// 	c.mutex.Unlock()
// }
