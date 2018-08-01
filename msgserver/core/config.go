package core

/*
	功能： 服务配置中心，定义RPC服务关联，Client消息分派。

 */
import (
	"steve/client_pb/msgid"
	"steve/msgserver/msg"
)


/////////////////////////////////////////[1.配置线程模型]////////////////////////////////////
// 是否采用单线程运行所有协程(goroutime)
var bSingleThread = true

/////////////////////////////[2.定义RPC服务实现]//////////////////////////////////////////
// PB文件中定义的RPC服务接口,如果不提供RPC服务，设置为nil
var pbService interface{} = nil
// PB定义的RPC服务的实现类,如果不提供RPC服务，设置为nil
var pbServerImp interface{} = nil

/////////////////////////////[3.处理Client消息]////////////////////////////////////////////
// 如果需要处理Client消息，需要在GateWay配置消息ID开始~ 结束区间 关联到当前服务名,GateWay才会把消息转发到此服务.

// 添加从GateWay过来的Client消息处理
var mapMsg  = map[msgid.MsgID] interface{} {
	msgid.MsgID_MSGSVR_GET_HORSE_RACE_REQ:msg.ProcessGetHorseRaceReq,
}







