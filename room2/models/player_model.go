package models

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/gutils"
	"steve/room2/contexts"
	"steve/room2/desk"
	"steve/room2/fixed"
	playerpkg "steve/room2/player"
	server_pb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type PlayerModel struct {
	BaseModel
	players []*playerpkg.Player
}

func (model PlayerModel) GetName() string {
	return fixed.Player
}
func (model *PlayerModel) Start() {
	model.players = make([]*playerpkg.Player, model.GetDesk().GetConfig().Num)
	ids := model.GetDesk().GetConfig().PlayerIds //GetModelManager().GetPlayerModel(model.GetDesk().GetUid()).GetDeskPlayerIDs()
	for i := 0; i < len(model.players); i++ {
		playerObj := playerpkg.GetPlayerMgr().GetPlayer(ids[i])
		playerObj.EnterDesk(model.GetDesk())
		model.players[i] = playerObj
	}
}
func (model PlayerModel) Stop() {
	playerMgr := playerpkg.GetPlayerMgr()

	playerIDs := make([]uint64, 0, len(model.players))
	for _, pla := range model.players {
		if !pla.IsDetached() {
			playerIDs = append(playerIDs, pla.GetPlayerID())
		}
	}
	playerMgr.UnbindPlayerRoomAddr(playerIDs)
}

func NewPlayertModel(desk *desk.Desk) DeskModel {
	result := &PlayerModel{}
	result.SetDesk(desk)
	return result
}

func (model *PlayerModel) PlayerEnter(player *playerpkg.Player) {
	// 判断行牌状态, 选项化后需修改
	player.EnterDesk(model.GetDesk())
	model.recoverGameForPlayer(player.GetPlayerID())
	model.setContextPlayerQuit(player, false)
	//d.playerQuitEnterDeskNtf(eqi.PlayerID, room.QuitEnterType_QET_ENTER)
	model.playerQuitEnterDeskNtf(player, room.QuitEnterType_QET_ENTER)
}

func (model *PlayerModel) recoverGameForPlayer(playerID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "recoverGameForPlayer",
		"playerID":  playerID,
	})
	ctx := model.GetDesk().GetConfig().Context.(*contexts.MjContext)
	mjContext := &ctx.MjContext
	bankerSeat := mjContext.GetZhuangjiaIndex()
	totalCardsNum := mjContext.GetCardTotalNum()
	gameStage := GetGameStage(mjContext.GetCurState())
	gameID := gutils.GameIDServer2Client(int(mjContext.GetGameId()))
	gameDeskInfo := room.GameDeskInfo{
		GameId:            &gameID,
		GameStage:         &gameStage,
		Players:           GetRecoverPlayerInfo(playerID, model.GetDesk()),
		Dices:             mjContext.GetDices(),
		BankerSeat:        &bankerSeat,
		EastSeat:          &bankerSeat,
		TotalCards:        &totalCardsNum,
		RemainCards:       proto.Uint32(uint32(len(mjContext.GetWallCards()))),
		CostTime:          proto.Uint32(GetStateCostTime(ctx.StateTime.Unix())),
		OperatePid:        GetOperatePlayerID(mjContext),
		NeedHsz:           proto.Bool(gutils.GameHasHszState(mjContext)),
		LastOutCard:       proto.Uint32(getLastOutCard(mjContext.GetLastOutCard())),
		LastOutCardPlayer: proto.Uint64(mjContext.GetLastChupaiPlayer()),
	}
	gameDeskInfo.HasZixun, gameDeskInfo.ZixunInfo = GetZixunInfo(playerID, mjContext)
	gameDeskInfo.HasWenxun, gameDeskInfo.WenxunInfo = GetWenxunInfo(playerID, mjContext)
	gameDeskInfo.HasQgh, gameDeskInfo.QghInfo = GetQghInfo(playerID, mjContext)

	_, gameDeskInfo.HuansanzhangInfo = getHuansanzhangInfo(playerID, mjContext)
	_, gameDeskInfo.DingqueInfo = getDingqueInfo(playerID, mjContext)
	if gameDeskInfo.GetHasZixun() {
		gameDeskInfo.DoorCard = GetDoorCard(mjContext)
	}
	rsp, err := proto.Marshal(&room.RoomResumeGameRsp{
		ResumeRes: room.RoomError_SUCCESS.Enum(),
		GameInfo:  &gameDeskInfo,
	})
	logEntry.WithField("desk_info", gameDeskInfo).Infoln("恢复数据")
	if err != nil {
		logEntry.WithError(err).Errorln("序列化失败")
		return
	}
	GetModelManager().GetMjEventModel(model.GetDesk().GetUid()).Reply([]server_pb.ReplyClientMessage{
		server_pb.ReplyClientMessage{
			Players: []uint64{playerID},
			MsgId:   int32(msgid.MsgID_ROOM_RESUME_GAME_RSP),
			Msg:     rsp,
		},
	})
}

