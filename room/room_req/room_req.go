package roomreq

import (
	"steve/client_pb/msgId"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
)

// RegisterRoomReqHandlers 注册牌桌请求处理函数
func RegisterRoomReqHandlers(e exchanger.Exchanger) {
	roomReqs := []msgid.MsgID{
		// TODO: 添加所有房间请求消息
		msgid.MsgID_room_huansanzhang_req,
	}

	for _, msg := range roomReqs {
		e.RegisterHandle(uint32(msg), handleRoomReq)
	}
}

func handleRoomReq(clientID uint64, header *steve_proto_gaterpc.Header, body []byte) (rspMsg []exchanger.ResponseMsg) {
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayerByClientID(clientID)
	if player == nil {
		return
	}
	deskMgr := global.GetDeskMgr()
	deskMgr.HandlePlayerRequest(player.GetID(), header, body)
	return []exchanger.ResponseMsg{}
}
