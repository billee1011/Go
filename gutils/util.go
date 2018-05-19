package gutils

import (
	"steve/client_pb/room"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// RoomCard2UInt32 Card 转 int
func RoomCard2UInt32(card *room.Card) uint32 {
	var color uint32
	if *card.Color == room.CardColor_CC_WAN {
		color = 1
	} else if *card.Color == room.CardColor_CC_TIAO {
		color = 2
	} else if *card.Color == room.CardColor_CC_TONG {
		color = 3
	} else if *card.Color == room.CardColor_CC_FENG {
		color = 4
	}
	value := color*10 + uint32(*card.Point)
	return value
}

// RoomCards2UInt32 Card 转 int
func RoomCards2UInt32(card []*room.Card) []uint32 {
	result := []uint32{}
	for _, c := range card {
		result = append(result, RoomCard2UInt32(c))
	}
	return result
}

//CardEqual 判断两张牌是否一样
func CardEqual(card1 *majongpb.Card, card2 *majongpb.Card) bool {
	return card1.GetColor() == card2.GetColor() && card1.GetPoint() == card2.GetPoint()
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

// ServerColor2ClientColor 服务端协议卡牌花色转客户端花色
func ServerColor2ClientColor(color majongpb.CardColor) room.CardColor {
	switch color {
	case majongpb.CardColor_ColorWan:
		{
			return room.CardColor_CC_WAN
		}
	case majongpb.CardColor_ColorTiao:
		{
			return room.CardColor_CC_TIAO
		}
	case majongpb.CardColor_ColorTong:
		{
			return room.CardColor_CC_TONG
		}
	case majongpb.CardColor_ColorFeng:
		{
			return room.CardColor_CC_FENG
		}
	}
	return room.CardColor(-1)
}

// MakeRoomCards 构造牌切片
func MakeRoomCards(card ...room.Card) []*room.Card {
	result := []*room.Card{}
	for i := range card {
		result = append(result, &card[i])
	}
	return result
}
