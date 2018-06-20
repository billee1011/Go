package utils

import (
	"errors"
	"fmt"
	"steve/client_pb/gate"
	"steve/client_pb/login"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/facade"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
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

// GenerateAccountName 生成账号名字
func GenerateAccountName(accountID uint64) string {
	return fmt.Sprintf("account_%v", accountID)
}

// RequestAuth 请求认证
func RequestAuth(client interfaces.Client, accountID uint64, accountName string, expireDuration time.Duration) (*login.LoginAuthRsp, error) {
	expire := time.Now().Add(expireDuration)
	request := &login.LoginAuthReq{
		AccountId:   proto.Uint64(accountID),
		AccountName: proto.String(accountName),
		Expire:      proto.Int64(expire.Unix()),
	}
	response := &login.LoginAuthRsp{}
	err := facade.Request(client, msgid.MsgID_LOGIN_AUTH_REQ, request, global.DefaultWaitMessageTime, msgid.MsgID_LOGIN_AUTH_RSP, response)
	return response, err
}

// RequestGateAuth 请求向网关服认证
func RequestGateAuth(client interfaces.Client, playerID uint64, expire int64, token string) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "RequestGateAuth",
		"player_id": playerID,
		"expire":    expire,
		"token":     token,
	})
	request := &gate.GateAuthReq{
		PlayerId: proto.Uint64(playerID),
		Expire:   proto.Int64(expire),
		Token:    proto.String(token),
	}
	response := &gate.GateAuthRsp{}
	err := facade.Request(client, msgid.MsgID_GATE_AUTH_REQ, request, global.DefaultWaitMessageTime, msgid.MsgID_GATE_AUTH_RSP, response)
	if err != nil {
		entry.WithError(err).Errorln("请求失败")
		return errors.New("请求失败")
	}
	if response.GetErrCode() != gate.ErrCode_SUCCESS {
		return fmt.Errorf("网关认证失败，错误码：%v", response.GetErrCode())
	}
	return nil
}
