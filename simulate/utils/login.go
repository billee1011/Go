package utils

import (
	"errors"
	"fmt"
	"steve/client_pb/gate"
	"steve/client_pb/login"
	msgid "steve/client_pb/msgId"
	"steve/simulate/config"
	"steve/simulate/connect"
	"steve/simulate/facade"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type clientPlayer struct {
	playerID  uint64
	coin      uint64
	client    interfaces.Client
	usrName   string
	accountID uint64
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

func (p *clientPlayer) GetAccountID() uint64 {
	return p.accountID
}

// LoginPlayer 登录玩家
func LoginPlayer(accountID uint64, accountName string) (interfaces.ClientPlayer, error) {
	loginClient := connect.NewTestClient(config.LoginServerAddr, config.ClientVersion)
	if loginClient == nil {
		return nil, errors.New("连接登录服失败")
	}
	loginResponse, err := RequestAuth(loginClient, accountID, accountName, time.Minute*5)
	if err != nil {
		return nil, fmt.Errorf("发起登录服认证请求失败：%v", err)
	}
	if loginResponse.GetErrCode() != login.ErrorCode_SUCCESS {
		return nil, fmt.Errorf("登录服认证失败， 错误码： %v", loginResponse.GetErrCode())
	}

	playerID := loginResponse.GetPlayerId()
	expire := loginResponse.GetExpire()
	token := loginResponse.GetGateToken()

	gateIP := loginResponse.GetGateIp()
	gatePort := loginResponse.GetGatePort()
	gateAddr := fmt.Sprintf("%s:%d", gateIP, gatePort)

	gateClient := connect.NewTestClient(gateAddr, config.ClientVersion)
	if gateClient == nil {
		return nil, fmt.Errorf("连接网关服失败，网关地址：%v", gateAddr)
	}
	err = RequestGateAuth(gateClient, playerID, expire, token)
	if err != nil {
		return nil, fmt.Errorf("网关认证失败: %v", err)
	}
	return &clientPlayer{
		playerID:  playerID,
		coin:      0, // TODO, 从服务器加载数据
		client:    gateClient,
		usrName:   "", // TODO: delete
		accountID: accountID,
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
