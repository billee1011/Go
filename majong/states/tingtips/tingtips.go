package tingtips

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/cardtype"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// calcHuTimes 计算胡牌倍数
func calcHuTimes(card utils.Card, player *majongpb.Player, gameID int) uint32 {
	calcor := &cardtype.ScxlCardTypeCalculator{}
	pengCards := []*majongpb.Card{}
	gangCards := []*majongpb.Card{}
	for _, pcard := range player.GetPengCards() {
		pengCards = append(pengCards, pcard.GetCard())
	}
	for _, gcard := range player.GetGangCards() {
		gangCards = append(gangCards, gcard.GetCard())
	}
	huCard, _ := utils.IntToCard(int32(card))

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
	player := utils.GetMajongPlayer(playerID, mjContext)
	playerCards := player.GetHandCards()

	tingCards := utils.FastCheckTingV2(utils.CardsToUtilCards(playerCards), map[utils.Card]bool{}) // TODO, 目前没有包括特殊牌型

	ntf := room.RoomTingInfoNtf{}
	for _, card := range tingCards {
		times := calcHuTimes(card, player, int(mjContext.GetGameId()))
		tingCardInfo := &room.TingCardInfo{
			TingCard: proto.Uint32(uint32(card)),
			Times:    proto.Uint32(times),
		}
		ntf.TingCardInfos = append(ntf.TingCardInfos, tingCardInfo)
	}
	flow.PushMessages([]uint64{playerID}, interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_TINGINFO_NTF),
		Msg:   &ntf,
	})
}
