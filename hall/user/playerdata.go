package user

import (
	"context"
	"fmt"
	"steve/entity/db"
	"steve/hall/data"
	"steve/server_pb/user"
	"time"

	"github.com/Sirupsen/logrus"
)

// playerDataService 实现 user.PlayerServer
type playerDataService struct{}

var _ user.PlayerDataServer = new(playerDataService)

// GetPlayerByAccount 根据账号获取玩家 ID
func (pds *playerDataService) GetPlayerByAccount(ctx context.Context, req *user.GetPlayerByAccountReq) (rsp *user.GetPlayerByAccountRsp, err error) {
	logrus.Debugln("GetPlayerByAccount req", *req)
	rsp, err = &user.GetPlayerByAccountRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
	}, nil

	accID := req.GetAccountId()
	exist, playerID, err := data.GetPlayerIDByAccountID(accID)
	if exist && err == nil {
		rsp.PlayerId, rsp.ErrCode = playerID, int32(user.ErrCode_EC_SUCCESS)
		return
	}

	var err2 error
	playerID, err2 = createPlayer(accID)
	if err2 != nil {
		logrus.WithField("account_id", accID).Errorln(err2)
		return
	}
	rsp.PlayerId, rsp.ErrCode = playerID, int32(user.ErrCode_EC_SUCCESS)
	return
}

// GetPlayerInfo 获取玩家基本信息
func (pds *playerDataService) GetPlayerInfo(ctx context.Context, req *user.GetPlayerInfoReq) (rsp *user.GetPlayerInfoRsp, err error) {
	logrus.Debugln("GetPlayerInfo req", *req)

	rsp, err = &user.GetPlayerInfoRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
	}, nil
	playerID := req.GetPlayerId()
	player, err := data.GetPlayerInfoByPlayerID(playerID)
	if err == nil {
		rsp.PlayerId, rsp.ErrCode = playerID, int32(user.ErrCode_EC_SUCCESS)
		rsp.NickName, rsp.Avatar = player.NickName, player.Avatar
	}
	return
}

// SetPlayerInfo 设置玩家信息
func (pds *playerDataService) UpdatePlayerInfo(ctx context.Context, req *user.UpdatePlayerInfoReq) (rsp *user.UpdatePlayerInfoRsp, err error) {
	logrus.Debugln("SetPlayerInfo req", *req)

	rsp, err = &user.UpdatePlayerInfoRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		Result:  false,
	}, nil

	exist, result, err := data.UpdatePlayerInfo(req.GetPlayerId(), req.GetNickName(), req.GetAvatar())
	rsp.Result = result
	if exist {
		rsp.ErrCode = int32(user.ErrCode_EC_SUCCESS)
	}
	return
}

// GetPlayerState 获取玩家状态
func (pds *playerDataService) GetPlayerState(ctx context.Context, req *user.GetPlayerStateReq) (rsp *user.GetPlayerStateRsp, err error) {
	logrus.Debugln("GetPlayerState req", *req)

	rsp, err = &user.GetPlayerStateRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		State:   user.PlayerState_PS_IDIE,
	}, nil
	state, err := data.GetPlayerState(req.GetPlayerId())
	if err == nil {
		rsp.State, rsp.ErrCode = user.PlayerState(state), int32(user.ErrCode_EC_SUCCESS)
	}
	return
}

// SetPlayerState 设置玩家状态
func (pds *playerDataService) UpdatePlayerState(ctx context.Context, req *user.UpdatePlayerStateReq) (rsp *user.UpdatePlayerStateRsp, err error) {
	logrus.Debugln("SetPlayerState req", *req)

	rsp, err = &user.UpdatePlayerStateRsp{
		ErrCode: int32(user.ErrCode_EC_FAIL),
		Result:  false,
	}, nil
	playerID := req.GetPlayerId()
	oldState := uint64(req.GetNewState())
	newState := uint64(req.GetNewState())
	serverType := int32(req.GetServerType())
	result, err := data.UpdatePlayerState(playerID, oldState, newState, serverType, req.GetServerAddr())
	if result && err == nil {
		rsp.Result, rsp.ErrCode = true, int32(user.ErrCode_EC_SUCCESS)
	}
	return
}

// GetGameListInfo 获取玩家游戏列表信息
func (pds *playerDataService) GetGameListInfo(ctx context.Context, req *user.GetGameListInfoReq) (rsp *user.GetGameListInfoRsp, err error) {
	logrus.Debugln("GetGameListInfo req", *req)

	rsp, err = &user.GetGameListInfoRsp{
		ErrCode:  int32(user.ErrCode_EC_FAIL),
		GameInfo: []*user.GameInfo{},
	}, nil
	rsp.GameInfo, err = data.GetGameInfoList()
	if err == nil {
		rsp.ErrCode = int32(user.ErrCode_EC_SUCCESS)
	}
	return
}

// createPlayer 创建玩家
func createPlayer(accID uint64) (uint64, error) {
	playerID := data.AllocPlayerID()

	if playerID == 0 {
		return 0, fmt.Errorf("分配玩家 ID 失败")
	}
	if err := data.InitPlayerData(db.TPlayer{
		Accountid:    int64(accID),
		Playerid:     int64(playerID),
		Type:         1,
		Channelid:    0,                                 // TODO ，渠道 ID
		Nickname:     fmt.Sprintf("player%d", playerID), // TODO,昵称
		Gender:       1,
		Avatar:       "", // TODO , 头像
		Provinceid:   0,  // TODO， 省ID
		Cityid:       0,  // TODO 市ID
		Name:         "", // TODO: 真实姓名
		Phone:        "", // TODO: 电话
		Idcard:       "", // TODO 身份证
		Iswhitelist:  0,
		Zipcode:      0,
		Shippingaddr: "",
		Status:       1,
		Remark:       "",
		Createtime:   time.Now(),
		Createby:     "",
		Updatetime:   time.Now(),
		Updateby:     "",
	}); err != nil {
		return 0, fmt.Errorf("初始化玩家(%d)数据失败: %v", playerID, err)
	}
	if err := data.InitPlayerCoin(db.TPlayerCurrency{
		Playerid:       int64(playerID),
		Coins:          10000,
		Ingots:         0,
		Keycards:       0,
		Obtainingots:   0,
		Obtainkeycards: 0,
		Costingots:     0,
		Costkeycards:   0,
		Remark:         "",
		Createtime:     time.Now(),
		Createby:       "",
		Updatetime:     time.Now(),
		Updateby:       "",
	}); err != nil {
		return playerID, fmt.Errorf("初始化玩家(%d)金币数据失败: %v", playerID, err)
	}
	if err := data.InitPlayerState(int64(playerID)); err != nil {
		return playerID, fmt.Errorf("初始化玩家(%d)状态失败: %v", playerID, err)
	}
	return playerID, nil
}
