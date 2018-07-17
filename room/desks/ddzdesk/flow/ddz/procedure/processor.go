package procedure

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/desks/ddzdesk/flow/ddz/ddzmachine"
	"steve/room/desks/ddzdesk/flow/ddz/states"
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/room/interfaces/global"
	"steve/server_pb/ddz"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// Result 处理牌局事件的结果
type Result struct {
	Context           ddz.DDZContext // 最新现场
	HasAutoEvent      bool
	AutoEventID       int
	AutoEventContext  []byte
	AutoEventDuration time.Duration
	Succeed           bool // 是否成功
}

// Params 处理牌局事件的参数
type Params struct {
	Context      ddz.DDZContext           // 牌局现场
	Sender       ddzmachine.MessageSender // 消息发送器， 拆分后要修改
	EventID      int                      // 事件 ID
	EventContext []byte                   // 事件现场
}

// HandleEvent 处理牌局事件
func HandleEvent(params Params) (result Result) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleEvent",
		"params":    params,
	})

	cloneContext := *proto.Clone(&params.Context).(*ddz.DDZContext)

	result = Result{
		Context:      cloneContext,
		Succeed:      false,
		HasAutoEvent: false,
	}
	m := ddzmachine.CreateDDZMachine(&cloneContext, states.NewFactory(), params.Sender)

	// 处理恢复对局的请求
	if params.EventID == int(ddz.EventID_event_resume_request) {
		if dealResumeRequest(params.EventContext, m, &cloneContext) == false {
			logEntry.Errorln("处理恢复对局失败")
		}
		return
	}

	err := m.ProcessEvent(machine.Event{
		EventID:   params.EventID,
		EventData: params.EventContext,
	})
	if err != nil {
		logEntry.WithError(err).Errorln("处理事件失败")
		return
	}
	result.Context = *m.GetDDZContext()
	e, d := m.GetAutoEvent()
	if e != nil {
		result.HasAutoEvent = true
		result.AutoEventID = e.EventID
		result.AutoEventContext = e.EventData
		result.AutoEventDuration = d
	} else {
		result.HasAutoEvent = false
	}
	result.Succeed = true
	return
}