func (model *PlayerModel) needTuoguan() bool {
	mjContext := model.GetGameContext().(*contexts.MjContext)
	state := mjContext.MjContext.GetCurState()
	switch state {
	case server_pb.StateID_state_init,
		server_pb.StateID_state_fapai,
		server_pb.StateID_state_huansanzhang,
		server_pb.StateID_state_dingque:
		return false
	}
	return true
}

func (model *PlayerModel) PlayerQuit(player *playerpkg.Player) {
	player.QuitDesk(model.GetDesk(), model.needTuoguan())
	//d.setMjPlayerQuitDesk(eqi.PlayerID, true)
	model.setContextPlayerQuit(player, true)
	//d.playerQuitEnterDeskNtf(eqi.PlayerID, room.QuitEnterType_QET_QUIT)
	model.playerQuitEnterDeskNtf(player, room.QuitEnterType_QET_QUIT)

	mjContext := model.GetGameContext().(*contexts.MjContext)
	majongPlayer := gutils.GetMajongPlayer(player.GetPlayerID(), &mjContext.MjContext)
	if !gutils.IsPlayerContinue(majongPlayer.GetXpState(), &mjContext.MjContext) {
		playerMgr := playerpkg.GetPlayerMgr()
		playerMgr.UnbindPlayerRoomAddr([]uint64{player.GetPlayerID()})
		player.SetDetached(true)
	}
}

func (model *PlayerModel) playerQuitEnterDeskNtf(player *playerpkg.Player, qeType room.QuitEnterType) {
	if player == nil {
		return
	}
	roomPlayer := TranslateToRoomPlayer(player)
	playerId := player.GetPlayerID()
	ntf := room.RoomDeskQuitEnterNtf{
		PlayerId:   &playerId,
		Type:       &qeType,
		PlayerInfo: &roomPlayer,
	}
	GetModelManager().GetMessageModel(model.GetDesk().GetUid()).BroadCastDeskMessageExcept([]uint64{playerId}, true, msgid.MsgID_ROOM_DESK_QUIT_ENTER_NTF, &ntf)
}

func (model *PlayerModel) setContextPlayerQuit(player *playerpkg.Player, value bool) {
	for _, p := range model.GetDesk().GetConfig().Context.(*contexts.MjContext).MjContext.Players {
		if p.GetPalyerId() == player.GetPlayerID() {
			p.IsQuit = value
		}
	}
}

func (model *PlayerModel) GetDeskPlayers() []*playerpkg.Player {
	return model.players
}

// GetDeskPlayerIDs 获取牌桌玩家 ID 列表， 座号作为索引
func (model *PlayerModel) GetDeskPlayerIDs() []uint64 {
	players := model.GetDeskPlayers()
	result := make([]uint64, len(players))
	for _, player := range players {
		result[player.GetSeat()] = player.GetPlayerID()
	}
	return result
}

// GetTuoguanPlayers 获取牌桌所有托管玩家
func (model *PlayerModel) GetTuoguanPlayers() []uint64 {
	players := GetModelManager().GetPlayerModel(model.GetDesk().GetUid()).GetDeskPlayers()
	result := make([]uint64, 0, len(players))
	for _, pla := range players {
		if pla.IsTuoguan() {
			result = append(result, pla.GetPlayerID())
		}
	}
	return result
}
