package server

import (
	"fmt"
	"steve/alms/data"
	client_alms "steve/client_pb/alms"
	"steve/client_pb/msgid"
	"steve/entity/constant"
	"steve/external/gateclient"
	"steve/external/hallclient"
	"steve/server_pb/alms"
	"steve/server_pb/user"
	"steve/structs"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	exposer := structs.GetGlobalExposer()
	if err := exposer.Subscriber.Subscribe(constant.PlayerLogin, "alms", &playerLoginHandler{}); err != nil {
		logrus.WithError(err).Panicln("订阅登录消息失败")
	}
	if err := exposer.Subscriber.Subscribe(constant.UpdateConfig, "alms", &almsConfigChangeHandler{}); err != nil {
		logrus.WithError(err).Panicln("订阅救济金改变配置失败")
	}
}

// 玩家登陆订阅救济金配置
type playerLoginHandler struct {
}

func (plh *playerLoginHandler) HandleMessage(message *nsq.Message) error {
	logrus.Debugln("玩家登陆订阅救济金配置")
	loginPb := user.PlayerLogin{}
	if err := proto.Unmarshal(message.Body, &loginPb); err != nil {
		logrus.WithError(err).Errorln("消息反序列化失败")
		return fmt.Errorf("消息反序列化失败：%v", err)
	}
	if err := getPlayerAlmsConfigInfo(loginPb.PlayerId); err != nil {
		logrus.WithError(err).Panicln("发送玩家登陆订阅救济金配置失败")
	}
	return nil
}

// 获取玩家救济配置信息,并发送请求
func getPlayerAlmsConfigInfo(playerID uint64) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "getPlayerAlmsConfigInfo",
		"playerID":  playerID,
	})
	// 判断玩家是否在线
	wgadder, err := hallclient.GetGateAddr(playerID)
	if err != nil || wgadder == "" {
		logrus.WithError(err).Errorln("判断玩家是否在线")
		return err
	}
	ac, err := data.GetAlmsConfigByPlayerID(playerID)
	if ac == nil {
		entry.WithError(err).Errorf("根据玩家ID获取救济金配置失败 playerID(%v)", playerID)
		return err
	}
	almsConfig := &client_alms.AlmsConfig{
		AlmsGetNorm:      proto.Int64(ac.GetNorm),                             // 救济线
		AlmsGetTimes:     proto.Int32(int32(ac.GetTimes)),                     // 救济次数
		AlmsGetNumber:    proto.Int64(ac.GetNumber),                           // 领取数量
		AlmsCountDonw:    proto.Int32(int32(ac.AlmsCountDonw)),                // 救济倒计时
		DepositCountDonw: proto.Int32(int32(ac.DepositCountDonw)),             // 快冲倒计时
		GameLeveIsOpen:   dataToClentPbGameLeveIsOpen(ac.GemeLeveIsOpentAlms), // 游戏场次对应的是否开启
		Version:          proto.Int32(int32(ac.Version)),                      // 版本
	}
	ntf := &client_alms.AlmsConfigNtf{
		AlmsConfig:     almsConfig,
		PlayerGotTimes: proto.Int32(int32(ac.PlayerGotTimes)), // 玩家已经领取的数量
	}
	entry.WithFields(logrus.Fields{
		"ntf": ntf,
	}).Infoln("玩家救济金配置信息")
	return gateclient.SendPackageByPlayerID(playerID, uint32(msgid.MsgID_ALMS_LOGIN_GOLD_CONFIG_NTF), ntf)
}

// 救济配置变化订阅
type almsConfigChangeHandler struct {
}

func (acch *almsConfigChangeHandler) HandleMessage(message *nsq.Message) error {
	logrus.Debugln("救济配置变化订阅")
	almsConfig := alms.AlmsConfig{}
	if err := proto.Unmarshal(message.Body, &almsConfig); err != nil {
		logrus.WithError(err).Errorln("消息反序列化失败")
		return fmt.Errorf("消息反序列化失败：%v", err)
	}
	almsConfigChange(almsConfig) // 修改redis
	return nil
}

//订阅消息配置改变,只改变redis,场次是否开启救济，在游戏场次配置
func almsConfigChange(a alms.AlmsConfig) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":       "almsConfigChange",
		"alms.AlmsConfig": a,
	})
	changeConfig := make(map[string]interface{})
	norm := a.GetAlmsGetNorm()
	if norm > 0 { // 救济线
		changeConfig[data.AlmsGetNorm] = norm
	}
	times := a.GetAlmsGetTimes()
	if times > 0 { // 救济金领取次数
		changeConfig[data.AlmsGetTimes] = times
	}
	number := a.GetAlmsGetNumber()
	if number > 0 { // 救济金币领取数量
		changeConfig[data.AlmsGetNumber] = number
	}
	aCountDonw := a.GetAlmsCountDonw()
	if aCountDonw > 0 { // 救济倒计时
		changeConfig[data.AlmsCountDonw] = aCountDonw
	}
	dCountDonw := a.GetDepositCountDonw()
	if dCountDonw > 0 { // 快冲倒计时
		changeConfig[data.DepositCountDonw] = dCountDonw
	}
	if len(changeConfig) == 0 {
		entry.Errorln("没有配置发生变化")
		return
	}
	// 每次发生变化版本号加1
	if err := data.UpdataAlmsConfigVersion(); err != nil {
		entry.Errorln("版本号变化失败")
		return
	}
	// 修改redis 救济金配置
	if err := data.SetAlmsConfigWatchs(changeConfig, data.RedisTimeOut); err != nil {
		entry.WithError(err).Errorf("redis 救济金配置 改变失败 changeConfig(%v) \n", changeConfig)
		return
	}
	entry.WithFields(logrus.Fields{
		"changeConfig": changeConfig,
	}).Infoln("救济金配置变化通知成功")
}
