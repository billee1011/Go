package common

import (
	"steve/client_pb/room"
	majongpb "steve/entity/majong"

	"steve/room/majong/utils"

	"github.com/Sirupsen/logrus"
)

//CartoonFinishData 动画完成数据
type CartoonFinishData struct {
	CurState        majongpb.StateID
	NextState       majongpb.StateID
	NeedCartoonType room.CartoonType
	EventContext    interface{}
}

// OnCartoonFinish 在某个状态上， 动画播放完成
func OnCartoonFinish(cartoonFinishData CartoonFinishData, mjContext *majongpb.MajongContext) (newState majongpb.StateID, err error) {
	curState := cartoonFinishData.CurState
	nextState := cartoonFinishData.NextState
	needCartoonType := cartoonFinishData.NeedCartoonType
	req := (cartoonFinishData.EventContext).(*majongpb.CartoonFinishRequestEvent)
	reqPlayerID := req.GetPlayerId()
	cartoonReqPlayerIDs := AddCartoonPlayerID(mjContext, reqPlayerID) //玩家发送请求，添加到临时数据中
	if len(cartoonReqPlayerIDs) < len(mjContext.GetPlayers()) {
		return curState, nil
	}
	reqCartoonType := req.GetCartoonType()
	if reqCartoonType != int32(needCartoonType) {
		return curState, nil
	}
	logrus.WithFields(logrus.Fields{"req_cartoon_type": "reqCartoonType",
		"cur_state":   curState,
		"next_state":  nextState,
		"cfPlayerIDs": cartoonReqPlayerIDs, "req": req}).Infoln("动画结束请求信息")
	return nextState, nil
}

// AddHuCard 添加胡的牌
func AddHuCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64, huType majongpb.HuType, isReal bool) {
	player.HuCards = append(player.GetHuCards(), &majongpb.HuCard{
		Card:      card,
		Type:      huType,
		SrcPlayer: srcPlayerID,
		IsReal:    isReal,
	})
}

// calcLastHuIndex 计算一次胡操作中，最后一个玩家索引
// startPlayer 从哪个索引开始算起
// players 哪些玩家胡了
// totalCount 总玩家数量
// return 最后一个胡的玩家索引
func calcLastHuIndex(startPlayer int, players []int, totalCount int) int {
	if len(players) == 0 {
		panic("胡的玩家为空")
	}
	maxStepPlayer := players[0]
	maxStep := calcStep(startPlayer, players[0], totalCount)

	for i := 1; i < len(players); i++ {
		step := calcStep(startPlayer, players[i], totalCount)
		if step > maxStep {
			maxStep = step
			maxStepPlayer = players[i]
		}
	}
	return maxStepPlayer
}

// calcStep 计算从 src 到 dest 的步骤数
func calcStep(src int, dest int, total int) int {
	if dest < src {
		return dest + total - src
	}
	return dest - src
}

// CalcMopaiPlayer 计算摸牌玩家 ID
// huPlayers 胡的玩家 ID 列表
// srcPlayer 原玩家
// players 全部玩家
func CalcMopaiPlayer(logEntry *logrus.Entry, huPlayers []uint64, srcPlayer uint64, players []*majongpb.Player) uint64 {
	huIndexs := []int{}
	for _, player := range huPlayers {
		index, err := utils.GetPlayerIndex(player, players)
		if err != nil {
			logEntry.WithError(err).Errorln("获取胡玩家索引失败")
			return srcPlayer
		}
		huIndexs = append(huIndexs, index)
	}
	srcIndex, err := utils.GetPlayerIndex(srcPlayer, players)
	if err != nil {
		logEntry.WithError(err).Errorln("获取源玩家索引失败")
		return srcPlayer
	}
	mopaiIndex := (calcLastHuIndex(srcIndex, huIndexs, len(players)) + 1) % len(players)
	return players[mopaiIndex].GetPlayerId()
}

func removeLastCard(logEntry *logrus.Entry, srcCards []*majongpb.Card, rmCard *majongpb.Card) []*majongpb.Card {
	pos := len(srcCards) - 1
	if pos >= 0 && (srcCards[pos].GetColor() == rmCard.GetColor()) &&
		(srcCards[pos].GetPoint() == rmCard.GetPoint()) {
		srcCards = srcCards[0:pos]
	} else {
		logEntry.Errorln("最后一张牌与目标牌不同，无法移除")
	}
	return srcCards
}

// AddCartoonPlayerID 添加动画完成请求玩家
func AddCartoonPlayerID(mjContext *majongpb.MajongContext, currPlayerID uint64) []uint64 {
	cfPlayerIDs := mjContext.GetTempData().GetCartoonReqPlayerIDs()
	for _, playerID := range cfPlayerIDs {
		if playerID == currPlayerID { // 用于判断玩家是否发过动画完成请求
			return cfPlayerIDs
		}
	}
	mjContext.GetTempData().CartoonReqPlayerIDs = append(cfPlayerIDs, currPlayerID)
	return mjContext.GetTempData().GetCartoonReqPlayerIDs()
}
