package core

import (
	"steve/hall/login"
	"steve/structs/exchanger"
	"steve/structs/proto/msg"

	"github.com/Sirupsen/logrus"
)

func registerHandles(e exchanger.Exchanger) error {

	panicRegister := func(msgID steve_proto_msg.MsgID, h interface{}) {
		if err := e.RegisterHandle(uint32(msgID), h); err != nil {
			logrus.WithField("msg_id", msgID).Panic(err)
		}
	}

	panicRegister(steve_proto_msg.MsgID_hall_login, login.HandleLogin)
	return nil
}
