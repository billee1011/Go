package core

import (
	"context"
	"steve/configuration/data"
	"steve/server_pb/config"

	"github.com/Sirupsen/logrus"
)

type configServer struct {
}

func (cs *configServer) GetConfig(ctx context.Context, req *config.GetConfigReq) (*config.GetConfigRsp, error) {
	key, subkey := req.GetKey(), req.GetSubkey()

	entry := logrus.WithFields(logrus.Fields{
		"key":    key,
		"subkey": subkey,
	})

	value, err := data.GetConfig(key, subkey)
	rsp := &config.GetConfigRsp{
		Value: value,
	}
	if err != nil {
		rsp.ErrCode = int32(-99)
		entry.Errorln("获取配置失败")
	} else {
		rsp.ErrCode = int32(0)
	}
	return rsp, nil
}
