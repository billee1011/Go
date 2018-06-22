package desks

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type joinApplyManager struct {
	applyChannel chan uint64
	applyXueZhan chan uint64
}

var gJoinApplyMgr *joinApplyManager
var once sync.Once

func getJoinApplyMgr() *joinApplyManager {
	once.Do(initApplyMgr)
	return gJoinApplyMgr
}

func initApplyMgr() {
	gJoinApplyMgr = newApplyMgr(true)
}

func newApplyMgr(runChecker bool) *joinApplyManager {
	mgr := &joinApplyManager{
		applyChannel: make(chan uint64, 1024),
		applyXueZhan: make(chan uint64, 1024),
	}
	if runChecker {
		go mgr.checkMatch()
	}
	return mgr
}

func (jam *joinApplyManager) getApplyChannel(gameID room.GameId) chan uint64 {
	switch gameID {
	case room.GameId_GAMEID_XUELIU:
		return jam.applyChannel
	case room.GameId_GAMEID_XUEZHAN:
		return jam.applyXueZhan
	default:
		return nil
	}
}

func (jam *joinApplyManager) joinPlayer(playerID uint64, gameID room.GameId) room.RoomError {
	// TODO: 检测玩家状态
	ch := jam.getApplyChannel(gameID)
	if ch == nil {
		return room.RoomError_FAILED
	}
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
	go jam.doApply(room.GameId_GAMEID_XUELIU)
	go jam.doApply(room.GameId_GAMEID_XUEZHAN)

}

func (jam *joinApplyManager) doApply(gameid room.GameId) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "checkMatch",
	})
	deskFactory := global.GetDeskFactory()
	deskMgr := global.GetDeskMgr()
	applyPlayers := make([]uint64, 0, 4)
	ch := jam.getApplyChannel(gameid)

	for {
		playerID, ok := <-ch
		logEntry.WithField("player_id", playerID).Debugln("accept player")
		if !ok {
			break
		}

		if jam.replicateApplyProc(applyPlayers, playerID) {
			continue
		}
		applyPlayers = append(applyPlayers, playerID)
		applyPlayers = jam.removeOfflinePlayer(applyPlayers)

		for len(applyPlayers) >= 4 {
			players := applyPlayers[:4]
			applyPlayers = applyPlayers[4:]
			infos := handleLocationInfos(players)
			//TODO:gameID这里是写死的,需要从协议里面拿
			result, err := deskFactory.CreateDesk(players, int(gameid), interfaces.CreateDeskOptions{}, infos)
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

func handleLocationInfos(players []uint64) map[uint64][]*room.GeographicalLocation {
	info := make(map[uint64][]*room.GeographicalLocation, 4)

	for _, pID := range players {
		infos, ok := locationInfos.Load(pID)
		if ok {
			info[pID] = infos.([]*room.GeographicalLocation)
			locationInfos.Delete(pID)
		}
	}
	return info
}

func (jam *joinApplyManager) replicateApplyProc(applyPlayers []uint64, newPlayerID uint64) bool {
	for _, playerID := range applyPlayers {
		if playerID == newPlayerID {
			header := &steve_proto_gaterpc.Header{
				MsgId: uint32(msgid.MsgID_ROOM_JOIN_DESK_RSP),
			}
			rsp := &room.RoomJoinDeskRsp{
				ErrCode: room.RoomError_DESK_ALREADY_APPLIED.Enum(),
			}
			SendMessageByPlayerID(playerID, header, rsp)
			return true
		}
	}
	return false
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

// 暂时用来存放地理信息
var locationInfos sync.Map

// HandleRoomJoinDeskReq 处理器玩家申请加入请求
// 	将玩家加入到申请列表中， 并且回复；
//TODO:RoomJoinDeskReq消息中需要加入gameID字段
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
	if _, exist := ExistInDesk(player.GetID()); exist {
		rsp.ErrCode = room.RoomError_DESK_GAME_PLAYING.Enum()
		return
	}
	rsp.ErrCode = getJoinApplyMgr().joinPlayer(player.GetID(), req.GetGameId()).Enum()
	locationInfos.Store(player.GetID(), req.GetLocation())
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

	rsp.ErrCode = getJoinApplyMgr().joinPlayer(player.GetID(), 1).Enum()
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
	playerID := player.GetID()
	deskMgr := global.GetDeskMgr()
	desk, err := deskMgr.GetRunDeskByPlayerID(playerID)
	if err != nil {
		return
	}
	desk.PlayerQuit(playerID)
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
func HandleResumeGameReq(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomCancelTuoGuanReq) (ret []exchanger.ResponseMsg) {
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayerByClientID(clientID)
	playerID := player.GetID()

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

// SendMessageByPlayerID 获取到playerID发送消息
func SendMessageByPlayerID(playerID uint64, head *steve_proto_gaterpc.Header, body proto.Message) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":   "sendMessageFromRoom",
		"newPlayerID": playerID,
		"head":        msgid.MsgID_name[int32(head.MsgId)],
	})
	playerMgr := global.GetPlayerMgr()
	p := playerMgr.GetPlayer(playerID)
	if p != nil {
		logEntry.Errorln("获取player失败")
		return
	}
	clientID := p.GetClientID()
	ms := global.GetMessageSender()
	err := ms.SendPackage(clientID, head, body)
	if err != nil {
		logEntry.WithError(err).Errorln("发送消息失败")
	}
}
