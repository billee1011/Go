package facade

import (
	msgid "steve/client_pb/msgId"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// BroadcaseMessage 将消息广播给牌桌所有玩家
func BroadcaseMessage(flow interfaces.MajongFlow, msgID msgid.MsgID, msg proto.Message) {
	mjContext := flow.GetMajongContext()
	players := []uint64{}

	for _, player := range mjContext.GetPlayers() {
		players = append(players, player.GetPalyerId())
	}
	flow.PushMessages(players, interfaces.ToClientMessage{
		MsgID: int(msgID),
		Msg:   msg,
	})
}

// CalculateCardValue 计算牌型倍数,根数
func CalculateCardValue(ctc interfaces.CardTypeCalculator, cardParams interfaces.CardCalcParams) (cardValue, genCount uint32) {
	types, gen := ctc.Calculate(cardParams)
	cardValue, genCount = ctc.CardTypeValue(cardParams.GameID, types, gen)
	return
}

// SettleGang 作杠结算
func SettleGang(factory interfaces.GameSettlerFactory, gameID int, params interfaces.GangSettleParams) *majongpb.SettleInfo {
	f := factory.CreateSettlerFactory(gameID)
	settler := f.CreateGangSettler()
	return settler.Settle(params)
}

// SettleHu 作胡结算
func SettleHu(factory interfaces.GameSettlerFactory, gameID int, params interfaces.HuSettleParams) []*majongpb.SettleInfo {
	f := factory.CreateSettlerFactory(gameID)
	settler := f.CreateHuSettler()
	return settler.Settle(params)
}

// SettleRound 作单局结算
func SettleRound(factory interfaces.GameSettlerFactory, gameID int, params interfaces.RoundSettleParams) ([]*majongpb.SettleInfo, []uint64) {
	f := factory.CreateSettlerFactory(gameID)
	settler := f.CreateRoundSettle()
	return settler.Settle(params)
}
