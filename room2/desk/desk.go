package desk

import (
	"steve/room2/desk/models"
	"github.com/Sirupsen/logrus"
	"steve/room2/desk/models/public"
	"steve/room2/desk/player"
	"context"
	"steve/client_pb/room"
	"steve/client_pb/msgid"
)

type Desk struct {
	uid       uint64
	gameID    int
	config    *DeskConfig
	models    map[string]models.DeskModel
	playerIds []uint64
	Context context.Context
	Cancel    context.CancelFunc // 取消事件处理
}

func NewDesk(uid uint64, gameId int,playerIds []uint64,config *DeskConfig) Desk {
	desk := Desk{uid: uid,
		gameID: gameId,
		config:config,
		playerIds : playerIds,
	}

	return desk
}

func (desk Desk) GetPlayerIds() []uint64{
	return desk.playerIds
}

func (desk Desk) InitModel(){
	desk.models = make(map[string]models.DeskModel,len(desk.config.Models))
	for _,name := range desk.config.Models{
		model := models.CreateModel(name,&desk)
		if model == nil{
			logrus.Error("创建Model失败["+name+"]")
			continue
		}
		desk.models[name] = model
	}
}

func (desk Desk) GetPlayer(playerId uint64) *player.Player {
	players := desk.GetDeskPlayers()
	for _,player := range players{
		if player.GetPlayerID()==playerId {
			return player
		}
	}
	return nil
}

func (desk Desk) GetDeskPlayers() []*player.Player {
	players := desk.GetModel(models.Player).(public.PlayerModel).GetDeskPlayers()
	return players
}

func (desk Desk) GetDeskPlayerIDs() []uint64{
	players := desk.GetModel(models.Player).(public.PlayerModel).GetDeskPlayerIDs()
	return players
}

func (desk Desk) GetModel(name string) models.DeskModel{
	return desk.models[name]
}

func (desk Desk) GetUid() uint64 {
	return desk.uid
}

func (desk Desk) GetGameId() int {
	return desk.gameID
}

func (desk Desk) Start() {
	desk.Context, desk.Cancel = context.WithCancel(context.Background())
	for _,v := range desk.models{
		v.Start()
	}
}

func (desk Desk) Stop() {
	desk.Cancel()
	for _,v := range desk.models{
		v.Stop()
	}
	ntf := room.RoomDeskDismissNtf{}
	desk.GetModel(models.Message).(public.MessageModel).BroadCastDeskMessage(nil, msgid.MsgID_ROOM_DESK_DISMISS_NTF, &ntf, true)
}

func (desk Desk) GetConfig() *DeskConfig {
	return desk.config
}