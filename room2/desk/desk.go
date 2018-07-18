package desk

import (
	"steve/room2/desk/models"
	"github.com/Sirupsen/logrus"
	"steve/room2/desk/models/public"
	"steve/room2/desk/player"
)

type Desk struct {
	uid            uint64
	gameID         int
	config *DeskConfig
	models map[string]models.DeskModel
}

func NewDesk(uid uint64, gameId int,config *DeskConfig) Desk {
	desk := Desk{uid: uid,
		gameID: gameId,
		config:config,
	}

	return desk
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
	for _,v := range desk.models{
		v.Start()
	}
}

func (desk Desk) Stop() {
	for _,v := range desk.models{
		v.Stop()
	}
}

func (desk Desk) GetConfig() *DeskConfig {
	return desk.config
}