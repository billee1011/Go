package models

import (
	"errors"
	"reflect"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/room2/player"
	"steve/room2/util"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type msgTranslator struct {
	f        interface{}
	bodyType reflect.Type
}

type translator struct {
	msgTranslators map[msgid.MsgID]msgTranslator
}

var errTranslatorNotExists = errors.New("转换器不存在")
var errUnmarshalReqFailed = errors.New("反序列化请求消息体失败")

func (t *translator) Translate(playerID uint64, header *steve_proto_gaterpc.Header, bodyData []byte) (eventID int, eventContext interface{}, err error) {
	f, ok := t.msgTranslators[msgid.MsgID(header.GetMsgId())]
	if !ok {
		err = errTranslatorNotExists
		return
	}
	return t.callTranslator(f, playerID, header, bodyData)
}

func (t *translator) callTranslator(msgTranslator msgTranslator, playerID uint64,
	header *steve_proto_gaterpc.Header, bodyData []byte) (eventID int, eventContext interface{}, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "translator.callTranslator",
		"player_id": playerID,
		"msg_id":    header.GetMsgId(),
	})

	bodyMsg := reflect.New(msgTranslator.bodyType).Interface()
	if err = proto.Unmarshal(bodyData, bodyMsg.(proto.Message)); err != nil {
		logEntry.WithError(err).Errorln(errUnmarshalReqFailed)
		err = errUnmarshalReqFailed
		return
	}

	f := reflect.ValueOf(msgTranslator.f)
	callResults := f.Call([]reflect.Value{
		reflect.ValueOf(playerID),
		reflect.ValueOf(header),
		reflect.ValueOf(bodyMsg).Elem(),
	})

	eventID = callResults[0].Interface().(int)
	eventContext = callResults[1].Interface()
	errInterface := callResults[2].Interface()
	if errInterface == nil {
		err = nil
	} else {
		err = errInterface.(error)
	}
	return
}

func (t *translator) addTranslator(msgID msgid.MsgID, f interface{}) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "translator.addTranslator",
		"msg_id":    msgID,
	})

	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		logEntry.Panic("需要函数类型")
	}
	if fType.NumIn() != 3 || fType.In(0).Kind() != reflect.Uint64 ||
		fType.In(1) != reflect.TypeOf(&steve_proto_gaterpc.Header{}) {
		logEntry.Panic("处理函数的参数错误")
	}
	bodyType := fType.In(2)
	msg := reflect.New(bodyType)
	if _, ok := msg.Interface().(proto.Message); !ok {
		logEntry.Panic("处理函数的第 3 个参数必须是 proto.Message 类型")
	}

	typeOfErr := reflect.TypeOf(errors.New(""))
	if fType.NumOut() != 3 || fType.Out(0).Kind() != reflect.Int || fType.Out(2).Name() != "error" {
		logEntry.WithFields(logrus.Fields{
			"num_out":     fType.NumOut(),
			"out_0_type":  fType.Out(0),
			"out_2_type":  fType.Out(2),
			"type_of_err": typeOfErr,
		}).Panic("处理函数的返回值类型错误")
	}
	t.msgTranslators[msgID] = msgTranslator{
		f:        f,
		bodyType: bodyType,
	}
}

func (t *translator) addTranslators() {
	// majong
	t.addTranslator(msgid.MsgID_ROOM_HUANSANZHANG_REQ, util.TranslateHuansanzhangReq)
	t.addTranslator(msgid.MsgID_ROOM_XINGPAI_ACTION_REQ, util.TranslateXingpaiActionReq)
	t.addTranslator(msgid.MsgID_ROOM_CHUPAI_REQ, util.TranslateChupaiReq)
	t.addTranslator(msgid.MsgID_ROOM_DINGQUE_REQ, util.TranslateDingqueReq)
	t.addTranslator(msgid.MsgID_ROOM_CARTOON_FINISH_REQ, util.TranslateCartoonFinishReq)

	// 斗地主
	t.addTranslator(msgid.MsgID_ROOM_DDZ_GRAB_LORD_REQ, util.TranslateGrabRequest)
	t.addTranslator(msgid.MsgID_ROOM_DDZ_DOUBLE_REQ, util.TranslateDoubleRequest)
	t.addTranslator(msgid.MsgID_ROOM_DDZ_PLAY_CARD_REQ, util.TranslatePlayCardRequest)
	t.addTranslator(msgid.MsgID_ROOM_DDZ_RESUME_REQ, util.TranslateResumeRequest)
}

var tran *translator

func GetTranslator() *translator {
	return tran
}

func init() {
	t := &translator{
		msgTranslators: make(map[msgid.MsgID]msgTranslator, 1),
	}
	t.addTranslators()
	tran = t
}

// TranslateToRoomPlayer 将 deskPlayer 转换成 RoomPlayerInfo
func translateToRoomPlayer(player *player.Player) room.RoomPlayerInfo {
	coin := player.GetCoin()
	return room.RoomPlayerInfo{
		PlayerId: proto.Uint64(player.GetPlayerID()),
		Name:     proto.String(""), // TODO
		Coin:     proto.Uint64(coin),
		Seat:     proto.Uint32(uint32(player.GetSeat())),
		// Location: TODO 没地方拿
	}
}
