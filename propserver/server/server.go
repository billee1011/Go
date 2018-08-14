package server

/*
	功能： RPC服务实现类，完成当前模块服务的所有RPC接口实现和处理
    作者： SkyWang
    日期： 2018-7-24

*/

import (
	"context"
	"github.com/Sirupsen/logrus"
	"steve/server_pb/propserver"
	"steve/propserver/logic"
)

// PropsServer 实现 props.PropsServer
type PropsServer struct{}

var _ props.PropsServer = new(PropsServer)

// 获取玩家道具列表
func (gs *PropsServer) GetUserProps(ctx context.Context, request *props.GetPropsReq) (response *props.GetPropsRsp, err error) {
	logrus.Debugln("GetUserProps req", *request)
	response = &props.GetPropsRsp{}
	response.ErrCode = props.ResultStat_FAILED

	// 参数检查
	if request.GetUid()== 0 {
		response.ErrCode = props.ResultStat_ERR_ARG
		logrus.Errorln("GetUserProps resp", *response)
		return response, nil
	}

	// 调用逻辑实现代码
	m, err := logic.GetMyLogic().GetUserProps(request.GetUid(), request.GetPropsId())

	// 逻辑代码返回错误
	if err != nil {
		response.ErrCode = props.ResultStat_FAILED
		logrus.WithError(err).Errorln("GetGold resp", *response)
		return response, nil
	}
	// 设置返回值
	retList := make([]*props.GetItem, 0,  len(m))
	for k, v := range  m {
		a := new(props.GetItem)
		a.PropsId = k
		a.PropsNum = v
		retList = append(retList, a)
	}
	response.PropsList = retList


	response.ErrCode = props.ResultStat_SUCCEED
	logrus.Debugln("GetUserProps resp", *response)
	return response, nil
}

// 添加用户道具
func (gs *PropsServer) AddUserProps(ctx context.Context, request *props.AddPropsReq) (response *props.AddPropsRsp, err error) {
	logrus.Debugln("AddUserProps req", *request)

	response = &props.AddPropsRsp{}
	response.ErrCode = props.ResultStat_FAILED

	// 参数检查
	if request.GetUid() == 0 {
		response.ErrCode = props.ResultStat_ERR_ARG
		logrus.Errorln("AddUserProps resp", *response)
		return response, nil
	}

	// 调用逻辑实现代码
	m := make(map[uint64]int64, len(request.PropsList))
	for _, prop := range  request.PropsList {
		m[prop.PropsId] = prop.AddNum
	}
	err = logic.GetMyLogic().AddUserProps(request.GetUid(), m,
		request.GetSeq(), request.GetFuncId(), request.GetChannel(), request.GetTime(), request.GetGameId(), request.GetLevel())

	// 逻辑代码返回错误
	if err != nil {
		response.ErrCode = props.ResultStat_FAILED
		logrus.Errorln("AddUserProps resp", *response)
		return response, nil
	}
	// 设置返回值


	response.ErrCode = props.ResultStat_SUCCEED
	logrus.Debugln("AddUserProps resp", *response)
	return response, nil
}
