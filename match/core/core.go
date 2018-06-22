package core

import (
	"context"
	"errors"
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/server_pb/room_mgr"
	"steve/structs"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
)

type matchCore struct {
	e *structs.Exposer
}

// NewService 创建服务
func NewService() service.Service {
	return new(matchCore)
}

func (c *matchCore) Init(e *structs.Exposer, param ...string) error {
	entry := logrus.WithField("name", "matchCore.Init")

	c.e = e
	if err := c.registerHandles(e.Exchanger); err != nil {
		entry.WithError(err).Error("注册消息处理器失败")
		return err
	}
	return nil
}

func (c *matchCore) Start() error {
	return nil
}

func (c *matchCore) registerHandles(e exchanger.Exchanger) error {
	registe := func(id msgid.MsgID, handler interface{}) {
		err := e.RegisterHandle(uint32(id), handler)
		if err != nil {
			panic(err)
		}
	}

	registe(msgid.MsgID_MATCH_REQ, c.handleMatch)
	return nil
}

func (c *matchCore) handleMatch(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomJoinDeskReq) (ret []exchanger.ResponseMsg) {
	response := &room.RoomJoinDeskRsp{
		ErrCode: room.RoomError_SUCCESS.Enum(),
	}
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_RSP),
		Body:  response,
	}}

	playerID := header.GetPlayerId()

	//TODO 匹配玩家

	//TODO 匹配成功，发起创建房间调用
	err := c.work(playerID)
	if err != nil {
		fmt.Println("call work failed")
		return
	}

	return
}

func (c *matchCore) work(playerID uint64) error {
	cc, err := c.e.RPCClient.GetConnectByServerName("room")
	if err != nil {
		return fmt.Errorf("Get client connection failed:%v", err)
	}
	if cc == nil {
		return errors.New("no service named room. ensure your consul agent is running and configed room")
	}

	client := roommgr.NewRoomMgrClient(cc)
	resp, err := client.CreateDesk(context.Background(), &roommgr.RoomMgrRequest{
		PlayerId: playerID,
	})

	if err != nil {
		return fmt.Errorf("call HelloRoom failed: %v", err)
	}

	fmt.Println("receive response from server:", resp.GetErrCode())
	return nil
}
