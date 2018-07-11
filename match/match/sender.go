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
func (s *Sender) createDesk(players []DeskPlayer, gameID int) (resp *roommgr.CreateDeskResponse, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "Sender::createDesk()",
	})
	e := structs.GetGlobalExposer()

	rs, err := e.RPCClient.GetConnectByServerName("room")
	if err != nil {
		logEntry.WithError(err).Errorln("get 'room' service failed!!!")
	}
	createPlayers := []*roommgr.DeskPlayer{}
	for _, player := range players {
		createPlayers = append(createPlayers, &roommgr.DeskPlayer{
			PlayerId:   player.GetPlayerID(),
			RobotLevel: int32(player.GetRobotLv()),
		})
	}

	roomMgrClient := roommgr.NewRoomMgrClient(rs)
	resp, err = roomMgrClient.CreateDesk(context.Background(), &roommgr.CreateDeskRequest{
		Players: createPlayers,
		GameId:  uint32(gameID),
	})

	if err != nil {
		logEntry.WithError(err).Errorln("create desk failed!!!")
		return
	}

	logEntry.Debugln("create desk success.")
	return
}
