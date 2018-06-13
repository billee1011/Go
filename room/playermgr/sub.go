package playermgr

import (
	"fmt"
	"steve/gutils/topics"
	"steve/room/interfaces/global"
	userpb "steve/server_pb/user"
	"steve/structs"
	"steve/structs/common"
	"steve/structs/pubsub"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	nsq "github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

// subscribeClientDisconnect 订阅客户端连接断开消息
func subscribeClientDisconnect() {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "subscribeClientDisconnect",
	})
	sub := getSubscriber()
	if err := sub.Subscribe(topics.ClientDisconnect, getChannelName(), &handler{}); err != nil {
		logEntry.WithError(err).Errorln("订阅消息失败")
	}
}

func getChannelName() string {
	addr := viper.GetString("rpc_addr")
	port := viper.GetInt("rpc_port")
	return fmt.Sprintf("%s_%s_%d", common.RoomServiceName, addr, port)
}

func getSubscriber() pubsub.Subscriber {
	exposer := structs.GetGlobalExposer()
	return exposer.Subscriber
}

type handler struct {
}

func (h *handler) HandleMessage(message *nsq.Message) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "handler.HandleMessage",
	})

	pb := userpb.ClientDisconnect{}
	if err := proto.Unmarshal(message.Body, &pb); err != nil {
		logEntry.WithError(err).Errorln("消息反序列化失败")
		return nil
	}
	playerMgr := global.GetPlayerMgr()
	playerMgr.OnClientDisconnect(pb.GetClientId())
	logEntry.WithField("client_id", pb.GetClientId()).Debugln("客户端断开连接")
	return nil
}

func init() {
	subscribeClientDisconnect()
}
