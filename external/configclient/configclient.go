package configclient

import (
	"context"
	"errors"
	"fmt"
	"steve/entity/constant"
	"steve/server_pb/config"
	"steve/structs"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	nsq "github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

// ConfigCliGetter config client 获取
var ConfigCliGetter func() (config.ConfigClient, error)

// GetConfig 获取配置
func GetConfig(key, subkey string) (string, error) {
	configCli, err := getConfigCli()
	if err != nil {
		return "", err
	}
	rsp, err := configCli.GetConfig(context.Background(), &config.GetConfigReq{
		Key: key, Subkey: subkey,
	})
	if err != nil {
		return "", fmt.Errorf("rpc 调用失败:%v", err.Error())
	}
	val := rsp.GetValue()
	errCode := rsp.GetErrCode()
	if errCode != 0 {
		return "", fmt.Errorf("获取配置失败，错误码：%d", errCode)
	}
	return val, nil
}

// GetConfigUntilSucc 获取配置直到成功
func GetConfigUntilSucc(key, subkey string, maxRetry int, retryInterval time.Duration) (string, error) {
	entry := logrus.WithFields(logrus.Fields{
		"key":            key,
		"subkey":         subkey,
		"max_retry":      maxRetry,
		"retry_interval": retryInterval,
	})

	var configCli config.ConfigClient
	var err error
	curRetry := 0
	for {
		configCli, err = getConfigCli()
		if err == nil && configCli != nil {
			break
		}
		curRetry++
		if curRetry < maxRetry {
			entry.Infoln("获取配置等待重试")
			time.Sleep(retryInterval)
			continue
		}
		break
	}
	if err != nil || configCli == nil {
		return "", err
	}
	rsp, err := configCli.GetConfig(context.Background(), &config.GetConfigReq{
		Key: key, Subkey: subkey,
	})
	if err != nil {
		return "", fmt.Errorf("rpc 调用失败:%v", err.Error())
	}
	val := rsp.GetValue()
	errCode := rsp.GetErrCode()
	if errCode != 0 {
		return "", fmt.Errorf("获取配置失败，错误码：%d", errCode)
	}
	return val, nil
}

// ConfigChangeHandle 配置变化处理器
type ConfigChangeHandle func(key, subkey, val string) error

type nsqHandler struct {
	handlerFunc ConfigChangeHandle
}

func (nh *nsqHandler) HandleMessage(message *nsq.Message) error {
	cfg := config.ConfigUpdate{}
	if err := proto.Unmarshal(message.Body, &cfg); err != nil {
		logrus.WithError(err).Errorln("反序列化失败")
		return fmt.Errorf("反序列化失败:%s", err.Error())
	}
	return nh.handlerFunc(cfg.GetKey(), cfg.GetSubkey(), cfg.GetVal())
}

// SubConfigChangeCustom 使用自定义通道订阅配置变化
func SubConfigChangeCustom(key, subkey, channel string, handle ConfigChangeHandle) error {
	exposer := structs.GetGlobalExposer()
	return exposer.Subscriber.Subscribe(constant.UpdateConfig, channel, &nsqHandler{
		handlerFunc: handle,
	})
}

// SubConfigChange 订阅配置改变
// 使用 [rpc_server_name]+[node] 作为 channel， 对于没有配置 rpc_server_name 的服务，要使用 SubConfigChangeCustom 来订阅
func SubConfigChange(key, subkey string, handle ConfigChangeHandle) error {
	exposer := structs.GetGlobalExposer()
	channel := fmt.Sprintf("%s_%d", viper.GetString("rpc_server_name"), viper.GetInt("node"))

	return exposer.Subscriber.Subscribe(constant.UpdateConfig, channel, &nsqHandler{
		handlerFunc: handle,
	})
}

// 根据金币服的路由策略生成服务连接获取方式
func getConfigCli() (config.ConfigClient, error) {
	if ConfigCliGetter != nil {
		return ConfigCliGetter()
	}
	e := structs.GetGlobalExposer()
	cc, err := e.RPCClient.GetConnectByServerName("configuration")
	if err != nil || cc == nil {
		return nil, errors.New("no connection")
	}
	return config.NewConfigClient(cc), nil
}
