package core

import (
	"context"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/peipai"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/room/loader_balancer"
	"steve/room/registers"
	"steve/server_pb/room"
	"steve/structs"
	"steve/structs/net"
	"steve/structs/proto/gate_rpc"
	"steve/structs/service"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
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

type RoomService struct {
}

type joinApplyManager struct {
	applyChannel chan uint64
}

var gJoinApplyMgr *joinApplyManager
var once sync.Once

func getJoinApplyMgr() *joinApplyManager {
	once.Do(initApplyMgr)
	return gJoinApplyMgr
}

func initApplyMgr() {
	gJoinApplyMgr = newApplyMgr(true)
}

func newApplyMgr(runChecker bool) *joinApplyManager {
	mgr := &joinApplyManager{
		applyChannel: make(chan uint64, 1024),
	}
	if runChecker {
		go mgr.checkMatch()
	}
	return mgr
}

func (jam *joinApplyManager) getApplyChannel() chan uint64 {
	return jam.applyChannel
}

func (jam *joinApplyManager) joinPlayer(playerID uint64) room.RoomError {
	// TODO: 检测玩家状态
	ch := jam.getApplyChannel()
	ch <- playerID
	return room.RoomError_SUCCESS
}

func (jam *joinApplyManager) replicateApplyProc(applyPlayers []uint64, newPlayerID uint64) bool {
	for _, playerID := range applyPlayers {
		if playerID == newPlayerID {
			header := &steve_proto_gaterpc.Header{
				MsgId: uint32(msgid.MsgID_ROOM_JOIN_DESK_RSP),
			}
			rsp := &room.RoomJoinDeskRsp{
				ErrCode: room.RoomError_DESK_ALREADY_APPLIED.Enum(),
			}
			SendMessageByPlayerID(playerID, header, rsp)
			return true
		}
	}
	return false
}

// SendMessageByPlayerID 获取到playerID发送消息
func SendMessageByPlayerID(playerID uint64, head *steve_proto_gaterpc.Header, body proto.Message) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":   "sendMessageFromRoom",
		"newPlayerID": playerID,
		"head":        msgid.MsgID_name[int32(head.MsgId)],
	})
	playerMgr := global.GetPlayerMgr()
	p := playerMgr.GetPlayer(playerID)
	if p != nil {
		logEntry.Errorln("获取player失败")
		return
	}
	clientID := p.GetClientID()
	ms := global.GetMessageSender()
	err := ms.SendPackage(clientID, head, body)
	if err != nil {
		logEntry.WithError(err).Errorln("发送消息失败")
	}
}

func (jam *joinApplyManager) removeOfflinePlayer(playerIDs []uint64) []uint64 {
	result := make([]uint64, 0, len(playerIDs))
	playerMgr := global.GetPlayerMgr()
	for _, playerID := range playerIDs {
		player := playerMgr.GetPlayer(playerID)
		if player != nil && player.GetClientID() != 0 {
			result = append(result, playerID)
		} else {
			logrus.WithField("player_id", playerID).Debugln("玩家不在线，移除")
		}
	}
	return result
}

func (jam *joinApplyManager) checkMatch() {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "checkMatch",
	})
	deskFactory := global.GetDeskFactory()
	deskMgr := global.GetDeskMgr()
	applyPlayers := make([]uint64, 0, 4)

	ch := jam.getApplyChannel()

	for {
		playerID, ok := <-ch
		logEntry.WithField("player_id", playerID).Debugln("accept player")
		if !ok {
			break
		}

		if jam.replicateApplyProc(applyPlayers, playerID) {
			continue
		}
		applyPlayers = append(applyPlayers, playerID)
		applyPlayers = jam.removeOfflinePlayer(applyPlayers)

		for len(applyPlayers) >= 4 {
			players := applyPlayers[:4]
			applyPlayers = applyPlayers[4:]
			result, err := deskFactory.CreateDesk(players, 1, interfaces.CreateDeskOptions{})
			if err != nil {
				logEntry.WithFields(
					logrus.Fields{
						"players": players,
						"result":  result,
					},
				).WithError(err).Errorln("创建房间失败")
				continue
			}
			notifyDeskCreate(result.Desk)
			deskMgr.RunDesk(result.Desk)
		}
	}
}

func notifyDeskCreate(desk interfaces.Desk) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "notifyDeskCreate",
	})
	players := desk.GetPlayers()
	clientIDs := []uint64{}

	playerMgr := global.GetPlayerMgr()
	for _, player := range players {
		playerID := player.GetPlayerId()
		p := playerMgr.GetPlayer(playerID)
		if p != nil {
			clientIDs = append(clientIDs, p.GetClientID())
		}
	}
	ntf := room.RoomDeskCreatedNtf{
		Players: desk.GetPlayers(),
	}
	head := &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid.MsgID_ROOM_DESK_CREATED_NTF)}
	ms := global.GetMessageSender()

	ms.BroadcastPackage(clientIDs, head, &ntf)
	logEntry.WithField("ntf_context", ntf).Debugln("广播创建房间")
}

func (hws *RoomService) CreateDesk(ctx context.Context, req *matchroom.MatchRoomRequest) (rsp *matchroom.MatchRoomResponse, err error) {
	//TODO 接到创建房间请求，创建房间
	playerID := req.GetPlayerId()
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayer(playerID)

	rsp = &matchroom.MatchRoomResponse{
		ErrCode: matchroom.RoomError_SUCCESS,
	}

	if player == nil {
		rsp.ErrCode = matchroom.RoomError_FAILED
		return
	}

	getJoinApplyMgr().joinPlayer(playerID)
	return

	//TODO 创建房间成功，广播房间消息到gateway
}

func (c *roomCore) Init(e *structs.Exposer, param ...string) error {
	logrus.Info("room init")
	c.e = e
	global.SetMessageSender(e.Exchanger)
	registers.RegisterHandlers(e.Exchanger)
	registerLbReporter(e)

	rpcServer := e.RPCServer
	err := rpcServer.RegisterService(matchroom.RegisterMatchRoomServer, &RoomService{})
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
	peipaiAddr := viper.GetString(ListenPeipaiAddr)
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
