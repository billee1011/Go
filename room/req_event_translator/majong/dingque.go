package majong

import (
	"steve/client_pb/room"
	server_pb "steve/entity/majong"
	"steve/structs/proto/gate_rpc"
)

// TranslateDingqueReq 转换定缺请求
func TranslateDingqueReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomDingqueReq) (eventID int, eventContext interface{}, err error) {

	eventHeader := translateHeader(playerID, header, &req)

	cardColor := translateClientCardColor(req.GetColor())
	eventContext = server_pb.DingqueRequestEvent{
		Head:  &eventHeader,
		Color: cardColor,
	}
	eventID = int(server_pb.EventID_event_dingque_request)
	return
}
