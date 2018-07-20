package pubsub

import (
	"steve/structs/pubsub"

	"github.com/Sirupsen/logrus"
	nsq "github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

type publisher struct {
	producer *nsq.Producer
}

// Publish 发布消息
func (pub *publisher) Publish(topic string, data []byte) error {
	return pub.producer.Publish(topic, data)
}

// CreatePublisher 创建 Publisher
func CreatePublisher() pubsub.Publisher {
	nsqAddr := viper.GetString("nsqd_addr")
	return &publisher{
		producer: createNsqProducer(nsqAddr),
	}
}

// createNsqProducer 创建 nsq 的生产者
func createNsqProducer(addr string) *nsq.Producer {
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		logrus.WithError(err).Panicln("创建 NSQ 生产者失败")
	}
	if err := producer.Ping(); err != nil {
		// 暂时改为 Error，后续还原成 Panic
		logrus.WithError(err).Errorln("连接 NSQ 失败")
	}
	return producer
}
