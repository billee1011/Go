package core

import (
	"encoding/json"
	"fmt"
	"steve/back/logic"
	"steve/entity/gamelog"
	"steve/gutils/topics"
	"steve/structs"

	"github.com/Sirupsen/logrus"
	nsq "github.com/nsqio/go-nsq"
)

type gameSummaryHandler struct {
}

type gameDetailHandler struct {
}

func (plh *gameSummaryHandler) HandleMessage(message *nsq.Message) error {
	logrus.Infoln("获取Summary")
	gameSummary := gamelog.TGameSummary{}
	if err := json.Unmarshal(message.Body, &gameSummary); err != nil {
		logrus.WithError(err).Errorln("消息反序列化失败")
		return fmt.Errorf("消息反序列化失败：%v", err)
	}
	logic.SaveSummaryInfo(gameSummary)
	return nil
}

func (plh *gameDetailHandler) HandleMessage(message *nsq.Message) error {
	logrus.Infoln("获取Detail")
	gameDetail := gamelog.TGameDetail{}
	if err := json.Unmarshal(message.Body, &gameDetail); err != nil {
		logrus.WithError(err).Errorln("消息反序列化失败")
		return fmt.Errorf("消息反序列化失败：%v", err)
	}
	logic.SaveDetailInfo(gameDetail)
	return nil
}

func init() {
	exposer := structs.GetGlobalExposer()
	if err := exposer.Subscriber.Subscribe(topics.GameSummaryRecord, "room", &gameSummaryHandler{}); err != nil {
		logrus.WithError(err).Panicln("订阅单局总消息失败")
	}
	if err := exposer.Subscriber.Subscribe(topics.GameDetailRecord, "room", &gameDetailHandler{}); err != nil {
		logrus.WithError(err).Panicln("订阅单局玩家明细消息失败")
	}
}
