package tests

import (
	"steve/client_pb/room"

	"github.com/golang/protobuf/proto"
)

var (
	// Card1W 1 万
	Card1W = room.Card{Color: room.CardColor_CC_WAN.Enum(), Point: proto.Int32(1)}
	// Card2W 2 万
	Card2W = room.Card{Color: room.CardColor_CC_WAN.Enum(), Point: proto.Int32(2)}
	// Card3W 3 万
	Card3W = room.Card{Color: room.CardColor_CC_WAN.Enum(), Point: proto.Int32(3)}
	// Card4W 4 万
	Card4W = room.Card{Color: room.CardColor_CC_WAN.Enum(), Point: proto.Int32(4)}
	// Card5W 5 万
	Card5W = room.Card{Color: room.CardColor_CC_WAN.Enum(), Point: proto.Int32(5)}
	// Card6W 6 万
	Card6W = room.Card{Color: room.CardColor_CC_WAN.Enum(), Point: proto.Int32(6)}
	// Card7W 7 万
	Card7W = room.Card{Color: room.CardColor_CC_WAN.Enum(), Point: proto.Int32(7)}
	// Card8W 8 万
	Card8W = room.Card{Color: room.CardColor_CC_WAN.Enum(), Point: proto.Int32(8)}
	// Card9W 9 万
	Card9W = room.Card{Color: room.CardColor_CC_WAN.Enum(), Point: proto.Int32(9)}

	// Card1T 1 条
	Card1T = room.Card{Color: room.CardColor_CC_TIAO.Enum(), Point: proto.Int32(1)}
	// Card2T 2 条
	Card2T = room.Card{Color: room.CardColor_CC_TIAO.Enum(), Point: proto.Int32(2)}
	// Card3T 3 条
	Card3T = room.Card{Color: room.CardColor_CC_TIAO.Enum(), Point: proto.Int32(3)}
	// Card4T 4 条
	Card4T = room.Card{Color: room.CardColor_CC_TIAO.Enum(), Point: proto.Int32(4)}
	// Card5T 5 条
	Card5T = room.Card{Color: room.CardColor_CC_TIAO.Enum(), Point: proto.Int32(5)}
	// Card6T 6 条
	Card6T = room.Card{Color: room.CardColor_CC_TIAO.Enum(), Point: proto.Int32(6)}
	// Card7T 7 条
	Card7T = room.Card{Color: room.CardColor_CC_TIAO.Enum(), Point: proto.Int32(7)}
	// Card8T 8 条
	Card8T = room.Card{Color: room.CardColor_CC_TIAO.Enum(), Point: proto.Int32(8)}
	// Card9T 9 条
	Card9T = room.Card{Color: room.CardColor_CC_TIAO.Enum(), Point: proto.Int32(9)}

	// Card1B 1 筒
	Card1B = room.Card{Color: room.CardColor_CC_TONG.Enum(), Point: proto.Int32(1)}
	// Card2B 2 筒
	Card2B = room.Card{Color: room.CardColor_CC_TONG.Enum(), Point: proto.Int32(2)}
	// Card3B 3 筒
	Card3B = room.Card{Color: room.CardColor_CC_TONG.Enum(), Point: proto.Int32(3)}
	// Card4B 4 筒
	Card4B = room.Card{Color: room.CardColor_CC_TONG.Enum(), Point: proto.Int32(4)}
	// Card5B 5 筒
	Card5B = room.Card{Color: room.CardColor_CC_TONG.Enum(), Point: proto.Int32(5)}
	// Card6B 6 筒
	Card6B = room.Card{Color: room.CardColor_CC_TONG.Enum(), Point: proto.Int32(6)}
	// Card7B 7 筒
	Card7B = room.Card{Color: room.CardColor_CC_TONG.Enum(), Point: proto.Int32(7)}
	// Card8B 8 筒
	Card8B = room.Card{Color: room.CardColor_CC_TONG.Enum(), Point: proto.Int32(8)}
	// Card9B 9 筒
	Card9B = room.Card{Color: room.CardColor_CC_TONG.Enum(), Point: proto.Int32(9)}
)
