package pubsub

import "github.com/nsqio/go-nsq"

// Publisher 消息发布
type Publisher interface {
	Publish(topic string, data []byte) error
}

// Subscriber 消息订阅
type Subscriber interface {
	Subscribe(topic string, channel string, handler nsq.Handler) error
}
