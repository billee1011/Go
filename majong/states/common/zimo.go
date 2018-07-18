package common

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// ZimoState 自摸状态
// 进入状态时，执行自摸动作，并广播给玩家
// 自摸完成事件，进入下家摸牌状态
type ZimoState struct {
}

var _ interfaces.MajongState = new(ZimoState)

// ProcessEvent 处理事件
func (s *ZimoState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_zimo_finish {
		return majongpb.StateID_state_zimo_settle, nil
	}
	return majongpb.StateID_state_zimo, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *ZimoState) OnEntry(flow interfaces.MajongFlow) {
	s.doZimo(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_zimo_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *ZimoState) OnExit(flow interfaces.MajongFlow) {

}

// doZimo 执行自摸操作
func (s *ZimoState) doZimo(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ZimoState.doZimo",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	player, card, err := s.getZimoInfo(mjContext)
	if err != nil {
		logEntry.Errorln(err)
		return
	}
	mjContext.LastHuPlayers = []uint64{player.GetPalyerId()}
	huType := s.calcHuType(player.GetPalyerId(), flow)
	s.notifyHu(card, huType, player.GetPalyerId(), flow)
	gutils.SetNextZhuangIndex(mjContext.GetLastHuPlayers(), player.GetPalyerId(), mjContext)
	player.HandCards, _ = utils.RemoveCards(player.GetHandCards(), card, 1)
	AddHuCard(card, player, player.GetPalyerId(), huType, true)

	// 玩家胡状态
	player.XpState = player.GetXpState() | majongpb.XingPaiState_hu
}

// isAfterGang 判断是否为杠开
// 杠后摸牌则为杠开
func (s *ZimoState) isAfterGang(mjContext *majongpb.MajongContext) bool {
	return mjContext.GetMopaiType() == majongpb.MopaiType_MT_GANG
}

// calcHuType 计算胡牌类型
func (s *ZimoState) calcHuType(huPlayerID uint64, flow interfaces.MajongFlow) majongpb.HuType {
	mjContext := flow.GetMajongContext()
	afterGang := s.isAfterGang(mjContext)
	isLast := !utils.HasAvailableWallCards(flow)
	if afterGang && isLast {
		return majongpb.HuType_hu_gangshanghaidilao
	} else if afterGang {
		return majongpb.HuType_hu_gangkai
	} else if isLast {
		return majongpb.HuType_hu_haidilao
	}
	huPlayer := utils.GetMajongPlayer(huPlayerID, mjContext)
	if len(huPlayer.PengCards) == 0 && len(huPlayer.GangCards) == 0 && len(huPlayer.HuCards) == 0 {
		if huPlayer.ZixunCount == 1 && huPlayerID == mjContext.Players[mjContext.ZhuangjiaIndex].GetPalyerId() {
			return majongpb.HuType_hu_tianhu
		}
		if huPlayer.MopaiCount == 1 && huPlayerID != mjContext.Players[mjContext.ZhuangjiaIndex].GetPalyerId() {
			return majongpb.HuType_hu_dihu
		}
	}
	return majongpb.HuType_hu_zimo
}

// notifyHu 广播胡
func (s *ZimoState) notifyHu(card *majongpb.Card, huType majongpb.HuType, playerID uint64, flow interfaces.MajongFlow) {
	// mjContext := flow.GetMajongContext()
	rhuType := s.huType2RoomHuType(huType)
	body := room.RoomHuNtf{
		Players:      []uint64{playerID},
		FromPlayerId: proto.Uint64(playerID),
		Card:         proto.Uint32(uint32(utils.ServerCard2Number(card))),
		HuType:       rhuType.Enum(),
		RealPlayerId: proto.Uint64(playerID),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HU_NTF, &body)
}

// getZimoInfo 获取自摸信息
func (s *ZimoState) getZimoInfo(mjContext *majongpb.MajongContext) (player *majongpb.Player, card *majongpb.Card, err error) {
	playerID := mjContext.GetLastMopaiPlayer()
	players := mjContext.GetPlayers()
	player = utils.GetPlayerByID(players, playerID)

	// 没有上个摸牌的玩家，是为天胡， 取庄家作为胡牌玩家
	if player.GetZixunCount() == 1 && player.GetPalyerId() == mjContext.Players[int(mjContext.GetZhuangjiaIndex())].GetPalyerId() {
		xpOption := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId()))
		switch xpOption.TianhuCardType {
		case mjoption.MostTingsCard:
			_, card = utils.CalcTianHuCardNum(mjContext, playerID)
		case mjoption.RightCard:
			card = player.HandCards[len(player.GetHandCards())-1]
		case mjoption.MoCard:
			card = mjContext.GetLastMopaiCard()
		}
	} else {
		card = mjContext.GetLastMopaiCard()
	}
	mjContext.LastHuPlayers = []uint64{playerID}
	return
}

func (s *ZimoState) huType2RoomHuType(huType majongpb.HuType) room.HuType {
	return map[majongpb.HuType]room.HuType{
		majongpb.HuType_hu_zimo:              room.HuType_HT_ZIMO,
		majongpb.HuType_hu_gangkai:           room.HuType_HT_GANGKAI,
		majongpb.HuType_hu_haidilao:          room.HuType_HT_HAIDILAO,
		majongpb.HuType_hu_gangshanghaidilao: room.HuType_HT_GANGSHANGHAIDILAO,
		majongpb.HuType_hu_tianhu:            room.HuType_HT_TIANHU,
		majongpb.HuType_hu_dihu:              room.HuType_HT_DIHU,
	}[huType]
}
