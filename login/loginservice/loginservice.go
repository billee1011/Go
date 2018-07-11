package loginservice

import (
	"context"
	"steve/login/auth"
	"steve/server_pb/login"
)

// LoginService 实现 login.LoginServiceServer
type LoginService struct{}

var _ login.LoginServiceServer = new(LoginService)

// Login 登录
func (ls *LoginService) Login(ctx context.Context, request *login.LoginRequest) (response *login.LoginResponse, err error) {
	response = &login.LoginResponse{}
	playerID := auth.HandleLoginRequest(request.GetAccountId())
	// TODO : 处理失败的情况
	response.ErrCode = 0
	response.PlayerId = playerID
	return response, nil
}
