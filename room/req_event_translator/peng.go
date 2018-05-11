package rtoet

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

func translatePengReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomPengReq) (eventID server_pb.EventID, eventContext proto.Message, err error) {

	eventHeader := translateHeader(playerID, header, &req)
	eventContext = &server_pb.PengRequestEvent{
		Head: &eventHeader,
	}
	eventID = server_pb.EventID_event_peng_request
	return
}
