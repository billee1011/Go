package auth

import (
	"context"
	"steve/client_pb/login"
	msgid "steve/client_pb/msgId"
	"steve/common/auth"
	"steve/common/data/player"
	"steve/login/config"
	"steve/login/facade"
	"steve/login/global"
	"steve/server_pb/gateway"
	"steve/structs"
	"steve/structs/common"
	"steve/structs/proto/base"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"

	"github.com/Sirupsen/logrus"
)

// OnAuthRequest 客户端请求认证
func OnAuthRequest(clientID uint64, header *steve_proto_base.Header, body []byte) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "OnAuthRequest",
		"client_id": clientID,
	})
	request, err := translateRequest(body)
	if err != nil {
		entry.WithError(err).Infoln("请求数据转换失败")
		return
	}
	entry.WithFields(logrus.Fields{
		"account_id":   request.GetAccountId(),
		"account_name": request.GetAccountName(),
	}).Infoln("处理认证请求")
	handleAuthRequest(clientID, header, request)
}

// translateRequest 反序列化请求信息
func translateRequest(body []byte) (login.LoginAuthReq, error) {
	var request login.LoginAuthReq
	err := proto.Unmarshal(body, &request)
	return request, err
}

// handleAuthRequest 处理认证请求
// step 1. 验证认证时间是否超时
// step 2. 验证账号信息是否合法
// step 3. 分配网关
// step 4. 判断用户是否存在，如果存在直接回复现有用户信息
// step 5. 如果用户不存在， 生成用户信息并回复客户端
func handleAuthRequest(clientID uint64, header *steve_proto_base.Header, request login.LoginAuthReq) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":    "handleAuthRequest",
		"client_id":    clientID,
		"account_id":   request.GetAccountId(),
		"account_name": request.GetAccountName(),
		"token":        request.GetToken(),
		"expire":       request.GetExpire(),
	})
	response := &login.LoginAuthRsp{
		ErrCode: login.ErrorCode_ABNORMAL.Enum(),
	}
	defer responseAuth(clientID, header, response)
	// 验证请求信息是否合法
	if !checkRequest(clientID, &request, response) {
		return
	}
	// 分配网关服
	if !allocGate(response) {
		return
	}
	// 查询玩家信息
	if !aquirePlayer(&request, response) {
		return
	}
	// 生成 token
	if !generateToken(response) {
		return
	}
	response.ErrCode = login.ErrorCode_SUCCESS.Enum()
	entry.WithFields(logrus.Fields{
		"player_id": response.GetPlayerId(),
		"gate_ip":   response.GetGateIp(),
		"gate_port": response.GetGatePort(),
		"token":     response.GetGateToken(),
	}).Infoln("认证成功")
}

// responseAuth 认证应答
func responseAuth(clientID uint64, header *steve_proto_base.Header, response *login.LoginAuthRsp) {
	respHeader := &steve_proto_base.Header{
		MsgId:  proto.Uint32(uint32(msgid.MsgID_LOGIN_AUTH_RSP)),
		RspSeq: proto.Uint64(header.GetSendSeq()),
	}
	facade.SendPackage(global.GetMessageSender(), clientID, respHeader, response)
}

// checkRequest 检查请求是否合法，包括 token 到期时间和 token
func checkRequest(clientID uint64, request *login.LoginAuthReq, response *login.LoginAuthRsp) bool {
	if time.Now().Unix() > request.GetExpire() {
		response.ErrCode = login.ErrorCode_ERR_EXPIRE_TOKEN.Enum()
		return false
	}
	if !checkAccountToken(request) {
		response.ErrCode = login.ErrorCode_ERR_INVALID_TOKEN.Enum()
		return false
	}
	return true
}

// checkAccountToken 检查 token
func checkAccountToken(request *login.LoginAuthReq) bool {
	// TODO: 等待账号系统完善
	return true
}

// allocGate 分配网关
func allocGate(response *login.LoginAuthRsp) bool {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "allocGate",
	})

	rpcClient := structs.GetGlobalExposer().RPCClient
	cc, err := rpcClient.GetConnectByServerName(common.GateServiceName)
	if err != nil {
		entry.WithError(err).Errorln("获取网关服连接失败")
		return false
	}
	client := gateway.NewGateServiceClient(cc)
	gateResp, err := client.GetGatewayAddress(context.Background(), &gateway.GetGatewayAddressRequest{})
	if err != nil {
		entry.WithError(err).Errorln("调用网关服接口失败")
		return false
	}
	addr := gateResp.GetAddr()
	response.GateIp = proto.String(addr.GetIp())
	response.GatePort = proto.Int32(addr.GetPort())

	entry.WithFields(logrus.Fields{
		"gate_ip":   response.GetGateIp(),
		"gate_port": response.GetGatePort(),
	}).Debugln("分配网关地址")
	return true
}

// aquirePlayer 查询玩家信息
func aquirePlayer(request *login.LoginAuthReq, response *login.LoginAuthRsp) bool {
	accountID := request.GetAccountId()
	playerID := player.GetAccountPlayerID(accountID)
	if playerID == 0 {
		playerID = newPlayer(accountID)
		if playerID == 0 {
			return false
		}
	}
	response.PlayerId = proto.Uint64(playerID)
	return true
}

// newPlayer 创建玩家
func newPlayer(accountID uint64) uint64 {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":  "newPlayer",
		"account_id": accountID,
	})
	playerID, err := player.AllocPlayerID()
	if err != nil {
		entry.WithError(err).Errorln("分配玩家 ID 失败")
		return 0
	}
	if err := player.NewPlayer(accountID, playerID); err != nil {
		entry.WithError(err).Errorln("创建玩家失败")
		return 0
	}
	return playerID
}

// generateToken 生成认证码
func generateToken(response *login.LoginAuthRsp) bool {
	expire := (time.Now().Add(time.Minute * 5)).Unix()
	response.Expire = proto.Int64(expire)

	key := viper.GetString(config.AuthKey)
	token := auth.GenerateAuthToken(response.GetPlayerId(), response.GetGateIp(), int(response.GetGatePort()), expire, key)
	response.GateToken = proto.String(token)
	return true
}
