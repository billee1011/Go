package utils

import (
	"steve/client_pb/room"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

//GetPlayerByID 根据玩家id获取玩家
func GetPlayerByID(players []*majongpb.Player, id uint64) *majongpb.Player {
	for _, player := range players {
		if player.PalyerId == id {
			return player
		}
	}
	return nil
}

//GetNextPlayerByID 根据玩家id获取下个玩家
func GetNextPlayerByID(players []*majongpb.Player, id uint64) *majongpb.Player {
	for k, player := range players {
		if player.PalyerId == id {
			index := (k + 1) % len(players)
			return players[index]
		}
	}
	return nil
}

//CheckHasDingQueCard 检查牌里面是否含有定缺的牌
func CheckHasDingQueCard(cards []*majongpb.Card, color majongpb.CardColor) bool {
	for _, card := range cards {
		if card.Color == color {
			return true
		}
	}
	return false
}

//DeleteIntCardFromLast 从int32类型的牌组中，找到对应的目标牌，并且移除
func DeleteIntCardFromLast(cards []int32, targetCard int32) ([]int32, bool) {
	index := -1
	l := len(cards)
	if l == 0 {
		return cards, false
	}
	for i := l - 1; i >= 0; i-- {
		if targetCard == cards[i] {
			index = i
			break
		}
	}
	if index != -1 {
		cards = append(cards[:index], cards[index+1:]...)
	}
	return cards, index != -1
}

//CardEqual 判断两张牌是否一样
func CardEqual(card1 *majongpb.Card, card2 *majongpb.Card) bool {
	return card1.GetColor() == card2.GetColor() && card1.GetPoint() == card2.GetPoint()
}

//DeleteCardFromLast 从majongpb.Card类型的牌组中，找到对应的目标牌，并且移除
func DeleteCardFromLast(cards []*majongpb.Card, targetCard *majongpb.Card) ([]*majongpb.Card, bool) {
	index := -1
	l := len(cards)
	if l == 0 {
		return cards, false
	}
	for i := l - 1; i >= 0; i-- {
		if CardEqual(targetCard, cards[i]) {
			index = i
			break
		}
	}
	if index != -1 {
		cards = append(cards[:index], cards[index+1:]...)
	}
	return cards, index != -1
}

//CardToRoomCard majongpb.card类型转room.Card类型
func CardToRoomCard(card *majongpb.Card) (*room.Card, error) {
	var color room.CardColor
	if card.Color.String() == room.CardColor_CC_WAN.String() {
		color = room.CardColor_CC_WAN
	}
	if card.Color.String() == room.CardColor_CC_TIAO.String() {
		color = room.CardColor_CC_TIAO
	}
	if card.Color.String() == room.CardColor_CC_TONG.String() {
		color = room.CardColor_CC_TONG
	}

	return &room.Card{
		Color: color.Enum(),
		Point: proto.Int32(card.Point),
	}, nil
}

// ServerCard2Number 服务器的 Card 转换成数字
func ServerCard2Number(card *majongpb.Card) int {
	var color int
	if card.Color == majongpb.CardColor_ColorWan {
		color = 1
	} else if card.Color == majongpb.CardColor_ColorTiao {
		color = 2
	} else if card.Color == majongpb.CardColor_ColorTong {
		color = 3
	} else if card.Color == majongpb.CardColor_ColorFeng {
		color = 4
	}
	value := color*10 + int(card.Point)
	return value
}

// ServerCards2Numbers 服务器的 Card 数组转 int 数组
func ServerCards2Numbers(cards []*majongpb.Card) []int {
	result := []int{}
	for _, c := range cards {
		result = append(result, ServerCard2Number(c))
	}
	return result
}

// CardsToRoomCards 将Card转换为room package中的Card
func CardsToRoomCards(cards []*majongpb.Card) []*room.Card {
	var rCards []*room.Card
	for i := 0; i < len(cards); i++ {
		rCards = append(rCards, &room.Card{
			Color: room.CardColor(cards[i].Color).Enum(),
			Point: &cards[i].Point,
		})
	}
	return rCards
}

// ContainCard 验证card是否存在于玩家手牌中，存在返回true,否则返回false
func ContainCard(cards []*majongpb.Card, card *majongpb.Card) bool {
	for i := 0; i < len(cards); i++ {
		if CardEqual(cards[i], card) {
			return true
		}
	}
	return false
}
