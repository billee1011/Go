package rtoet

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

func translateCartoonFinishReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomCartoonFinishReq) (eventID server_pb.EventID, eventContext proto.Message, err error) {
	eventContext = &server_pb.CartoonFinishRequestEvent{
		CartoonType: int32(req.GetCartoonType()),
	}
	eventID = server_pb.EventID_event_cartoon_finish_request
	return
}
