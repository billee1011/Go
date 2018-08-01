package majong

import (
	"steve/client_pb/room"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// translateHeader 生成事件头
func translateHeader(playerID uint64, header *steve_proto_gaterpc.Header, body proto.Message) server_pb.RequestEventHead {
	return server_pb.RequestEventHead{
		PlayerId: playerID,
	}
}

// translateClientCardColor 转换卡牌花色
func translateClientCardColor(color room.CardColor) server_pb.CardColor {
	switch color {
	case room.CardColor_CC_WAN:
		{
			return server_pb.CardColor_ColorWan
		}
	case room.CardColor_CC_TIAO:
		{
			return server_pb.CardColor_ColorTiao
		}
	case room.CardColor_CC_TONG:
		{
			return server_pb.CardColor_ColorTong
		}
	case room.CardColor_CC_ZI:
		{
			return server_pb.CardColor_ColorZi
		}
	}
	return server_pb.CardColor(-1)
}

// translateCardColor 转换卡牌花色
func translateCardColor(cardVal uint32) server_pb.CardColor {
	switch cardVal / 10 {
	case 1:
		{
			return server_pb.CardColor_ColorWan
		}
	case 2:
		{
			return server_pb.CardColor_ColorTiao
		}
	case 3:
		{
			return server_pb.CardColor_ColorTong
		}
	case 4:
		{
			return server_pb.CardColor_ColorZi
		}
	}
	return server_pb.CardColor(-1)
}

// translateCard 转换卡牌
func translateCard(card uint32) server_pb.Card {
	return server_pb.Card{
		Color: translateCardColor(card),
		Point: int32(card % 10),
	}
}

// translateCards 转换多个卡牌
func translateCards(cards []uint32) []*server_pb.Card {
	result := []*server_pb.Card{}

	for _, card := range cards {
		serverCard := translateCard(card)
		result = append(result, &serverCard)
	}
	return result
}

func translateTingAction(tingAction *room.TingAction) *server_pb.TingAction {
	serTingAction := &server_pb.TingAction{}
	serTingAction.EnableTing = tingAction.GetEnableTing()
	serTingAction.TingType = translateTingType(tingAction.GetTingType())
	return serTingAction
}

func translateTingType(tingType room.TingType) server_pb.TingType {
	switch tingType {
	case room.TingType_TT_NORMAL_TING:
		return server_pb.TingType_TT_NORMAL_TING
	case room.TingType_TT_TIAN_TING:
		return server_pb.TingType_TT_TIAN_TING
	}
	return server_pb.TingType_TT_NORMAL_TING
}
