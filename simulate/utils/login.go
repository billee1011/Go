package utils

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"

	"github.com/Sirupsen/logrus"
)

type clientPlayer struct {
	playerID  uint64
	coin      uint64
	client    interfaces.Client
	usrName   string
	expectors map[msgid.MsgID]interfaces.MessageExpector
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

func (p *clientPlayer) AddExpectors(msgIDs ...msgid.MsgID) {
	for _, msgID := range msgIDs {
		p.expectors[msgID], _ = p.client.ExpectMessage(msgID)
	}
}

func (p *clientPlayer) GetExpector(msgID msgid.MsgID) interfaces.MessageExpector {
	return p.expectors[msgID]
}

func createPlayer(playerID uint64, coin uint64, client interfaces.Client, userName string) *clientPlayer {
	return &clientPlayer{
		playerID:  playerID,
		coin:      coin,
		client:    client,
		usrName:   userName,
		expectors: make(map[msgid.MsgID]interfaces.MessageExpector),
	}
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
	return createPlayer(rsp.GetPlayerId(), rsp.GetCoin(), client, ""), nil
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
	return createPlayer(rsp.GetPlayerId(), rsp.GetCoin(), client, ""), nil
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
