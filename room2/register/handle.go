package registers

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
	"steve/room2/util"
	player2 "steve/room2/player"
	"github.com/Sirupsen/logrus"
	modelmanager "steve/room2/models"
)

func HandleRoomChatReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskChatReq) (ret []exchanger.ResponseMsg){
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		return
	}
	modelmanager.GetModelManager().GetChatModel(player.GetDesk().GetUid()).RoomChatMsgReq(player,header,req)
	return
}

// HandleRoomDeskQuitReq 处理玩家退出桌面请求
// 失败先不回复
func HandleRoomDeskQuitReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskQuitReq) (rspMsg []exchanger.ResponseMsg) {
	response := room.RoomDeskQuitRsp{
		UserData: proto.Uint32(req.GetUserData()),
		ErrCode:  room.RoomError_DESK_NO_GAME_PLAYING.Enum(),
	}
	defer util.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_DESK_QUIT_RSP, &response)

	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if player == nil {
		return
	}
	desk := player.GetDesk()
	response.ErrCode = room.RoomError_SUCCESS.Enum()
	modelmanager.GetModelManager().GetPlayerModel(desk.GetUid()).PlayerQuit(player)
	player.SetTuoguan(true,true)
	return
}

// HandleResumeGameReq 恢复对局请求
func HandleResumeGameReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomCancelTuoGuanReq) (ret []exchanger.ResponseMsg) {
	player := player2.GetPlayerMgr().GetPlayer(playerID)
	if !(player==nil) {
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

	modelmanager.GetModelManager().GetPlayerModel(player.GetDesk().GetUid()).PlayerEnter(player)
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