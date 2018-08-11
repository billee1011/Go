package exchanger

import (
	"context"
	"errors"
	"steve/external/hallclient"
	"steve/structs"
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

// broadcastBare 广播消息给玩家
// step 1. 获取所有玩家所在的网关，并按照网关地址分类
// step 2. 利用 gate_rpc 所提供的服务，给玩家发送消息
func (s *sender) broadcastBare(playerIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "sender.broadcastBare",
		"player_ids": playerIDs,
		"msg_id":     head.GetMsgId(),
	})
	gates := s.classifyPlayers(playerIDs)
	for cc, clis := range gates {
		if err := s.gateBraodcast(cc, clis, head, bodyData); err != nil {
			logEntry.WithField("failed_clients", clis).WithError(err).Warningln("广播消息失败")
		}
	}
	return nil
}

// broadcast 广播消息
func (s *sender) broadcast(playerIDs []uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "sender.broadcast",
		"player_ids": playerIDs,
		"msg_id":     head.GetMsgId(),
	})
	bodyData, err := proto.Marshal(body)
	if err != nil {
		logEntry.WithError(err).Errorln(errBodyMarshal)
		return errBodyMarshal
	}
	return s.broadcastBare(playerIDs, head, bodyData)
}

// sendBare 发送消息
func (s *sender) sendBare(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	cc := s.aquirePlayerGate(playerID)
	if cc == nil {
		return errNoClient
	}
	return s.gateBraodcast(cc, []uint64{playerID}, head, bodyData)
}

// send 向玩家发送消息
func (s *sender) send(playerID uint64, head *steve_proto_gaterpc.Header, body proto.Message) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "sender.send",
		"player_id": playerID,
		"msg_id":    head.GetMsgId(),
	})
	bodyData, err := proto.Marshal(body)
	if err != nil {
		logEntry.WithError(err).Errorln(errBodyMarshal)
		return errBodyMarshal
	}
	return s.sendBare(playerID, head, bodyData)
}

// gateBraodcast 通过 gate 提供的 rpc 服务，向玩家广播消息
func (s *sender) gateBraodcast(cc *grpc.ClientConn, playerIDs []uint64, head *steve_proto_gaterpc.Header, bodyData []byte) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "sender.gateBraodcast",
		"player_ids": playerIDs,
		"msg_id":     head.GetMsgId(),
	})

	mc := steve_proto_gaterpc.NewMessageSenderClient(cc)
	r, err := mc.SendMessage(context.Background(), &steve_proto_gaterpc.SendMessageRequest{
		PlayerId: playerIDs,
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

// classifyPlayers 将玩家 id 按照所在网关分类
func (s *sender) classifyPlayers(playerIDs []uint64) map[*grpc.ClientConn][]uint64 {
	result := map[*grpc.ClientConn][]uint64{}
	for _, playerID := range playerIDs {
		cc := s.aquirePlayerGate(playerID)
		if cc == nil {
			continue
		}
		if result[cc] == nil {
			result[cc] = make([]uint64, 0, len(playerIDs))
		}
		result[cc] = append(result[cc], playerID)
	}
	return result
}

// aquireClientGate 查询玩家所在的网关服
func (s *sender) aquirePlayerGate(playerID uint64) *grpc.ClientConn {
	entry := logrus.WithFields(logrus.Fields{
		"player_id": playerID,
	})

	g := structs.GetGlobalExposer()
	gateAddr, _ := hallclient.GetGateAddr(playerID)
	if gateAddr == "" {
		// entry.Debugln("玩家不在线")
		return nil
	}
	cc, err := g.RPCClient.GetConnectByAddr(gateAddr)
	if err != nil {
		entry.WithError(err).Infoln("获取网关连接失败")
		return nil
	}
	return cc
}
