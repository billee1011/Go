package pubsub

import (
	"errors"
	"steve/structs/pubsub"

	"github.com/Sirupsen/logrus"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

// Subscriber 消息订阅者
type subscriber struct {
	nsqLookupdAddrs []string
}

var errCreateNsqConsumer = errors.New("创建 NSQ 消费者失败")
var errConnectNSQLookupds = errors.New("连接 nsqlookupd 失败")

// Subscribe 订阅消息
func (sub *subscriber) Subscribe(topic string, channel string, handler nsq.Handler) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"topic":   topic,
		"channel": channel,
	})
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		logEntry.WithError(err).Errorln(errCreateNsqConsumer)
		return errCreateNsqConsumer
	}
	consumer.AddHandler(handler)
	if err := consumer.ConnectToNSQLookupds(sub.nsqLookupdAddrs); err != nil {
		logEntry.WithError(err).Errorln(errConnectNSQLookupds)
		return errConnectNSQLookupds
	}
	logEntry.Infoln("订阅消息成功")
	return nil
}

// CreateSubscriber 创建消息订阅者
func CreateSubscriber() pubsub.Subscriber {
	addrs := viper.GetStringSlice("nsqlookupd_addrs")
	if len(addrs) == 0 {
		return nil
	}
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "CreateSubscriber",
	})

	logEntry = logEntry.WithField("nsqlookupd_addrs", addrs)
	logEntry.Infoln("创建消息订阅者")
	return &subscriber{
		nsqLookupdAddrs: addrs,
	}
}
