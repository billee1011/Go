package auth

import (
	"context"
	"fmt"
	"steve/client_pb/login"
	"steve/client_pb/msgId"
	"steve/common/data/player"
	"steve/gateway/config"
	"steve/gateway/connection"
	"steve/gateway/gateservice"
	"steve/server_pb/gateway"
	server_login_pb "steve/server_pb/login"
	"steve/structs"
	"steve/structs/common"
	"steve/structs/proto/base"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HandleLoginRequest 处理登录请求
func HandleLoginRequest(clientID uint64, reqHeader *steve_proto_base.Header, body []byte) (responses []*steve_proto_gaterpc.ResponseMessage) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "handleLoginRequest",
		"client_id": clientID,
	})
	request := login.LoginAuthReq{}
	if err := proto.Unmarshal(body, &request); err != nil {
		entry.WithError(err).Warningln("反序列化失败")
		return nil
	}

	response := execLogin(clientID, request)
	body, err := proto.Marshal(&response)
	if err != nil {
		entry.WithError(err).Errorln("序列化失败")
		return nil
	}

	responses = []*steve_proto_gaterpc.ResponseMessage{{
		Header: &steve_proto_gaterpc.Header{
			MsgId: uint32(msgid.MsgID_LOGIN_AUTH_RSP),
		},
		Body: body,
	}}
	return responses
}

// execLogin 执行登录
func execLogin(clientID uint64, clientRequest login.LoginAuthReq) (clientResponse login.LoginAuthRsp) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":  "execLogin",
		"client_id":  clientID,
		"account_id": clientRequest.GetAccountId(),
		"token":      clientRequest.GetToken(),
	})
	clientResponse = login.LoginAuthRsp{
		ErrCode:  login.ErrorCode_ABNORMAL.Enum(),
		PlayerId: proto.Uint64(0),
	}
	cm := connection.GetConnectionMgr()
	connection := cm.GetConnection(clientID)
	if connection == nil || connection.GetPlayerID() != 0 {
		entry.Warningln("客户端已经登录")
		return
	}
	// TODO : 校验 token
	playerID, err := callLoginService(clientRequest.GetAccountId())
	if err != nil {
		entry.Errorln(err)
		return
	}

	clientResponse.ErrCode = login.ErrorCode_SUCCESS.Enum()
	clientResponse.PlayerId = proto.Uint64(playerID)

	checkAnother(playerID)
	connection.AttachPlayer(playerID)
	return
}

// callLoginService 调用登录服务
func callLoginService(accountID uint64) (playerID uint64, err error) {
	exposer := structs.GetGlobalExposer()
	cc, err := exposer.RPCClient.GetConnectByServerName(common.LoginServiceName)
	if err != nil {
		return
	}
	loginClient := server_login_pb.NewLoginServiceClient(cc)
	request := server_login_pb.LoginRequest{
		AccountId: accountID,
	}
	response, err := loginClient.Login(context.Background(), &request)
	if err != nil {
		return
	}
	playerID = response.GetPlayerId()
	return playerID, nil
}

// checkAnother 顶号检查
func checkAnother(playerID uint64) {
	gateAddr := player.GetPlayerGateAddr(playerID)
	if gateAddr == "" {
		return
	}
	localGateAddr := fmt.Sprintf("%s:%d", config.GetRPCAddr(), config.GetRPCPort())

	entry := logrus.WithFields(logrus.Fields{
		"func_name":       "checkAnother",
		"player_id":       playerID,
		"gate_addr":       gateAddr,
		"local_gate_addr": localGateAddr,
	})
	entry.Infoln("顶号登录")
	if gateAddr == localGateAddr {
		// 玩家原本在此网关登录
		gateservice.AnotherLogin(playerID)
	} else {
		// 玩家原本在其他网关服登录
		exposer := structs.GetGlobalExposer()
		cc, err := exposer.RPCClient.GetConnectByServerName(common.GateServiceName)
		if err != nil || cc == nil {
			entry.WithError(err).Warningln("发起顶号通知失败")
			return
		}
		client := gateway.NewGateServiceClient(cc)
		client.AnotherLogin(context.Background(), &gateway.AnotherLoginRequest{PlayerId: playerID})
	}
}
