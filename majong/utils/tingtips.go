package utils

import (
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// calcHuTimes 计算胡牌倍数
func calcHuTimes(card *majongpb.Card, player *majongpb.Player, mjContext *majongpb.MajongContext) uint32 {
	calcor := global.GetFanTypeCalculator()
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
	value, _, _ := facade.CalculateCardValue(calcor, mjContext, params)
	return uint32(value)
}

// NotifyTingCards 通知玩家当前听的牌
func NotifyTingCards(flow interfaces.MajongFlow, playerID uint64) {
	mjContext := flow.GetMajongContext()
	player := GetMajongPlayer(playerID, mjContext)
	playerCards := player.GetHandCards()

	// 不存在定缺牌
	if gutils.CheckHasDingQueCard(playerCards, player.GetDingqueColor()) {
		return
	}
	// wuhongwei 增加七对提示
	tingCards, _ := GetTingCards(playerCards, nil) // TODO, 目前没有包括特殊牌型

	ntf := room.RoomTingInfoNtf{}
	for _, utilscard := range tingCards {
		card, _ := IntToCard(int32(utilscard))
		// 胡提示不能是定缺牌
		if card.GetColor() != player.GetDingqueColor() {
			newCard, _ := CardToInt(*card)
			times := calcHuTimes(card, player, mjContext)
			tingCardInfo := &room.TingCardInfo{
				TingCard: proto.Uint32(uint32(*newCard)),
				Times:    proto.Uint32(times),
			}
			ntf.TingCardInfos = append(ntf.TingCardInfos, tingCardInfo)
			// 记录听牌信息
			mjTingInfo := &majongpb.TingCardInfo{
				TingCard: uint32(*newCard),
				Times:    times,
			}
			player.TingCardInfo = append(player.TingCardInfo, mjTingInfo)
		}
	}
	flow.PushMessages([]uint64{playerID}, interfaces.ToClientMessage{
		MsgID: int(msgId.MsgID_ROOM_TINGINFO_NTF),
		Msg:   &ntf,
	})
}
