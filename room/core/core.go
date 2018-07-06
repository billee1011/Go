package core

import (
	"context"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/config"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/room/loader_balancer"
	"steve/room/peipai"
	"steve/room/registers"
	"steve/server_pb/room_mgr"
	"steve/structs"
	"steve/structs/net"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	_ "steve/room/autoevent" // 引入 autoevent 包，设置工厂
	_ "steve/room/desks"
	_ "steve/room/playermgr"
	_ "steve/room/req_event_translator"
	_ "steve/room/settle"
)

type roomCore struct {
	e   *structs.Exposer
	dog net.WatchDog
}

// NewService 创建服务
func NewService() service.Service {
	return new(roomCore)
}

// RoomService room房间RPC服务
type RoomService struct {
}

func notifyDeskCreate(desk interfaces.Desk) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "notifyDeskCreate",
	})
	ntf := room.RoomDeskCreatedNtf{
		Players: desk.GetPlayers(),
	}
	facade.BroadCastDeskMessage(desk, nil, msgid.MsgID_ROOM_DESK_CREATED_NTF, &ntf, true)
	logEntry.WithField("ntf_context", ntf).Debugln("广播创建房间")
}

// CreateDesk 创建牌桌
func (hws *RoomService) CreateDesk(ctx context.Context, req *roommgr.CreateDeskRequest) (rsp *roommgr.CreateDeskResponse, err error) {

	// 日志
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "RoomService::CreateDesk",
	})

	logEntry.WithField("players", req.GetPlayerId()).Debugln("RoomService::CreateDesk()")

	// 回复match服的消息
	rsp = &roommgr.CreateDeskResponse{
		ErrCode: roommgr.RoomError_FAILED, // 默认是失败的
	}

	// 请求的玩家ID数组
	playersID := req.GetPlayerId()

	// 个数须为4
	if len(playersID) < 4 {
		logEntry.WithField("len(playersID):", len(playersID)).Errorln("players数组长度不为4")
		return
	}

	deskFactory := global.GetDeskFactory()

	deskMgr := global.GetDeskMgr()

	//playerMgr := global.GetPlayerMgr()

	// 创建桌子
	result, err := deskFactory.CreateDesk(playersID, int(req.GetGameId()), interfaces.CreateDeskOptions{})
	if err != nil {
		logEntry.WithFields(
			logrus.Fields{
				"players": playersID,
				"result":  result,
			},
		).WithError(err).Errorln("创建桌子失败")

		return
	}

	logEntry.WithField("players", req.GetPlayerId()).Debugln("创建桌子成功")

	// 回复match：创建桌子成功
	rsp.ErrCode = roommgr.RoomError_SUCCESS

	// 通知该桌子的所有人
	notifyDeskCreate(result.Desk)

	// 桌子开始运行
	deskMgr.RunDesk(result.Desk)

	// 添加进playerManager
	// todo

	return
}

func (c *roomCore) Init(e *structs.Exposer, param ...string) error {
	logrus.Info("room init")
	c.e = e
	global.SetMessageSender(e.Exchanger)
	registers.RegisterHandlers(e.Exchanger)
	registerLbReporter(e)

	rpcServer := e.RPCServer
	err := rpcServer.RegisterService(roommgr.RegisterRoomMgrServer, &RoomService{})
	if err != nil {
		return err
	}

	return nil
}

func (c *roomCore) Start() error {
	go startPeipai()
	return nil
}

func startPeipai() error {
	peipaiAddr := viper.GetString(config.ListenPeipaiAddr)
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "startPeipai",
		"addr":      peipaiAddr,
	})
	if peipaiAddr != "" {
		logEntry.Info("启动配牌服务")
		err := peipai.Run(peipaiAddr)
		if err != nil {
			logEntry.WithError(err).Panic("配牌服务启动失败")
		}
		return err
	}
	logEntry.Info("未配置配牌")
	return nil
}

func registerLbReporter(exposer *structs.Exposer) {
	if err := lb.RegisterLBReporter(exposer.RPCServer); err != nil {
		logrus.WithError(err).Panicln("注册负载上报服务失败")
	}
}
