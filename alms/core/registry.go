package core

import (
	"steve/alms/almsserver"
	"steve/client_pb/msgid"
	"steve/structs/exchanger"

	"github.com/Sirupsen/logrus"
)

func registerHandles(e exchanger.Exchanger) error {

	panicRegister := func(msgID msgid.MsgID, h interface{}) {
		if err := e.RegisterHandle(uint32(msgID), h); err != nil {
			logrus.WithField("msg_id", msgID).Panic(err)
		}
	}
	panicRegister(msgid.MsgID_ALMS_GET_GOLD_REQ, almsserver.HandleGetAlmsReq)
	return nil
}
