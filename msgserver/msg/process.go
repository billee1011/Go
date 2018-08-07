package msg

import (
	"steve/structs/proto/gate_rpc"
	"steve/client_pb/msgserver"
	"steve/structs/exchanger"
	"github.com/golang/protobuf/proto"
	"steve/client_pb/common"
	"steve/client_pb/msgid"
	"steve/msgserver/logic"
	"github.com/Sirupsen/logrus"
)

/*
 功能：
		1. 完成从GateWay(网关）过来的所有Client的请求消息的处理。
 		2. 通过core.coreConfig配置需要处理的消息列表。
		3. 需要在GateWay配置消息ID开始~ 结束区间 关联到当前服务名,GateWay才会把消息转发到此服务
*/


// 处理获取跑马灯请求
func ProcessGetHorseRaceReq(playerID uint64, header *steve_proto_gaterpc.Header, req msgserver.MsgSvrGetHorseRaceReq) (ret []exchanger.ResponseMsg) {

	logrus.Debugln("ProcessGetHorseRaceReq req", req)

	response := &msgserver.MsgSvrGetHorseRaceRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	//
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MSGSVR_GET_HORSE_RACE_RSP),
		Body:  response,
	}}

	list, tick, sleep, err := logic.GetMsgMgr().GetHorseRace(playerID)
	if err != nil {
		response.ErrCode = proto.Int32(int32(common.ErrCode_EC_FAIL))
		response.ErrDesc = proto.String("失败")
		logrus.Debugln("ProcessGetHorseRaceReq err:", err)
		return nil
	}
	response.Content = list
	response.Tick = &tick
	response.Sleep = &sleep
	logrus.Debugln("ProcessGetHorseRaceReq resp", response)
	return ret
}

