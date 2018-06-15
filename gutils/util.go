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
func ServerCard2Number(card *majongpb.Card) uint32 {
	var color uint32
	if card.Color == majongpb.CardColor_ColorWan {
		color = 1
	} else if card.Color == majongpb.CardColor_ColorTiao {
		color = 2
	} else if card.Color == majongpb.CardColor_ColorTong {
		color = 3
	} else if card.Color == majongpb.CardColor_ColorFeng {
		color = 4
	}
	value := color*10 + uint32(card.Point)
	return value
}

// ServerCards2Numbers 服务器的 Card 数组转 int 数组
func ServerCards2Numbers(cards []*majongpb.Card) []uint32 {
	result := []uint32{}
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

// GangTypeSvr2Client server的杠类型转换为恢复牌局的麻将组类型，server_pb-->client_pb
func GangTypeSvr2Client(gangType majongpb.GangType) (groupType room.CardsGroupType) {
	switch gangType {
	case majongpb.GangType_gang_angang:
		groupType = room.CardsGroupType_CGT_ANGANG
	case majongpb.GangType_gang_minggang:
		groupType = room.CardsGroupType_CGT_MINGGANG
	case majongpb.GangType_gang_bugang:
		groupType = room.CardsGroupType_CGT_BUGANG
	}
	return
}

// HuTypeSvr2Client 胡类型转换，server_pb-->client_pb
func HuTypeSvr2Client(recordHuType majongpb.HuType) *room.HuType {
	var huType room.HuType
	switch recordHuType {
	case majongpb.HuType_hu_dianpao:
		huType = room.HuType_HT_DIANPAO
	case majongpb.HuType_hu_dihu:
		huType = room.HuType_HT_DIHU
	case majongpb.HuType_hu_ganghoupao:
		huType = room.HuType_HT_GANGHOUPAO
	case majongpb.HuType_hu_gangkai:
		huType = room.HuType_HT_GANGKAI
	case majongpb.HuType_hu_gangshanghaidilao:
		huType = room.HuType_HT_GANGSHANGHAIDILAO
	case majongpb.HuType_hu_haidilao:
		huType = room.HuType_HT_HAIDILAO
	case majongpb.HuType_hu_qiangganghu:
		huType = room.HuType_HT_QIANGGANGHU
	case majongpb.HuType_hu_tianhu:
		huType = room.HuType_HT_TIANHU
	case majongpb.HuType_hu_zimo:
		huType = room.HuType_HT_ZIMO
	default:
		return nil
	}
	return &huType
}

// CanTingCardInfoSvr2Client 玩家停牌信息转换，server_pb-->client_pb
func CanTingCardInfoSvr2Client(minfos []*majongpb.CanTingCardInfo) []*room.CanTingCardInfo {
	rinfos := []*room.CanTingCardInfo{}
	for _, minfo := range minfos {
		rinfo := &room.CanTingCardInfo{}
		rinfo.OutCard = proto.Uint32(minfo.GetOutCard())
		rinfo.TingCardInfo = tingCardInfoSvr2Client(minfo.GetTingCardInfo())
		rinfos = append(rinfos, rinfo)
	}
	return rinfos
}

// tingCardInfoSvr2Client 具体听牌信息转换，server_pb-->client_pb
func tingCardInfoSvr2Client(minfos []*majongpb.TingCardInfo) []*room.TingCardInfo {
	rinfos := []*room.TingCardInfo{}
	for _, minfo := range minfos {
		rinfo := &room.TingCardInfo{}
		rinfo.TingCard = proto.Uint32(minfo.GetTingCard())
		rinfo.Times = proto.Uint32(minfo.GetTimes())
		rinfos = append(rinfos, rinfo)
	}
	return rinfos
}

// GetCardsGroup 获取玩家牌组信息
func GetCardsGroup(player *majongpb.Player) []*room.CardsGroup {
	cardsGroupList := make([]*room.CardsGroup, 0)
	// 碰牌
	for _, pengCard := range player.PengCards {
		card := ServerCard2Number((*pengCard).Card)
		cardsGroup := &room.CardsGroup{
			Pid:   proto.Uint64(player.PalyerId),
			Type:  room.CardsGroupType_CGT_PENG.Enum(),
			Cards: []uint32{uint32(card)},
		}
		cardsGroupList = append(cardsGroupList, cardsGroup)
	}
	// 杠牌
	var groupType *room.CardsGroupType
	for _, gangCard := range player.GangCards {
		if gangCard.Type == majongpb.GangType_gang_angang {
			groupType = room.CardsGroupType_CGT_ANGANG.Enum()
		}
		if gangCard.Type == majongpb.GangType_gang_minggang {
			groupType = room.CardsGroupType_CGT_MINGGANG.Enum()
		}
		if gangCard.Type == majongpb.GangType_gang_bugang {
			groupType = room.CardsGroupType_CGT_BUGANG.Enum()
		}
		card := ServerCard2Number((*gangCard).Card)
		cardsGroup := &room.CardsGroup{
			Pid:   proto.Uint64(player.PalyerId),
			Type:  groupType,
			Cards: []uint32{uint32(card)},
		}
		cardsGroupList = append(cardsGroupList, cardsGroup)
	}
	// 手牌
	handCards := ServerCards2Numbers(player.HandCards)
	cards := make([]uint32, 0)
	for _, handCard := range handCards {
		cards = append(cards, uint32(handCard))
	}
	cardsGroup := &room.CardsGroup{
		Pid:   proto.Uint64(player.PalyerId),
		Type:  room.CardsGroupType_CGT_HAND.Enum(),
		Cards: cards,
	}
	cardsGroupList = append(cardsGroupList, cardsGroup)
	return cardsGroupList
}
