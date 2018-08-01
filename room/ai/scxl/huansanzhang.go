package scxlai

import (
	"steve/client_pb/room"
	"steve/gutils"
	"steve/server_pb/majong"

	"steve/room/ai"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type huansanzhangStateAI struct {
}

// GenerateAIEvent 生成 换三张AI 事件
// 无论是超时、托管还是机器人，若已存在换三张的牌，则直接换该三张牌，否则取花色最少的三张手牌换三张， 并且产生相应的事件
func (h *huansanzhangStateAI) GenerateAIEvent(params ai.AIEventGenerateParams) (result ai.AIEventGenerateResult, err error) {
	result, err = ai.AIEventGenerateResult{
		Events: []ai.AIEvent{},
	}, nil

	mjContext := params.MajongContext
	player := gutils.GetMajongPlayer(params.PlayerID, mjContext)
	if !player.GetHuansanzhangSure() { // 玩家没选换牌情况下,生成换牌请求
		if event := h.huansanzhang(player); event != nil {
			result.Events = append(result.Events, *event)
		}
		return
	}
	if params.AIType != ai.RobotAI { //不是机器人不能发送动画完成请求
		return
	}
	finished := mjContext.GetExcutedHuansanzhang() // 是否所有人已经执行换牌
	if finished {
		crPlayerIDs := mjContext.GetTempData().GetCartoonReqPlayerIDs()
		if len(crPlayerIDs) == len(mjContext.GetPlayers()) { //所有玩家都发送过动画完成请求
			return
		}
		for _, playerID := range crPlayerIDs {
			if playerID == player.GetPalyerId() { // 当前玩家已经发送过
				return
			}
		}
		// 发送动画完成请求
		if event := CartoonFinsh(player, int32(room.CartoonType_CTNT_HUANSANZHANG)); event != nil {
			result.Events = append(result.Events, *event)
		}
	}
	return
}

// getHszCards 获取换三张的牌
func (h *huansanzhangStateAI) getHszCards(player *majong.Player) (hszCards []*majong.Card) {
	if len(player.GetHuansanzhangCards()) == 3 { //超时获取存的换三张
		return player.GetHuansanzhangCards()
	}
	player.HuansanzhangCards = gutils.GetRecommedHuanSanZhang(player.GetHandCards())
	hszCards = player.GetHuansanzhangCards()
	logrus.WithFields(logrus.Fields{"player_id": player.GetPalyerId(), "hszCards": gutils.FmtMajongpbCards(hszCards)}).Infoln("服务器推荐换三张")
	return hszCards
}

// huansanzhang 生成换三张请求事件
func (h *huansanzhangStateAI) huansanzhang(player *majong.Player) *ai.AIEvent {
	hszCards := h.getHszCards(player)

	mjContext := majong.HuansanzhangRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
		Cards: hszCards,
		Sure:  true,
	}

	data, err := proto.Marshal(&mjContext)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func_name": "huansanzhangStateAI.huansanzhang",
			"player_id": player.GetPalyerId(),
			"hszCards":  hszCards,
		}).Errorln("事件序列化失败")
		return nil
	}
	return &ai.AIEvent{
		ID:      int32(majong.EventID_event_huansanzhang_request),
		Context: data,
	}
}
