package rtoet

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

func translateChupaiReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomChupaiReq) (eventID server_pb.EventID, eventContext proto.Message, err error) {

	eventHeader := translateHeader(playerID, header, &req)

	card := translateCard(req.GetCard())
	eventContext = &server_pb.ChupaiRequestEvent{
		Head:       &eventHeader,
		Cards:      &card,
		TingAction: translateTingAction(req.GetTingAction()),
	}
	eventID = server_pb.EventID_event_chupai_request
	return
}
