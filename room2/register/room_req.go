package registers

import (
	"steve/client_pb/msgid"
	modelmanager "steve/room2/models"
	"steve/room2/player"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
)

// RegisterRoomReqHandlers 注册牌桌请求处理函数
func RegisterRoomReqHandlers(e exchanger.Exchanger) {
	roomReqs := []msgid.MsgID{
		msgid.MsgID_ROOM_HUANSANZHANG_REQ,
		msgid.MsgID_ROOM_XINGPAI_ACTION_REQ,
		msgid.MsgID_ROOM_DINGQUE_REQ,
		msgid.MsgID_ROOM_CHUPAI_REQ,
		msgid.MsgID_ROOM_CARTOON_FINISH_REQ,
		//斗地主请求
		msgid.MsgID_ROOM_DDZ_GRAB_LORD_REQ,
		msgid.MsgID_ROOM_DDZ_DOUBLE_REQ,
		msgid.MsgID_ROOM_DDZ_PLAY_CARD_REQ,
		msgid.MsgID_ROOM_DDZ_RESUME_REQ,
	}

	for _, msg := range roomReqs {
		e.RegisterHandle(uint32(msg), handleRoomReq)
	}
}

func handleRoomReq(playerID uint64, header *steve_proto_gaterpc.Header, body []byte) (rspMsg []exchanger.ResponseMsg) {
	player := player.GetPlayerMgr().GetPlayer(playerID)
	if player.GetDesk() == nil {
		logrus.WithField("player_id", playerID).Debugln("not found player desk")
		return
	}
	deskID := player.GetDesk().GetUid()
	modelManager := modelmanager.GetModelManager()
	requestModel := modelManager.GetRequestModel(deskID)
	requestModel.HandlePlayerRequest(playerID, header, body)
	return []exchanger.ResponseMsg{}
}
