package deskbase

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/golang/protobuf/proto"
)

// HandleRoomDeskQuitReq 处理玩家退出桌面请求
// 失败先不回复
func HandleRoomDeskQuitReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskQuitReq) (rspMsg []exchanger.ResponseMsg) {
	response := room.RoomDeskQuitRsp{
		UserData: proto.Uint32(req.GetUserData()),
		ErrCode:  room.RoomError_SUCCESS.Enum(),
	}
	defer facade.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_DESK_QUIT_RSP, &response)

	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayer(playerID)
	if player == nil {
		return
	}
	deskMgr := global.GetDeskMgr()
	desk, err := deskMgr.GetRunDeskByPlayerID(playerID)
	if err != nil {
		return
	}
	response.ErrCode = room.RoomError_SUCCESS.Enum()
	<-desk.PlayerQuit(playerID) // 等到退出完成后返回
	return
}

// ExistInDesk 是否在游戏中
func ExistInDesk(playerID uint64) (interfaces.Desk, bool) {
	deskMgr := global.GetDeskMgr()
	desk, _ := deskMgr.GetRunDeskByPlayerID(playerID)
	if desk == nil {
		return nil, false
	}
	return desk, true
}

// HandleResumeGameReq 恢复对局请求
func HandleResumeGameReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomResumeGameReq) (ret []exchanger.ResponseMsg) {
	desk, exist := ExistInDesk(playerID)
	if !exist {
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

	desk.PlayerEnter(playerID)
	return
}
