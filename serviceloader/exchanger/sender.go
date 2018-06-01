package exchanger

import (
	"context"
	"errors"
	"steve/structs"
	"steve/structs/common"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

type sender struct{}

var errBodyMarshal = errors.New("消息序列化失败")
var errCallRPCFailed = errors.New("调用 RPC 服务失败")
var errGateSendFailed = errors.New("网关转发消息失败")
var errNoClient = errors.New("客户端连接不存在")

// broadcast 广播消息给客户端
// step 1. 获取所有客户端所在的网关，并按照网关地址分类
// step 2. 利用 gate_rpc 所提供的服务，给客户端发送消息
func (s *sender) broadcastBare(clientIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "sender.broadcastBare",
		"client_ids": clientIDs,
		"msg_id":     head.GetMsgId(),
	})
	gates := s.classify(clientIDs)
	for cc, clis := range gates {
		if err := s.gateBraodcast(cc, clis, head, bodyData); err != nil {
			logEntry.WithField("failed_clients", clis).WithError(err).Warningln("广播消息失败")
		}
	}
	return nil
}

// broadcast 广播消息
func (s *sender) broadcast(clientIDs []uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "sender.broadcast",
		"client_ids": clientIDs,
		"msg_id":     head.GetMsgId(),
	})
	bodyData, err := proto.Marshal(body)
	if err != nil {
		logEntry.WithError(err).Errorln(errBodyMarshal)
		return errBodyMarshal
	}
	return s.broadcastBare(clientIDs, head, bodyData)
}

// sendBare 发送消息
func (s *sender) sendBare(clientID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	cc := s.aquireClientGate(clientID)
	if cc == nil {
		return errNoClient
	}
	return s.gateBraodcast(cc, []uint64{clientID}, head, bodyData)
}

// send 发送消息
func (s *sender) send(clientID uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "sender.send",
		"client_id": clientID,
		"msg_id":    head.GetMsgId(),
	})
	bodyData, err := proto.Marshal(body)
	if err != nil {
		logEntry.WithError(err).Errorln(errBodyMarshal)
		return errBodyMarshal
	}
	return s.sendBare(clientID, head, bodyData)
}

// gateBraodcast 通过 gate 提供的 rpc 服务，向客户端广播消息
func (s *sender) gateBraodcast(cc *grpc.ClientConn, clientIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "sender.gateBraodcast",
		"client_ids": clientIDs,
		"msg_id":     head.GetMsgId(),
	})

	mc := steve_proto_gaterpc.NewMessageSenderClient(cc)
	r, err := mc.SendMessage(context.Background(), &steve_proto_gaterpc.SendMessageRequest{
		ClientId: clientIDs,
		Header:   head,
		Data:     bodyData,
	})
	if err != nil {
		logEntry.WithError(err).Errorln(errCallRPCFailed)
		return errCallRPCFailed
	}
	if !r.Ok {
		return errGateSendFailed
	}
	return nil
}

// classify 将客户端 id 按照所在网关分类
func (s *sender) classify(clientIDs []uint64) map[*grpc.ClientConn][]uint64 {
	result := map[*grpc.ClientConn][]uint64{}
	for _, clientID := range clientIDs {
		cc := s.aquireClientGate(clientID)
		if cc == nil {
			continue
		}
		if result[cc] == nil {
			result[cc] = make([]uint64, 0, len(clientIDs))
		}
		result[cc] = append(result[cc], clientID)
	}
	return result
}

// aquireClientGate 查询客户端连接所在的网关服
func (s *sender) aquireClientGate(clientID uint64) *grpc.ClientConn {
	g := structs.GetGlobalExposer()
	// TODO : 暂时先用任意网关代替
	if cc, err := g.RPCClient.GetConnectByServerName(common.GateServiceName); err == nil {
		return cc
	}
	return nil
}
