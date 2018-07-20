package public

import (
	"steve/room2/desk/models"
	"github.com/Sirupsen/logrus"
	"steve/client_pb/room"
	"github.com/golang/protobuf/proto"
	"steve/client_pb/msgid"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"steve/room2/desk/player"
)

type ChatModel struct {
	BaseModel
}

func (model ChatModel) GetName() string {
	return models.Chat
}
func (model ChatModel) Start() {

}
func (model ChatModel) Stop() {

}

// 所有聊天类型
var chatTypeAll map[room.ChatType]string

func init() {
	chatTypeAll = map[room.ChatType]string{
		room.ChatType_CT_EXPRESSION: "表情",
		room.ChatType_CT_QUICK:      "快捷语",
		room.ChatType_CT_VOICE:      "语音",
		room.ChatType_CT_WRITE:      "打字",
	}
}

//RoomChatMsgReq 房间玩家的聊天信息请求
func (model ChatModel) RoomChatMsgReq(player *player.Player, header *steve_proto_gaterpc.Header, req room.RoomDeskChatReq) (ret []exchanger.ResponseMsg) {
	playerID := player.GetPlayerID()
	// 日志
	logentry := logrus.WithFields(logrus.Fields{
		"func_name": "RoomChatMsgReq",
		"sourceID":  playerID,
	})
	// 聊天类型
	chatType := req.GetChatType()
	// 聊天信息
	chatInfo := req.GetChatInfo()
	// 响应消息
	ntf := &room.RoomDeskChatNtf{
		PlayerId: &playerID,
		ChatType: &chatType,
		ChatInfo: &chatInfo,
	}
	// 聊天类型是否存在
	strChatType, isExist := chatTypeAll[chatType]
	if !isExist {
		logentry.WithFields(logrus.Fields{
			"chatType": chatType,
		}).Infoln("---玩家聊天：聊天类型不存在---")
		return
	}
	// 广播聊天通知
	err := model.broadChatNotify(player, ntf)
	if err != nil {
		logentry.WithFields(logrus.Fields{
			"err": err,
		}).Infoln("---广播聊天失败---")
	}
	//日志信息
	logentry.WithFields(logrus.Fields{
		"chatType": strChatType,
		"chatInfo": chatInfo,
	}).Infoln("---玩家聊天---")
	return
}

// 广播聊天通知
func (model ChatModel) broadChatNotify(player *player.Player, ntf *room.RoomDeskChatNtf) error {
	// 聊天通知序列化
	msgBody, err := proto.Marshal(ntf)
	if err != nil {
		return err
	}
	// 广播聊天消息([]uint64{}为所有玩家，true为退出玩家不发送聊天消息)
	player.GetDesk().GetModel(models.Message).(MessageModel).BroadcastMessage([]uint64{}, msgid.MsgID_ROOM_CHAT_NTF, msgBody, true)
	return nil
}