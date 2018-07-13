package public

import "steve/room2/desk/models"

type RequestModel struct {
	BaseModel
}
func (model RequestModel) GetName() string{
	return models.Request
}
func (model RequestModel) Start(){

}
func (model RequestModel) Stop(){

}