// 处理恢复对局的请求
// eventContext : 事件体
// machine		: 状态机
// ddzContext	: 斗地主牌局信息
// bool 		: 成功/失败
func dealResumeRequest(eventContext []byte, machine *ddzmachine.DDZMachine, ddzContext *ddz.DDZContext) bool {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "dealResumeRequest",
	})

	message := &ddz.ResumeRequestEvent{}
	err := proto.Unmarshal(eventContext, message)
	if err != nil {
		logEntry.WithError(err).Errorln("处理恢复对局事件失败")
		return false
	}

	// 请求的玩家ID
	reqPlayerID := message.GetHead().GetPlayerId()

	bExist := false

	// 找到这个玩家
	for _, player := range ddzContext.GetPlayers() {
		if player.GetPlayerId() == reqPlayerID {
			bExist = true
		}
	}

	// 存在的话则发送游戏信息
	if bExist {
		playersInfo := []*room.DDZPlayerInfo{}

		for _, player := range ddzContext.GetPlayers() {

			// Player转为RoomPlayer
			roomPlayerInfo := TranslateDDZPlayerToRoomPlayer(*player)
			lord := player.GetLord()
			//double := player.GetIsDouble()
			tuoguan := false // TODO

			ddzPlayerInfo := room.DDZPlayerInfo{}

			ddzPlayerInfo.PlayerInfo = &roomPlayerInfo
			ddzPlayerInfo.OutCards = player.GetOutCards()

			// 只发送自己的手牌，其他人的手牌为空
			if player.GetPlayerId() == reqPlayerID {
				ddzPlayerInfo.HandCards = player.GetHandCards()
			} else {
				ddzPlayerInfo.HandCards = []uint32{}
			}

			ddzPlayerInfo.Lord = &lord
			ddzPlayerInfo.Tuoguan = &tuoguan
			//ddzPlayerInfo.Grablord =

			// 叫/抢地主
			if ddzContext.CurState == ddz.StateID_state_grab {

				grabLord := room.GrabLordType_GLT_CALLLORD      // 叫地主
				notGrabLord := room.GrabLordType_GLT_NOTCALLORD // 不叫
				grab := room.GrabLordType_GLT_GRAB              // 抢地主
				notGrab := room.GrabLordType_GLT_NOTGRAB        // 不抢
				noneOpe := room.GrabLordType_GLT_NONE           // 未操作

				// 首个叫地主玩家的playerID
				firstPlayerID := ddzContext.GetFirstGrabPlayerId()

				// 自己是否操作过
				bOpera := false
				allOpePlayers := ddzContext.GetGrabbedPlayers()
				for i := 0; i < len(allOpePlayers); i++ {
					if allOpePlayers[i] == player.GetPlayerId() {
						bOpera = true
						break
					}
				}

				// 没人叫过地主
				if firstPlayerID == 0 {
					if player.GetGrab() {
						logEntry.Error("出现错误！firstPlayerID为0，但是player.GetGrab()为true")
						ddzPlayerInfo.GrabLord = &grabLord
					} else {
						// 操作过说明不叫，否则说明未操作
						if bOpera {
							ddzPlayerInfo.GrabLord = &notGrabLord
						} else {
							ddzPlayerInfo.GrabLord = &noneOpe
						}
					}
				} else { // 已经有人叫过地主了

					// 自己叫/抢过
					if player.GetGrab() {

						// 首次的是自己，说明是叫地主; 首次的不是自己，说明是抢地主
						if ddzContext.GetFirstGrabPlayerId() == player.GetPlayerId() {
							ddzPlayerInfo.GrabLord = &grabLord
						} else {
							ddzPlayerInfo.GrabLord = &grab
						}
					} else { // 自己没叫/抢过
						// 操作过说明是不叫/不抢，否则说明未操作
						if bExist {
							// 指定叫地主的是自己，说明是不叫;否则说明是不抢
							if ddzContext.GetCallPlayerId() == player.GetPlayerId() {
								ddzPlayerInfo.GrabLord = &notGrabLord
							} else {
								ddzPlayerInfo.GrabLord = &notGrab
							}
						} else {
							ddzPlayerInfo.GrabLord = &noneOpe
						}
					}
				}
			}

			// 加倍状态
			if ddzContext.CurState == ddz.StateID_state_double {

				double := room.DoubleType_DT_DOUBLE
				notDouble := room.DoubleType_DT_NOTDOUBLE
				noneOpera := room.DoubleType_DT_NONE

				// 加倍
				if player.GetIsDouble() {
					ddzPlayerInfo.Double = &double
				} else {
					// 是否已操作过
					bExist := false
					allOpePlayers := ddzContext.GetDoubledPlayers()
					for i := 0; i < len(allOpePlayers); i++ {
						if allOpePlayers[i] == player.GetPlayerId() {
							bExist = true
							break
						}
					}

					// 存在说明不加倍，不存在说明未操作
					if bExist {
						ddzPlayerInfo.Double = &notDouble
					} else {
						ddzPlayerInfo.Double = &noneOpera
					}
				}
			}

			playersInfo = append(playersInfo, &ddzPlayerInfo)
		}

		var errCode uint32 = 0
		errDesc := ""

		// 开始时间
		startTime := time.Time{}
		startTime.UnmarshalBinary(ddzContext.StartTime)

		// 限制时间
		duration := time.Second * time.Duration(ddzContext.Duration)

		// 剩余时间
		leftTime := duration - time.Now().Sub(startTime)

		if leftTime < 0 {
			leftTime = 0
		}

		curStage := room.DDZStage(int32(ddzContext.CurStage))

		// 发送游戏信息
		machine.SendMessage([]uint64{reqPlayerID}, msgid.MsgID_ROOM_DDZ_RESUME_RSP, &room.DDZResumeGameRsp{
			Result: &room.Result{ErrCode: &errCode, ErrDesc: &errDesc},
			GameInfo: &room.DDZDeskInfo{
				Players: playersInfo, // 每个人的信息
				Stage: &room.NextStage{
					Stage: &curStage,
					Time:  proto.Uint32(uint32(leftTime)),
				},
				CurPlayerId: proto.Uint64(ddzContext.GetCurrentPlayerId()), // 当前操作的玩家
				Dipai:       ddzContext.GetDipai(),
			},
		})
	}

	return true
}

// TranslateDDZPlayerToRoomPlayer 将 ddzPlayer 转换成 RoomPlayerInfo
func TranslateDDZPlayerToRoomPlayer(ddzPlayer ddz.Player) room.RoomPlayerInfo {
	playerMgr := global.GetPlayerMgr()
	playerID := ddzPlayer.GetPlayerId()
	player := playerMgr.GetPlayer(playerID)
	var coin uint64
	if player != nil {
		coin = player.GetCoin()
	}
	return room.RoomPlayerInfo{
		PlayerId: proto.Uint64(playerID),
		Name:     proto.String(""), // TODO
		Coin:     proto.Uint64(coin),
		Seat:     proto.Uint32(0), // TODO
		// Location: TODO 没地方拿
	}
}
