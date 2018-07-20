package models

import (
	"sync"
	"sync/atomic"
	"steve/room2/contexts"
	"steve/server_pb/room_mgr"
	"context"
	"steve/client_pb/room"
	"steve/client_pb/msgid"
	"steve/room2/player"
	deskpkg "steve/room2/desk"
)

type DeskMgr struct {
	deskMap sync.Map // deskID -> *desk
	maxID   uint64
}

const (
	GameId_GAMEID_XUELIU   = 1
	GameId_GAMEID_XUEZHAN  = 2
	GameId_GAMEID_DOUDIZHU = 3
	GameId_GAMEID_ERRENMJ  = 4
)

var deskMgr DeskMgr

func init() {
	deskMgr = DeskMgr{maxID: 0}
}

func GetDeskMgr() DeskMgr {
	return deskMgr
}

func (mgr DeskMgr) CreateDesk(ctx context.Context, req *roommgr.CreateDeskRequest) (rsp *roommgr.CreateDeskResponse, err error) {
	players := req.GetPlayers()
	// 回复match服的消息
	rsp = &roommgr.CreateDeskResponse{
		ErrCode: roommgr.RoomError_FAILED, // 默认是失败的
	}
	length := len(players)
	var playerIds []uint64
	var robotLvs []int
	for _, pbPlayer := range players {
		playerIds = append(playerIds,pbPlayer.GetPlayerId())
		robotLvs = append(robotLvs,int(pbPlayer.GetRobotLevel()))
	}
	desk,err := mgr.CreateDeskObj(length, playerIds[:], int(req.GetGameId()), robotLvs[:])
	if err != nil{
		rsp.ErrCode = roommgr.RoomError_FAILED // 默认是失败的
		return
	}

	rsp.ErrCode = roommgr.RoomError_SUCCESS
	pbPlayers := []*room.RoomPlayerInfo{}
	//通知玩家
	for _, tempPlayer := range GetModelManager().GetPlayerModel(desk.GetUid()).GetDeskPlayers() {
		roomPlayer := TranslateToRoomPlayer(tempPlayer)
		pbPlayers = append(pbPlayers, &roomPlayer)
	}
	ntf := room.RoomDeskCreatedNtf{
		GameId:  room.GameId(desk.GetGameId()).Enum(),
		Players: pbPlayers,
	}
	GetModelManager().GetMessageModel(desk.GetUid()).BroadCastDeskMessage( nil, msgid.MsgID_ROOM_DESK_CREATED_NTF, &ntf, true)

	desk.Start(nil)
	//mgr.deskMap.Store(desk.GetUid(),desk)
	return
}

//创建桌子并初始化所有model
func (mgr DeskMgr) CreateDeskObj(length int,players []uint64, gameID int, robotLvs []int) (*deskpkg.Desk,error) {
	var config deskpkg.DeskConfig
	var context interface{}
	id, _ := mgr.allocDeskID()
	desk := deskpkg.NewDesk(id, gameID,players[:], &config)
	playerSli := players[:]
	var err error = nil
	var ctx interface{} = nil
	switch gameID {
	case GameId_GAMEID_DOUDIZHU:
		config = deskpkg.NewDDZMDeskCreateConfig(context, length)
	default:
		config = deskpkg.NewMjDeskCreateConfig(context, NewMajongSettle(),length)
		ctx,err = contexts.CreateMajongContext(playerSli,gameID)
	}
	if err != nil{
		return nil,err
	}
	desk.GetConfig().Context = ctx
	player.GetPlayerMgr().InitDeskData(players, gameID, robotLvs)
	GetModelManager().InitDeskModel(desk.GetUid(),desk.GetConfig().Models,&desk)
	return &desk,nil
}

func (mgr DeskMgr) allocDeskID() (uint64, error) {
	return atomic.AddUint64(&mgr.maxID, 1), nil
}