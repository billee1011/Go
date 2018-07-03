package ddz

import (
	"steve/client_pb/room"
	"steve/server_pb/ddz"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// translateHeader 生成事件头
func translateHeader(playerID uint64, header *steve_proto_gaterpc.Header, body proto.Message) ddz.RequestEventHead {
	return ddz.RequestEventHead{
		PlayerId: playerID,
	}
}

// TranslateGrabRequest 转换抢地主请求
func TranslateGrabRequest(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.DDZGrabLordReq) (eventID int, eventContext proto.Message, err error) {

	head := translateHeader(playerID, header, &req)
	eventContext = &ddz.GrabRequestEvent{
		Head: &head,
	}
	eventID = int(ddz.EventID_event_grab_request)
	return
}
