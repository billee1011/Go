package rtoet

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

func translateDingqueReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomDingqueReq) (eventID server_pb.EventID, eventContext proto.Message, err error) {

	eventHeader := translateHeader(playerID, header, &req)

	cardColor := translateCardColor(req.GetColor())
	eventContext = &server_pb.DingqueRequestEvent{
		Head:  &eventHeader,
		Color: cardColor,
	}
	eventID = server_pb.EventID_event_dingque_request
	return
}
