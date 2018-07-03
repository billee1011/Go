package ddzdesk

import (
	"steve/room/desks/deskbase"
	"steve/room/interfaces"
	"steve/structs/proto/gate_rpc"
)

type playerRequest struct {
	playerID uint64
	head     *steve_proto_gaterpc.Header
	body     []byte
}

// desk 斗地主牌桌
type desk struct {
	deskbase.DeskBase
	requestChannel chan playerRequest
}

// Start 启动牌桌逻辑
// finish : 当牌桌逻辑完成时调用
func (d *desk) Start(finish func()) error {
	d.requestChannel = make(chan playerRequest)

	go func() {
		d.processRequests()
		finish()
	}()
	return nil
}

// Stop 停止牌桌
func (d *desk) Stop() error {
	close(d.requestChannel)
	return nil
}

// PushRequest 压入玩家请求
func (d *desk) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	d.requestChannel <- playerRequest{
		playerID: playerID,
		head:     head,
		body:     bodyData,
	}
}

// PushEvent 压入事件
func (d *desk) PushEvent(event interfaces.Event) {
	return
}

// processRequests 处理请求
func (d *desk) processRequests() {

forstart:
	for {
		select {
		case request, ok := <-d.requestChannel:
			{
				if !ok {
					break forstart
				}
				d.processRequest(&request)
			}
		}
	}
}

func (d *desk) processRequest(request *playerRequest) {

}
