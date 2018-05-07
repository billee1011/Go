package connect

import (
	"reflect"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"

	"github.com/Sirupsen/logrus"
)

// 响应消息解析列表
var metaByID = map[uint32]*MessageMeta{}

// RegisterResponseMessageMeta 注册响应消息
func RegisterResponseMessageMeta(msgID uint32, msgType reflect.Type) {
	entry := logrus.WithField("name", "client.RegisterMessageMeta")

	meta := &MessageMeta{
		Type: msgType,
		ID:   msgID,
	}
	if _, ok := metaByID[msgID]; ok {
		entry.WithFields(logrus.Fields{"msgID": msgID}).Fatalf("重复消息注册")
	}

	metaByID[msgID] = meta
}

func init() {
	RegisterResponseMessageMeta(uint32(msgid.MsgID_room_login_req), reflect.TypeOf((*room.RoomLoginReq)(nil)).Elem())
}
