package matchtests

import (
	"steve/client_pb/common"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/simulate/cheater"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClearAllMatch 清空所有的匹配
// 第一步:登录一个玩家a，发起匹配请求
// 第二步:清空所有的匹配
// 第三步:登录一个玩家b,发起匹配

// 期望:第三步执行完后,玩家b匹配
func TestClearAllMatch(t *testing.T) {
	// 先清空所有的匹配
	cheater.ClearAllMatch()

	cancelMatchRspNtfExpectors := map[int]interfaces.MessageExpector{}

	// 登录一个玩家
	for i := 0; i < 1; i++ {
		// 登录用户
		player, err := utils.LoginNewPlayer()
		assert.Nil(t, err)
		assert.NotNil(t, player)
		client := player.GetClient()

		// 取消匹配的期待
		cancelMatchRspNtfExpector, err := client.ExpectMessage(msgid.MsgID_CANCEL_MATCH_RSP)
		assert.Nil(t, err)
		cancelMatchRspNtfExpectors[i] = cancelMatchRspNtfExpector

		// 发起匹配
		_, err = utils.RequestMatch(player, uint32(common.GameId_GAMEID_XUELIU), uint32(1))
		assert.Nil(t, err)
	}

	// 清空所有的匹配
	cheater.ClearAllMatch()

	// 期望收到取消匹配的回复
	for _, e := range cancelMatchRspNtfExpectors {
		ntf := &match.CancelMatchRsp{}
		assert.Nil(t, e.Recv(global.DefaultWaitMessageTime, ntf))

		// 期望取消成功
		assert.Equal(t, int32(match.MatchError_EC_SUCCESS), ntf.GetErrCode())
	}
}
