package models

import (
	"context"
	"fmt"
	"steve/room/contexts"
	deskpkg "steve/room/desk"
	"steve/room/player"
	"steve/server_pb/room_mgr"
	"sync"
	"sync/atomic"

	"steve/common/data/redis"
	"steve/entity/cache"

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
	desk, err := mgr.CreateDeskObj(req.GetDeskId(), length, playerIDs, int(req.GetGameId()), int32(req.GetLevelId()), robotLvs, req)
	if err != nil {
		entry.WithError(err).Errorln("创建桌子失败")
		rsp.ErrCode = roommgr.RoomError_FAILED // 默认是失败的
		return
	}
	deskID := desk.GetUid()

	rsp.ErrCode = roommgr.RoomError_SUCCESS

	modelMgr := GetModelManager()

	/* 	roomPlayers := []*room.RoomPlayerInfo{}
	   	deskPlayers := modelMgr.GetPlayerModel(deskID).GetDeskPlayers()
	   	for _, player := range deskPlayers {
	   		roomPlayer := TranslateToRoomPlayer(player)
	   		roomPlayers = append(roomPlayers, &roomPlayer)
	   	}

	   	ntf := room.RoomDeskCreatedNtf{
	   		GameId:  room.GameId(desk.GetGameId()).Enum(),
	   		Players: roomPlayers,
	   	}
	   	messageModel := modelMgr.GetMessageModel(deskID)
	   	if messageModel != nil {
	   		messageModel.BroadCastDeskMessage(nil, msgid.MsgID_ROOM_DESK_CREATED_NTF, &ntf, true)
	   	} */
	if err = modelMgr.StartDeskModel(deskID); err != nil {
		entry.WithError(err).Errorln("牌桌启动失败")
		rsp.ErrCode = roommgr.RoomError_FAILED // 默认是失败的
		return
	}

	entry.Infoln("牌桌创建成功")

	reportKey := cache.FmtGameReportKey(int(req.GetGameId()) ,int(desk.GetLevel())) //临时0
	redisCli := redis.GetRedisClient()
	redisCli.IncrBy(reportKey, int64(length))
	return
}

func createDeskContext(gameID int, players []uint64, zhuang int, fixzhuang bool) (interface{}, error) {
	switch gameID {
	case GameId_GAMEID_DOUDIZHU:
		return contexts.CreateInitDDZContext(players), nil
	default:
		deskcontext, err := contexts.CreateMajongContext(players, gameID, uint32(zhuang), fixzhuang)
		if err != nil {
			return nil, fmt.Errorf("创建麻将现场失败:%v", err)
		}
		return deskcontext, nil
	}
}

func createDeskSettler(gameID int) deskpkg.DeskSettler {
	switch gameID {
	case GameId_GAMEID_DOUDIZHU:
		{
			return nil
		}
	default:
		return NewMajongSettle()
	}
}

func (mgr *DeskManager) createDeskConfig(gameID int, players []uint64, req *roommgr.CreateDeskRequest) (deskpkg.DeskConfig, error) {
	context, err := createDeskContext(gameID, players, 0, false)
	if err != nil {
		return deskpkg.DeskConfig{}, fmt.Errorf("创建牌桌现场失败：%v", err)
	}

	var config deskpkg.DeskConfig

	switch gameID {
	case GameId_GAMEID_DOUDIZHU:
		config = deskpkg.NewDDZMDeskCreateConfig(context, len(players))
	default:
		config = deskpkg.NewMjDeskCreateConfig(context, NewMajongSettle(), len(players))
	}
	config.MinScore = req.GetMinCoin()
	config.MaxScore = req.GetMaxCoin()
	config.BaseScore = req.GetBaseCoin()
	return config, nil
}

// CreateDeskObj 创建桌子并初始化所有model
func (mgr *DeskManager) CreateDeskObj(deskID uint64, length int, players []uint64, gameID int, levelID int32, robotLvs []int, req *roommgr.CreateDeskRequest) (*deskpkg.Desk, error) {
	config, err := mgr.createDeskConfig(gameID, players, req)
	if err != nil {
		return nil, fmt.Errorf("create desk config failed:%v", err)
	}

	config.PlayerIds = players
	desk := deskpkg.NewDesk(deskID, gameID,levelID, players, &config)

	player.GetPlayerMgr().InitDeskData(players, 2, robotLvs)
	player.GetPlayerMgr().BindPlayerRoomAddr(players, gameID, int(levelID))
	GetModelManager().InitDeskModel(desk.GetUid(), desk.GetConfig().Models, &desk)
	return &desk, nil
}

func (mgr *DeskManager) allocDeskID() (uint64, error) {
	return atomic.AddUint64(&mgr.maxID, 1), nil
}
