package loginservice

import (
	"context"
	"crypto/md5"
	"fmt"
	"steve/client_pb/common"
	"steve/gutils"
	"steve/login/data"
	"steve/server_pb/login"
	"steve/server_pb/user"
	"steve/structs"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"steve/external/datareportclient"
	"steve/datareport/fixed"
)

// tokenSetter for test mock
var tokenSetter = data.SetPlayerToken

// tokenGetter for test mock
var tokenGetter = data.GetPlayerToken

// playerIDGetter for test mock
var playerIDGetter = func(accID uint64) (uint64, int) {
	entry := logrus.WithField("account_id", accID)
	rpcCli := structs.GetGlobalExposer().RPCClient
	hallCli, err := rpcCli.GetConnectByServerName("hall")
	if err != nil {
		entry.WithError(err).Errorln("获取大厅服连接失败")
		return 0, int(common.ErrCode_EC_FAIL)
	}
	playerDataService := user.NewPlayerDataClient(hallCli)
	rsp, err := playerDataService.GetPlayerByAccount(context.Background(),
		&user.GetPlayerByAccountReq{
			AccountId: accID,
		})
	if err != nil {
		entry.WithError(err).Errorln("请求玩家 ID 失败")
		return 0, int(common.ErrCode_EC_FAIL)
	}
	return rsp.GetPlayerId(), int(rsp.GetErrCode())
}

// ---------------------------------------------------------------------------------------------------

// LoginService 实现 login.LoginServiceServer
type LoginService struct {
	idAllocNode *gutils.Node
}

var defaultLoginService *LoginService

// Default return default object
func Default() *LoginService {
	return defaultLoginService
}

// Login 登录
func (ls *LoginService) Login(ctx context.Context, request *login.LoginRequest) (
	response *login.LoginResponse, err error) {

	entry := logrus.WithFields(logrus.Fields{
		"player_id": request.GetPlayerId(),
		"token":     request.GetToken(),
	})
	entry.Debugln("收到登录请求")
	response = &login.LoginResponse{
		ErrCode:  uint32(common.ErrCode_EC_FAIL),
		PlayerId: 0,
		Token:    "",
	}
	// auth by token success
	if authByToken(request.GetToken(), request.GetPlayerId()) {
		entry.Debugln("通过 token 认证成功")
		response.ErrCode = uint32(common.ErrCode_EC_SUCCESS)
		response.PlayerId = request.GetPlayerId()
		response.Token = generateToken(request.GetPlayerId())
		return
	}

	accID := request.GetAccountId()
	if accID == 0 {
		accID = uint64(ls.idAllocNode.Generate())
	}
	playerID, errCode := playerIDGetter(accID)
	if errCode != int(common.ErrCode_EC_SUCCESS) {
		response.ErrCode = uint32(errCode)
		return
	}

	response.ErrCode = uint32(common.ErrCode_EC_SUCCESS)
	response.PlayerId = playerID
	response.Token = generateToken(playerID)

	datareportclient.DataReport(fixed.LOG_TYPE_ACT,0,0,0,playerID,"1")

	return
}

// ---------------------------------------------------------------------------------------------------

// 生成并缓存 token
func generateToken(playerID uint64) string {
	entry := logrus.WithField("player_id", playerID)

	tokenSrc := fmt.Sprintf("%d%d%s", playerID,
		time.Now().Nanosecond(),
		viper.GetString("auth_key"))

	token := fmt.Sprintf("%x", md5.Sum([]byte(tokenSrc)))
	if err := tokenSetter(playerID, token, 24*time.Hour); err != nil {
		entry.WithError(err).Errorln("存储玩家 token 失败")
		return ""
	}
	entry.WithField("token", token).Debugln("生成token")
	return token
}

// authByToken token 认证
func authByToken(token string, playerID uint64) bool {
	if token == "" || playerID == 0 {
		return false
	}
	entry := logrus.WithFields(logrus.Fields{
		"token":     token,
		"player_id": playerID,
	})

	saveToken, err := tokenGetter(playerID)
	if err != nil {
		entry.WithError(err).Debugln("获取玩家 token 失败")
		return false
	}
	entry.WithField("save_token", saveToken).Debugln("token 认证")
	return saveToken == token
}

func init() {
	viper.SetDefault("auth_key", "some-secret-key")

	idAllocNode, err := gutils.NewNode(viper.GetInt64("node"))
	if err != nil {
		logrus.Panicln(err)
	}
	defaultLoginService = &LoginService{
		idAllocNode: idAllocNode,
	}
}
