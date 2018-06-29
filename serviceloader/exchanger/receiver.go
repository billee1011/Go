package exchanger

import (
	"context"
	"errors"
	"reflect"
	iexchanger "steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

var byteSliceType = reflect.TypeOf([]byte{})

type receiver struct {
	handlerMgr iexchanger.HandlerMgr
}

// HandleClientMessage 处理客户端消息
func (r *receiver) HandleClientMessage(ctx context.Context, msg *steve_proto_gaterpc.ClientMessage) (*steve_proto_gaterpc.HandleResult, error) {
	header := msg.GetHeader()
	msgID := header.GetMsgId()
	playerID := msg.GetPlayerId()
	handler := r.handlerOf(msgID)
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "receiver.HandleClientMessage",
		"msg_id":    msgID,
		"player_id": playerID,
	})
	if handler == nil {
		logEntry.Warnln("没有对应的消息处理器")
		return &steve_proto_gaterpc.HandleResult{}, nil
	}
	responses := r.callHandler(logEntry, handler, playerID, header, msg.GetRequestData())
	return r.packResults(responses)
}

// callHandler 根据消息类型反序列化消息体和回调处理器
func (r *receiver) callHandler(logEntry *logrus.Entry, handler *iexchanger.Handler, playerID uint64,
	header *steve_proto_gaterpc.Header, body []byte) []iexchanger.ResponseMsg {

	result, _ := iexchanger.CallHandler(handler, playerID, header, body)
	return result
}

// packResults 将应答消息打包返回
func (r *receiver) packResults(responses []iexchanger.ResponseMsg) (*steve_proto_gaterpc.HandleResult, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "receiver.packResults",
	})
	resultMessages := []*steve_proto_gaterpc.ResponseMessage{}
	for _, resp := range responses {
		bodyData, err := proto.Marshal(resp.Body)
		if err != nil {
			logEntry.WithField("msg_id", resp.MsgID).Errorln("消息序列化失败")
			return nil, errors.New("消息序列化失败")
		}
		resultMessages = append(resultMessages, &steve_proto_gaterpc.ResponseMessage{
			Header: &steve_proto_gaterpc.Header{
				MsgId: resp.MsgID,
			},
			Body: bodyData,
		})
	}
	return &steve_proto_gaterpc.HandleResult{
		Responses: resultMessages,
	}, nil
}

func (r *receiver) handlerOf(msgID uint32) *iexchanger.Handler {
	return r.handlerMgr.GetHandler(msgID)
}
