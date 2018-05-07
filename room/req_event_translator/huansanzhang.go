package rtoet

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

func translateHuansanzhangReq(playerID uint64, header *steve_proto_gaterpc.Header, body proto.Message) (eventID server_pb.EventID, eventContext []byte, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "translateHuansanzhangReq",
		"player_id": playerID,
	})

	req, ok := body.(*room.RoomHuansanzhangReq)
	if !ok {
		err = errMessageTypeNotMatch
		return
	}

	eventHeader := translateHeader(playerID, header, body)

	event := server_pb.HuansanzhangRequestEvent{
		Head:  &eventHeader,
		Cards: translateCards(req.GetCards()),
	}
	eventID = server_pb.EventID_event_huansanzhang_request
	if err = proto.Unmarshal(eventContext, &event); err != nil {
		logEntry.WithError(err).Errorln(err)
		err = errUnmarshalEvent
		return
	}
	return
}
