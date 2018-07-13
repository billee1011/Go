package desk

import (
	"steve/room2/desk/models"
	"github.com/Sirupsen/logrus"
)

type Desk struct {
	uid            uint64
	gameID         int
	config *DeskConfig
	models []models.DeskModel
}

func NewDesk(uid uint64, gameId int,config *DeskConfig) Desk {
	desk := Desk{uid: uid,
		gameID: gameId,
		config:config,
	}

	return desk
}

func (desk Desk) InitModel(){
	desk.models = make([]models.DeskModel,len(desk.config.Models))
	for index,name := range desk.config.Models{
		model := models.CreateModel(name,&desk)
		if model == nil{
			logrus.Error("创建Model失败["+name+"]")
			continue
		}
		desk.models[index] = model
	}
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