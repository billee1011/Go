package game

import (
	"context"
	"github.com/golang/protobuf/proto"
	"runtime/debug"
	"steve/client_pb/msgid"
	mahjonghandler "steve/majong/export/process"
	"steve/room/interfaces/global"
	mahjong "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"
	"time"
)

// deskEvent 房间事件
type deskEvent struct {
	id          mahjong.EventID // 事件 ID
	context     []byte          // 事件现场
	eventType   int             // 事件类型
	playerID    uint64          // 针对哪个玩家的事件
	stateNumber int
}

type deskContext struct {
	context     interface{} // 牌局现场
	stateNumber int         // 状态序号
	stateTime   time.Time   // 状态时间
}

type Desk struct {
	deskID  uint64
	gameID  int
	players []*Player

	option *DeskOption

	dContext *deskContext // 牌桌现场

	event chan deskEvent

	cancel context.CancelFunc
}

func NewDesk(deskID uint64, gameID int, option *DeskOption) *Desk {
	// 1. 根据option创建桌子
	return &Desk{
		deskID:  deskID,
		gameID:  gameID,
		players: []*Player{},

		option: option,

		event: make(chan deskEvent, 16), //todo 16???
	}
}

func (d *Desk) Start() error {
	d.initGameContext(d.gameID)

	var ctx context.Context
	ctx, d.cancel = context.WithCancel(context.Background())
	go func() {
		d.process(ctx)
	}()

	d.event <- deskEvent{
		id:          mahjong.EventID_event_start_game,
		context:     []byte{},
		eventType:   1,
		stateNumber: d.dContext.stateNumber,
	}
	return nil
}

func (d *Desk) initGameContext(gameID int) {
	gameContext := InitGameContext(gameID)

	d.dContext = &deskContext{
		context:     gameContext,
		stateNumber: 0,
		stateTime:   time.Now(),
	}
}

func (d *Desk) Stop() {
	d.cancel()

}

func (d *Desk) process(ctx context.Context) {
	defer func() {
		if x := recover(); x != nil {
			debug.PrintStack()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case event := <-d.event:
			{
				result, ok := d.handleEvent(event.id, event.context)
				if !ok {

				}

				d.Reply(result.ReplyMsgs)

			}
		}
	}
}

func (d *Desk) handleEvent(eventID mahjong.EventID, eventContext []byte) (result mahjonghandler.HandleMajongEventResult, ok bool) {
	ok = true
	stateNumber, gameContext, stateTime := d.dContext.stateNumber, d.dContext.context.(mahjong.MajongContext), d.dContext.stateTime
	oldState := gameContext.GetCurState()

	if result = mahjonghandler.HandleMajongEvent(mahjonghandler.HandleMajongEventParams{
		MajongContext: gameContext,
		EventID:       eventID,
		EventContext:  eventContext,
	}); !result.Succeed {
		ok = false
		return
	}

	newContext := result.NewContext
	newState := newContext.GetCurState()
	if newState != oldState {
		stateNumber++
		stateTime = time.Now()
	}

	// dContext 的每次修改都是一块新内存，用来确保并发安全。
	d.dContext = &deskContext{
		context:     newContext,
		stateNumber: stateNumber,
		stateTime:   stateTime,
	}

	return
}

func (d *Desk) Reply(replyMsgs []mahjong.ReplyClientMessage) {
	if replyMsgs == nil {
		return
	}

	for _, msg := range replyMsgs {
		d.BroadcastMessage(msg.GetPlayers(), msgid.MsgID(msg.GetMsgId()), msg.GetMsg())
	}
}

func (d *Desk) BroadcastMessage(playerIDs []uint64, msgID msgid.MsgID, body []byte) {
	if playerIDs == nil || len(playerIDs) == 0 {
		for _, player := range d.players {
			if !DefaultPlayManager.GetPlayer(player.playerID).isQuit {
				playerIDs = append(playerIDs, player.playerID)
			}
		}
	}

	if len(playerIDs) == 0 {
		return
	}
	DefaultSender.BroadCastMessageBare(playerIDs, msgID, body)
}

func (d *Desk) HandlePlayerRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	trans := global.GetReqEventTranslator()
	eventID, eventContext, err := trans.Translate(playerID, head, bodyData)
	if err != nil {
		return
	}
	eventMessage, ok := eventContext.(proto.Message)
	if !ok {
		return
	}
	eventConetxt, err := proto.Marshal(eventMessage)
	if err != nil {
		return
	}

	d.event <- deskEvent{
		id:          mahjong.EventID(eventID),
		context:     eventConetxt,
		eventType:   1,
		playerID:    playerID,
		stateNumber: d.dContext.stateNumber,
	}
}

func (d *Desk) Player(playerID uint64) *Player {
	for _, player := range d.players {
		if player.playerID == playerID {
			return player
		}
	}
	return nil
}
