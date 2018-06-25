package utils

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// IsPlayerContinue   玩家的状态在麻将可行牌数组中包含则返回true
func IsPlayerContinue(playerStater majongpb.XingPaiState, option *majongpb.MajongCommonOption) bool {
	// 麻将可行牌数组
	xpStates := option.GetValidXpStateSet()
	flag := false
	for _, state := range xpStates {
		if state == playerStater {
			flag = true
			break
		}
	}
	logrus.WithFields(logrus.Fields{"playerStater": playerStater,
		"xpStates": xpStates, "isCanXp": flag}).Info("判断玩家是否可以继续")
	return flag
}

//GetNextXpPlayerByID 获取下一个行牌玩家
func GetNextXpPlayerByID(srcPlayerID uint64, players []*majongpb.Player, option *majongpb.MajongCommonOption) (nextPalyer *majongpb.Player) {
	curPlayerID, i := srcPlayerID, 0
	for i < 4 {
		nextPalyer = GetNextPlayerByID(players, curPlayerID)
		// 当前下个玩家可以继续，退出循环
		if IsPlayerContinue(nextPalyer.GetXpState(), option) {
			break
		}
		curPlayerID = nextPalyer.GetPalyerId()
	}
	logrus.WithFields(logrus.Fields{"playerID": nextPalyer.GetPalyerId(),
		"playerStatus": nextPalyer.GetXpState()}).Info("获取下个正常状态的玩家")
	return nextPalyer
}

//GetXpPlayers 获取行牌玩家数组
func GetXpPlayers(players []*majongpb.Player, option *majongpb.MajongCommonOption) []*majongpb.Player {
	newPlalyers := make([]*majongpb.Player, 0)
	for _, player := range players {
		// 不是正常行牌的玩家，不能检查胡，碰，杠，摸牌。。。
		if !IsPlayerContinue(player.GetXpState(), option) {
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
func IsXpPlayerInsufficient(players []*majongpb.Player, option *majongpb.MajongCommonOption) bool {
	conut := 0
	for _, player := range players {
		if IsPlayerContinue(player.GetXpState(), option) {
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
		player.XpState = majongpb.XingPaiState_give_up
	}
}

// GetNextState 下一状态获取
func GetNextState(mjContext *majongpb.MajongContext) majongpb.StateID {
	// 正常玩家<=1,游戏结束
	if IsXpPlayerInsufficient(mjContext.GetPlayers(), mjContext.GetOption()) {
		return majongpb.StateID_state_gameover
	}
	return majongpb.StateID_state_mopai
}
