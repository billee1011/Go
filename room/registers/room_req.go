package registers

import (
	"steve/client_pb/msgId"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
)

// RegisterRoomReqHandlers 注册牌桌请求处理函数
func RegisterRoomReqHandlers(e exchanger.Exchanger) {
	roomReqs := []msgid.MsgID{
		msgid.MsgID_ROOM_HUANSANZHANG_REQ,
		msgid.MsgID_ROOM_XINGPAI_ACTION_REQ,
		msgid.MsgID_ROOM_DINGQUE_REQ,
		msgid.MsgID_ROOM_CHUPAI_REQ,
		msgid.MsgID_ROOM_CARTOON_FINISH_REQ,
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
