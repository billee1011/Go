package gutils

import (
	"fmt"
	"steve/client_pb/room"
	"steve/common/mjoption"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
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
	} else if *card.Color == room.CardColor_CC_ZI {
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
	if card.Color.String() == room.CardColor_CC_HUA.String() {
		color = room.CardColor_CC_HUA
	}
	return &room.Card{
		Color: color.Enum(),
		Point: proto.Int32(card.Point),
	}, nil
}

// ServerCard2Number 服务器的 Card 转换成数字
func ServerCard2Number(card *majongpb.Card) uint32 {
	var color uint32
	if card == nil {
		return 0
	}
	if card.Color == majongpb.CardColor_ColorWan {
		color = 1
	} else if card.Color == majongpb.CardColor_ColorTiao {
		color = 2
	} else if card.Color == majongpb.CardColor_ColorTong {
		color = 3
	} else if card.Color == majongpb.CardColor_ColorZi {
		color = 4
	} else if card.Color == majongpb.CardColor_ColorHua {
		color = 5
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

// ServerCards2Int32 服务器的 Card 数组转 int 数组
func ServerCards2Int32(cards []*majongpb.Card) []int32 {
	result := []int32{}
	for _, c := range cards {
		result = append(result, int32(ServerCard2Number(c)))
	}
	return result
}

// CardsToRoomCards 将Card转换为room package中的Card
func CardsToRoomCards(cards []*majongpb.Card) []*room.Card {
	var rCards []*room.Card
	for i := 0; i < len(cards); i++ {
		rCards = append(rCards, &room.Card{
			Color: ServerColor2ClientColor(cards[i].Color).Enum(),
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
	case majongpb.CardColor_ColorZi:
		{
			return room.CardColor_CC_ZI
		}
	}
	return room.CardColor(-1)
}

// ServerFanType2ClientHuType fanType获取hutype
func ServerFanType2ClientHuType(cardTypeOptionID int, fanTypes []int) int32 {
	cardTypeOption := mjoption.GetCardTypeOption(cardTypeOptionID)
	for _, fanType := range fanTypes {
		if _, ok := cardTypeOption.FanType2HuType[fanType]; ok {
			return int32(cardTypeOption.FanType2HuType[fanType].ID)
		}
	}
	return -1
}

// GetShowFan 获取实际显示的番型，移除番型中的胡牌类型及结算类型
func GetShowFan(cardTypeOptionID int, fanTypes []int) []int64 {
	cardTypeOption := mjoption.GetCardTypeOption(cardTypeOptionID)
	showFan := make([]int64, 0)
	for _, fanType := range fanTypes {
		if !cardTypeOption.EnableFanTypeDeal { // 胡类型是否从番型拿出
			showFan = append(showFan, int64(fanType))
			continue
		}
		_, isHuType := cardTypeOption.FanType2HuType[fanType]
		_, isSettleType := cardTypeOption.FanType2Settle[fanType]
		if !isHuType && !isSettleType {
			showFan = append(showFan, int64(fanType))
		}
	}
	return showFan
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

// TingTypeSvr2Client 听类型转换，server_pb-->client_pb
func TingTypeSvr2Client(recordTingType majongpb.TingType) *room.TingType {
	var tingType room.TingType
	switch recordTingType {
	case majongpb.TingType_TT_NORMAL_TING:
		tingType = room.TingType_TT_NORMAL_TING
	case majongpb.TingType_TT_TIAN_TING:
		tingType = room.TingType_TT_TIAN_TING
	}
	return &tingType
}

// CanTingCardInfoSvr2Client 玩家停牌信息转换，server_pb-->client_pb
func CanTingCardInfoSvr2Client(minfos []*majongpb.CanTingCardInfo) []*room.CanTingCardInfo {
	rinfos := []*room.CanTingCardInfo{}
	for _, minfo := range minfos {
		rinfo := &room.CanTingCardInfo{}
		rinfo.OutCard = proto.Uint32(minfo.GetOutCard())
		rinfo.TingCardInfo = TingCardInfoSvr2Client(minfo.GetTingCardInfo())
		rinfos = append(rinfos, rinfo)
	}
	return rinfos
}

// TingCardInfoSvr2Client 具体听牌信息转换，server_pb-->client_pb
func TingCardInfoSvr2Client(minfos []*majongpb.TingCardInfo) []*room.TingCardInfo {
	rinfos := []*room.TingCardInfo{}
	for _, minfo := range minfos {
		rinfo := &room.TingCardInfo{}
		rinfo.TingCard = proto.Uint32(minfo.GetTingCard())
		rinfo.Times = proto.Uint32(minfo.GetTimes())
		rinfos = append(rinfos, rinfo)
	}
	return rinfos
}

// CardsGroupSvr2Client server的牌组类型转换，server_pb-->client_pb
func CardsGroupSvr2Client(cardsGroups []*majongpb.CardsGroup) (cardsGroupList []*room.CardsGroup) {
	cardsGroupList = make([]*room.CardsGroup, 0)
	for _, cardsGroup := range cardsGroups {
		cardsGroupList = append(cardsGroupList, &room.CardsGroup{
			Pid:   proto.Uint64(cardsGroup.Pid),
			Type:  room.CardsGroupType(int32(cardsGroup.GetType())).Enum(),
			Cards: cardsGroup.Cards,
		})
	}
	return
}

//FmtPlayerInfo 打印玩家信息
func FmtPlayerInfo(player *majongpb.Player) logrus.Fields {
	fields := logrus.Fields{
		"玩家ID":      player.GetPalyerId(),
		"手牌":        FmtMajongpbCards(player.HandCards),
		"问询下可以有的操作": player.PossibleActions,
		"杠过的牌":      FmtGangCards(player.GangCards),
		"胡过的牌":      FmtHuCards(player.HuCards),
		"碰过的牌":      FmtPengCards(player.PengCards),
		"出过的牌":      FmtMajongpbCards(player.OutCards),
	}
	return fields
}

//FmtMajongpbCards 打印牌组
func FmtMajongpbCards(cards []*majongpb.Card) string {
	results := ""
	for _, card := range cards {
		if card != nil {
			results += fmt.Sprintf("%v%v ", card.Point, getColor(card.Color))
		}
	}
	return results
}

//FmtGangCards 打印gangCards
func FmtGangCards(gangCards []*majongpb.GangCard) string {
	result := ""
	for _, gangCard := range gangCards {
		result += fmt.Sprintf("杠的类型:%v ", gangCard.Type.String())
		result += fmt.Sprintf("杠的牌:%v%v ", gangCard.Card.Point, getColor(gangCard.Card.Color))
		result += fmt.Sprintf("来自玩家:%v ", gangCard.SrcPlayer)
	}
	return result
}

//FmtPengCards 打印pengCards
func FmtPengCards(pengCards []*majongpb.PengCard) string {
	result := ""
	for _, pengCard := range pengCards {
		result += fmt.Sprintf("碰的牌:%v%v ", pengCard.Card.Point, getColor(pengCard.Card.Color))
		result += fmt.Sprintf("来自玩家:%v; ", pengCard.SrcPlayer)
	}
	return result
}

//FmtHuCards 打印hucards
func FmtHuCards(huCards []*majongpb.HuCard) string {
	result := ""
	for _, huCard := range huCards {
		result += fmt.Sprintf("胡的类型:%v ", huCard.Type.String())
		result += fmt.Sprintf("胡的牌:%v%v ", huCard.Card.Point, getColor(huCard.Card.Color))
		result += fmt.Sprintf("来自玩家:%v ", huCard.SrcPlayer)
	}
	return result
}

func getColor(srcColor majongpb.CardColor) string {
	if srcColor == majongpb.CardColor_ColorWan {
		return "w"
	}
	if srcColor == majongpb.CardColor_ColorTiao {
		return "t"
	}
	if srcColor == majongpb.CardColor_ColorTong {
		return "b"
	}
	if srcColor == majongpb.CardColor_ColorZi {
		return "z"
	}
	if srcColor == majongpb.CardColor_ColorHua {
		return "h"
	}
	return "none"
}

//FmtMajongContxt 打印麻将现场
func FmtMajongContxt(context *majongpb.MajongContext) logrus.Fields {

	return logrus.Fields{
		"LastGangPlayer":   context.GetLastGangPlayer(),
		"LastChupaiPlayer": context.GetLastChupaiPlayer(),
		"LastOutCard":      FmtMajongpbCards([]*majongpb.Card{context.LastOutCard}),
		"LastMopaiPlayer":  context.GetLastMopaiPlayer(),
		"LastMopaiCard":    FmtMajongpbCards([]*majongpb.Card{context.LastMopaiCard}),
		"LastPengPlayer":   context.GetLastPengPlayer(),
		"MopaiPlayer":      context.GetMopaiPlayer(),
	}
}

//CheckHasDingQueCard 检查牌里面是否含有定缺的牌
func CheckHasDingQueCard(context *majongpb.MajongContext, player *majongpb.Player) bool {
	xpOption := mjoption.GetXingpaiOption(int(context.GetXingpaiOptionId()))
	cards, color, hasDq := player.GetHandCards(), player.GetDingqueColor(), xpOption.EnableDingque
	if !hasDq {
		return false
	}
	for _, card := range cards {
		if card.GetColor() == color {
			return true
		}
	}
	return false
}

//IsDingQueCard 当前的牌是不是定缺牌
func IsDingQueCard(context *majongpb.MajongContext, dqColor majongpb.CardColor, card *majongpb.Card) bool {
	xpOption := mjoption.GetXingpaiOption(int(context.GetXingpaiOptionId()))
	if xpOption.EnableDingque && card.GetColor() == dqColor {
		return true
	}
	return false
}

// GetCardsGroup 获取玩家牌组信息
func GetCardsGroup(player *majongpb.Player) []*room.CardsGroup {
	cardsGroupList := make([]*room.CardsGroup, 0)
	// 手牌组
	cltHandCard := ServerCards2Numbers(player.GetHandCards())
	handCardGroup := &room.CardsGroup{
		Cards: cltHandCard,
		Type:  room.CardsGroupType_CGT_HAND.Enum(),
	}
	cardsGroupList = append(cardsGroupList, handCardGroup)
	// 吃牌组
	var chiCardGroups []*room.CardsGroup
	for _, chiCard := range player.GetChiCards() {
		srcPlayerID := chiCard.GetSrcPlayer()
		card := ServerCard2Number(chiCard.GetCard())
		chiCardGroup := &room.CardsGroup{
			Cards: []uint32{card, card + 1, card + 2},
			Type:  room.CardsGroupType_CGT_CHI.Enum(),
			Pid:   &srcPlayerID,
		}
		chiCardGroups = append(chiCardGroups, chiCardGroup)
	}
	cardsGroupList = append(cardsGroupList, chiCardGroups...)
	// 碰牌组
	var pengCardGroups []*room.CardsGroup
	for _, pengCard := range player.GetPengCards() {
		srcPlayerID := pengCard.GetSrcPlayer()
		cards := []uint32{ServerCard2Number(pengCard.GetCard())}
		pengCardGroup := &room.CardsGroup{
			Cards: append(cards, cards[0], cards[0]),
			Type:  room.CardsGroupType_CGT_PENG.Enum(),
			Pid:   &srcPlayerID,
		}
		pengCardGroups = append(pengCardGroups, pengCardGroup)
	}
	cardsGroupList = append(cardsGroupList, pengCardGroups...)
	// 杠牌组
	var gangCardGroups []*room.CardsGroup
	for _, gangCard := range player.GetGangCards() {
		groupType := GangTypeSvr2Client(gangCard.GetType())
		srcPlayerID := gangCard.GetSrcPlayer()
		cards := []uint32{ServerCard2Number(gangCard.GetCard())}
		gangCardGroup := &room.CardsGroup{
			Cards: append(cards, cards[0], cards[0], cards[0]),
			Type:  &groupType,
			Pid:   &srcPlayerID,
		}
		gangCardGroups = append(gangCardGroups, gangCardGroup)
	}
	cardsGroupList = append(cardsGroupList, gangCardGroups...)
	// 胡牌组
	var huCardGroups []*room.CardsGroup
	for _, huCard := range player.GetHuCards() {
		srcPlayerID := huCard.GetSrcPlayer()
		huCardGroup := &room.CardsGroup{
			Cards:  []uint32{ServerCard2Number(huCard.GetCard())},
			Type:   room.CardsGroupType_CGT_HU.Enum(),
			Pid:    &srcPlayerID,
			IsReal: proto.Bool(huCard.GetIsReal()),
		}
		huCardGroups = append(huCardGroups, huCardGroup)
	}
	cardsGroupList = append(cardsGroupList, huCardGroups...)
	// 花牌组
	var huaCardGroups []*room.CardsGroup
	for _, huaCard := range player.GetHuaCards() {
		huaCardGroup := &room.CardsGroup{
			Cards: []uint32{ServerCard2Number(huaCard)},
			Type:  room.CardsGroupType_CGT_HUA.Enum(),
		}
		huaCardGroups = append(huaCardGroups, huaCardGroup)
	}
	cardsGroupList = append(cardsGroupList, huaCardGroups...)
	return cardsGroupList
}

// DeleteHuType 移除番型中的胡类型
func DeleteHuType(cardTypeOptionID int, fanTypes []int) []int {
	cardTypeOption := mjoption.GetCardTypeOption(cardTypeOptionID)
	showFan := make([]int, 0)
	for _, fanType := range fanTypes {
		_, isHuType := cardTypeOption.FanType2HuType[fanType]
		_, isSettleType := cardTypeOption.FanType2Settle[fanType]
		//TODO:建议此处选项化不要以胡类型结算类型来排除行牌过程的胡牌提示，添加一个在行牌阶段可以查番的番型列表或者不可查番型的番型列表
		//将天胡，报听一发这些在行牌阶段不确定的番型进行归类，查的时候直接排除不去查就行了，这里暂时先将报听一发写死在代码里进行排除，后面
		//会统一在番型选项化中对这些不确定番型进行排除
		if !isHuType && !isSettleType && fanType != int(room.FanType_FT_BAOTINGYIFA) {
			showFan = append(showFan, fanType)
		}
	}
	if len(fanTypes) != len(showFan) {
		logrus.WithFields(logrus.Fields{
			"func_name": "DeleteHuType",
			"fanTypes":  fanTypes,
			"showFan":   showFan,
		}).Error("移除番型中的胡类型")
	}
	return showFan
}
