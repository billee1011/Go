package public

import (
	"steve/room2/desk/models"
	"github.com/Sirupsen/logrus"
	a "steve/room/interfaces"
	b "steve/room/interfaces/global"
	"github.com/golang/protobuf/proto"
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
	desk.GetModel(models.Event).(mj.MjEventModel).
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

// PushRequest 压入玩家请求
func (d *desk) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "desk.PushRequest",
		"desk_uid":   d.GetUID(),
		"game_id":    d.GetGameID(),
		"player_id":  playerID,
		"message_id": head.GetMsgId(),
	})

	trans := global.GetReqEventTranslator()
	eventID, eventContext, err := trans.Translate(playerID, head, bodyData)
	if err != nil {
		logEntry.WithError(err).Errorln("消息转事件失败")
		return
	}
	eventMessage, ok := eventContext.(proto.Message)
	if !ok {
		logEntry.Errorln("转换事件函数返回值类型错误")
		return
	}
	eventConetxtByte, err := proto.Marshal(eventMessage)
	if err != nil {
		logEntry.WithError(err).Errorln("序列化事件现场失败")
	}

	d.PushEvent(interfaces.Event{
		ID:        server_pb.EventID(eventID),
		Context:   eventConetxtByte,
		EventType: interfaces.NormalEvent,
		PlayerID:  playerID,
	})
}