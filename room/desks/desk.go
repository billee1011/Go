package desks

import (
	"context"
	"errors"
	"runtime/debug"
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

var errInitMajongContext = errors.New("初始化麻将现场失败")
var errAllocDeskIDFailed = errors.New("分配牌桌 ID 失败")
var errPlayerNotExist = errors.New("玩家不存在")

type deskPlayer struct {
	playerID uint64
	seat     uint32 // 座号
}

// deskEvent 房间事件
type deskEvent struct {
	eventID      server_pb.EventID
	eventContext []byte
}

type desk struct {
	deskUID   uint64
	gameID    int
	mjContext server_pb.MajongContext
	settler   interfaces.DeskSettler   // 结算器
	players   map[uint32]deskPlayer    // Seat -> player
	event     chan deskEvent           // 牌桌事件通道
	autoEvent chan server_pb.AutoEvent // 自动事件通道
	cancel    context.CancelFunc       // 取消事件处理
}

func makeDeskPlayers(logEntry *logrus.Entry, players []uint64) (map[uint32]deskPlayer, error) {
	playerMgr := global.GetPlayerMgr()
	deskPlayers := make(map[uint32]deskPlayer, 4)
	seat := uint32(0)
	for _, playerID := range players {
		player := playerMgr.GetPlayer(playerID)
		if player == nil {
			logEntry.WithField("player_id", playerID).Errorln(errPlayerNotExist)
			return nil, errPlayerNotExist
		}
		deskPlayers[seat] = deskPlayer{
			playerID: playerID,
			seat:     seat,
		}
		seat++
	}
	return deskPlayers, nil
}

func newDesk(players []uint64, gameID int, opt interfaces.CreateDeskOptions) (result interfaces.CreateDeskResult, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "newDesk",
		"game_id":   gameID,
		"players":   players,
	})
	alloc := global.GetDeskIDAllocator()
	id, err := alloc.AllocDeskID()
	if err != nil {
		logEntry.Errorln(errAllocDeskIDFailed)
		err = errAllocDeskIDFailed
		return
	}
	logEntry = logEntry.WithField("desk_uid", id)
	deskPlayers, err := makeDeskPlayers(logEntry, players)
	if err != nil {
		return
	}
	// TODO: mjOption
	mjContext := server_pb.MajongContext{}

	return interfaces.CreateDeskResult{
		Desk: &desk{
			deskUID:   id,
			gameID:    gameID,
			mjContext: mjContext,
			settler:   global.GetDeskSettleFactory().CreateDeskSettler(gameID),
			players:   deskPlayers,
			event:     make(chan deskEvent, 16),
			autoEvent: make(chan server_pb.AutoEvent, 1),
		},
	}, nil
}

// GetUID 获取牌桌 UID
func (d *desk) GetUID() uint64 {
	return d.deskUID
}

// GetGameID 获取游戏 ID
func (d *desk) GetGameID() int {
	return d.gameID
}

// GetPlayers 获取牌桌玩家数据
func (d *desk) GetPlayers() []*room.RoomPlayerInfo {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.GetPlayers",
		"desk_uid":  d.deskUID,
		"game_id":   d.gameID,
	})
	result := []*room.RoomPlayerInfo{}

	playerMgr := global.GetPlayerMgr()

	for seat, deskPlayer := range d.players {
		player := playerMgr.GetPlayer(deskPlayer.playerID)
		if player == nil {
			logEntry.WithField("player_id", deskPlayer.playerID).Errorln("玩家不存在")
			return nil
		}
		result = append(result, &room.RoomPlayerInfo{
			PlayerId: proto.Uint64(deskPlayer.playerID),
			Coin:     proto.Uint64(player.GetCoin()),
			Seat:     proto.Uint32(uint32(seat)),
		})
	}
	return result
}

// Start 启动牌桌逻辑
// finish : 当牌桌逻辑完成时调用
// step 1. 初始化牌桌现场
// step 2. 启动发送事件的 goroutine
// step 3. 写入开始游戏事件
func (d *desk) Start(finish func()) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.Start",
		"desk_uid":  d.deskUID,
		"game_id":   d.gameID,
	})

	if err := d.initMajongContext(); err != nil {
		logEntry.WithError(err).Errorln(errInitMajongContext)
		return errInitMajongContext
	}
	var ctx context.Context
	ctx, d.cancel = context.WithCancel(context.Background())

	go func() {
		d.processEvents(ctx)
		finish()
	}()

	d.event <- deskEvent{
		eventID:      server_pb.EventID_event_start_game,
		eventContext: []byte{},
	}
	return nil
}

// PushRequest 压入玩家请求
func (d *desk) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, body proto.Message) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "desk.PushRequest",
		"desk_uid":   d.deskUID,
		"game_id":    d.gameID,
		"player_id":  playerID,
		"message_id": head.GetMsgId(),
	})

	trans := global.GetReqEventTranslator()
	eventID, eventContext, err := trans.Translate(playerID, head, body)
	if err != nil {
		logEntry.WithError(err).Errorln("消息转事件失败")
		return
	}

	d.event <- deskEvent{
		eventID:      eventID,
		eventContext: eventContext,
	}
}

func (d *desk) initMajongContext() error {
	return nil // TODO
}

func (d *desk) processEvents(ctx context.Context) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.processEvent",
		"desk_uid":  d.deskUID,
		"game_id":   d.gameID,
	})
	defer func() {
		if x := recover(); x != nil {
			logEntry.Errorln(x)
			debug.PrintStack()
		}
	}()

	select {
	case <-ctx.Done():
		{
			return
		}
	case autoEvent := <-d.autoEvent: // 需要确保 autoEvent 通道有 1 个缓冲区
		{
			d.processEvent(&deskEvent{
				eventID:      autoEvent.GetEventId(),
				eventContext: autoEvent.GetEventContext(),
			})
		}
	case event := <-d.event:
		{
			d.processEvent(&event)
		}
	}
}

// processEvent 处理单个事件
// step 1. 调用麻将逻辑的接口来处理事件(返回最新麻将现场, 自动事件， 发送给玩家的消息)， 并且更新 mjContext
// step 2. 将消息发送给玩家
// step 3. 调用 room 的结算逻辑来处理结算
// step 4. 如果有自动事件， 将自动事件写入自动事件通道
// step 5. 如果当前状态是游戏结束状态， 调用 cancel 终止游戏
func (d *desk) processEvent(event *deskEvent) {
	// TODO
	var mjContext server_pb.MajongContext
	var autoEvent *server_pb.AutoEvent

	// TODO 发送消息给玩家

	d.mjContext = mjContext
	d.settler.Settle(d, d.mjContext)
	if autoEvent != nil {
		d.autoEvent <- *autoEvent
	}

	// 游戏结束
	if d.mjContext.GetCurState() == server_pb.StateID_state_gameover {
		d.cancel()
	}
}
