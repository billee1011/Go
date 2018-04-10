package core

import (
	"context"
	"steve/structs/exchanger"
	"steve/structs/net"
	"steve/structs/proto/base"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
)

type messageObserver struct {
	core *gatewayCore
}

var _ net.MessageObserver = new(messageObserver)

func (o *messageObserver) OnRecv(clientID uint64, header *steve_proto_base.Header, body []byte) {
	msgID := header.GetMsgId()

	logEntry := logrus.WithFields(logrus.Fields{
		"name":      "messageObserver.OnRecv",
		"client_id": clientID,
		"msg_id":    msgID,
	})
	logEntry.Debug("收到客户端消息")

	handleServer := exchanger.GetMessageServer(msgID)
	if handleServer == "" {
		logEntry.Error("消息没有对应的处理服务")
		return
	}
	logEntry = logEntry.WithField("server_name", handleServer)

	// TODO 处理服务绑定
	cc, err := o.core.e.RPCClient.GetClientConnByServerName(handleServer)
	if err != nil {
		logEntry.WithError(err).Error("获取服务失败")
	}

	client := steve_proto_gaterpc.NewMessageHandlerClient(cc)
	handleResult, err := client.HandleClientMessage(context.Background(), &steve_proto_gaterpc.ClientMessage{
		ClientId: clientID,
		Header: &steve_proto_gaterpc.Header{
			MsgId: header.GetMsgId(),
		},
		RequestData: body,
	})
	if err != nil {
		logEntry.WithError(err).Error("调用 HandleClientMessage 失败")
		return
	}
	respDatas := handleResult.GetResponseDatas()
	header.RspSeq = header.SendSeq

	for _, rspData := range respDatas {
		if err := o.core.dog.SendPackage(clientID, header, rspData.GetValue()); err != nil {
			logEntry.WithError(err).Error("回复数据失败")
			return
		}
	}
}
