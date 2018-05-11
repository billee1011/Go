package rtoet

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

func translateQiReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomQiReq) (eventID server_pb.EventID, eventContext proto.Message, err error) {

	eventHeader := translateHeader(playerID, header, &req)
	eventContext = &server_pb.QiRequestEvent{
		Head: &eventHeader,
	}
	eventID = server_pb.EventID_event_qi_request
	return
}
