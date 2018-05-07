package rtoet

import (
	"errors"
	"steve/client_pb/msgId"
	"steve/room/interfaces/global"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

type msgTranslator func(playerID uint64, header *steve_proto_gaterpc.Header, body proto.Message) (eventID server_pb.EventID, eventContext []byte, err error)

type translator struct {
	msgTranslators map[msgid.MsgID]msgTranslator
}

var errTranslatorNotExists = errors.New("转换器不存在")

func (t *translator) Translate(playerID uint64, header *steve_proto_gaterpc.Header, body proto.Message) (eventID server_pb.EventID, eventContext []byte, err error) {
	f, ok := t.msgTranslators[msgid.MsgID(header.GetMsgId())]
	if !ok {
		err = errTranslatorNotExists
		return
	}
	return f(playerID, header, body)
}

func (t *translator) addTranslator(msgID msgid.MsgID, f msgTranslator) {
	t.msgTranslators[msgID] = f
}

func (t *translator) addTranslators() {
	// TODO 添加所有请求转事件表
	t.addTranslator(msgid.MsgID_room_huansanzhang_req, translateHuansanzhangReq)

}

func init() {
	t := &translator{
		msgTranslators: make(map[msgid.MsgID]msgTranslator, 1),
	}
	t.addTranslators()

	global.SetReqEventTranslator(t)
}
