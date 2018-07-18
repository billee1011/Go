package public

import (
	"steve/room2/desk/models"
	"steve/room2/desk/player"
	"steve/client_pb/room"
	"steve/room2/desk/contexts"
	"steve/client_pb/msgid"
	"steve/room2/util"
	server_pb "steve/server_pb/majong"
	"github.com/golang/protobuf/proto"
	"steve/gutils"
	"github.com/Sirupsen/logrus"
)

type PlayerModel struct {
	BaseModel
	players []*player.Player
}
func (model PlayerModel) GetName() string{
	return models.Player
}
func (model PlayerModel) Start(){
	model.players = make([]*player.Player,model.GetDesk().GetConfig().Num)
}
func (model PlayerModel) Stop(){

}

func (model PlayerModel) PlayerEnter(player *player.Player,seat uint32){
	player.SetSeat(seat)
	player.EnterDesk(model.GetDesk())

	// 判断行牌状态, 选项化后需修改
	context := player.GetDesk().GetConfig().Context.(contexts.MjContext).MjContext
	mjPlayer := util.GetMajongPlayer(player.PlayerID, &context)
	// 非主动退出，再进入后取消托管；主动退出再进入不取消托管
	// 胡牌后没有托管，但是在客户端退出时，需要托管来自动胡牌,重新进入后把托管取消
	if !player.IsQuit() || mjPlayer.GetXpState() != server_pb.XingPaiState_normal {
		player.SetTuoguan(false, false)
	}
	player.EnterDesk(model.GetDesk())
	model.recoverGameForPlayer(eqi.PlayerID)
	d.setMjPlayerQuitDesk(eqi.PlayerID, false)
	d.playerQuitEnterDeskNtf(eqi.PlayerID, room.QuitEnterType_QET_ENTER)
	logEntry.Debugln("玩家进入")
}

func (model PlayerModel) recoverGameForPlayer(playerID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "recoverGameForPlayer",
		"playerID":  playerID,
	})

	mjContext := &d.dContext.mjContext
	bankerSeat := mjContext.GetZhuangjiaIndex()
	totalCardsNum := mjContext.GetCardTotalNum()
	gameStage := getGameStage(mjContext.GetCurState())
	gameID := gutils.GameIDServer2Client(int(mjContext.GetGameId()))
	gameDeskInfo := room.GameDeskInfo{
		GameId:      &gameID,
		GameStage:   &gameStage,
		Players:     getRecoverPlayerInfo(playerID, d),
		Dices:       mjContext.GetDices(),
		BankerSeat:  &bankerSeat,
		EastSeat:    &bankerSeat,
		TotalCards:  &totalCardsNum,
		RemainCards: proto.Uint32(uint32(len(mjContext.GetWallCards()))),
		CostTime:    proto.Uint32(getStateCostTime(d.dContext.stateTime.Unix())),
		OperatePid:  getOperatePlayerID(mjContext),
		DoorCard:    getDoorCard(mjContext),
		NeedHsz:     proto.Bool(gutils.GameHasHszState(mjContext)),
	}
	gameDeskInfo.HasZixun, gameDeskInfo.ZixunInfo = getZixunInfo(playerID, mjContext)
	gameDeskInfo.HasWenxun, gameDeskInfo.WenxunInfo = getWenxunInfo(playerID, mjContext)
	gameDeskInfo.HasQgh, gameDeskInfo.QghInfo = getQghInfo(playerID, mjContext)
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
	d.reply([]server_pb.ReplyClientMessage{
		server_pb.ReplyClientMessage{
			Players: []uint64{playerID},
			MsgId:   int32(msgid.MsgID_ROOM_RESUME_GAME_RSP),
			Msg:     rsp,
		},
	})
}





func (model PlayerModel) PlayerQuit(player *player.Player){
	player.QuitDesk(model.GetDesk())
	//d.setMjPlayerQuitDesk(eqi.PlayerID, true)
	model.setContextPlayerQuit(player,true)
	//d.playerQuitEnterDeskNtf(eqi.PlayerID, room.QuitEnterType_QET_QUIT)
	model.playerQuitEnterDeskNtf(player,room.QuitEnterType_QET_QUIT)
}

func (model PlayerModel) playerQuitEnterDeskNtf(player *player.Player, qeType room.QuitEnterType) {
	if player == nil {
		return
	}
	roomPlayer := util.TranslateToRoomPlayer(player)
	playerId := player.GetPlayerID()
	ntf := room.RoomDeskQuitEnterNtf{
		PlayerId:   &playerId,
		Type:       &qeType,
		PlayerInfo: &roomPlayer,
	}
	player.GetDesk().GetModel(models.Message).(MessageModel).BroadCastDeskMessageExcept([]uint64{playerId}, true, msgid.MsgID_ROOM_DESK_QUIT_ENTER_NTF, &ntf)
}

func (model PlayerModel) setContextPlayerQuit(player *player.Player,value bool){
	for _,p:= range model.GetDesk().GetConfig().Context.(contexts.MjContext).MjContext.Players{
		if p.GetPalyerId() == player.GetPlayerID(){
			p.IsQuit = value
		}
	}
}

func (model PlayerModel) GetDeskPlayers() []*player.Player {
	return model.players
}

// GetDeskPlayerIDs 获取牌桌玩家 ID 列表， 座号作为索引
func (model PlayerModel) GetDeskPlayerIDs() []uint64 {
	players := model.GetDeskPlayers()
	result := make([]uint64, len(players))
	for _, player := range players {
		result[player.GetSeat()] = player.GetPlayerID()
	}
	return result
}