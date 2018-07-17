package game

import (
	"context"
	"steve/server_pb/room_mgr"
)

// RoomService room房间RPC服务
type GameService struct {
}

// CreateDesk 创建牌桌
func (gs *GameService) HandleMatchRequest(ctx context.Context, req *roommgr.CreateDeskRequest) (rsp *roommgr.CreateDeskResponse, err error) {
	players := req.GetPlayers()

	rsp = &roommgr.CreateDeskResponse{
		ErrCode: roommgr.RoomError_SUCCESS,
	}

	playerIDs := []uint64{}
	for _, player := range players {
		playerIDs = append(playerIDs, player.GetPlayerId())
	}

	option := &DeskOption{}
	desk, err := DefaultDeskManager.CreateDesk(playerIDs, int(req.GetGameId()), option)
	if err != nil {
		rsp.ErrCode = roommgr.RoomError_FAILED
		return
	}

	rsp.ErrCode = roommgr.RoomError_SUCCESS
	DefaultDeskManager.RunDesk(desk.deskID)
	return
}
