package core

import (
	"context"
	"errors"
	msgid "steve/client_pb/msgId"
	"steve/gateway/msgrange"
	"steve/structs"
	"steve/structs/net"
	"steve/structs/proto/base"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

type receiver struct {
	core *gatewayCore
}

var _ net.MessageObserver = new(receiver)

var errNoCorresponseServer = errors.New("消息没有对应的处理服务")
var errCallServiceFailed = errors.New("调用服务失败")
var errGetConnectByServerName = errors.New("根据服务名称获取连接失败")

// getConnection 根据消息 ID 和客户端 ID 获取处理服务器的 RPC 连接
func (o *receiver) getConnection(msgID uint32, clientID uint64) (*grpc.ClientConn, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"name":      "receiver.getConnection",
		"client_id": clientID,
		"msg_id":    msgID,
	})
	server := msgrange.GetMessageServer(msgID)
	if server == "" {
		logEntry.Error(errNoCorresponseServer)
		return nil, errNoCorresponseServer
	}
	logEntry = logEntry.WithField("server_name", server)
	e := structs.GetGlobalExposer()
	// TODO 处理服务绑定
	cc, err := e.RPCClient.GetConnectByServerName(server)
	if cc == nil {
		logEntry.WithError(err).Errorln(errGetConnectByServerName)
		return nil, errGetConnectByServerName
	}
	return cc, nil
}

// handle 通过 RPC 服务处理消息
func (o *receiver) handle(cc *grpc.ClientConn, clientID uint64, msgID uint32, body []byte) ([]*steve_proto_gaterpc.ResponseMessage, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"name":      "receiver.handle",
		"client_id": clientID,
		"msg_id":    msgID,
	})
	client := steve_proto_gaterpc.NewMessageHandlerClient(cc)
	handleResult, err := client.HandleClientMessage(context.Background(), &steve_proto_gaterpc.ClientMessage{
		ClientId: clientID,
		Header: &steve_proto_gaterpc.Header{
			MsgId: msgID,
		},
		RequestData: body,
	})
	if err != nil {
		logEntry.WithError(err).Error(errCallServiceFailed)
		return nil, errCallServiceFailed
	}
	return handleResult.GetResponses(), nil
}

// response 将 RPC 服务处理消息的结果回复给客户端
func (o *receiver) response(clientID uint64, reqHeader *steve_proto_base.Header, responses []*steve_proto_gaterpc.ResponseMessage) {
	logEntry := logrus.WithFields(logrus.Fields{
		"name":       "receiver.response",
		"client_id":  clientID,
		"req_msg_id": msgid.MsgID(reqHeader.GetMsgId()),
		"send_seq":   reqHeader.GetSendSeq(),
	})
	dog := o.core.dog
	for _, response := range responses {
		rspMsgID := response.GetHeader().GetMsgId()
		header := &steve_proto_base.Header{
			RspSeq: proto.Uint64(reqHeader.GetSendSeq()),
			MsgId:  proto.Uint32(rspMsgID),
		}
		newLogEntry := logEntry.WithFields(logrus.Fields{
			"rsp_msg_id": msgid.MsgID(rspMsgID),
			"rsp_seq":    header.GetRspSeq(),
		})
		if err := dog.SendPackage(clientID, header, response.GetBody()); err != nil {
			newLogEntry.WithError(err).Errorln("发送消息失败")
		}
	}
}

// OnRecv 收到消息后的处理
func (o *receiver) OnRecv(clientID uint64, header *steve_proto_base.Header, body []byte) {
	msgID := header.GetMsgId()

	logEntry := logrus.WithFields(logrus.Fields{
		"name":      "receiver.OnRecv",
		"client_id": clientID,
		"msg_id":    msgid.MsgID(msgID),
	})
	cc, err := o.getConnection(msgID, clientID)
	if err != nil {
		return
	}
	responses, err := o.handle(cc, clientID, msgID, body)
	if err != nil {
		logEntry.WithError(err).Errorln("处理消息失败")
		return
	}
	o.response(clientID, header, responses)
}
