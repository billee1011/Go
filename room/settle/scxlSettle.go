package settle

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/utils"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	majongpb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// scxlSettle 血流麻将结算
type scxlSettle struct {
	// 每条setttleInfo中每个玩家实际输赢分 key:settleId value:playerCoin
	settleMap map[uint64]playerCoin
	// 汇总setttleInfo中每个玩家输赢总分 key:playerID value:score
	roundScore map[uint64]int64
	// setttleInfo处理情况 		key:settleId value:true为已处理，false为未处理
	handleSettle map[uint64]bool
}

// newScxlSettle 创建四川血流结算
func newScxlSettle() *scxlSettle {
	return &scxlSettle{
		settleMap:    make(map[uint64]playerCoin),
		handleSettle: make(map[uint64]bool),
		roundScore:   make(map[uint64]int64),
	}
}

// playerCoin 玩家实际输赢分   key:playerID value:score
type playerCoin map[uint64]int64

// Settle 结算信息扣分并通知客户端
func (s *scxlSettle) Settle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	if len(mjContext.SettleInfos) != 0 {
		for _, settleInfo := range mjContext.SettleInfos {
			if !s.handleSettle[settleInfo.Id] {
				playerCoin := make(playerCoin)
				billplayerInfos := make([]*room.BillPlayerInfo, 0)
				for _, player := range desk.GetPlayers() {
					pid := *player.PlayerId
					billplayerInfo := &room.BillPlayerInfo{
						Pid:      player.PlayerId,
						BillType: room.BillType(settleInfo.SettleType).Enum(),
					}
					score := settleInfo.Scores[pid]

					if score < 0 && (-score) >= int64(*player.Coin) {
						playerCoin[pid] = int64(*player.Coin)
						billplayerInfo.Score = proto.Int64(int64(*player.Coin))
						*player.Coin = 0
					} else {
						playerCoin[pid] = score
						billplayerInfo.Score = proto.Int64(score)
						*player.Coin = uint64(int64(*player.Coin) + score)
					}

					s.settleMap[settleInfo.Id] = playerCoin
					s.roundScore[pid] = s.roundScore[pid] + playerCoin[pid]
					billplayerInfo.CurrentScore = proto.Int64(int64(*player.Coin))
					billplayerInfos = append(billplayerInfos, billplayerInfo)
				}
				s.handleSettle[settleInfo.Id] = true
				// 广播即时结算消息
				instantSettle := room.RoomSettleInstantRsp{
					BillPlayersInfo: billplayerInfos,
				}
				notifyDeskMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &instantSettle)
			}
		}
	}
}

func notifyDeskMessage(desk interfaces.Desk, msgid msgid.MsgID, message proto.Message) {
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
		MsgId: uint32(msgid)}
	ms := global.GetMessageSender()

	logrus.WithFields(logrus.Fields{
		"msg": message.String(),
	}).Debugln("通知立即结算")

	ms.BroadcastPackage(clientIDs, head, message)
}

// RoundSettle 单局结算信息
func (s *scxlSettle) RoundSettle(desk interfaces.Desk, mjContext majongpb.MajongContext) {
	balanceRsp := new(room.RoomBalanceInfoRsp)
	for _, roomPlayer := range desk.GetPlayers() {
		pid := *roomPlayer.PlayerId
		for _, settleInfo := range mjContext.SettleInfos {
			billDetail := s.getBillDetail(pid, settleInfo)
			if billDetail != nil {
				balanceRsp.BillDetail = append(balanceRsp.BillDetail, billDetail)
			}
		}
		balanceRsp.Pid = proto.Uint64(pid)
		balanceRsp.BillPlayersInfo = s.getBillPlayerInfo(pid, mjContext)
	}
	notifyDeskMessage(desk, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp)
}

// getBillDetail 单次结算详情，包括番型，分数，倍数，以及输赢玩家
func (s *scxlSettle) getBillDetail(palyerID uint64, settleInfo *majongpb.SettleInfo) *room.BillDetail {
	if settleInfo.Scores[palyerID] != 0 {
		billDetail := &room.BillDetail{
			SetleType: room.SettleType(settleInfo.SettleType).Enum(),
			HuType:    room.HuType(settleInfo.HuType).Enum(),
			FanValue:  proto.Uint32(settleInfo.CardValue),
			Score:     proto.Int64(s.settleMap[settleInfo.Id][palyerID]),
		}
		fanTypes := make([]room.FanType, 0)
		for _, cardType := range settleInfo.CardType {
			fanTypes = append(fanTypes, room.FanType(cardType))
		}
		billDetail.FanType = fanTypes
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

// getBillPlayerInfo 单局结算玩家详情,包括玩家自己牌型,输赢分数，以及其余每个玩家的输赢分数
func (s *scxlSettle) getBillPlayerInfo(playerID uint64, context majongpb.MajongContext) []*room.BillPlayerInfo {
	billPlayerInfos := make([]*room.BillPlayerInfo, 0)
	billPlayerInfo := new(room.BillPlayerInfo)
	for _, player := range context.Players {
		billPlayerInfo.Pid = proto.Uint64(player.GetPalyerId())
		billPlayerInfo.Score = proto.Int64(s.roundScore[player.GetPalyerId()])
		if player.PalyerId == playerID {
			billPlayerInfo.CardsGroup = getCardsGroup(utils.GetPlayerByID(context.Players, playerID))
		}
		billPlayerInfos = append(billPlayerInfos, billPlayerInfo)
	}
	return billPlayerInfos
}

// getCardsGroup 获取玩家牌组信息
func getCardsGroup(player *majongpb.Player) []*room.CardsGroup {
	cardsGroupList := make([]*room.CardsGroup, 0)
	// 碰牌
	for _, pengCard := range player.PengCards {
		card, _ := utils.CardToInt(*pengCard.Card)
		cardsGroup := &room.CardsGroup{
			Pid:   proto.Uint64(player.PalyerId),
			Type:  room.CardsGroupType_CGT_PENG.Enum(),
			Cards: []uint32{uint32(*card)},
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
		Type:  room.CardsGroupType_CGT_HAND.Enum(),
		Cards: cards,
	}
	cardsGroupList = append(cardsGroupList, cardsGroup)
	return cardsGroupList
}
