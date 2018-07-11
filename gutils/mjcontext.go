package gutils

import (
	"steve/common/mjoption"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// 牌对应都值
const (
	//W1 1万
	W1 = 11
	//W2 2万
	W2 = 12
	//W3 3万
	W3 = 13
	//W4 4万
	W4 = 14
	//W5 5万
	W5 = 15
	//W6 6万
	W6 = 16
	//W7 7万
	W7 = 17
	//W8 8万
	W8 = 18
	//W9 9万
	W9 = 19

	//T1 1条
	T1 = 21
	//T2 2条
	T2 = 22
	//T3 3条
	T3 = 23
	//T4 4条
	T4 = 24
	//T5 5条
	T5 = 25
	//T6 6Xi
	T6 = 26
	//T7 7Xi
	T7 = 27
	//T8 8条
	T8 = 28
	//T9 9条
	T9 = 29

	//B1 1筒
	B1 = 31
	//B2 2筒
	B2 = 32
	//B3 3筒
	B3 = 33
	//B4 4筒
	B4 = 34
	//B5 5筒
	B5 = 35
	//B6 6筒
	B6 = 36
	//B7 7筒
	B7 = 37
	//B8 8筒
	B8 = 38
	//B9 9筒
	B9 = 39

	//Dong 东风
	Dong = 41
	//Nan 南风
	Nan = 42
	//Xi 西风
	Xi = 43
	//Bei 北风
	Bei = 44
	//Zhong 红中
	Zhong = 45
	//Fa 发财
	Fa = 46
	//Bai 白板
	Bai = 47
)

// GetMajongPlayer 从 MajongContext 中根据玩家 ID 获取玩家
func GetMajongPlayer(playerID uint64, mjContext *majongpb.MajongContext) *majongpb.Player {
	for _, player := range mjContext.GetPlayers() {
		if player.GetPalyerId() == playerID {
			return player
		}
	}
	return nil
}

// GetPlayerIndex 获取玩家索引
func GetPlayerIndex(playerID uint64, players []*majongpb.Player) int {
	for index, player := range players {
		if player.GetPalyerId() == playerID {
			return index
		}
	}
	return -1
}

// GetPlayerAndIndex 获取玩家索引
func GetPlayerAndIndex(playerID uint64, players []*majongpb.Player) (int, *majongpb.Player) {
	for index, player := range players {
		if player.GetPalyerId() == playerID {
			return index, player
		}
	}
	return -1, nil
}

// IsPlayerContinue   玩家的状态在麻将不可行牌数组中包含则返回false
func IsPlayerContinue(playerState majongpb.XingPaiState, mjContext *majongpb.MajongContext) bool {
	// 麻将不可行牌数组
	xpOption := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId()))
	flag := xpOption.PlayerNoNormalStates&int32(playerState) == 0
	logrus.WithFields(logrus.Fields{
		"playerStater":   playerState,
		"canNotXpStates": xpOption.PlayerNoNormalStates,
		"isCanXp":        flag,
	}).Info("判断玩家是否可以继续")
	return flag
}

// GetPlayerSeat 获取玩家所在的座位
func GetPlayerSeat(renNum, playerIndex int) int {
	indexs := map[int][]int{
		2: []int{Dong, Xi},           // 东西
		3: []int{Dong, Nan, Xi},      // 东南西
		4: []int{Dong, Nan, Xi, Bei}, // 东南西北
	}[renNum]
	return indexs[playerIndex]
}
