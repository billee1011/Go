package rtoet

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

func translateHuRequest(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomHuReq) (eventID server_pb.EventID, eventContext proto.Message, err error) {

	eventHeader := translateHeader(playerID, header, &req)
	eventContext = &server_pb.HuRequestEvent{
		Head: &eventHeader,
	}
	eventID = server_pb.EventID_event_hu_request
	return
}
