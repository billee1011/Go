package majong

import (
	msgid "steve/client_pb/msgId"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// BroadCastMessage 广播通知牌桌所有玩家
func BroadCastMessage(desk interfaces.Desk, msgid msgid.MsgID, message proto.Message) {
	clientIDs := []uint64{}

	for _, player := range desk.GetDeskPlayers() {
		p := global.GetPlayerMgr().GetPlayer(player.GetPlayerID())
		if !player.IsQuit() {
			clientIDs = append(clientIDs, p.GetClientID())
		}
	}
	ms := global.GetMessageSender()
	ms.BroadcastPackage(clientIDs, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid)}, message)

	logrus.WithFields(logrus.Fields{
		"deskId": desk.GetUID(),
		"msg":    message.String(),
	}).Debugln("room广播通知牌桌所有玩家")
}

// SendMessageToPlayers 通知消息给指定玩家
func SendMessageToPlayers(desk interfaces.Desk, playerIds []uint64, msgid msgid.MsgID, message proto.Message) {

	clientIds := make([]uint64, 0)
	for _, playerID := range playerIds {
		clientID := global.GetPlayerMgr().GetPlayer(playerID).GetClientID()
		clientIds = append(clientIds, clientID)
	}

	ms := global.GetMessageSender()

	ms.BroadcastPackage(clientIds, &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid)}, message)

	logrus.WithFields(logrus.Fields{
		"deskId":    desk.GetUID(),
		"playerIds": playerIds,
		"msg":       message.String(),
	}).Debugln("room通知消息给指定玩家")
}
