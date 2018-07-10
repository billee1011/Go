package connection

import (
	"steve/client_pb/gate"
	"steve/client_pb/msgId"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// HandleHeartBeat 处理心跳
func HandleHeartBeat(clientID uint64, header *steve_proto_gaterpc.Header, req gate.GateHeartBeatReq) (ret []exchanger.ResponseMsg) {
	connection := GetConnectionMgr().GetConnection(clientID)
	if connection == nil {
		return
	}
	connection.HeartBeat()

	response := gate.GateHeartBeatRsp{
		TimeStamp: proto.Uint64(req.GetTimeStamp()),
	}
	return []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_GATE_HEART_BEAT_RSP),
			Body:  &response,
		},
	}
}