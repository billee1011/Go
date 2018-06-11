package registers

import (
	"steve/client_pb/msgId"
	"steve/room/desks"
	"steve/room/login"
	"steve/structs/exchanger"
)

// RegisterHandlers 注册消息处理器
func RegisterHandlers(e exchanger.Exchanger) {
	registe := func(id msgid.MsgID, handler interface{}) {
		err := e.RegisterHandle(uint32(id), handler)
		if err != nil {
			panic(err)
		}
	}

	registe(msgid.MsgID_ROOM_LOGIN_REQ, login.HandleLogin)                     // 登录请求
	registe(msgid.MsgID_ROOM_VISITOR_LOGIN_REQ, login.HandleVisitorLogin)      // 游客登录请求
	registe(msgid.MsgID_ROOM_JOIN_DESK_REQ, desks.HandleRoomJoinDeskReq)       // 加入牌桌请求
	registe(msgid.MsgID_ROOM_DESK_QUIT_REQ, desks.HandleRoomDeskQuitReq)       // 退出牌桌请求
	registe(msgid.MsgID_ROOM_DESK_CONTINUE_REQ, desks.HandleRoomContinueReq)   // 续局请求
	registe(msgid.MsgID_ROOM_CANCEL_TUOGUAN_REQ, desks.HandleCancelTuoGuanReq) // 取消托管请求
	registe(msgid.MsgID_ROOM_RESUME_GAME_REQ, desks.HandleResumeGameReq)       // 恢复对局请求
	RegisterRoomReqHandlers(e)
}
