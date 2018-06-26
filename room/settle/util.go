package settle

import (
	msgid "steve/client_pb/msgId"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	majongpb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// NotifyMessage 将消息广播给牌桌所有玩家
func NotifyMessage(desk interfaces.Desk, msgid msgid.MsgID, message proto.Message) {
	players := desk.GetDeskPlayers()
	clientIDs := []uint64{}

	playerMgr := global.GetPlayerMgr()
	for _, player := range players {
		playerID := player.GetPlayerID()
		p := playerMgr.GetPlayer(playerID)
		if p != nil && !player.IsQuit() {
			clientIDs = append(clientIDs, p.GetClientID())
		}
	}
	head := &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid)}
	ms := global.GetMessageSender()

	logrus.WithFields(logrus.Fields{
		"msg": message.String(),
	}).Debugln("room消息通知desk")

	ms.BroadcastPackage(clientIDs, head, message)
}

// NotifyPlayersMessage 将消息广播给牌桌指定playerIds[]中的玩家
func NotifyPlayersMessage(desk interfaces.Desk, playerIds []uint64, msgid msgid.MsgID, message proto.Message) {

	head := &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid)}
	ms := global.GetMessageSender()

	clientIds := make([]uint64, 0)
	for _, playerID := range playerIds {
		clientID := global.GetPlayerMgr().GetPlayer(playerID).GetClientID()
		clientIds = append(clientIds, clientID)
	}

	logrus.WithFields(logrus.Fields{
		"msg":       message.String(),
		"playerIds": playerIds,
	}).Debugln("room消息通知Player")

	ms.BroadcastPackage(clientIds, head, message)
}

// GetDeskPlayer 获取指定id的room Player
func GetDeskPlayer(deskPlayers []interfaces.DeskPlayer, pid uint64) interfaces.DeskPlayer {
	for _, p := range deskPlayers {
		if p.GetPlayerID() == pid {
			return p
		}
	}
	return nil
}

// GetSettleInfoBySid 根据settleID获取对应settleInfo的下标index
func GetSettleInfoBySid(settleInfos []*majongpb.SettleInfo, ID uint64) int {
	for index, s := range settleInfos {
		if s.Id == ID {
			return index
		}
	}
	return -1
}
