package settle

import (
	"steve/client_pb/room"
	"steve/majong/utils"
	"steve/server_pb/majong"

	"github.com/gogo/protobuf/proto"
)

// RoundSettle 单局结算
type RoundSettle struct {
}

// SettleRound 单局结算
func (roundSettle *RoundSettle) SettleRound(context *majong.MajongContext) {
	// 退税 TODO

	//resp := new(room.RoomBalanceInfoRsp)
	for _, settleInfo := range context.SettleInfos {
		for _, player := range context.Players {
			billDetail := createRoomBalanceInfoRsBillDetail(player.PalyerId, settleInfo)
		}
	}
}

// createRoomBalanceInfoRsBillDetail 生成结算详情
func createRoomBalanceInfoRsBillDetail(palyerID uint64, settleInfo *majong.SettleInfo) *room.RoomBalanceInfoRsp_BillDetail {
	billDetail := &room.RoomBalanceInfoRsp_BillDetail{
		SettleTyoe: proto.String(string(settleInfo.SettleType)),
		FanType:    proto.String(settleInfo.FanType),
		Times:      proto.Int32(settleInfo.Times),
		Score:      proto.Int64(settleInfo.Scores[palyerID]),
	}
	if palyerID == settleInfo.GetPalyerId() {
		for uid, score := range settleInfo.Scores {
			if (uid != palyerID) && score != 0 {
				billDetail.RelatedPid = append(billDetail.RelatedPid, settleInfo.PalyerId)
			}
		}
	} else {
		if settleInfo.Scores[palyerID] != 0 {
			billDetail.RelatedPid = append(billDetail.RelatedPid, settleInfo.PalyerId)
		}
	}
	return billDetail

}

// createRoomBillPlayerInfo 生成玩家结算信息
func createRoomBillPlayerInfo(player *majong.Player, context majong.MajongContext) *room.BillPlayerInfo {
	total := int64(0)
	for _, settleInfo := range context.SettleInfos {
		total := settleInfo.Scores[player.PalyerId] + total
	}
	return &room.BillPlayerInfo{
		Pid:        proto.Uint64(player.PalyerId),
		Score:      proto.Int64(total),
		Name:       player.Name,
		CardsGroup: getCardsGroup(),
	}
	return nil
}

// getCardsGroup 获取玩家牌组信息
func getCardsGroup(player *majong.Player) []*room.CardsGroup {
	cardsGroupList := make([]*room.CardsGroup, 0)
	// 碰牌
	for _, pengCard := range player.PengCards {
		card, _ := utils.CardToInt(*pengCard.Card)
		cardsGroup := &room.CardsGroup{
			Pid:   proto.Uint64(player.PalyerId),
			Type:  room.CardsGroupType_CardGroupType_Peng.Enum(),
			Cards: []uint32{uint32(*card)},
		}
		cardsGroupList = append(cardsGroupList, cardsGroup)
	}
	// 杠牌
	var groupType *room.CardsGroupType
	for _, gangCard := range player.GangCards {
		if gangCard.Type == server_pb.GangType_gang_angang {
			groupType = room.CardsGroupType_CardGroupType_AnGang.Enum()
		}
		if gangCard.Type == server_pb.GangType_gang_minggang {
			groupType = room.CardsGroupType_CardGroupType_MingGang.Enum()
		}
		if gangCard.Type == server_pb.GangType_gang_bugang {
			groupType = room.CardsGroupType_CardGroupType_BuGang.Enum()
		}
		card, _ := utils.CardToInt(*gangCard.Card)
		cardsGroup := &room.CardsGroup{
			Pid:   proto.Uint64(player.PalyerId),
			Type:  groupType,
			Cards: []uint32{uint32(*card)},
		}
		cardsGroupList = append(cardsGroupList, cardsGroup)
	}
	// 手牌
	handCards, _ := utils.CardsToInt(player.HandCards)
	cards := make([]uint32, 0)
	for _, handCard := range handCards {
		cards = append(cards, uint32(handCard))
	}
	cardsGroup := &room.CardsGroup{
		Pid:   proto.Uint64(player.PalyerId),
		Type:  room.CardsGroupType_CardGroupType_Hand.Enum(),
		Cards: cards,
	}
	cardsGroupList = append(cardsGroupList, cardsGroup)
	return cardsGroupList
}
