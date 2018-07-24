package majong

import (
	"steve/client_pb/room"
	server_pb "steve/entity/majong"
	"steve/structs/proto/gate_rpc"
)

// TranslateHuansanzhangReq 转换换三张请求
func TranslateHuansanzhangReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomHuansanzhangReq) (eventID int, eventContext interface{}, err error) {

	eventHeader := translateHeader(playerID, header, &req)

	eventContext = server_pb.HuansanzhangRequestEvent{
		Head:  &eventHeader,
		Cards: translateCards(req.GetCards()),
		Sure:  req.GetSure(),
	}
	eventID = int(server_pb.EventID_event_huansanzhang_request)
	return
}
