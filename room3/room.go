package main

import (
	"github.com/Sirupsen/logrus"
	"steve/client_pb/msgid"
	"steve/room3/game"
	"steve/server_pb/room_mgr"
	"steve/structs"
	"steve/structs/exchanger"
	"steve/structs/net"
	"steve/structs/service"
)

type room struct {
	e   *structs.Exposer
	dog net.WatchDog
}

func (r *room) registerHandlers(e exchanger.Exchanger) {
	register := func(id msgid.MsgID, handler interface{}) {
		err := e.RegisterHandle(uint32(id), handler)
		if err != nil {
			panic(err)
		}
	}

	register(msgid.MsgID_ROOM_DESK_QUIT_REQ, game.DefaultDeskManager.HandleExitRequest)               // 退出牌桌请求
	register(msgid.MsgID_ROOM_CANCEL_TUOGUAN_REQ, game.DefaultDeskManager.HandleCancelTuoGuanRequest) // 取消托管请求

	// mahjong
	register(msgid.MsgID_ROOM_HUANSANZHANG_REQ, game.DefaultDeskManager.HandleDeskRequest)
	register(msgid.MsgID_ROOM_XINGPAI_ACTION_REQ, game.DefaultDeskManager.HandleDeskRequest)
	register(msgid.MsgID_ROOM_DINGQUE_REQ, game.DefaultDeskManager.HandleDeskRequest)
	register(msgid.MsgID_ROOM_CHUPAI_REQ, game.DefaultDeskManager.HandleDeskRequest)
	register(msgid.MsgID_ROOM_CARTOON_FINISH_REQ, game.DefaultDeskManager.HandleDeskRequest)
}

func (r *room) Init(e *structs.Exposer, param ...string) error {
	logrus.Info("room init")
	r.e = e

	game.DefaultSender.SetSender(e.Exchanger)

	r.registerHandlers(e.Exchanger)

	rpcServer := e.RPCServer
	err := rpcServer.RegisterService(roommgr.RegisterRoomMgrServer, &game.GameService{})
	if err != nil {
		return err
	}

	return nil
}

func (r *room) Start() error {
	return nil
}

func GetService() service.Service {
	return new(room)
}

func main() {}
