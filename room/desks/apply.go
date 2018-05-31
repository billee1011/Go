package desks

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
)

type joinApplyManager struct {
	applyChannel chan uint64
}

var gJoinApplyMgr = newApplyMgr(true)

func newApplyMgr(runChecker bool) *joinApplyManager {
	mgr := &joinApplyManager{
		applyChannel: make(chan uint64, 1024),
	}
	if runChecker {
		go mgr.checkMatch()
	}
	return mgr
}

func (jam *joinApplyManager) getApplyChannel() chan uint64 {
	return jam.applyChannel
}

func (jam *joinApplyManager) joinPlayer(playerID uint64) room.RoomError {
	// TODO: 检测玩家状态
	ch := jam.getApplyChannel()
	ch <- playerID
	return room.RoomError_SUCCESS
}

func (jam *joinApplyManager) removeOfflinePlayer(playerIDs []uint64) []uint64 {
	result := make([]uint64, 0, len(playerIDs))
	playerMgr := global.GetPlayerMgr()
	for _, playerID := range playerIDs {
		player := playerMgr.GetPlayer(playerID)
		if player != nil && player.GetClientID() != 0 {
			result = append(result, playerID)
		} else {
			logrus.WithField("player_id", playerID).Debugln("玩家不在线，移除")
		}
	}
	return result
}

func (jam *joinApplyManager) checkMatch() {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "checkMatch",
	})
	deskFactory := global.GetDeskFactory()
	deskMgr := global.GetDeskMgr()
	applyPlayers := make([]uint64, 0, 4)

	ch := jam.getApplyChannel()

	for {
		playerID, ok := <-ch
		logEntry.WithField("player_id", playerID).Debugln("accept player")
		if !ok {
			break
		}
		applyPlayers = append(applyPlayers, playerID)
		applyPlayers = jam.removeOfflinePlayer(applyPlayers)

		for len(applyPlayers) >= 4 {
			players := applyPlayers[:4]
			applyPlayers = applyPlayers[4:]
			result, err := deskFactory.CreateDesk(players, 1, interfaces.CreateDeskOptions{})
			if err != nil {
				logEntry.WithFields(
					logrus.Fields{
						"players": players,
						"result":  result,
					},
				).WithError(err).Errorln("创建房间失败")
				continue
			}
			notifyDeskCreate(result.Desk)
			deskMgr.RunDesk(result.Desk)
		}
	}
}

func notifyDeskCreate(desk interfaces.Desk) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "notifyDeskCreate",
	})
	players := desk.GetPlayers()
	clientIDs := []uint64{}

	playerMgr := global.GetPlayerMgr()
	for _, player := range players {
		playerID := player.GetPlayerId()
		p := playerMgr.GetPlayer(playerID)
		if p != nil {
			clientIDs = append(clientIDs, p.GetClientID())
		}
	}
	ntf := room.RoomDeskCreatedNtf{
		Players: desk.GetPlayers(),
	}
	head := &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid.MsgID_ROOM_DESK_CREATED_NTF)}
	ms := global.GetMessageSender()

	ms.BroadcastPackage(clientIDs, head, &ntf)
	logEntry.WithField("ntf_context", ntf).Debugln("广播创建房间")
}

// HandleRoomJoinDeskReq 处理器玩家申请加入请求
// 	将玩家加入到申请列表中， 并且回复；
func HandleRoomJoinDeskReq(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomJoinDeskReq) (rspMsg []exchanger.ResponseMsg) {
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayerByClientID(clientID)

	rsp := &room.RoomJoinDeskRsp{
		ErrCode: room.RoomError_SUCCESS.Enum(),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_ROOM_JOIN_DESK_RSP),
			Body:  rsp,
		},
	}

	if player == nil {
		rsp.ErrCode = room.RoomError_NOT_LOGIN.Enum()
		return
	}
	if ExsitInDesk(player.GetID()) {
		rsp.ErrCode = room.RoomError_desk_already_applied.Enum()
		return
	}
	rsp.ErrCode = gJoinApplyMgr.joinPlayer(player.GetID()).Enum()
	return
}

// HandleRoomContinueReq 玩家申请续局
func HandleRoomContinueReq(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskContinueReq) (rspMsg []exchanger.ResponseMsg) {
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayerByClientID(clientID)
	rsp := &room.RoomDeskContinueRsp{
		ErrCode: room.RoomError_SUCCESS.Enum(),
	}
	rspMsg = []exchanger.ResponseMsg{
		exchanger.ResponseMsg{
			MsgID: uint32(msgid.MsgID_ROOM_DESK_CONTINUE_RSP),
			Body:  rsp,
		},
	}

	if player == nil {
		rsp.ErrCode = room.RoomError_NOT_LOGIN.Enum()
		return
	}
	if ExsitInDesk(player.GetID()) {
		rsp.ErrCode = room.RoomError_desk_already_applied.Enum()
		return
	}
	rsp.ErrCode = gJoinApplyMgr.joinPlayer(player.GetID()).Enum()
	return
}

// HandleRoomDeskQuitReq 处理玩家退出桌面请求
// 失败先不回复
func HandleRoomDeskQuitReq(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomDeskQuitReq) (rspMsg []exchanger.ResponseMsg) {
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayerByClientID(clientID)
	if player == nil {
		return
	}
	deskMgr := global.GetDeskMgr()
	desk, err := deskMgr.GetRunDeskByPlayerID(player.GetID())
	if err != nil {
		return
	}
	err = desk.Stop()
	if err != nil {
		return
	}
	return
}

// ExsitInDesk 是否在游戏中
func ExsitInDesk(playerID uint64) bool {
	deskMgr := global.GetDeskMgr()
	desk, _ := deskMgr.GetRunDeskByPlayerID(playerID)
	if desk == nil {
		return false
	}
	return true
}
