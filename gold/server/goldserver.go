package server

/*
	功能： RPC服务实现类，完成当前模块服务的所有RPC接口实现和处理
    作者： SkyWang
    日期： 2018-7-24

*/

import (
	"context"
	"steve/gold/logic"
	"steve/server_pb/gold"

	"github.com/Sirupsen/logrus"
)

// GoldService 实现 gold.GoldServer
type GoldServer struct{}

var _ gold.GoldServer = new(GoldServer)

// 获取玩家金币
func (gs *GoldServer) GetGold(ctx context.Context, request *gold.GetGoldReq) (response *gold.GetGoldRsp, err error) {
	logrus.Debugln("GetGold req", *request)
	response = &gold.GetGoldRsp{}
	response.ErrCode = gold.ResultStat_FAILED

	// 参数检查
	if request.GetItem() == nil {
		response.ErrCode = gold.ResultStat_ERR_ARG
		logrus.Errorln("GetGold resp", *response)
		return response, nil
	}

	// 调用逻辑实现代码
	item := request.GetItem()
	value, err := logic.GetGoldMgr().GetGold(item.GetUid(), int16(item.GetGoldType()))

	// 逻辑代码返回错误
	if err != nil {
		response.ErrCode = gold.ResultStat_FAILED
		logrus.WithError(err).Errorln("GetGold resp", *response)
		return response, nil
	}
	// 设置返回值
	item.Value = value
	response.Item = item

	response.ErrCode = gold.ResultStat_SUCCEED
	logrus.Debugln("GetGold resp", *response)
	return response, nil
}

// 加玩家金币
func (gs *GoldServer) AddGold(ctx context.Context, request *gold.AddGoldReq) (response *gold.AddGoldRsp, err error) {
	logrus.Debugln("AddGold req", *request)

	response = &gold.AddGoldRsp{}
	response.ErrCode = gold.ResultStat_FAILED

	// 参数检查
	if request.GetItem() == nil {
		response.ErrCode = gold.ResultStat_ERR_ARG
		logrus.Errorln("AddGold resp", *response)
		return response, nil
	}

	// 调用逻辑实现代码
	item := request.GetItem()
	after, err := logic.GetGoldMgr().AddGold(item.GetUid(), int16(item.GetGoldType()), item.GetChangeValue(),
		item.GetSeq(), item.GetFuncId(), item.GetChannel(), item.GetTime(), item.GetGameId(), item.GetLevel())

	// 逻辑代码返回错误
	if err != nil {
		response.ErrCode = gold.ResultStat_FAILED
		logrus.Errorln("AddGold resp", *response)
		return response, nil
	}
	// 设置返回值
	response.CurValue = after

	response.ErrCode = gold.ResultStat_SUCCEED
	logrus.Debugln("AddGold resp", *response)
	return response, nil
}
