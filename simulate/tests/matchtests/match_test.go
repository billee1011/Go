package matchtests

import (
	"steve/client_pb/msgId"
	"steve/server_pb/match"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	accountID := global.AllocAccountID()
	accountName := utils.GenerateAccountName(accountID)
	player, err := utils.LoginPlayer(accountID, accountName)
	assert.Nil(t, err)
	assert.NotNil(t, player)

	req := match.MatchRequest{
		Name: "hello",
	}

	rsp := match.MatchResponse{
		Echo: "world",
	}
	client := player.GetClient()
	client.Request(utils.CreateMsgHead(msgid.MsgID_MATCH_REQ), &req, global.DefaultWaitMessageTime, uint32(msgid.MsgID_MATCH_RSP), &rsp)
}
