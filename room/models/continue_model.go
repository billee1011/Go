package models

import (
	client_match_pb "steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/external/goldclient"
	"steve/room/desk"
	"steve/room/fixed"
	playerpkg "steve/room/player"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type continueRequestInfo struct {
	playerID uint64
	request  *client_match_pb.MatchDeskContinueReq
	response *client_match_pb.MatchDeskContinueRsp
	finish   chan struct{}
}

// ContinueModel 续局 model
type ContinueModel struct {
	BaseModel

	requestChannel chan continueRequestInfo
	readyPlayers   []uint64
	started        bool          // 是否已经启动续局逻辑
	startChannel   chan struct{} // 续局逻辑开启
	stopChannel    chan struct{} // model 停止
	continueTime   time.Time     // 启动续局逻辑的时间

	zhuang    int
	fixzhuang bool
}

// NewContinueModel 创建续局 model
func NewContinueModel(desk *desk.Desk) DeskModel {
	result := &ContinueModel{
		requestChannel: make(chan continueRequestInfo, 4),
		startChannel:   make(chan struct{}),
		stopChannel:    make(chan struct{}),
		started:        false,
		readyPlayers:   make([]uint64, 0, 4),
	}
	result.SetDesk(desk)
	return result
}

// GetName 获取 model 名称
func (model *ContinueModel) GetName() string {
	return fixed.ChatModelName
}

// Active 激活 model
func (model *ContinueModel) Active() {}

// Start 启动 model
func (model *ContinueModel) Start() {
	go model.run()
}

// startContinue 开始续局逻辑
func (model *ContinueModel) startContinue() {
	model.startChannel <- struct{}{}
}

func (model *ContinueModel) run() {
	ticker := time.NewTicker(time.Second * 1)
forstart:
	for {
		select {
		case <-model.startChannel:
			{
				model.started = true
				model.continueTime = time.Now()
			}
		case requestInfo := <-model.requestChannel:
			{
				model.handlePlayerContinueRequest(&requestInfo)
				close(requestInfo.finish)
			}
		case <-ticker.C:
			{
				model.checkDismiss()
			}
		case <-model.stopChannel:
			{
				break forstart
			}
		}
	}
}

// checkDismiss 超过20s解散牌桌
func (model *ContinueModel) checkDismiss() {
	if !model.started {
		return
	}
	if time.Now().Sub(model.continueTime) <= time.Second*20 {
		return
	}
	modelMgr := GetModelManager()
	deskID := model.GetDesk().GetUid()

	modelMgr.GetMessageModel(deskID).BroadCastDeskMessageExcept(nil, true,
		msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF,
		&client_match_pb.MatchContinueDeskDimissNtf{})

	modelMgr.StopDeskModel(model.GetDesk())
}

// handleCancelRequest 处理取消续局请求
func (model *ContinueModel) handleCancelRequest(requestInfo *continueRequestInfo) {
	modelMgr := GetModelManager()
	deskID := model.GetDesk().GetUid()
	modelMgr.GetMessageModel(deskID).BroadCastDeskMessageExcept([]uint64{requestInfo.playerID}, true,
		msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF, &client_match_pb.MatchContinueDeskDimissNtf{})

	modelMgr.StopDeskModel(model.GetDesk())
}

func (model *ContinueModel) handlePlayerContinueRequest(requestInfo *continueRequestInfo) {
	playerID := requestInfo.playerID
	entry := logrus.WithField("player_id", playerID)

	response := requestInfo.response
	if !model.started {
		response.ErrCode = proto.Int32(1)
		response.ErrDesc = proto.String("当前不能续局")
		return
	}

	if requestInfo.request.GetCancel() {
		model.handleCancelRequest(requestInfo)
		return
	}

	for _, _playerID := range model.readyPlayers {
		if _playerID == playerID {
			response.ErrCode = proto.Int32(0)
			response.ErrDesc = proto.String("")
			return
		}
	}
	response.ErrCode = proto.Int32(1)
	response.ErrDesc = proto.String("续局失败")

	playerCoin, err := goldclient.GetGold(playerID, 1)
	if err != nil {
		entry.WithError(err).Errorln("获取玩家金币数失败")
		return
	}
	desk := model.GetDesk()
	if uint64(playerCoin) < desk.GetConfig().MinScore {
		response.ErrDesc = proto.String("金豆不足")
		return
	}
	response.ErrCode = proto.Int32(0)
	response.ErrDesc = proto.String("")

	model.readyPlayers = append(model.readyPlayers, playerID)

	playerModel := GetModelManager().GetPlayerModel(desk.GetUid())

	allplayer := playerModel.GetDeskPlayerIDs()
	for _, _playerID := range allplayer {
		exist := false
		for _, readyPlayerID := range model.readyPlayers {
			if _playerID == readyPlayerID {
				exist = true
			}
		}
		if !exist {
			return
		}
	}
	model.startNextRound()
}

// startNextRound 开始下一局
func (model *ContinueModel) startNextRound() {
	model.started = false
	model.readyPlayers = model.readyPlayers[0:0]

	desk := model.GetDesk()
	var err error
	desk.GetConfig().Context, err = createDeskContext(desk.GetGameId(), desk.GetPlayerIds(), model.zhuang, model.fixzhuang)
	if err != nil {
		logrus.WithField("players", model.readyPlayers).Errorln("初始化牌桌现场失败")
		return
	}
	eventModel := GetEventModel(model.GetDesk().GetUid())
	eventModel.StartProcessEvents()
}

// Stop 停止 model
func (model *ContinueModel) Stop() {
	model.stopChannel <- struct{}{}
}

// PushContinueRequest 处理玩家续局请求
func (model *ContinueModel) PushContinueRequest(playerID uint64, request *client_match_pb.MatchDeskContinueReq) (response client_match_pb.MatchDeskContinueRsp) {
	finish := make(chan struct{})
	model.requestChannel <- continueRequestInfo{
		playerID: playerID,
		request:  request,
		response: &response,
		finish:   finish,
	}
	select {
	case <-finish:
		break
	case <-time.NewTimer(time.Second).C: // 1s 后还没处理，直接返回
		logrus.WithField("player_id", playerID).Warningln("处理超时")
		break
	}
	return
}

// ContinueDesk 开始续局逻辑
func (model *ContinueModel) ContinueDesk(fixBanker bool, bankerSeat int, settleMap map[uint64]int64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":   "DeskBase.ContinueDesk",
		"fix_banker":  fixBanker,
		"banker_seat": bankerSeat,
	})

	model.fixzhuang = fixBanker
	model.zhuang = bankerSeat

	desk := model.GetDesk()
	playerIDs := desk.GetPlayerIds()
	playerMgr := playerpkg.GetPlayerMgr()

	for _, playerID := range playerIDs {
		player := playerMgr.GetPlayer(playerID)
		if player == nil || player.GetDesk() != desk || player.IsQuit() { // 玩家已经退出牌桌，不续局
			entry.WithFields(logrus.Fields{
				"player_id": playerID,
				"quited":    player.IsQuit(),
			}).Debugln("玩家不满足续局条件")
			messageModel := GetModelManager().GetMessageModel(desk.GetUid())
			messageModel.BroadCastDeskMessage(nil, msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF,
				&client_match_pb.MatchContinueDeskDimissNtf{}, true)
			return
		}
	}

	model.startContinue()
}
