package models

import (
	"context"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/room/contexts"
	deskpkg "steve/room/desk"
	"steve/room/player"
	"steve/server_pb/room_mgr"
	"sync"
	"sync/atomic"

	"github.com/Sirupsen/logrus"
)

type DeskManager struct {
	deskMap sync.Map // deskID -> *desk
	maxID   uint64
}

const (
	GameId_GAMEID_XUELIU   = 1
	GameId_GAMEID_XUEZHAN  = 2
	GameId_GAMEID_DOUDIZHU = 3
	GameId_GAMEID_ERRENMJ  = 4
)

var deskMgr *DeskManager

func init() {
	deskMgr = &DeskManager{maxID: 0}
}

func GetDeskMgr() *DeskManager {
	return deskMgr
}

func (mgr *DeskManager) CreateDesk(ctx context.Context, req *roommgr.CreateDeskRequest) (rsp *roommgr.CreateDeskResponse, err error) {
	entry := logrus.WithField("request", req.String())
	players := req.GetPlayers()
	// 回复match服的消息
	rsp = &roommgr.CreateDeskResponse{
		ErrCode: roommgr.RoomError_FAILED, // 默认是失败的
	}
	length := len(players)
	playerIDs := make([]uint64, length, length)
	robotLvs := make([]int, length, length)
	for _, pbPlayer := range players {
		seat := int(pbPlayer.GetSeat())
		if seat >= length || seat < 0 {
			entry.Errorln("座号错误")
			return
		}
		if playerIDs[seat] != 0 {
			entry.Errorln("座号重复")
			return
		}
		playerIDs[seat] = pbPlayer.GetPlayerId()
		robotLvs[seat] = int(pbPlayer.GetRobotLevel())
	}
	desk, err := mgr.CreateDeskObj(length, playerIDs, int(req.GetGameId()), robotLvs, req)
	if err != nil {
		rsp.ErrCode = roommgr.RoomError_FAILED // 默认是失败的
		return
	}

	rsp.ErrCode = roommgr.RoomError_SUCCESS

	roomPlayers := []*room.RoomPlayerInfo{}
	deskPlayers := GetModelManager().GetPlayerModel(desk.GetUid()).GetDeskPlayers()
	for _, player := range deskPlayers {
		roomPlayer := TranslateToRoomPlayer(player)
		roomPlayers = append(roomPlayers, &roomPlayer)
	}
	ntf := room.RoomDeskCreatedNtf{
		GameId:  room.GameId(desk.GetGameId()).Enum(),
		Players: roomPlayers,
	}

	GetModelManager().GetMessageModel(desk.GetUid()).BroadCastDeskMessage(nil, msgid.MsgID_ROOM_DESK_CREATED_NTF, &ntf, true)
	return
}

// CreateDeskObj 创建桌子并初始化所有model
func (mgr *DeskManager) CreateDeskObj(length int, players []uint64, gameID int, robotLvs []int, req *roommgr.CreateDeskRequest) (*deskpkg.Desk, error) {
	var config deskpkg.DeskConfig
	var context interface{}
	id, _ := mgr.allocDeskID()
	desk := deskpkg.NewDesk(id, gameID, players[:], &config)
	playerSli := players[:]
	var err error
	switch gameID {
	case GameId_GAMEID_DOUDIZHU:
		context = contexts.CreateInitDDZContext(playerSli)
		config = deskpkg.NewDDZMDeskCreateConfig(context, length)
	default:
		context, err = contexts.CreateMajongContext(playerSli, gameID, req.GetBankerSeat(), req.GetFixBanker())
		config = deskpkg.NewMjDeskCreateConfig(context, NewMajongSettle(), length)
	}
	if err != nil {
		return nil, err
	}
	desk.GetConfig().Context = context
	desk.GetConfig().PlayerIds = players
	player.GetPlayerMgr().InitDeskData(players, 2, robotLvs)
	player.GetPlayerMgr().BindPlayerRoomAddr(players, gameID)
	GetModelManager().InitDeskModel(desk.GetUid(), desk.GetConfig().Models, &desk)
	return &desk, nil
}

func (mgr *DeskManager) allocDeskID() (uint64, error) {
	return atomic.AddUint64(&mgr.maxID, 1), nil
}
