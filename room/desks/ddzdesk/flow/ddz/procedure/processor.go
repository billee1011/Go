package procedure

import (
	msgid "steve/client_pb/msgid"
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
		message := &ddz.ResumeRequestEvent{}
		err := proto.Unmarshal(params.EventContext, message)
		if err != nil {
			logEntry.WithError(err).Errorln("处理恢复对局事件失败")
			return
		}

		// 请求的玩家ID
		reqPlayerID := message.GetHead().GetPlayerId()

		bExist := false

		// 找到这个玩家
		for _, player := range cloneContext.GetPlayers() {
			if player.GetPalyerId() == reqPlayerID {
				bExist = true
			}
		}

		// 存在的话则发送游戏信息
		if bExist {
			playersInfo := []*room.DDZPlayerInfo{}

			for _, player := range cloneContext.GetPlayers() {

				// Player转为RoomPlayer
				roomPlayerInfo := TranslateDDZPlayerToRoomPlayer(*player)
				lord := player.GetLord()
				// double := player.GetIsDouble()
				tuoguan := false // TODO

				ddzPlayerInfo := room.DDZPlayerInfo{}
				ddzPlayerInfo.PlayerInfo = &roomPlayerInfo
				ddzPlayerInfo.OutCards = player.GetOutCards()
				ddzPlayerInfo.HandCards = player.GetHandCards()
				ddzPlayerInfo.Lord = &lord
				// ddzPlayerInfo.IsDouble = &double
				ddzPlayerInfo.Tuoguan = &tuoguan

				playersInfo = append(playersInfo, &ddzPlayerInfo)
			}

			var errCode uint32 = 0
			errDesc := ""
			//stage := room.DDZStage_DDZ_STAGE_PLAYING

			// 发送游戏信息
			m.SendMessage([]uint64{reqPlayerID}, msgid.MsgID_ROOM_DDZ_RESUME_REQ, &room.DDZResumeGameRsp{
				Result: &room.Result{ErrCode: &errCode, ErrDesc: &errDesc},
				GameInfo: &room.DDZDeskInfo{
					Players: playersInfo,
					//Stage: , TODO
				},
			})
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

// TranslateDDZPlayerToRoomPlayer 将 ddzPlayer 转换成 RoomPlayerInfo
func TranslateDDZPlayerToRoomPlayer(ddzPlayer ddz.Player) room.RoomPlayerInfo {
	playerMgr := global.GetPlayerMgr()
	playerID := ddzPlayer.GetPalyerId()
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
