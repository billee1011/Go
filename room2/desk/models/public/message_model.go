package public

import (
	"steve/room/interfaces"
	"steve/room2/desk/models"
)

var gMessageSender interfaces.MessageSender

type MessageModel struct {
	BaseModel
}

func (model MessageModel) GetName() string {
	return models.Message
}
func (model MessageModel) Start() {

}
func (model MessageModel) Stop() {

}
