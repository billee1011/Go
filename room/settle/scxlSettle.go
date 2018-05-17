package settle

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	majongpb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// scxlSettle 血流麻将结算
type scxlSettle struct {
	// setttleInfo中每个玩家实际输赢分 key:settleId value:playerCoin
	settleMap map[uint64]playerCoin
	// setttleInfo处理情况 		key:settleId value:true为已处理，false为未处理
	handleSettle map[uint64]bool
}

// newScxlSettle 创建四川血流结算
func newScxlSettle() *scxlSettle {
	return &scxlSettle{
		settleMap:    make(map[uint64]playerCoin),
		handleSettle: make(map[uint64]bool),
	}
}

// playerCoin 玩家实际输赢分   key:playerID value:score
type playerCoin map[uint64]int64

// Settle 结算信息扣分并通知客户端
func (s *scxlSettle) Settle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	if len(mjContext.SettleInfos) != 0 {
		lastSettleInfo := mjContext.SettleInfos[len(mjContext.SettleInfos)-1]
		if !s.handleSettle[lastSettleInfo.Id] {
			playerCoin := make(playerCoin)
			billplayerInfos := make([]*room.BillPlayerInfo, 0)
			for _, player := range desk.GetPlayers() {
				pid := *player.PlayerId
				billplayerInfo := &room.BillPlayerInfo{
					Pid:      player.PlayerId,
					BillType: room.BillType(lastSettleInfo.SettleType).Enum(),
				}
				score := lastSettleInfo.Scores[pid]

				if score < 0 && (-score) >= int64(*player.Coin) {
					playerCoin[pid] = int64(*player.Coin)
					billplayerInfo.Score = proto.Int64(int64(*player.Coin))
					*player.Coin = 0
				} else {
					playerCoin[pid] = score
					billplayerInfo.Score = proto.Int64(score)
					*player.Coin = uint64(int64(*player.Coin) + score)
				}

				s.settleMap[lastSettleInfo.Id] = playerCoin
				billplayerInfo.CurrentScore = proto.Int64(int64(*player.Coin))
				billplayerInfos = append(billplayerInfos, billplayerInfo)
			}
			s.handleSettle[lastSettleInfo.Id] = true
			// 广播即时结算消息
			instantSettle := room.RoomSettleInstantRsp{
				BillPlayersInfo: billplayerInfos,
			}
			notifyDeskMessage(desk, &instantSettle)
		}
	}
}

func notifyDeskMessage(desk interfaces.Desk, message proto.Message) {
	players := desk.GetPlayers()
	clientIDs := []uint64{}

	playerMgr := global.GetPlayerMgr()
	for _, player := range players {
		playerID := player.GetPlayerId()
		p := playerMgr.GetPlayer(playerID)
		if p != nil {
			clientIDs = append(clientIDs, p.GetClientID())
		}
	}
	head := &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid.MsgID_ROOM_INSTANT_SETTLE)}
	ms := global.GetMessageSender()

	ms.BroadcastPackage(clientIDs, head, message)
}

// RoundSettle 单局结算信息
func (s *scxlSettle) RoundSettle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	// balanceRsp := new(room.RoomBalanceInfoRsp)
	// for _, roomPlayer := range desk.GetPlayers() {
	// 	for _, settleInfo := range mjContext.SettleInfos {
	// 		billDetail := s.createBillDetail(*roomPlayer.PlayerId, settleInfo)
	// 		if billDetail != nil {
	// 			balanceRsp.BillDetail = append(balanceRsp.BillDetail, billDetail)
	// 		}
	// 	}
	// 	balanceRsp.BillPlayerInfo = s.createBillPlayerInfo(roomPlayer, mjContext)
	// }

}

