package public

import (
	"steve/room2/desk/models"
	"github.com/Sirupsen/logrus"
	"steve/structs/proto/gate_rpc"
	"steve/room2/desk/models/mj"
	player2 "steve/room2/desk/player"
)

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

// HandlePlayerRequest 处理玩家请求
func (model RequestModel) HandlePlayerRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "deskMgr.HandlePlayerRequest",
		"player_id":  playerID,
		"message_id": head.GetMsgId(),
	})

	//iDeskID, ok := dm.playerDeskMap.Load(playerID)
	player := player2.GetRoomPlayerMgr().GetPlayer(playerID)
	desk := player.GetDesk()
	if !(desk==nil) {
		logEntry.Infoln("玩家不在牌桌上")
		return
	}
	desk.GetModel(models.Event).(mj.MjEventModel).PushRequest(playerID,head,bodyData) //TODO 临时
	/*deskID := iDeskID.(uint64)
	logEntry = logEntry.WithField("desk_id", deskID)

	iDesk, ok := dm.deskMap.Load(deskID)
	if !ok {
		logEntry.Infoln("牌桌可能已经结束")
		return
	}
	desk := iDesk.(interfaces.Desk)
	desk.PushRequest(playerID, head, bodyData)*/

}