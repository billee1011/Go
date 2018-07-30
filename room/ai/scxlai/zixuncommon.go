package scxlai

import (
	"steve/common/mjoption"
	"steve/gutils"
	"steve/room/interfaces"
	"steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

func (h *zixunStateAI) getNormalZiXunAIEvent(player *majong.Player, mjContext *majong.MajongContext) (aiEvent interfaces.AIEvent) {
	zxRecord := player.GetZixunRecord()
	handCards := player.GetHandCards()
	canHu := zxRecord.GetEnableZimo()
	if (gutils.IsHu(player) || gutils.IsTing(player)) && canHu {
		aiEvent = h.hu(player)
		return
	}
	// 优先出定缺牌
	if gutils.CheckHasDingQueCard(mjContext, player) {
		for i := len(handCards) - 1; i >= 0; i-- {
			hc := handCards[i]
			if mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId())).EnableDingque &&
				hc.GetColor() == player.GetDingqueColor() {
				aiEvent = h.chupai(player, hc)
				return
			}
		}
	}

	// 正常出牌
	if player.GetMopaiCount() == 0 {
		aiEvent = h.chupai(player, handCards[len(handCards)-1])
	} else {
		aiEvent = h.chupai(player, mjContext.GetLastMopaiCard())
	}
	return
}

func (h *zixunStateAI) getNormalZiXunTingStateAIEvent(player *majong.Player, mjContext *majong.MajongContext) (aiEvent interfaces.AIEvent) {
	zxRecord := player.GetZixunRecord()
	if gutils.IsTing(player) {
		canHu := zxRecord.GetEnableZimo()
		if !canHu {
			aiEvent = h.chupai(player, mjContext.GetLastMopaiCard())
		}
	}
	return
}

func (h *zixunStateAI) getNormalZiXunHuStateAIEvent(player *majong.Player, mjContext *majong.MajongContext) (aiEvent interfaces.AIEvent) {
	zxRecord := player.GetZixunRecord()
	if gutils.IsHu(player) {
		canHu := zxRecord.GetEnableZimo()
		if canHu {
			aiEvent = h.hu(player)
		} else {
			aiEvent = h.chupai(player, mjContext.GetLastMopaiCard())
		}
	}
	return
}

func (h *zixunStateAI) chupai(player *majong.Player, card *majong.Card) interfaces.AIEvent {
	eventContext := majong.ChupaiRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
		Cards: card,
	}
	data, err := proto.Marshal(&eventContext)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func_name": "zixunStateAI.chupai",
			"player_id": player.GetPalyerId(),
		}).Errorln("事件序列化失败")
	}
	return interfaces.AIEvent{
		ID:      int32(majong.EventID_event_chupai_request),
		Context: data,
	}
}

func (h *zixunStateAI) hu(player *majong.Player) interfaces.AIEvent {
	eventContext := majong.HuRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
	}
	data, err := proto.Marshal(&eventContext)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func_name": "zixunStateAI.hu",
			"player_id": player.GetPalyerId(),
		}).Errorln("事件序列化失败")
	}
	return interfaces.AIEvent{
		ID:      int32(majong.EventID_event_hu_request),
		Context: data,
	}
}
