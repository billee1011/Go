package utils

import (
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"

	"github.com/Sirupsen/logrus"
)

type clientPlayer struct {
	playerID uint64
	coin     uint64
	client   interfaces.Client
	usrName  string
}

func (p *clientPlayer) GetID() uint64 {
	return p.playerID
}
func (p *clientPlayer) GetCoin() uint64 {
	return p.coin
}

func (p *clientPlayer) GetClient() interfaces.Client {
	return p.client
}

func (p *clientPlayer) GetUsrName() string {
	return p.usrName
}

// LoginUser 登录用户
func LoginUser(client interfaces.Client, userName string) (interfaces.ClientPlayer, error) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "LoginUser",
		"user_name": userName,
	})

	rsp := room.RoomLoginRsp{}
	err := client.Request(interfaces.SendHead{
		Head: interfaces.Head{
			MsgID: uint32(msgid.MsgID_ROOM_LOGIN_REQ),
		},
	}, &room.RoomLoginReq{
		UserName: &userName,
	}, global.DefaultWaitMessageTime, uint32(msgid.MsgID_ROOM_LOGIN_RSP), &rsp)

	if err != nil {
		logEntry.WithError(err).Errorln(errRequestFailed)
		return nil, err
	}
	return &clientPlayer{
		playerID: rsp.GetPlayerId(),
		coin:     rsp.GetCoin(),
		client:   client,
		usrName:  userName,
	}, nil
}

// LoginVisitor 登录游客
func LoginVisitor(client interfaces.Client, RoomVisitorLoginReq *room.RoomVisitorLoginReq) (interfaces.ClientPlayer, error) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "LoginVisitor",
	})

	rsp := room.RoomVisitorLoginRsp{}
	err := client.Request(interfaces.SendHead{
		Head: interfaces.Head{
			MsgID: uint32(msgid.MsgID_ROOM_VISITOR_LOGIN_REQ),
		},
	}, RoomVisitorLoginReq, global.DefaultWaitMessageTime, uint32(msgid.MsgID_ROOM_VISITOR_LOGIN_RSP), &rsp)

	if err != nil {
		logEntry.WithError(err).Errorln(errRequestFailed)
		return nil, err
	}
	return &clientPlayer{
		playerID: rsp.GetPlayerId(),
		coin:     rsp.GetCoin(),
		client:   client,
	}, nil
}

func UpdatePlayerClientInfo(client interfaces.Client, player interfaces.ClientPlayer, deskData *DeskData) {
	oldPlayer, exist := deskData.Players[player.GetID()]
	if !exist {
		return
	}

	newPlayer := DeskPlayer{
		Player:    player,
		Seat:      oldPlayer.Seat,
		Expectors: createPlayerExpectors(player.GetClient()),
	}
	deskData.Players[player.GetID()] = newPlayer
	return
}
