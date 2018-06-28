package core

import (
	"context"
	"errors"
	msgid "steve/client_pb/msgId"
	"steve/gateway/global"
	"steve/gateway/msgrange"
	"steve/structs"
	"steve/structs/common"
	"steve/structs/exchanger"
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

var errNoCorrespondeServer = errors.New("消息没有对应的处理服务")
var errCallServiceFailed = errors.New("调用服务失败")
var errGetConnectByServerName = errors.New("根据服务名称获取连接失败")

// getConnection 根据服务名称和客户端 ID 获取处理服务器的 RPC 连接
func (o *receiver) getConnection(serverName string, clientID uint64) (*grpc.ClientConn, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":   "receiver.getConnection",
		"client_id":   clientID,
		"server_name": serverName,
	})

	logEntry = logEntry.WithField("server_name", serverName)
	e := structs.GetGlobalExposer()
	// TODO 处理服务绑定
	cc, err := e.RPCClient.GetConnectByServerName(serverName)
	if cc == nil {
		logEntry.WithError(err).Errorln(errGetConnectByServerName)
		return nil, errGetConnectByServerName
	}
	return cc, nil
}

func (o *receiver) getPlayerID(clientID uint64) uint64 {
	cm := global.GetConnectionManager()
	connection := cm.GetConnection(clientID)
	if connection == nil {
		return 0
	}
	return connection.GetPlayerID()
}

// handle 通过 RPC 服务处理消息
func (o *receiver) handle(cc *grpc.ClientConn, clientID uint64, playerID uint64, msgID uint32, body []byte) ([]*steve_proto_gaterpc.ResponseMessage, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"name":      "receiver.handle",
		"client_id": clientID,
		"player_id": playerID,
		"msg_id":    msgID,
	})
	client := steve_proto_gaterpc.NewMessageHandlerClient(cc)
	handleResult, err := client.HandleClientMessage(context.Background(), &steve_proto_gaterpc.ClientMessage{
		ClientId: clientID,
		Header: &steve_proto_gaterpc.Header{
			MsgId:    msgID,
			PlayerId: playerID,
		},
		RequestData: body,
	})
	if err != nil {
		logEntry.WithError(err).Error(errCallServiceFailed)
		return nil, errCallServiceFailed
	}
	return handleResult.GetResponses(), nil
}

// responseRPCMessage 将 RPC 服务处理消息的结果回复给客户端
func (o *receiver) responseRPCMessage(clientID uint64, reqHeader *steve_proto_base.Header, responses []*steve_proto_gaterpc.ResponseMessage) {
	for _, response := range responses {
		rspMsgID := response.GetHeader().GetMsgId()
		o.response(clientID, reqHeader, rspMsgID, response.GetBody())
	}
}

// responseLocalMessage 回复本地消息处理器返回的结果
func (o *receiver) responseLocalMessage(clientID uint64, reqHeader *steve_proto_base.Header, responses []exchanger.ResponseMsg) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":    "receiver.responseLocalMessage",
		"req_send_seq": reqHeader.GetSendSeq(),
		"client_id":    clientID,
	})
	for _, response := range responses {
		body, err := proto.Marshal(response.Body)
		if err != nil {
			entry.WithField("msg_id", response.MsgID).WithError(err).Errorln("消息序列化失败")
			continue
		}
		o.response(clientID, reqHeader, response.MsgID, body)
	}
}

func (o *receiver) response(clientID uint64, reqHeader *steve_proto_base.Header, rspMsgID uint32, body []byte) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":    "receiver.response",
		"rsp_msg_id":   msgid.MsgID(rspMsgID),
		"req_send_seq": reqHeader.GetSendSeq(),
		"client_id":    clientID,
	})

	header := &steve_proto_base.Header{
		RspSeq: proto.Uint64(reqHeader.GetSendSeq()),
		MsgId:  proto.Uint32(rspMsgID),
	}
	dog := o.core.dog
	if err := dog.SendPackage(clientID, header, body); err != nil {
		entry.WithError(err).Errorln("发送消息失败")
	}
}

// OnRecv 收到消息后的处理
func (o *receiver) OnRecv(clientID uint64, header *steve_proto_base.Header, body []byte) {
	msgID := header.GetMsgId()

	playerID := o.getPlayerID(clientID)
	logEntry := logrus.WithFields(logrus.Fields{
		"name":      "receiver.OnRecv",
		"client_id": clientID,
		"msg_id":    msgid.MsgID(msgID),
		"player_id": playerID,
	})
	serverName := msgrange.GetMessageServer(msgID)
	if serverName == "" {
		logEntry.Errorln(errNoCorrespondeServer)
		return
	}
	logEntry = logEntry.WithField("server_name", serverName)

	if serverName == common.GateServiceName {
		// 发往网关服的，调用本地消息处理器
		o.callLocalHandler(clientID, playerID, header, body)
	} else {
		// 发往其他服务器的，调用远程消息处理器
		o.callRemoteHandler(clientID, playerID, header, body, serverName)
	}
}

func (o *receiver) callRemoteHandler(clientID uint64, playerID uint64, reqHeader *steve_proto_base.Header, body []byte, serverName string) {
	msgID := reqHeader.GetMsgId()
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "receiver.callRemoteHandler",
		"client_id": clientID,
		"player_id": playerID,
		"msg_id":    msgid.MsgID(msgID),
	})
	cc, err := o.getConnection(serverName, clientID)
	if err != nil {
		return
	}
	responses, err := o.handle(cc, clientID, playerID, msgID, body)
	if err != nil {
		entry.WithError(err).Errorln("处理消息失败")
		return
	}
	o.responseRPCMessage(clientID, reqHeader, responses)
}

// getLocalHandler 获取本地消息处理器
func (o *receiver) getLocalHandler(msgID uint32) *exchanger.Handler {
	exposer := structs.GetGlobalExposer()
	return exposer.Exchanger.GetHandler(msgID)
}

func (o *receiver) callLocalHandler(clientID uint64, playerID uint64, reqHeader *steve_proto_base.Header, body []byte) {
	msgID := reqHeader.GetMsgId()
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "receiver.callLocalHandler",
		"client_id": clientID,
		"player_id": playerID,
		"msg_id":    msgID,
	})
	handler := o.getLocalHandler(msgID)
	if handler == nil {
		entry.Infoln("不存在对应的消息处理器")
		return
	}
	responses, err := exchanger.CallHandler(handler, clientID, &steve_proto_gaterpc.Header{
		MsgId:    msgID,
		PlayerId: playerID,
	}, body)
	if err != nil {
		entry.Errorln("调用消息处理器失败")
		return
	}
	o.responseLocalMessage(clientID, reqHeader, responses)
}
