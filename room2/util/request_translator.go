package util

import (
	"steve/client_pb/room"
	server_pb "steve/entity/majong"
	"steve/structs/proto/gate_rpc"
)

// TranslateXingpaiActionReq 转换行牌动作请求
func TranslateXingpaiActionReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomXingpaiActionReq) (eventID int, eventContext interface{}, err error) {

	eventHeader := translateHeader(playerID, header, &req)
	switch req.GetActionId() {
	case room.XingpaiAction_XA_CHI:
		{
			eventID = int(server_pb.EventID_event_chi_request)
			cards := translateCards(req.GetChiCards())
			eventContext = &server_pb.ChiRequestEvent{
				Head:  &eventHeader,
				Cards: cards,
			}
		}
	case room.XingpaiAction_XA_PENG:
		{
			eventID = int(server_pb.EventID_event_peng_request)
			eventContext = &server_pb.PengRequestEvent{
				Head: &eventHeader,
			}
		}
	case room.XingpaiAction_XA_GANG:
		{
			eventID = int(server_pb.EventID_event_gang_request)
			card := translateCard(req.GetGangCard())
			eventContext = &server_pb.GangRequestEvent{
				Head: &eventHeader,
				Card: &card,
			}
		}
	case room.XingpaiAction_XA_HU:
		{
			eventID = int(server_pb.EventID_event_hu_request)
			eventContext = &server_pb.HuRequestEvent{
				Head: &eventHeader,
			}
		}
	case room.XingpaiAction_XA_QI:
		{
			eventID = int(server_pb.EventID_event_qi_request)
			eventContext = &server_pb.QiRequestEvent{
				Head: &eventHeader,
			}
		}
	}
	return
}

// TranslateDingqueReq 转换定缺请求
func TranslateDingqueReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomDingqueReq) (eventID int, eventContext interface{}, err error) {

	eventHeader := translateHeader(playerID, header, &req)

	cardColor := translateClientCardColor(req.GetColor())
	eventContext = &server_pb.DingqueRequestEvent{
		Head:  &eventHeader,
		Color: cardColor,
	}
	eventID = int(server_pb.EventID_event_dingque_request)
	return
}

// TranslateChupaiReq 转换出牌请求
func TranslateChupaiReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomChupaiReq) (eventID int, eventContext interface{}, err error) {

	eventHeader := translateHeader(playerID, header, &req)

	card := translateCard(req.GetCard())
	eventContext = &server_pb.ChupaiRequestEvent{
		Head:       &eventHeader,
		Cards:      &card,
		TingAction: translateTingAction(req.GetTingAction()),
	}
	eventID = int(server_pb.EventID_event_chupai_request)
	return
}

// TranslateCartoonFinishReq 转换动画完成请求
func TranslateCartoonFinishReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomCartoonFinishReq) (eventID int, eventContext interface{}, err error) {
	eventContext = &server_pb.CartoonFinishRequestEvent{
		CartoonType: int32(req.GetCartoonType()),
	}
	eventID = int(server_pb.EventID_event_cartoon_finish_request)
	return
}

// TranslateHuansanzhangReq 转换换三张请求
func TranslateHuansanzhangReq(playerID uint64, header *steve_proto_gaterpc.Header,
	req room.RoomHuansanzhangReq) (eventID int, eventContext interface{}, err error) {

	eventHeader := translateHeader(playerID, header, &req)

	eventContext = &server_pb.HuansanzhangRequestEvent{
		Head:  &eventHeader,
		Cards: translateCards(req.GetCards()),
		Sure:  req.GetSure(),
	}
	eventID = int(server_pb.EventID_event_huansanzhang_request)
	return
}

// translateHeader 生成事件头
func translateHeader(playerID uint64, header *steve_proto_gaterpc.Header, body interface{}) server_pb.RequestEventHead {
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
