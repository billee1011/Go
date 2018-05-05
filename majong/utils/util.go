package utils

import (
	"fmt"
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

//CardsToUtilCards 用来辅助查ting胡的工具,将Card转为适合查胡的工具
func CardsToUtilCards(cards []*majongpb.Card) []Card {
	cardsCard := make([]Card, 0)
	for _, v := range cards {
		cardInt, _ := CardToInt(*v)
		cardsCard = append(cardsCard, Card(*cardInt))
	}
	return cardsCard
}

//HuCardsToUtilCards 用来辅助查ting胡的工具,将Card转为适合查胡的工具
func HuCardsToUtilCards(cards []*majongpb.HuCard) []Card {
	cardsCard := make([]Card, 0)
	for _, v := range cards {

		cardInt, _ := CardToInt(*v.Card)
		cardsCard = append(cardsCard, Card(*cardInt))
	}
	return cardsCard
}

//CardToInt 将Card转换成牌值（int32）
func CardToInt(card majongpb.Card) (*int32, error) {
	var color int32
	switch card.GetColor() {
	case majongpb.CardColor_ColorWan:
		color = 1
	case majongpb.CardColor_ColorTiao:
		color = 2
	case majongpb.CardColor_ColorTong:
		color = 3
	default:
		return &color, fmt.Errorf("cant trans card ")
	}
	tValue := card.Point
	value := color*10 + tValue
	return &value, nil
}

//CardsToInt 将Card转换成牌值（int32）
func CardsToInt(cards []*majongpb.Card) ([]int32, error) {
	var cardsInt []int32
	var color int32
	for _, card := range cards {
		switch card.GetColor() {
		case majongpb.CardColor_ColorWan:
			color = 1
		case majongpb.CardColor_ColorTiao:
			color = 2
		case majongpb.CardColor_ColorTong:
			color = 3
		default:
			return nil, fmt.Errorf("cant trans card ")
		}
		tValue := card.Point
		value := color*10 + tValue
		cardsInt = append(cardsInt, value)
	}
	return cardsInt, nil
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

//IntToCard int32类型转majongpb.Card类型
func IntToCard(cardValue int32) (*majongpb.Card, error) {
	colorValue := cardValue / 10
	value := cardValue % 10
	var color majongpb.CardColor
	switch colorValue {
	case 1:
		color = majongpb.CardColor_ColorWan
	case 2:
		color = majongpb.CardColor_ColorTiao
	case 3:
		color = majongpb.CardColor_ColorTong
	default:
		return nil, fmt.Errorf("cant trans card %d", cardValue)
	}
	return &majongpb.Card{
		Color: color,
		Point: value,
	}, nil
}

//IntToRoomCard int32类型转room.Card类型
func IntToRoomCard(cardValue int32) (*room.Card, error) {
	colorValue := cardValue / 10
	value := cardValue % 10
	var color room.CardColor
	switch colorValue {
	case 1:
		color = room.CardColor_ColorWan
	case 2:
		color = room.CardColor_ColorTiao
	case 3:
		color = room.CardColor_ColorTong
	default:
		return nil, fmt.Errorf("cant trans card %d", cardValue)
	}
	return &room.Card{
		Color: color.Enum(),
		Point: proto.Int32(value),
	}, nil
}

//CardToRoomCard majongpb.card类型转room.Card类型
func CardToRoomCard(card *majongpb.Card) (*room.Card, error) {
	return &room.Card{
		Color: room.CardColor(card.Color).Enum(),
		Point: proto.Int32(card.Point),
	}, nil
}

//IntToUtilCard uint32类型的数组强转成类型
func IntToUtilCard(cards []int32) []Card {
	cardsCard := make([]Card, 0, 0)
	for _, v := range cards {

		utilCard := Card(v)
		cardsCard = append(cardsCard, utilCard)
	}
	return cardsCard
}

//ContainHuCards 判断当前可以胡的牌中是否包含已经胡过的所有牌
func ContainHuCards(targetHuCards []Card, HuCards []Card) bool {
	sameHuCards := 0
	for _, tagetCard := range targetHuCards {
		for _, Card := range HuCards {
			if tagetCard == Card {
				sameHuCards++
			}
		}
	}
	if len(HuCards) == sameHuCards {
		return true
	}
	return false
}

//CheckHu 用来辅助胡牌查胡工具 cards玩家的所有牌，huCard点炮的牌（自摸时huCard为0）
func CheckHu(cards []*majongpb.Card, huCard uint32) bool {
	cardsCard := CardsToUtilCards(cards)
	if huCard > 0 {
		cardsCard = append(cardsCard, Card(huCard))
	}
	// flag, _ := util.FastCheckHuV1(cardsCard) // 检测玩家能否推倒胡
	laizi := make(map[Card]bool)
	flag := FastCheckHuV2(cardsCard, laizi) // 检测玩家能否推倒胡
	if flag != true {
		flag = FastCheckQiDuiHu(cardsCard) // 检测玩家能否七对胡
	}
	return flag
}
