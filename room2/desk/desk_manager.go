package desk

import (
	"sync"
	"sync/atomic"
	"steve/room2/desk/contexts"
	"steve/server_pb/room_mgr"
	"context"
	"steve/client_pb/room"
	"steve/client_pb/msgid"
	"steve/room2/desk/models"
	"steve/room2/util"
	"steve/room2/desk/settle"
	"steve/room2/desk/player"
)

type DeskMgr struct {
	deskMap sync.Map // deskID -> desk
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
	var playerIds [length]uint64
	var robotLvs [length]int
	for index, pbPlayer := range players {
		playerIds[index] = pbPlayer.GetPlayerId()
		robotLvs[index] = int(pbPlayer.GetRobotLevel())
	}
	desk,err := mgr.CreateDeskObj(length, playerIds, int(req.GetGameId()), robotLvs)
	if err != nil{
		rsp.ErrCode = roommgr.RoomError_FAILED // 默认是失败的
		return
	}

	rsp.ErrCode = roommgr.RoomError_SUCCESS
	pbPlayers := []*room.RoomPlayerInfo{}
	//通知玩家
	for _, tempPlayer := range desk.GetModel(models.Player).(models.PlayerModel).GetDeskPlayers() {
		roomPlayer := util.TranslateToRoomPlayer(tempPlayer)
		pbPlayers = append(pbPlayers, &roomPlayer)
	}
	ntf := room.RoomDeskCreatedNtf{
		GameId:  room.GameId(desk.GetGameId()).Enum(),
		Players: pbPlayers,
	}
	desk.GetModel(models.Message).(models.MessageModel).BroadCastDeskMessage( nil, msgid.MsgID_ROOM_DESK_CREATED_NTF, &ntf, true)

	desk.Start()

	return
}

//创建桌子并初始化所有model
func (mgr DeskMgr) CreateDeskObj(len int, players [len]uint64, gameID int, robotLvs [len]int) (Desk,error) {
	var config DeskConfig
	var context interface{}
	id, _ := mgr.allocDeskID()
	desk := NewDesk(id, gameID,players[:], &config)
	var err error = nil
	var ctx interface{} = nil
	switch gameID {
	case GameId_GAMEID_DOUDIZHU:
		config = NewDDZMDeskCreateConfig(context, len)
	default:
		config = NewMjDeskCreateConfig(context, settle.NewMajongSettle(),len)
		ctx,err = contexts.CreateMajongContext(desk.GetDeskPlayerIDs(),gameID)
	}
	if err != nil{
		return desk,err
	}
	desk.GetConfig().Context = ctx
	player.GetRoomPlayerMgr().InitDeskData(len, players, gameID, robotLvs)
	desk.InitModel()
	return desk,nil
}

func (mgr DeskMgr) allocDeskID() (uint64, error) {
	return atomic.AddUint64(&mgr.maxID, 1), nil
}
