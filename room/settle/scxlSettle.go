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
	logEntry := logrus.WithFields(logrus.Fields{
		"name":        "Settle",
		"SettleInfos": mjContext.SettleInfos,
	})
	if len(mjContext.SettleInfos) != 0 {
		deskPlayers := desk.GetPlayers()
		for _, settleInfo := range mjContext.SettleInfos {
			if !s.handleSettle[settleInfo.Id] {
				playerCoin := make(playerCoin)
				billplayerInfos := make([]*room.BillPlayerInfo, 0)
				for i := 0; i < len(deskPlayers); i++ {
					pid := *deskPlayers[i].PlayerId
					score := settleInfo.Scores[pid]
					coin := int64(*deskPlayers[i].Coin)
					if score != 0 {
						billplayerInfo := &room.BillPlayerInfo{
							Pid:      deskPlayers[i].PlayerId,
							BillType: room.BillType(settleInfo.SettleType).Enum(),
						}
						if score < 0 && (-score) >= int64(*deskPlayers[i].Coin) {
							playerCoin[pid] = int64(*deskPlayers[i].Coin)
							billplayerInfo.Score = proto.Int64(int64(coin))
							deskPlayers[i].Coin = proto.Uint64(0)
						} else {
							playerCoin[pid] = score
							billplayerInfo.Score = proto.Int64(score)
							deskPlayers[i].Coin = proto.Uint64(uint64(coin + score))
						}
						global.GetPlayerMgr().GetPlayer(pid).SetCoin(*deskPlayers[i].Coin)
						s.settleMap[settleInfo.Id] = playerCoin
						s.roundScore[pid] = s.roundScore[pid] + playerCoin[pid]
						billplayerInfo.CurrentScore = proto.Int64(int64(*deskPlayers[i].Coin))
						billplayerInfos = append(billplayerInfos, billplayerInfo)
					}
				}
				s.handleSettle[settleInfo.Id] = true
				// 广播即时结算消息
				notifyDeskMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
					BillPlayersInfo: billplayerInfos,
				})
			}
		}
		if len(mjContext.RevertSettles) != 0 {
			billplayerInfos := make([]*room.BillPlayerInfo, 0)
			for i := 0; i < len(deskPlayers); i++ {
				pid := *deskPlayers[i].PlayerId
				coin := int64(*deskPlayers[i].Coin)
				billplayerInfo := &room.BillPlayerInfo{
					Pid:      deskPlayers[i].PlayerId,
					BillType: room.BillType_BILL_REFUND.Enum(),
					Score:    proto.Int64(0),
				}
				for _, revertSettle := range mjContext.RevertSettles {
					if score, ok := s.settleMap[revertSettle][pid]; ok && score != 0 {
						billplayerInfo.Score = proto.Int64(*billplayerInfo.Score + (-score))
						deskPlayers[i].Coin = proto.Uint64(uint64(int64(coin) + (-score)))
					}
				}
				global.GetPlayerMgr().GetPlayer(pid).SetCoin(*deskPlayers[i].Coin)
				billplayerInfo.CurrentScore = proto.Int64(int64(*deskPlayers[i].Coin))
				billplayerInfos = append(billplayerInfos, billplayerInfo)
			}
			// 广播即时结算消息
			notifyDeskMessage(desk, msgid.MsgID_ROOM_INSTANT_SETTLE, &room.RoomSettleInstantRsp{
				BillPlayersInfo: billplayerInfos,
			})
		}
	}
	logEntry.Debugln("room 结算")
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

func notifyDeskPlayerMessage(desk interfaces.Desk, playerID uint64, msgid msgid.MsgID, message proto.Message) {
	clientID := global.GetPlayerMgr().GetPlayer(playerID).GetClientID()

	head := &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid)}
	ms := global.GetMessageSender()

	ms.BroadcastPackage([]uint64{clientID}, head, message)
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
		//balanceRsp.Pid = proto.Uint64(roomPlayer.GetPlayerId())
		balanceRsp.BillPlayersInfo = s.getBillPlayerInfo(pid, mjContext)
		notifyDeskPlayerMessage(desk, pid, msgid.MsgID_ROOM_ROUND_SETTLE, balanceRsp)
	}
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
			billPlayerInfo.CardsGroup = utils.GetCardsGroup(utils.GetPlayerByID(context.Players, playerID))
		}
		billPlayerInfos = append(billPlayerInfos, billPlayerInfo)
	}
	return billPlayerInfos
}
