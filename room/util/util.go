package util

import (
	"github.com/golang/protobuf/proto"
	"steve/client_pb/room"
	majongpb "steve/entity/majong"
)

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

// ServerCards2Numbers 服务器的 Card 数组转 int 数组
func ServerCards2Numbers(cards []*majongpb.Card) []uint32 {
	result := []uint32{}
	for _, c := range cards {
		result = append(result, ServerCard2Number(c))
	}
	return result
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
	} else if card.Color == majongpb.CardColor_ColorZi {
		color = 4
	} else if card.Color == majongpb.CardColor_ColorHua {
		color = 5
	}
	value := color*10 + uint32(card.Point)
	return value
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

// IsTing 玩家是否是听的状态
func IsTing(player *majongpb.Player) bool {
	tingState := player.GetTingStateInfo()
	if tingState.GetIsTing() || tingState.GetIsTianting() {
		return true
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

func GetMajongPlayer(playerID uint64, mjContext *majongpb.MajongContext) *majongpb.Player {
	for _, player := range mjContext.GetPlayers() {
		if player.GetPlayerId() == playerID {
			return player
		}
	}
	return nil
}

//合并任意个数组
func MergeStringArray(strings [][]string) []string {
	return mergeStringArray(nil, 0, 0, 0, strings)
}

/**
递归合并数组，调用时result传入nil,lastIndex,maxIndex,rIndex=0
*/
func mergeStringArray(result []string, maxIndex int, lastIndex int, rIndex int, strings [][]string) []string {
	if result == nil {
		for _, v := range strings {
			maxIndex += len(v)
		}
		result = make([]string, maxIndex)
	}

	if maxIndex == rIndex {
		return result
	}

	for _, v := range strings[lastIndex] {
		result[rIndex] = v
		rIndex++
	}
	lastIndex++
	return mergeStringArray(result, maxIndex, lastIndex, rIndex, strings)
}
