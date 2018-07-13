package matchtests

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_PlayerStates 测试玩家空闲状态
func Test_PlayerStates_Idle(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)

	player.AddExpectors(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP)

	client := player.GetClient()
	client.SendPackage(utils.CreateMsgHead(msgid.MsgID_HALL_GET_PLAYER_STATE_REQ), &hall.HallGetPlayerStateReq{})

	expector := player.GetExpector(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP)
	response := hall.HallGetPlayerStateRsp{}
	expector.Recv(global.DefaultWaitMessageTime, &response)
	assert.Equal(t, response.GetPlayerState(), common.PlayerState_PS_IDLE)
}

// Test_PlayerStates_Gameing 测试玩家在游戏中的状态
func Test_PlayerStates_Gameing(t *testing.T) {
	params := global.NewCommonStartGameParams()
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)

	deskPlayer := utils.GetDeskPlayerBySeat(0, deskData)
	deskPlayer.Player.AddExpectors(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP)

	client := deskPlayer.Player.GetClient()
	client.SendPackage(utils.CreateMsgHead(msgid.MsgID_HALL_GET_PLAYER_STATE_REQ), &hall.HallGetPlayerStateReq{})

	expector := deskPlayer.Player.GetExpector(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP)
	response := hall.HallGetPlayerStateRsp{}
	expector.Recv(global.DefaultWaitMessageTime, &response)
	assert.Equal(t, response.GetPlayerState(), common.PlayerState_PS_GAMEING)
	assert.Equal(t, response.GetGameId(), common.GameId(params.GameID))
}
