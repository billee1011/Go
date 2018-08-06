package nsq

import (
	"steve/structs"
	"github.com/Sirupsen/logrus"
	"steve/structs/proto/gate_rpc"
	"github.com/nsqio/go-nsq"
	"github.com/golang/protobuf/proto"
	"fmt"
	"steve/gutils/topics"
	"steve/gateway/connection"
	"steve/gateway/watchdog"
	"steve/structs/proto/base"
)


type BroadcastMsgHandler struct {
}

func (plh *BroadcastMsgHandler) HandleMessage(message *nsq.Message) error {
	logrus.Debugln("recv nsq msg: ", message.Body)
	req := steve_proto_gaterpc.BroadcastMsgRequest{}
	if err := proto.Unmarshal(message.Body, &req); err != nil {
		logrus.WithError(err).Errorln("消息反序列化失败")
		return fmt.Errorf("消息反序列化失败：%v", err)
	}

	// 先向所有玩家发送广播
	nsqBroadMessage(req.GetHeader().GetMsgId(), req.GetData())

	return nil
}

// 通过NSQ消息队列向Client广播通知消息
func  nsqBroadMessage(msgID uint32, msgData []byte) error {
	connMgr := connection.GetConnectionMgr()
	clientIDs := connMgr.GetAllClientID()

	header := base.Header{
		MsgId:   proto.Uint32(msgID),
		Version: proto.String("1.0"),
	}
	if len(clientIDs) != 0 {
		dog := watchdog.Get()
		err := dog.BroadPackage(clientIDs, &header, msgData)
		if err != nil {
			return err
		}
	}
	return nil
}


func Init() {
	exposer := structs.GetGlobalExposer()
	if exposer.Subscriber == nil {
		return
	}
	if err := exposer.Subscriber.Subscribe(topics.BroadcastMsg, "gateway", &BroadcastMsgHandler{}); err != nil {
		logrus.WithError(err).Panicln("订阅登录消息失败")
	}
}