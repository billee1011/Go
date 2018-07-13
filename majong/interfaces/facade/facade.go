package facade

import (
	"steve/client_pb/msgid"
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
func CalculateCardValue(ctc interfaces.FantypeCalculator, context *majongpb.MajongContext, fanParams interfaces.FantypeParams) (cardValue uint64, gen, hua int) {
	types, gen, hua := ctc.Calculate(fanParams)
	cardValue = ctc.CardTypeValue(context, types, gen, hua)
	return
}
