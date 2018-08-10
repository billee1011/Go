package registers

import (
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/room/models"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	modelmanager "steve/room/models"
	player2 "steve/room/player"
	"steve/room/util"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HandleRoomChatReq 处理玩家聊天请求
func HandleRoomChatReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskChatReq) (ret []exchanger.ResponseMsg) {
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		return
	}
	modelmanager.GetModelManager().GetChatModel(player.GetDesk().GetUid()).RoomChatMsgReq(player, header, req)
	return
}

// HandleRoomDeskQuitReq 处理玩家退出桌面请求
// 失败先不回复
func HandleRoomDeskQuitReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskQuitReq) (rspMsg []exchanger.ResponseMsg) {
	response := room.RoomDeskQuitRsp{
		UserData: proto.Uint32(req.GetUserData()),
		ErrCode:  room.RoomError_FAILED.Enum(),
	}
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		util.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_DESK_QUIT_RSP, &response)
		return
	}
	desk := player.GetDesk()
	if desk == nil {
		util.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_DESK_QUIT_RSP, &response)
		return
	}
	modelmanager.GetModelManager().GetPlayerModel(desk.GetUid()).PlayerQuit(player)
	response.ErrCode = room.RoomError_SUCCESS.Enum()
	util.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_DESK_QUIT_RSP, &response)
	return
}

func noGamePlaying() []exchanger.ResponseMsg {
	body := &room.RoomResumeGameRsp{
		ResumeRes: room.RoomError_DESK_NO_GAME_PLAYING.Enum(),
	}
	return []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_ROOM_RESUME_GAME_RSP),
			Body:  body,
		},
	}
}

// HandleResumeGameReq 恢复对局请求
func HandleResumeGameReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomResumeGameReq) (ret []exchanger.ResponseMsg) {
	entry := logrus.WithField("player_id", playerID)
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		entry.Debugln("玩家不存在")
		return noGamePlaying()
	}
	desk := player.GetDesk()
	if desk == nil {
		entry.Debugln("没有对应的牌桌")
		return noGamePlaying()
	}
	modelmanager.GetModelManager().GetPlayerModel(desk.GetUid()).PlayerEnter(player)
	return
}

// HandleCancelTuoGuanReq 处理取消托管请求
func HandleCancelTuoGuanReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomCancelTuoGuanReq) (ret []exchanger.ResponseMsg) {
	ret = []exchanger.ResponseMsg{}

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleCancelTuoGuanReq",
		"player_id": playerID,
	})
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		logEntry.Debugln("获取玩家失败")
		return
	}
	desk := player.GetDesk()
	if desk == nil {
		logEntry.Debugln("玩家不在房间中")
		return
	}
	player.SetTuoguan(false, true)
	return
}

// HandleTuoGuanReq 处理取消托管请求
func HandleTuoGuanReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomTuoGuanReq) (ret []exchanger.ResponseMsg) {
	ret = []exchanger.ResponseMsg{}

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleTuoGuanReq",
		"player_id": playerID,
	})
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		logEntry.Debugln("获取玩家失败")
		return
	}
	desk := player.GetDesk()
	if desk == nil {
		logEntry.Debugln("玩家不在房间中")
		return
	}
	player.SetTuoguan(req.GetTuoguan(), true)
	return
}

// HandleContinueReq 处理续局请求
func HandleContinueReq(playerID uint64, header *steve_proto_gaterpc.Header, req match.MatchDeskContinueReq) (ret []exchanger.ResponseMsg) {
	entry := logrus.WithField("player_id", playerID)

	response := &match.MatchDeskContinueRsp{
		ErrCode: proto.Int32(0),
		ErrDesc: proto.String("成功"),
	}
	ret = []exchanger.ResponseMsg{{
		MsgID: uint32(msgid.MsgID_MATCH_CONTINUE_RSP),
		Body:  response,
	}}

	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		entry.Debugln("获取玩家失败")
		return
	}
	desk := player.GetDesk()
	if desk == nil {
		entry.Debugln("玩家不在房间")
		return
	}
	continueModel := models.GetContinueModel(desk.GetUid())
	*response = continueModel.PushContinueRequest(playerID, &req)
	return
}
