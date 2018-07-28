package loginservice

import (
	"context"
	"steve/client_pb/common"
	"steve/server_pb/login"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoginService_Login(t *testing.T) {
	const (
		ACCOUNTID uint64 = 1
		PLAYERID  uint64 = 100
	)

	playerIDGetter = func(accID uint64) (uint64, int) {
		assert.Equal(t, accID, ACCOUNTID)
		return PLAYERID, 0
	}

	var settedToken string
	tokenSetter = func(playerID uint64, token string, duration time.Duration) error {
		assert.Equal(t, playerID, PLAYERID)
		assert.NotEmpty(t, token)
		assert.NotZero(t, duration)
		settedToken = token
		return nil
	}

	tokenGetter = func(playerID uint64) (string, error) {
		return settedToken, nil
	}

	// 测试普通登录
	ls := LoginService{}
	response, err := ls.Login(context.Background(), &login.LoginRequest{
		AccountId: 1,
	})
	assert.Nil(t, err)
	assert.Equal(t, uint32(common.ErrCode_EC_SUCCESS), response.GetErrCode())
	assert.Equal(t, PLAYERID, response.GetPlayerId())
	assert.Equal(t, settedToken, response.GetToken())

	// 测试第二次使用 token 登录
	response, err = ls.Login(context.Background(), &login.LoginRequest{
		PlayerId: PLAYERID,
		Token:    settedToken,
	})
	assert.Nil(t, err)
	assert.Equal(t, uint32(common.ErrCode_EC_SUCCESS), response.GetErrCode())
	assert.Equal(t, PLAYERID, response.GetPlayerId())
	assert.Equal(t, settedToken, response.GetToken())
}
