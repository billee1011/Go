package match

import (
	"context"
	"steve/server_pb/room_mgr"
	"steve/structs"

	"github.com/Sirupsen/logrus"
)

type Sender struct {
}

func NewSender() *Sender {
	return &Sender{}
}

// 通知room服创建desk
func (s *Sender) createDesk(playersID []uint64, gameID int) (resp *roommgr.CreateDeskResponse, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "Sender::createDesk()",
	})
	e := structs.GetGlobalExposer()

	rs, err := e.RPCClient.GetConnectByServerName("room")
	if err != nil {
		logEntry.WithError(err).Errorln("get 'room' service failed!!!")
	}

	roomMgrClient := roommgr.NewRoomMgrClient(rs)
	resp, err = roomMgrClient.CreateDesk(context.Background(), &roommgr.CreateDeskRequest{
		PlayerId: playersID,
		GameId:   uint32(gameID),
	})

	if err != nil {
		logEntry.WithError(err).Errorln("create desk failed!!!")
		return
	}

	logEntry.Debugln("create desk success.")
	return
}
