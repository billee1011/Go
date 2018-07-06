package utils

import (
	"steve/gutils"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

//GetNextXpPlayerByID 获取下一个行牌玩家
func GetNextXpPlayerByID(srcPlayerID uint64, players []*majongpb.Player, mjContext *majongpb.MajongContext) (nextPalyer *majongpb.Player) {
	curPlayerID, i := srcPlayerID, 0
	for i < len(mjContext.Players) {
		nextPalyer = GetNextPlayerByID(players, curPlayerID)
		// 当前下个玩家可以继续，退出循环
		if gutils.IsPlayerContinue(nextPalyer.GetXpState(), mjContext) {
			break
		}
		curPlayerID = nextPalyer.GetPalyerId()
	}
	logrus.WithFields(logrus.Fields{"playerID": nextPalyer.GetPalyerId(),
		"playerStatus": nextPalyer.GetXpState()}).Info("获取下个正常状态的玩家")
	return nextPalyer
}

//GetXpPlayers 获取行牌玩家数组
func GetXpPlayers(players []*majongpb.Player, mjContext *majongpb.MajongContext) []*majongpb.Player {
	newPlalyers := make([]*majongpb.Player, 0)
	for _, player := range players {
		// 不是正常行牌的玩家，不能检查胡，碰，杠，摸牌。。。
		if !gutils.IsPlayerContinue(player.GetXpState(), mjContext) {
			continue
		}
		newPlalyers = append(newPlalyers, player)
	}
	// 正常的玩家ID
	playerID := make([]uint64, 0)
	for _, player := range newPlalyers {
		playerID = append(playerID, player.GetPalyerId())
	}
	logrus.WithFields(logrus.Fields{"PlayerIDs": playerID}).Info("获取正常玩家数组")
	return newPlalyers
}

//IsXpPlayerInsufficient 判断能行牌人数是否足够
func IsXpPlayerInsufficient(players []*majongpb.Player, mjContext *majongpb.MajongContext) bool {
	conut := 0
	for _, player := range players {
		if gutils.IsPlayerContinue(player.GetXpState(), mjContext) {
			conut++
		}
	}
	logrus.WithFields(logrus.Fields{"NormalPlayerConut": conut}).Infoln("正常状态玩家数量")
	return conut <= 1
}

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
		player := GetMajongPlayer(pid, mjContext)
		if player == nil {
			logEntry.WithField("playerID: ", pid).Errorln("找不到玩家")
			continue
		}
		player.XpState = player.GetXpState() | majongpb.XingPaiState_give_up
	}

}

// GetNextState 下一状态获取
func GetNextState(mjContext *majongpb.MajongContext) majongpb.StateID {
	// 正常玩家<=1,游戏结束
	if IsXpPlayerInsufficient(mjContext.GetPlayers(), mjContext) {
		return majongpb.StateID_state_gameover
	}
	return majongpb.StateID_state_mopai
}
