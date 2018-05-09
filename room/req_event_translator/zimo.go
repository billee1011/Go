package rtoet

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

func translateZimoReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomZimoReq) (eventID server_pb.EventID, eventContext proto.Message, err error) {

	eventHeader := translateHeader(playerID, header, &req)
	card := translateCard(*req.GetCard())
	eventContext = &server_pb.ZimoRequestEvent{
		Head:  &eventHeader,
		Cards: &card,
	}
	eventID = server_pb.EventID_event_zimo_request
	return
}