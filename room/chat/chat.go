package chat

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/gogo/protobuf/proto"
)

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
func RoomChatMsgReq(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskChatReq) (ret []exchanger.ResponseMsg) {
	// 获取聊天发起者ID
	playerID := global.GetPlayerMgr().GetPlayerByClientID(clientID).GetID()
	// 日志
	logentry := logrus.WithFields(logrus.Fields{
		"func_name": "RoomChatMsgReq",
		"client_id": clientID,
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
	broadChatNotify(playerID, ntf)
	//日志信息
	logentry.WithFields(logrus.Fields{
		"chatType": strChatType,
		"chatInfo": chatInfo,
	}).Infoln("---玩家聊天---")
	return
}

// 广播聊天通知
func broadChatNotify(playerID uint64, ntf *room.RoomDeskChatNtf) {
	// 获取桌面
	desk, err := global.GetDeskMgr().GetRunDeskByPlayerID(playerID)
	printErr("---广播聊天通知：获取当前桌面失败---", playerID, err)
	// err = facade.BroadCastDeskMessageExcept(desk, []uint64{playerID}, true, msgid.MsgID_ROOM_CHAT_NTF, ntf)
	// 聊天通知序列化
	msgBody, err := proto.Marshal(ntf)
	printErr("---广播聊天通知：序列化失败---", playerID, err)
	// 广播聊天消息([]uint64{}为所有玩家，true为退出玩家不发送聊天消息)
	desk.BroadcastMessage([]uint64{}, msgid.MsgID_ROOM_CHAT_NTF, msgBody, true)
}

// 打印错误日志
func printErr(str string, playerID uint64, err error) {
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func_name": "broadChatNotify",
			"playerID":  playerID,
			"err":       err,
		}).Info(str)
	}
}
