package settle

import (
	"steve/client_pb/room"
	"steve/majong/utils"
	"steve/room/interfaces"
	server_pb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// SInfoID 结算信息id
type SInfoID uint64

// PlayerID 玩家ID
type PlayerID uint64

// scxlSettle 血流麻将结算
type scxlSettle struct {
	settleMap   map[SInfoID]int64  //key->:结算id;value->:该条结算扣除金币数
	roundScore  map[PlayerID]int64 //key->:玩家id;value->:单局结算扣除金币数
	revertedMap map[SInfoID]bool   //退税SInfoID
}

// Settle 结算信息扣分
func (s *scxlSettle) Settle(desk interfaces.Desk, mjContext server_pb.MajongContext) {
	lastSettleInfo := mjContext.SettleInfos[len(mjContext.SettleInfos)-1]
	if !lastSettleInfo.Handle {
		for _, player := range desk.GetPlayers() {
			settleScore := lastSettleInfo.Scores[*player.PlayerId]
			*player.Coin = uint64(settleScore + int64(*player.Coin))
			if settleScore < 0 && *player.Coin < 0 {
				settleScore = int64(*player.Coin)
				*player.Coin = 0
			}
			s.settleMap[SInfoID(lastSettleInfo.Id)] = settleScore
			s.roundScore[PlayerID(*player.PlayerId)] = s.roundScore[PlayerID(*player.PlayerId)] + settleScore
		}
		lastSettleInfo.Handle = true
	}
}

// RoundSettle 单局结算信息
func (s *scxlSettle) RoundSettle(desk interfaces.Desk, mjContext server_pb.MajongContext) {
	balanceRsp := new(room.RoomBalanceInfoRsp)
	for _, roomPlayer := range desk.GetPlayers() {
		for _, settleInfo := range mjContext.SettleInfos {
			billDetail := s.createBillDetail(*roomPlayer.PlayerId, settleInfo)
			if billDetail != nil {
				balanceRsp.BillDetail = append(balanceRsp.BillDetail, billDetail)
			}
		}
		balanceRsp.BillPlayerInfo = s.createBillPlayerInfo(roomPlayer, mjContext)
	}

}

// createBillDetail 单次结算详情，包括番型，分数，倍数，以及输赢玩家
func (s *scxlSettle) createBillDetail(palyerID uint64, settleInfo *server_pb.SettleInfo) *room.RoomBalanceInfoRsp_BillDetail {
	if settleInfo.Scores[palyerID] != 0 {
		billDetail := &room.RoomBalanceInfoRsp_BillDetail{
			SetleType: proto.String(string(settleInfo.SettleType)),
			FanType:   proto.String(settleInfo.FanType),
			Times:     proto.Int32(settleInfo.Times),
			Score:     proto.Int64(int64(s.settleMap[SInfoID(settleInfo.Id)])),
		}
		if settleInfo.Scores[palyerID] > 0 {
			for pid, score := range settleInfo.Scores {
				if palyerID != pid && score != 0 {
					billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
				}
			}
		} else {
			for pid, score := range settleInfo.Scores {
				if palyerID != pid && score > 0 {
					billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
				}
			}
		}
		return billDetail
	}
	return nil
}

// createBillPlayerInfo 单局结算详情,包括玩家自己牌型,输赢分数，以及其余每个玩家的输赢分数
func (s *scxlSettle) createBillPlayerInfo(roomPlayer *room.RoomPlayerInfo, context server_pb.MajongContext) []*room.BillPlayerInfo {
	billPlayerInfos := make([]*room.BillPlayerInfo, 0)
	billPlayerInfo := new(room.BillPlayerInfo)
	for _, player := range context.Players {
		billPlayerInfo.Pid = roomPlayer.PlayerId
		billPlayerInfo.Score = proto.Int64(s.roundScore[PlayerID(player.PalyerId)])
		billPlayerInfo.Name = roomPlayer.Name
		if player.PalyerId == *roomPlayer.PlayerId {
			billPlayerInfo.CardsGroup = getCardsGroup(utils.GetPlayerByID(context.Players, *roomPlayer.PlayerId))
		}
		billPlayerInfos = append(billPlayerInfos, billPlayerInfo)
	}
	return billPlayerInfos
}

// getCardsGroup 获取玩家牌组信息
func getCardsGroup(player *server_pb.Player) []*room.CardsGroup {
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
