package majong

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// TranslateDingqueReq 转换定缺请求
func TranslateDingqueReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomDingqueReq) (eventID int, eventContext proto.Message, err error) {

	eventHeader := translateHeader(playerID, header, &req)

	cardColor := translateClientCardColor(req.GetColor())
	eventContext = &server_pb.DingqueRequestEvent{
		Head:  &eventHeader,
		Color: cardColor,
	}
	eventID = int(server_pb.EventID_event_dingque_request)
	return
}
