package rtoet

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

// translateCardColor 转换卡牌花色
func translateCardColor(color room.CardColor) server_pb.CardColor {
	switch color {
	case room.CardColor_ColorWan:
		{
			return server_pb.CardColor_ColorWan
		}
	case room.CardColor_ColorTiao:
		{
			return server_pb.CardColor_ColorTiao
		}
	case room.CardColor_ColorTong:
		{
			return server_pb.CardColor_ColorTong
		}
	case room.CardColor_ColorFeng:
		{
			return server_pb.CardColor_ColorFeng
		}
	}
	return server_pb.CardColor(-1)
}

// translateCard 转换卡牌
func translateCard(card room.Card) server_pb.Card {
	return server_pb.Card{
		Color: translateCardColor(card.GetColor()),
		Point: card.GetPoint(),
	}
}

// translateCards 转换多个卡牌
func translateCards(cards []*room.Card) []*server_pb.Card {
	result := []*server_pb.Card{}

	for _, card := range cards {
		serverCard := translateCard(*card)
		result = append(result, &serverCard)
	}
	return result
}
