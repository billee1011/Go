package util

import (
	"steve/client_pb/room"
	"steve/server_pb/ddz"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// translateDDZHeader 生成事件头
func translateDDZHeader(playerID uint64, header *steve_proto_gaterpc.Header, body proto.Message) ddz.RequestEventHead {
	return ddz.RequestEventHead{
		PlayerId: playerID,
	}
}

// TranslateGrabRequest 转换抢地主请求
func TranslateGrabRequest(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.DDZGrabLordReq) (eventID int, eventContext proto.Message, err error) {

	head := translateDDZHeader(playerID, header, &req)
	eventContext = &ddz.GrabRequestEvent{
		Head: &head,
		Grab: req.GetGrab(),
	}
	eventID = int(ddz.EventID_event_grab_request)
	return
}

// TranslateDoubleRequest 转换加倍请求
func TranslateDoubleRequest(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.DDZDoubleReq) (eventID int, eventContext proto.Message, err error) {

	head := translateDDZHeader(playerID, header, &req)
	eventContext = &ddz.DoubleRequestEvent{
		Head:     &head,
		IsDouble: req.GetIsDouble(),
	}
	eventID = int(ddz.EventID_event_double_request)
	return
}

// TranslatePlayCardRequest 转换出牌请求
func TranslatePlayCardRequest(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.DDZPlayCardReq) (eventID int, eventContext proto.Message, err error) {

	head := translateDDZHeader(playerID, header, &req)
	eventContext = &ddz.PlayCardRequestEvent{
		Head:     &head,
		Cards:    req.Cards,
		CardType: ddz.CardType(int32(req.GetCardType())),
	}
	eventID = int(ddz.EventID_event_chupai_request)
	return
}

// TranslateResumeRequest 转换恢复对局请求
func TranslateResumeRequest(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.DDZResumeGameReq) (eventID int, eventContext proto.Message, err error) {

	head := translateDDZHeader(playerID, header, &req)
	eventContext = &ddz.ResumeRequestEvent{
		Head: &head,
	}
	eventID = int(ddz.EventID_event_resume_request)
	return
}