// createBillDetail 单次结算详情，包括番型，分数，倍数，以及输赢玩家
func (s *scxlSettle) createBillDetail(palyerID uint64, settleInfo *majongpb.SettleInfo) *room.BillDetail {
	// if settleInfo.Scores[palyerID] != 0 {
	// 	billDetail := &room.RoomBalanceInfoRsp_BillDetail{
	// 		SetleType: proto.String(string(settleInfo.SettleType)),
	// 		// FanType:   proto.String(settleInfo.FanType),
	// 		// Times:     proto.Int32(settleInfo.Times),
	// 		Score: proto.Int64(int64(s.settleMap[SInfoID(settleInfo.Id)])),
	// 	}
	// 	if settleInfo.Scores[palyerID] > 0 {
	// 		for pid, score := range settleInfo.Scores {
	// 			if palyerID != pid && score != 0 {
	// 				billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
	// 			}
	// 		}
	// 	} else {
	// 		for pid, score := range settleInfo.Scores {
	// 			if palyerID != pid && score > 0 {
	// 				billDetail.RelatedPid = append(billDetail.RelatedPid, pid)
	// 			}
	// 		}
	// 	}
	// 	return billDetail
	// }
	return nil
}

// createBillPlayerInfo 单局结算详情,包括玩家自己牌型,输赢分数，以及其余每个玩家的输赢分数
func (s *scxlSettle) createBillPlayerInfo(roomPlayer *room.RoomPlayerInfo, context majongpb.MajongContext) []*room.BillPlayerInfo {
	// billPlayerInfos := make([]*room.BillPlayerInfo, 0)
	// billPlayerInfo := new(room.BillPlayerInfo)
	// for _, player := range context.Players {
	// 	billPlayerInfo.Pid = roomPlayer.PlayerId
	// 	billPlayerInfo.Score = proto.Int64(s.roundScore[PlayerID(player.PalyerId)])
	// 	billPlayerInfo.Name = roomPlayer.Name
	// 	if player.PalyerId == *roomPlayer.PlayerId {
	// 		billPlayerInfo.CardsGroup = getCardsGroup(utils.GetPlayerByID(context.Players, *roomPlayer.PlayerId))
	// 	}
	// 	billPlayerInfos = append(billPlayerInfos, billPlayerInfo)
	// }
	// return billPlayerInfos
	return nil
}

// getCardsGroup 获取玩家牌组信息
func getCardsGroup(player *majongpb.Player) []*room.CardsGroup {
	// cardsGroupList := make([]*room.CardsGroup, 0)
	// // 碰牌
	// for _, pengCard := range player.PengCards {
	// 	card, _ := utils.CardToInt(*pengCard.Card)
	// 	cardsGroup := &room.CardsGroup{
	// 		Pid:   proto.Uint64(player.PalyerId),
	// 		Type:  room.CardsGroupType_CardGroupType_Peng.Enum(),
	// 		Cards: []uint32{uint32(*card)},
	// 	}
	// 	cardsGroupList = append(cardsGroupList, cardsGroup)
	// }
	// // 杠牌
	// var groupType *room.CardsGroupType
	// for _, gangCard := range player.GangCards {
	// 	if gangCard.Type == server_pb.GangType_gang_angang {
	// 		groupType = room.CardsGroupType_CardGroupType_AnGang.Enum()
	// 	}
	// 	if gangCard.Type == server_pb.GangType_gang_minggang {
	// 		groupType = room.CardsGroupType_CardGroupType_MingGang.Enum()
	// 	}
	// 	if gangCard.Type == server_pb.GangType_gang_bugang {
	// 		groupType = room.CardsGroupType_CardGroupType_BuGang.Enum()
	// 	}
	// 	card, _ := utils.CardToInt(*gangCard.Card)
	// 	cardsGroup := &room.CardsGroup{
	// 		Pid:   proto.Uint64(player.PalyerId),
	// 		Type:  groupType,
	// 		Cards: []uint32{uint32(*card)},
	// 	}
	// 	cardsGroupList = append(cardsGroupList, cardsGroup)
	// }
	// // 手牌
	// handCards, _ := utils.CardsToInt(player.HandCards)
	// cards := make([]uint32, 0)
	// for _, handCard := range handCards {
	// 	cards = append(cards, uint32(handCard))
	// }
	// cardsGroup := &room.CardsGroup{
	// 	Pid:   proto.Uint64(player.PalyerId),
	// 	Type:  room.CardsGroupType_CardGroupType_Hand.Enum(),
	// 	Cards: cards,
	// }
	// cardsGroupList = append(cardsGroupList, cardsGroup)
	// return cardsGroupList
	return nil
}
