package facade

import (
	msgid "steve/client_pb/msgId"
	"steve/majong/interfaces"

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
