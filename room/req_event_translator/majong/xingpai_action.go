package majong

import (
	"errors"
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// TranslateXingpaiActionReq 转换行牌动作请求
func TranslateXingpaiActionReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomXingpaiActionReq) (eventID int, eventContext proto.Message, err error) {

	eventHeader := translateHeader(playerID, header, &req)
	switch req.GetActionId() {
	case room.XingpaiAction_XA_CHI:
		{
			err = errors.New("未实现")
		}
	case room.XingpaiAction_XA_PENG:
		{
			eventID = int(server_pb.EventID_event_peng_request)
			eventContext = &server_pb.PengRequestEvent{
				Head: &eventHeader,
			}
		}
	case room.XingpaiAction_XA_GANG:
		{
			eventID = int(server_pb.EventID_event_gang_request)
			card := translateCard(req.GetGangCard())
			eventContext = &server_pb.GangRequestEvent{
				Head: &eventHeader,
				Card: &card,
			}
		}
	case room.XingpaiAction_XA_HU:
		{
			eventID = int(server_pb.EventID_event_hu_request)
			eventContext = &server_pb.HuRequestEvent{
				Head: &eventHeader,
			}
		}
	case room.XingpaiAction_XA_QI:
		{
			eventID = int(server_pb.EventID_event_qi_request)
			eventContext = &server_pb.QiRequestEvent{
				Head: &eventHeader,
			}
		}
	}
	return
}
