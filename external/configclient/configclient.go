package configclient

import (
	"context"
	"errors"
	"fmt"
	"steve/server_pb/config"
	"steve/structs"
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
