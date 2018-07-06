package utils

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// calcHuTimes 计算胡牌倍数
func calcHuTimes(card Card, player *majongpb.Player, gameID int) uint32 {
	calcor := global.GetCardTypeCalculator()
	pengCards := []*majongpb.Card{}
	gangCards := []*majongpb.Card{}
	for _, pcard := range player.GetPengCards() {
		pengCards = append(pengCards, pcard.GetCard())
	}
	for _, gcard := range player.GetGangCards() {
		gangCards = append(gangCards, gcard.GetCard())
	}
	huCard, _ := IntToCard(int32(card))

	params := interfaces.CardCalcParams{
		HandCard: player.GetHandCards(),
		PengCard: pengCards,
		GangCard: gangCards,
		HuCard:   huCard,
		GameID:   gameID,
	}
	value, _ := facade.CalculateCardValue(calcor, params)
	return value
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
			times := calcHuTimes(Card(*newCard), player, int(mjContext.GetGameId()))
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
		MsgID: int(msgid.MsgID_ROOM_TINGINFO_NTF),
		Msg:   &ntf,
	})
}
