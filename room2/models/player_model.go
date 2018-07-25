package models

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	server_pb "steve/entity/majong"
	"steve/gutils"
	"steve/room2/contexts"
	"steve/room2/desk"
	"steve/room2/fixed"
	"steve/room2/player"
	"steve/room2/util"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type PlayerModel struct {
	BaseModel
	players []*player.Player
}

func (model PlayerModel) GetName() string {
	return fixed.Player
}
func (model *PlayerModel) Start() {
	model.players = make([]*player.Player, model.GetDesk().GetConfig().Num)
	ids := model.GetDesk().GetConfig().PlayerIds //GetModelManager().GetPlayerModel(model.GetDesk().GetUid()).GetDeskPlayerIDs()
	for i := 0; i < len(model.players); i++ {
		model.players[i] = player.GetPlayerMgr().GetPlayer(ids[i])
	}
}
func (model PlayerModel) Stop() {

}

func NewPlayertModel(desk *desk.Desk) DeskModel {
	result := &PlayerModel{}
	result.SetDesk(desk)
	return result
}

func (model *PlayerModel) PlayerEnter(player *player.Player) {
	// 判断行牌状态, 选项化后需修改
	context := player.GetDesk().GetConfig().Context.(*contexts.MjContext).MjContext
	mjPlayer := util.GetMajongPlayer(player.PlayerID, &context)
	// 非主动退出，再进入后取消托管；主动退出再进入不取消托管
	// 胡牌后没有托管，但是在客户端退出时，需要托管来自动胡牌,重新进入后把托管取消
	if !player.IsQuit() || mjPlayer.GetXpState() != server_pb.XingPaiState_normal {
		player.SetTuoguan(false, false)
	}
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
	ctx := model.GetDesk().GetConfig().Context.(*contexts.MjContext).MjContext
	mjContext := &ctx
	bankerSeat := mjContext.GetZhuangjiaIndex()
	totalCardsNum := mjContext.GetCardTotalNum()
	gameStage := GetGameStage(mjContext.GetCurState())
	gameID := gutils.GameIDServer2Client(int(mjContext.GetGameId()))
	gameDeskInfo := room.GameDeskInfo{
		GameId:      &gameID,
		GameStage:   &gameStage,
		Players:     GetRecoverPlayerInfo(playerID, model.GetDesk()),
		Dices:       mjContext.GetDices(),
		BankerSeat:  &bankerSeat,
		EastSeat:    &bankerSeat,
		TotalCards:  &totalCardsNum,
		RemainCards: proto.Uint32(uint32(len(mjContext.GetWallCards()))),
		CostTime:    proto.Uint32(GetStateCostTime(model.GetDesk().GetConfig().Context.(*contexts.MjContext).StateTime.Unix())),
		OperatePid:  GetOperatePlayerID(mjContext),
		DoorCard:    GetDoorCard(mjContext),
		NeedHsz:     proto.Bool(gutils.GameHasHszState(mjContext)),
	}
	gameDeskInfo.HasZixun, gameDeskInfo.ZixunInfo = GetZixunInfo(playerID, mjContext)
	gameDeskInfo.HasWenxun, gameDeskInfo.WenxunInfo = GetWenxunInfo(playerID, mjContext)
	gameDeskInfo.HasQgh, gameDeskInfo.QghInfo = GetQghInfo(playerID, mjContext)
	rsp, err := proto.Marshal(&room.RoomResumeGameRsp{
		ResumeRes: room.RoomError_SUCCESS.Enum(),
		GameInfo:  &gameDeskInfo,
	})
	logEntry.Infoln("恢复数据")
	logEntry.Infoln(gameDeskInfo)
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

func (model *PlayerModel) PlayerQuit(player *player.Player) {
	player.QuitDesk(model.GetDesk())
	//d.setMjPlayerQuitDesk(eqi.PlayerID, true)
	model.setContextPlayerQuit(player, true)
	//d.playerQuitEnterDeskNtf(eqi.PlayerID, room.QuitEnterType_QET_QUIT)
	model.playerQuitEnterDeskNtf(player, room.QuitEnterType_QET_QUIT)
}

func (model *PlayerModel) playerQuitEnterDeskNtf(player *player.Player, qeType room.QuitEnterType) {
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

func (model *PlayerModel) setContextPlayerQuit(player *player.Player, value bool) {
	for _, p := range model.GetDesk().GetConfig().Context.(*contexts.MjContext).MjContext.Players {
		if p.GetPalyerId() == player.GetPlayerID() {
			p.IsQuit = value
		}
	}
}

func (model *PlayerModel) GetDeskPlayers() []*player.Player {
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
