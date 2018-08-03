package core

import (
	"context"
	"errors"
	msgid "steve/client_pb/msgid"
	"steve/gateway/auth"
	"steve/gateway/connection"
	"steve/gateway/msgrange"
	"steve/gateway/router"
	"steve/gateway/watchdog"
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

type observer struct {
}

var _ net.MessageObserver = new(observer)

var errNoCorrespondeServer = errors.New("消息没有对应的处理服务")
var errCallServiceFailed = errors.New("调用服务失败")
var errGetConnectByServerName = errors.New("根据服务名称获取连接失败")

func (o *observer) getPlayerID(clientID uint64) uint64 {
	cm := connection.GetConnectionMgr()
	connection := cm.GetConnection(clientID)
	if connection == nil {
		return 0
	}
	return connection.GetPlayerID()
}

// handle 通过 RPC 服务处理消息
func (o *observer) handle(cc *grpc.ClientConn, clientID uint64, playerID uint64, msgID uint32, body []byte) ([]*steve_proto_gaterpc.ResponseMessage, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"name":      "receiver.handle",
		"client_id": clientID,
		"player_id": playerID,
		"msg_id":    msgID,
	})
	client := steve_proto_gaterpc.NewMessageHandlerClient(cc)
	handleResult, err := client.HandleClientMessage(context.Background(), &steve_proto_gaterpc.ClientMessage{
		PlayerId: playerID,
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

// responseRPCMessage 将 RPC 服务处理消息的结果回复给客户端
func (o *observer) responseRPCMessage(clientID uint64, reqHeader *base.Header, responses []*steve_proto_gaterpc.ResponseMessage) {
	for _, response := range responses {
		rspMsgID := response.GetHeader().GetMsgId()
		o.response(clientID, reqHeader, rspMsgID, response.GetBody())
	}
}

// responseLocalMessage 回复本地消息处理器返回的结果
func (o *observer) responseLocalMessage(clientID uint64, reqHeader *base.Header, responses []exchanger.ResponseMsg) {
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

func (o *observer) response(clientID uint64, reqHeader *base.Header, rspMsgID uint32, body []byte) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":    "receiver.response",
		"rsp_msg_id":   msgid.MsgID(rspMsgID),
		"req_send_seq": reqHeader.GetSendSeq(),
		"client_id":    clientID,
	})

	header := &base.Header{
		RspSeq: proto.Uint64(reqHeader.GetSendSeq()),
		MsgId:  proto.Uint32(rspMsgID),
	}
	dog := watchdog.Get()
	if err := dog.SendPackage(clientID, header, body); err != nil {
		entry.WithError(err).Errorln("发送消息失败")
	}
}

func heartBeat(clientID uint64) {
	conn := connection.GetConnectionMgr().GetConnection(clientID)
	if conn != nil {
		conn.HeartBeat()
	}
}

// AfterSend 消息发送后的回调
func (o *observer) AfterSend(clientID uint64, header *base.Header, body []byte, err error) {
	// 只要发送消息成功，就重置心跳时间
	if err == nil {
		heartBeat(clientID)
	}
}

// OnRecv 收到消息后的处理
func (o *observer) OnRecv(clientID uint64, header *base.Header, body []byte) {
	msgID := header.GetMsgId()
	// 收到消息时，非心跳消息，也计算一次心跳
	// 也就是说，只要收到消息就重置心跳时间
	if msgID != uint32(msgid.MsgID_GATE_HEART_BEAT_REQ) {
		heartBeat(clientID)
	}

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

	logEntry.Debugln("recv client msg")
	if serverName == common.GateServiceName {
		// 发往网关服的，调用本地消息处理器
		o.callLocalHandler(clientID, playerID, header, body)
	} else {
		// 发往其他服务器的，调用远程消息处理器
		o.callRemoteHandler(clientID, playerID, header, body, serverName)
	}
}

func (o *observer) callRemoteHandler(clientID uint64, playerID uint64, reqHeader *base.Header, body []byte, serverName string) {
	msgID := reqHeader.GetMsgId()
	entry := logrus.WithFields(logrus.Fields{
		"client_id":   clientID,
		"player_id":   playerID,
		"msg_id":      msgid.MsgID(msgID),
		"router":      reqHeader.GetRoutine(),
		"server_name": serverName,
	})
	// 未绑定玩家，处理登录消息
	if playerID == 0 {
		if msgID == uint32(msgid.MsgID_LOGIN_AUTH_REQ) {
			responses := auth.HandleLoginRequest(clientID, reqHeader, body)
			if responses != nil {
				o.responseRPCMessage(clientID, reqHeader, responses)
			}
		} else {
			entry.Warningln("未绑定玩家，不能调用远程处理器")
		}
		return
	}
	cc, err := router.GetConnection(serverName, playerID, reqHeader.GetRoutine())
	if err != nil {
		entry.WithError(err).Warningln("获取服务连接失败")
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
func (o *observer) getLocalHandler(msgID uint32) *exchanger.Handler {
	exposer := structs.GetGlobalExposer()
	return exposer.Exchanger.GetHandler(msgID)
}

func (o *observer) callLocalHandler(clientID uint64, playerID uint64, reqHeader *base.Header, body []byte) {
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
		MsgId: msgID,
	}, body)
	if err != nil {
		entry.Errorln("调用消息处理器失败")
		return
	}
	o.responseLocalMessage(clientID, reqHeader, responses)
}
