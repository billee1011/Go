package configclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	entityConf "steve/entity/config"
	"steve/server_pb/config"
	"steve/structs"

	"github.com/Sirupsen/logrus"
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

func ParseToGameLevelConfigMap(jsonStr string) (conf []entityConf.GameLevelConfig) {
	if err := json.Unmarshal([]byte(jsonStr), &conf); err != nil {
		logrus.Errorf("游戏配置数据反序列化失败：%s", err.Error())
	}
	return
}

//获取救济金配置
func GetAlmsConfigMap() (conf []entityConf.AlmsConfig,err error){
	almsStr, err := GetConfig("game", "alms")
	if err != nil {
		logrus.WithError(err).Errorln("获取救济金配置失败")
		return nil, err
	}
	if err := json.Unmarshal([]byte(almsStr), &conf); err != nil {
		logrus.WithError(err).Errorf("游戏配置数据反序列化失败：%s", err.Error())
		return nil, err
	}
	return
}

// GetGameConfigMap 获取游戏配置信息
func GetGameConfigMap() (gameConf []entityConf.GameConfig, err error) {
	gameStr, err := GetConfig("game", "config")
	if err != nil {
		logrus.WithError(err).Errorln("获取游戏配置失败")
		return nil, err
	}

	if err := json.Unmarshal([]byte(gameStr), &gameConf); err != nil {
		logrus.WithError(err).Errorf("游戏配置数据反序列化失败：%s", err.Error())
		return nil, err
	}

	return
}

// GetGameLevelConfigMap 获取游戏级别配置信息
func GetGameLevelConfigMap() (levelConf []entityConf.GameLevelConfig, err error) {
	levelStr, err := GetConfig("game", "levelconfig")
	if err != nil {
		logrus.WithError(err).Errorln("获取游戏级别配置失败")
		return nil, err
	}

	if err := json.Unmarshal([]byte(levelStr), &levelConf); err != nil {
		logrus.WithError(err).Errorf("游戏级别配置数据反序列化失败：%s", err.Error())
		return nil, err
	}

	return
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
