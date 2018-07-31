package utils

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/gutils"
	"steve/room/majong/interfaces"
	majongpb "steve/entity/majong"

	"steve/room/majong/bus"

	"github.com/golang/protobuf/proto"
)

// calcHuTimes 计算胡牌倍数
func calcHuTimes(card *majongpb.Card, player *majongpb.Player, mjContext *majongpb.MajongContext) uint32 {
	calcor := bus.GetFanTypeCalculator()
	pengCards := []*majongpb.Card{}
	gangCards := []*majongpb.GangCard{}
	for _, pcard := range player.GetPengCards() {
		pengCards = append(pengCards, pcard.GetCard())
	}
	for _, gcard := range player.GetGangCards() {
		gangCards = append(gangCards, gcard)
	}

	params := interfaces.FantypeParams{
		PlayerID:  player.GetPalyerId(),
		MjContext: mjContext,
		HandCard:  player.GetHandCards(),
		PengCard:  pengCards,
		GangCard:  gangCards,
		HuCard: &majongpb.HuCard{
			Card: card,
			Type: majongpb.HuType_hu_dianpao,
		},
	}
	types, gen, hua := calcor.Calculate(params)
	types = gutils.DeleteHuType(int(mjContext.GetCardtypeOptionId()), types) // 移除胡类型的番型
	value := calcor.CardTypeValue(mjContext, types, gen, hua)
	return uint32(value)
}

// NotifyTingCards 通知玩家当前听的牌
func NotifyTingCards(flow interfaces.MajongFlow, playerID uint64) {
	mjContext := flow.GetMajongContext()
	player := GetMajongPlayer(playerID, mjContext)
	playerCards := player.GetHandCards()
	//清除上一次听牌记录
	player.TingCardInfo = nil
	// 不存在定缺牌
	if gutils.CheckHasDingQueCard(mjContext, player) {
		return
	}
	// wuhongwei 增加七对提示
	tingCards, _ := GetTingCards(playerCards, nil) // TODO, 目前没有包括特殊牌型
	ntf := room.RoomTingInfoNtf{}
	for _, utilscard := range tingCards {
		card, _ := IntToCard(int32(utilscard))
		// 胡提示不能是定缺牌
		if mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId())).EnableDingque && card.GetColor() == player.GetDingqueColor() {
			continue
		}
		newCard := ServerCard2Uint32(card)
		times := calcHuTimes(card, player, mjContext)
		tingCardInfo := &room.TingCardInfo{
			TingCard: proto.Uint32(newCard),
			Times:    proto.Uint32(times),
		}
		ntf.TingCardInfos = append(ntf.TingCardInfos, tingCardInfo)
		// 记录听牌信息
		mjTingInfo := &majongpb.TingCardInfo{
			TingCard: newCard,
			Times:    times,
		}
		player.TingCardInfo = append(player.TingCardInfo, mjTingInfo)
	}
	flow.PushMessages([]uint64{playerID}, interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_TINGINFO_NTF),
		Msg:   &ntf,
	})
}
