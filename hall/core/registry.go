package core

import (
	"steve/client_pb/msgid"
	"steve/hall/user"
	"steve/structs/exchanger"

	"github.com/Sirupsen/logrus"
)

func registerHandles(e exchanger.Exchanger) error {

	panicRegister := func(msgID msgid.MsgID, h interface{}) {
		if err := e.RegisterHandle(uint32(msgID), h); err != nil {
			logrus.WithField("msg_id", msgID).Panic(err)
		}
	}
	panicRegister(msgid.MsgID_HALL_GET_PLAYER_INFO_REQ, user.HandleGetPlayerInfoReq)
	panicRegister(msgid.MsgID_HALL_GET_PLAYER_STATE_REQ, user.HandleGetPlayerStateReq)
	panicRegister(msgid.MsgID_HALL_GET_GAME_INFO_REQ, user.HandleGetGameInfoReq)
	panicRegister(msgid.MsgID_HALL_UPDATE_PLAYER_INFO_REQ, user.HandleUpdatePlayerInoReq)
	panicRegister(msgid.MsgID_HALL_REAL_NAME_REQ, user.HandleRealNameReq)
	panicRegister(msgid.MsgID_HALL_GET_PLAYER_GAME_INFO_REQ, user.HandleRealNameReq)
	return nil
}
