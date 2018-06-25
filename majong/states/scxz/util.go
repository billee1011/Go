package scxz

import (
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

//SettleOver 结算完成
func SettleOver(flow interfaces.MajongFlow, message *majongpb.SettleFinishEvent) {
	logEntry := logrus.WithFields(
		logrus.Fields{
			"func_name": "settleOver",
		},
	)
	mjContext := flow.GetMajongContext()
	playerIds := message.GetPlayerId()
	for _, pid := range playerIds {
		player := utils.GetMajongPlayer(pid, mjContext)
		if player == nil {
			logEntry.WithField("playerID: ", pid).Errorln("找不到玩家")
			continue
		}
		player.XpState = majongpb.XingPaiState_give_up
	}
}

// GetNextState 下一状态获取
func GetNextState(mjContext *majongpb.MajongContext) majongpb.StateID {
	// 正常玩家<=1,游戏结束
	if utils.IsNormalPlayerInsufficient(mjContext.GetPlayers()) {
		return majongpb.StateID_state_gameover
	}
	return majongpb.StateID_state_mopai
}
