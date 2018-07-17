package public

import (
	"steve/room2/desk/models"
	"steve/room2"
)

type PlayerModel struct {
	BaseModel
	players []*room2.RoomPlayer
}
func (model PlayerModel) GetName() string{
	return models.Player
}
func (model PlayerModel) Start(){
	model.players = make([]*room2.RoomPlayer,model.GetDesk().GetConfig().Num)
}
func (model PlayerModel) Stop(){

}

func (model PlayerModel) PlayerEnter(player *room2.RoomPlayer,seat uint32){
	player.SetSeat(seat)
	player.EnterDesk(model.GetDesk())
}

func (model PlayerModel) PlayerQuit(player room2.RoomPlayer){
	player.QuitDesk(model.GetDesk())
}

func (model PlayerModel) GetDeskPlayers() []*room2.RoomPlayer{
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