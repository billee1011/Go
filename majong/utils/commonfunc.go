package utils

import (
	"steve/common/mjoption"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// IsPlayerContinue   玩家的状态在麻将不可行牌数组中包含则返回false
func IsPlayerContinue(playerStater majongpb.XingPaiState, mjContext *majongpb.MajongContext) bool {
	// 麻将不可行牌数组
	xpOption := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId()))
	flag := xpOption.PlayerNoNormalStates&uint64(playerStater) == 0
	logrus.WithFields(logrus.Fields{"playerStater": playerStater,
		"canNotXpStates": xpOption.PlayerNoNormalStates, "isCanXp": flag}).Info("判断玩家是否可以继续")
	return flag
}

//GetNextXpPlayerByID 获取下一个行牌玩家
func GetNextXpPlayerByID(srcPlayerID uint64, players []*majongpb.Player, mjContext *majongpb.MajongContext) (nextPalyer *majongpb.Player) {
	curPlayerID, i := srcPlayerID, 0
	for i < len(mjContext.Players) {
		nextPalyer = GetNextPlayerByID(players, curPlayerID)
		// 当前下个玩家可以继续，退出循环
		if IsPlayerContinue(nextPalyer.GetXpState(), mjContext) {
			break
		}
		curPlayerID = nextPalyer.GetPalyerId()
	}
	logrus.WithFields(logrus.Fields{"playerID": nextPalyer.GetPalyerId(),
		"playerStatus": nextPalyer.GetXpState()}).Info("获取下个正常状态的玩家")
	return nextPalyer
}

//GetCanXpPlayers 获取能行牌玩家数组
func GetCanXpPlayers(players []*majongpb.Player, mjContext *majongpb.MajongContext) []*majongpb.Player {
	newPlalyers := make([]*majongpb.Player, 0)
	for _, player := range players {
		// 不是正常行牌的玩家，不能检查胡，碰，杠，摸牌。。。
		if !IsPlayerContinue(player.GetXpState(), mjContext) {
			logrus.WithFields(logrus.Fields{"PlayerIDs": player.GetPalyerId(), "PlayerState": player.GetXpState()}).Info("不正常玩家")
			continue
		}
		newPlalyers = append(newPlalyers, player)
	}
	return newPlalyers
}

//IsAreThereEnoughpeople 判断是否有足够多的人数
func IsAreThereEnoughpeople(players []*majongpb.Player, mjContext *majongpb.MajongContext) bool {
	count := 0
	for _, player := range players {
		if IsPlayerContinue(player.GetXpState(), mjContext) {
			count++
		}
	}
	logrus.WithFields(logrus.Fields{"NormalPlayerConut": count}).Infoln("正常状态玩家数量")
	return count <= 1
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

// IsGameOverReturnState 判断游戏是否结束返回状态
func IsGameOverReturnState(mjContext *majongpb.MajongContext) majongpb.StateID {
	// 正常玩家<=1,游戏结束
	if IsAreThereEnoughpeople(mjContext.GetPlayers(), mjContext) {
		return majongpb.StateID_state_gameover
	}
	return majongpb.StateID_state_mopai
}
