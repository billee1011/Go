package majong

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// TranslateCartoonFinishReq 转换动画完成请求
func TranslateCartoonFinishReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomCartoonFinishReq) (eventID int, eventContext proto.Message, err error) {
	eventContext = &server_pb.CartoonFinishRequestEvent{
		CartoonType: int32(req.GetCartoonType()),
		PlayerId:    playerID,
	}
	eventID = int(server_pb.EventID_event_cartoon_finish_request)
	return
}
