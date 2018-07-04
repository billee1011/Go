package core

import (
	"github.com/Sirupsen/logrus"
	"steve/server_pb/room_mgr"
	"fmt"
	"errors"
	"context"
	"steve/structs"
)

type Sender struct {
	e *structs.Exposer
}

func NewSender() *Sender {
	s:= &Sender{
		e: matchCore.e,
	}

	return s
}

// 通知room服创建desk
func (s *Sender) createDesk(playersID []uint64) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "sender::CreateDesk",
	})
	logEntry.Debugln("sender::CreateDesk()")

	roomConnection, roomErr := s.core.e.RPCClient.GetConnectByServerName("room")
	if roomErr != nil {
		logEntry.WithError(roomErr).Errorln("获取room服失败")
	}

	if roomConnection == nil {
		logEntry.Errorln("获取room服失败，room_connection == nil")
		return errors.New("获取room服失败，matchCore::NofityRoomCreateDesk() room_connection == nil")
	}

	// 建立一个新的连接
	roomMgrClient := roommgr.NewRoomMgrClient(roomConnection)
	deskResp, deskErr := roomMgrClient.CreateDesk(context.Background(), &roommgr.CreateDeskRequest{
		PlayerId: playersID,
	})

	if deskErr != nil {
		logEntry.WithError(deskErr).Errorln("调用room服的CreateDesk()失败")
		return fmt.Errorf("call room::CreateDesk() failed: %v", deskErr)
	}

	fmt.Println("收到room服 CreateDesk()返回消息 : ", deskResp.GetErrCode())
	return nil
}
