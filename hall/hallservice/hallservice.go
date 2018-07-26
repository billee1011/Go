package hallservice

import (
	"context"
	"steve/entity/cache"
	"steve/hall/player"
	"steve/server_pb/hall"
)

// Hallservice 大厅服务
type Hallservice struct {
}

var defaultObject = new(Hallservice)
var _ hall.HallServiceServer = Default()

// Default 默认对象
func Default() *Hallservice {
	return defaultObject
}

// PlayerLogin 初始化玩家登录
func (hs *Hallservice) PlayerLogin(ctx context.Context, request *hall.PlayerLoginReq) (*hall.PlayerLoginRsp, error) {
	playerID, err := player.Login(request.GetAccountId())
	if err != nil {
		return nil, err
	}
	response := &hall.PlayerLoginRsp{
		PlayerId: playerID,
	}
	return response, nil
}

// GetPlayerInfo 获取玩家信息
func (hs *Hallservice) GetPlayerInfo(ctx context.Context, request *hall.GetPlayerInfoReq) (*hall.GetPlayeInfoRsp, error) {
	hallPlayer, err := player.GetPlayerInfo(request.GetPlayerId())
	if err != nil {
		return nil, err
	}
	response := &hall.GetPlayeInfoRsp{
		PlayerId:    hallPlayer.PlayerID,
		NickName:    hallPlayer.NickName,
		HeadImage:   hallPlayer.HeadImage,
		Coin:        hallPlayer.Coin,
		PlayerState: hall.PlayerState(hallPlayer.State),
	}
	return response, nil
}

// SetPlayerInfo 设置玩家信息
func (hs *Hallservice) SetPlayerInfo(ctx context.Context, request *hall.SetPlayerInfoReq) (*hall.SetPlayerInfoRsp, error) {
	hallPlayer := cache.HallPlayer{
		PlayerID:  request.GetPlayerId(),
		NickName:  request.GetNickName(),
		HeadImage: request.GetHeadImage(),
	}
	response := &hall.SetPlayerInfoRsp{
		Result: player.UpdatePlayerInfo(hallPlayer),
	}
	return response, nil
}

// GetPlayerState 获取玩家状态
func (hs *Hallservice) GetPlayerState(ctx context.Context, request *hall.GetPlayerStateReq) (*hall.GetPlayerStateRsp, error) {
	state, err := player.GetPlayerState(request.GetPlayerId())
	if err != nil {
		return nil, err
	}
	response := &hall.GetPlayerStateRsp{
		State: hall.PlayerState(int32(state)),
	}
	return response, nil
}

// SetPlayerState 设置玩家状态
func (hs *Hallservice) SetPlayerState(ctx context.Context, request *hall.SetPlayerStateReq) (*hall.SetPlayerStateRsp, error) {
	playerID := request.GetPlayerId()
	return &hall.SetPlayerStateRsp{
		Result: player.UpdatePlayerState(playerID, uint64(request.GetOldState()), uint64(request.GetNewState()), int32(request.GetServerType()), request.GetServerAddr()),
	}, nil
}

// GetGameListInfo 获取游戏列表信息
func (hs *Hallservice) GetGameListInfo(ctx context.Context, request *hall.GetGameListInfoReq) (*hall.GetGameListInfoRsp, error) {
	response := &hall.GetGameListInfoRsp{}
	return response, nil
}
