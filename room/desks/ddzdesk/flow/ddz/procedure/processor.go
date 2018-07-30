package procedure

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/room/desks/ddzdesk/flow/ddz/ddzmachine"
	"steve/room/desks/ddzdesk/flow/ddz/states"
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/room/interfaces/global"
	"steve/server_pb/ddz"
	"time"

	"steve/room/interfaces"
	"steve/room/interfaces/facade"

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
	PlayerMgr    interfaces.DeskPlayerMgr // 是否托管
	Context      ddz.DDZContext           // 牌局现场
	Sender       ddzmachine.MessageSender // 消息发送器， 拆分后要修改
	EventID      int                      // 事件 ID
	EventContext []byte                   // 事件现场
}

// HandleEvent 处理牌局事件
func HandleEvent(params Params) (result Result) {
	start := time.Now()
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
		resumeErr := dealResumeRequest(&params, m, &cloneContext)
		if resumeErr != nil {
			logEntry.WithError(resumeErr).Errorln("处理恢复对局失败")
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

	end := time.Now()
	logrus.WithField("duration", end.Sub(start)).Debug("状态机从创建到退出")
	return
}

// 处理恢复对局的请求
// eventContext : 事件体
// machine		: 状态机
// ddzContext	: 斗地主牌局信息
// bool 		: 成功/失败
func dealResumeRequest(param *Params, machine *ddzmachine.DDZMachine, ddzContext *ddz.DDZContext) error {
	message := &ddz.ResumeRequestEvent{}
	err := proto.Unmarshal(param.EventContext, message)
	if err != nil {
		return err
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
		var playersInfo []*room.DDZPlayerInfo

		players := ddzContext.GetPlayers()
		for index := 0; index < len(players); index++ {
			player := players[index]

			// Player转为RoomPlayer
			roomPlayerInfo := TranslateDDZPlayerToRoomPlayer(*player, uint32(index))
			lord := player.GetLord()
			//double := player.GetIsDouble()
			deskPlayer := facade.GetDeskPlayerByID(param.PlayerMgr, player.GetPlayerId())
			tuoguan := deskPlayer.IsTuoguan()

			ddzPlayerInfo := room.DDZPlayerInfo{}

			ddzPlayerInfo.PlayerInfo = &roomPlayerInfo
			ddzPlayerInfo.OutCards = player.GetOutCards()

			// 只发送自己的手牌，其他人的手牌为空
			if player.GetPlayerId() == reqPlayerID {
				ddzPlayerInfo.HandCards = player.GetHandCards()
			} else {
				ddzPlayerInfo.HandCards = []uint32{}
			}
			ddzPlayerInfo.HandCardsCount = proto.Uint32(uint32(len(player.GetHandCards())))

			ddzPlayerInfo.Lord = &lord
			ddzPlayerInfo.Tuoguan = &tuoguan

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
						logrus.Error("出现错误！firstPlayerID为0，但是player.GetGrab()为true")
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

		// 开始时间
		startTime := time.Time{}
		startTime.UnmarshalBinary(ddzContext.StartTime)
		logrus.Debugf("startTime = %v", startTime)

		// 限制时间
		duration := time.Second * time.Duration(ddzContext.Duration)
		logrus.Debugf("duration = %v", duration)

		// 剩余时间
		leftTime := duration - time.Now().Sub(startTime)
		logrus.Debugf("leftTime = %v", leftTime)

		if leftTime < 0 {
			leftTime = 0
		}
		leftTimeInt32 := uint32(leftTime.Seconds())

		curStage := room.DDZStage(int32(ddzContext.CurStage))
		curCardType := room.CardType(ddzContext.GetCurCardType())

		totalGrab := ddzContext.GetTotalGrab()
		if totalGrab == 0 {
			totalGrab = 1
		}

		resumeMsg := &room.DDZResumeGameRsp{
			Result: &room.Result{ErrCode: proto.Uint32(0), ErrDesc: proto.String("")},
			GameInfo: &room.DDZDeskInfo{
				Players: playersInfo, // 每个人的信息
				Stage: &room.NextStage{
					Stage: &curStage,
					Time:  proto.Uint32(leftTimeInt32),
				},
				CurPlayerId:  proto.Uint64(ddzContext.GetCurrentPlayerId()), // 当前操作的玩家
				Dipai:        ddzContext.GetDipai(),
				TotalGrab:    &totalGrab,
				TotalDouble:  proto.Uint32(ddzContext.GetTotalDouble()),
				TotalBomb:    proto.Uint32(ddzContext.GetTotalBomb()),
				CurCardType:  &curCardType,
				CurCardPivot: proto.Uint32(ddzContext.GetCardTypePivot()),
				CurOutCards:  ddzContext.CurOutCards,
			},
		}
		// 发送游戏信息
		machine.SendMessage([]uint64{reqPlayerID}, msgid.MsgID_ROOM_DDZ_RESUME_RSP, resumeMsg)
		logrus.WithField("resumeMsg", resumeMsg).WithField("playerId", reqPlayerID).Infoln("斗地主发送恢复对局消息")
	}

	return nil
}

// TranslateDDZPlayerToRoomPlayer 将 ddzPlayer 转换成 RoomPlayerInfo
func TranslateDDZPlayerToRoomPlayer(ddzPlayer ddz.Player, seat uint32) room.RoomPlayerInfo {
	playerMgr := global.GetPlayerMgr()
	playerID := ddzPlayer.GetPlayerId()
	player := playerMgr.GetPlayer(playerID)
	var coin uint64
	if player != nil {
		coin = player.GetCoin()
	}

	return room.RoomPlayerInfo{
		PlayerId: proto.Uint64(playerID),
		Name:     proto.String("player" + string(playerID)),
		Coin:     proto.Uint64(coin),
		Seat:     proto.Uint32(seat),
		// Location: TODO 没地方拿
	}
}
