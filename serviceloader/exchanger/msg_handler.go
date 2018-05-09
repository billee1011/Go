package exchanger

import (
	"context"
	"fmt"
	"reflect"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type msgHandler struct {
	exchanger exchangerImpl
}

func (h *msgHandler) HandleClientMessage(ctx context.Context, msg *steve_proto_gaterpc.ClientMessage) (*steve_proto_gaterpc.HandleResult, error) {
	logEntry := logrus.WithField("name", "msgHandler.HandleClientMessage")

	header := msg.GetHeader()
	msgID := header.GetMsgId()
	handler := h.exchanger.getHandler(msgID)

	logEntry = logEntry.WithFields(logrus.Fields{
		"msg_id":    msgID,
		"client_id": msg.GetClientId(),
	})

	f := reflect.ValueOf(handler.handleFunc)
	bodyMsg := reflect.New(handler.msgType).Interface()
	if err := proto.Unmarshal(msg.GetRequestData(), bodyMsg.(proto.Message)); err != nil {
		logEntry.WithError(err).Errorln("反序列化消息体失败")
		return nil, fmt.Errorf("反序列化消息体失败 %v", err)
	}

	results := f.Call([]reflect.Value{
		reflect.ValueOf(msg.GetClientId()),
		reflect.ValueOf(header),
		reflect.ValueOf(bodyMsg).Elem(),
	})
	_ = results[0]
	handleResult := &steve_proto_gaterpc.HandleResult{}

	// TODO : 回复
	// retMessages, _ := result.Interface().([]proto.Message)
	// for _, retMessage := range retMessages {
	// 	anyMsg, err := ptypes.MarshalAny(retMessage)
	// 	if err != nil {
	// 		logEntry.WithError(err).Errorln("消息反序列化失败")
	// 		return handleResult, fmt.Errorf("消息序列化成 Any 失败: %v", err)
	// 	}
	// 	handleResult.ResponseDatas = append(handleResult.ResponseDatas, anyMsg)
	// }

	return handleResult, nil
}

// NewMessageHandlerServer 创建消息处理服务
// 	返回值 steve_proto_gaterpc.MessageHandlerServer 为消息处理服务
//  返回值 exchanger.Exchanger 为客户端交互器
func NewMessageHandlerServer() (steve_proto_gaterpc.MessageHandlerServer, exchanger.Exchanger) {
	s := &msgHandler{
		exchanger: exchangerImpl{},
	}
	return s, &s.exchanger
}
