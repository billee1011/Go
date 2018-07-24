package majong

import (
	"steve/client_pb/room"
	server_pb "steve/entity/majong"
	"steve/structs/proto/gate_rpc"
)

// TranslateChupaiReq 转换出牌请求
func TranslateChupaiReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomChupaiReq) (eventID int, eventContext interface{}, err error) {

	eventHeader := translateHeader(playerID, header, &req)

	card := translateCard(req.GetCard())
	eventContext = server_pb.ChupaiRequestEvent{
		Head:       &eventHeader,
		Cards:      &card,
		TingAction: translateTingAction(req.GetTingAction()),
	}
	eventID = int(server_pb.EventID_event_chupai_request)
	return
}
