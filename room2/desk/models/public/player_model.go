package public

import "steve/room2/desk/models"

type PlayerModel struct {
	BaseModel

}
func (model PlayerModel) GetName() string{
	return models.Player
}
func (model PlayerModel) Start(){

}
func (model PlayerModel) Stop(){

}