package models

import (
	"steve/room/desk"
	"steve/room/fixed"
	player2 "steve/room/player"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
)

type RequestModel struct {
	BaseModel
}

func (model *RequestModel) GetName() string {
	return fixed.RequestModelName
}

// Active 激活 model
func (model *RequestModel) Active() {}

func (model *RequestModel) Start() {

}
func (model *RequestModel) Stop() {

}

func NewRequestModel(desk *desk.Desk) DeskModel {
	result := &RequestModel{}
	result.SetDesk(desk)
	return result
}

// HandlePlayerRequest 处理玩家请求
func (model *RequestModel) HandlePlayerRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"player_id":  playerID,
		"message_id": head.GetMsgId(),
	})

	player := player2.GetPlayerMgr().GetPlayer(playerID)
	desk := player.GetDesk()
	if desk != model.GetDesk() {
		logEntry.Infoln("玩家不在牌桌上")
		return
	}
	eventModel := GetEventModel(desk.GetUid())
	if eventModel == nil {
		logEntry.Errorln("获取 event model 失败")
		return
	}
	eventModel.PushRequest(playerID, head, bodyData)
}
